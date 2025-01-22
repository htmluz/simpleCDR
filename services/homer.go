package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"radiusgo/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func convertTimestamp(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04", dateStr)
	if err != nil {
		return "", err
	}
	timestamp := strconv.FormatInt(t.UnixNano(), 10)
	return timestamp, nil
}

func BuildLokiQuery(filters *models.FilterParamsHomer) (string, error) {
	baseQuery := `{job="heplify-server"}`
	var conditions []string

	if filters.CalledPhone != "" {
		conditions = append(conditions, fmt.Sprintf(`|= "%s"`, filters.CalledPhone))
	}
	if filters.CallingPhone != "" {
		conditions = append(conditions, fmt.Sprintf(`|= "%s"`, filters.CallingPhone))
	}
	if filters.AnyPhone != "" {
		conditions = append(conditions, fmt.Sprintf(`|= "%s"`, filters.AnyPhone))
	}
	if filters.CallID != "" {
		baseQuery = fmt.Sprintf(`{job="heplify-server", call_id="%s"}`, filters.CallID)
	}
	if filters.Domain != "" {
		conditions = append(conditions, fmt.Sprintf(`|= "%s"`, filters.Domain))
	}

	if len(conditions) > 0 {
		return baseQuery + " " + strings.Join(conditions, " "), nil
	}
	return baseQuery, nil
}

func QueryLoki(filters *models.FilterParamsHomer) (*models.LokiResponse, error) {
	lokiURL := "http://10.90.0.58:3100/loki/api/v1/query_range"
	query, err := BuildLokiQuery(filters)
	if err != nil {
		return nil, fmt.Errorf("erro criando a query %v", err)
	}

	timeStart, err := convertTimestamp(filters.StartDate)
	if err != nil {
		return nil, err
	}
	timeEnd, err := convertTimestamp(filters.EndDate)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("start", timeStart)
	params.Add("end", timeEnd)

	resp, err := http.Get(lokiURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("erro fazendo a query pro loki %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro lendo a resposta do loki %v", err)
	}

	var lokiResp models.LokiResponse
	if err := json.Unmarshal(body, &lokiResp); err != nil {
		return nil, fmt.Errorf("erro parseando a resposta do loki %v", err)
	}
	return &lokiResp, nil
}

func ProcessLokiResponse(resp *models.LokiResponse) ([]*models.HomerCall, error) {
	calls := make(map[string]*models.HomerCall)

	for _, result := range resp.Data.Result {
		callID := result.Stream.CallID
		if callID == "" {
			continue
		}

		if _, exists := calls[callID]; !exists {
			calls[callID] = &models.HomerCall{
				CallID:     callID,
				Messages:   []models.LokiResult{},
				FromNumber: extractNumber(result.Stream.From),
				ToNumber:   extractNumber(result.Stream.To),
			}
		}
		calls[callID].Messages = append(calls[callID].Messages, result)

		timestamp, err := strconv.ParseInt(result.Values[0][0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %v", err)
		}
		timeStr := time.Unix(0, timestamp).Format(time.RFC3339)

		if calls[callID].StartTime == "" || timeStr < calls[callID].StartTime {
			calls[callID].StartTime = timeStr
		}
		if calls[callID].EndTime == "" || timeStr > calls[callID].EndTime {
			calls[callID].EndTime = timeStr
		}
	}
	result := make([]*models.HomerCall, 0, len(calls))
	for _, call := range calls {
		result = append(result, call)
	}
	return result, nil
}

func extractNumber(sipURI string) string {
	re := regexp.MustCompile(`sip:(\d+)@`)
	matches := re.FindStringSubmatch(sipURI)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
