package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"radiusgo/models"
	"sort"
	"strings"
	"time"
)

func GetBriefHomerCalls(filter *models.HomerFilterParams, db *sql.DB) ([]models.BriefHomerCall, error) {
	q := strings.Builder{}
	q.WriteString(`
		with cte as (
			select
				sid, data_header, data_header::json as data_header_json, create_date,
				row_number() over (partition by sid order by create_date asc) as row_num
			from hep_proto_1_call
			where 1=1
	`)
	args := []interface{}{}
	argPosition := 1

	if filter.StartTime != "" {
		q.WriteString(fmt.Sprintf(" and create_date >= $%d", argPosition))
		args = append(args, filter.StartTime)

		argPosition++
	}
	if filter.EndTime != "" {
		q.WriteString(fmt.Sprintf(" and create_date <= $%d", argPosition))
		args = append(args, filter.EndTime)
		argPosition++

	}

	q.WriteString(`
    )
    SELECT
        sid,
        data_header,
        MIN(create_date) AS start_time,
        MAX(create_date) AS end_time
    FROM
        cte
    WHERE
        row_num = 1
	`)

	if filter.RuriDomain != "" {
		q.WriteString(fmt.Sprintf(" AND data_header_json->>'ruri_domain' LIKE $%d", argPosition))
		args = append(args, "%"+filter.RuriDomain+"%")
		argPosition++
	}
	if filter.RuriUser != "" {
		q.WriteString(fmt.Sprintf(" AND data_header_json->>'ruri_user' LIKE $%d", argPosition))
		args = append(args, "%"+filter.RuriUser+"%")
		argPosition++
	}
	if filter.FromUser != "" {
		q.WriteString(fmt.Sprintf(" AND data_header_json->>'from_user' LIKE $%d", argPosition))
		args = append(args, "%"+filter.FromUser+"%")

		argPosition++
	}

	q.WriteString(" group by sid, data_header")
	log.Println(q)
	log.Println(q.String())
	rows, err := db.Query(q.String(), args...)
	if err != nil {
		log.Printf("Erro fazendo query no Homer %v", err)
		return nil, err
	}
	defer rows.Close()

	var calls []models.BriefHomerCall
	for rows.Next() {
		var call models.BriefHomerCall
		var dataHeader string
		err := rows.Scan(&call.Sid, &dataHeader, &call.StartTime, &call.EndTime)
		if err != nil {
			log.Printf("Erro parseando rows %v", err)
			return nil, err
		}
		call.CallInfo = parseDataHeader(dataHeader)
		calls = append(calls, call)
	}

	return calls, nil
}

func parseDataHeader(dataHeader string) models.BriefHomerCallInfo {
	var callInfo models.BriefHomerCallInfo
	err := json.Unmarshal([]byte(dataHeader), &callInfo)
	if err != nil {
		log.Printf("Erro parseando data_header %v", err)
	}
	return callInfo
}

func GetHomerMessages(callID string, db *sql.DB) models.HomerResponse {
	sipMsgs, err := getSIPMessages(callID, db)
	if err != nil {
		log.Fatalf("Erro msg sip %v", err)
	}

	rtcpMsgs, err := getRTCPMessages(callID, db)
	if err != nil {
		log.Fatalf("Erro msg rtcp %v", err)
	}

	allMessages := make([]interface{}, 0, len(sipMsgs)+len(rtcpMsgs))
	for _, msg := range sipMsgs {
		allMessages = append(allMessages, msg)
	}
	for _, msg := range rtcpMsgs {
		allMessages = append(allMessages, msg)
	}

	sort.Slice(allMessages, func(i, j int) bool {
		var dateI, dateJ time.Time

		switch msg := allMessages[i].(type) {
		case models.HomerSIPMessage:
			dateI = msg.CreateDate

		case models.HomerRTCPMessage:
			dateI = msg.CreateDate
		}

		switch msg := allMessages[j].(type) {
		case models.HomerSIPMessage:
			dateJ = msg.CreateDate
		case models.HomerRTCPMessage:
			dateJ = msg.CreateDate
		}

		return dateI.Before(dateJ)
	})

	messages := models.HomerResponse{
		CallID:   callID,
		Messages: allMessages,
	}
	return messages
}

func getSIPMessages(CallID string, db *sql.DB) ([]models.HomerSIPMessage, error) {
	q := strings.Builder{}
	q.WriteString(`
		select create_date, protocol_header, data_header, raw
		from hep_proto_1_call where sid = $1
	`)
	rows, err := db.Query(q.String(), CallID)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta SIP %w", err)
	}
	defer rows.Close()

	var messages []models.HomerSIPMessage
	for rows.Next() {
		var createDate time.Time
		var protocolHeader, dataHeader, raw string

		if err := rows.Scan(&createDate, &protocolHeader, &dataHeader, &raw); err != nil {
			return nil, fmt.Errorf("erro lendo sip message %w", err)
		}

		var protocolHeaderStruct models.HomerMessageProtocolHeader
		if err := json.Unmarshal([]byte(protocolHeader), &protocolHeaderStruct); err != nil {
			return nil, fmt.Errorf("erro parseando sip protocol header %w", err)
		}

		var dataHeaderStruct models.BriefHomerCallInfo
		if err := json.Unmarshal([]byte(dataHeader), &dataHeaderStruct); err != nil {
			return nil, fmt.Errorf("erro parseando sip data header %w", err)
		}

		messages = append(messages, models.HomerSIPMessage{
			CreateDate:     createDate,
			ProtocolHeader: protocolHeaderStruct,
			DataHeader:     dataHeaderStruct,
			Raw:            raw,
			Type:           "sip",
		})
	}
	return messages, nil
}

func getRTCPMessages(CallID string, db *sql.DB) ([]models.HomerRTCPMessage, error) {
	q := strings.Builder{}
	q.WriteString(`
		select create_date, protocol_header, data_header, raw
		from hep_proto_5_default where sid = $1
	`)
	rows, err := db.Query(q.String(), CallID)
	if err != nil {
		return nil, fmt.Errorf("Erro fazendo a consulta messages RTCP %w", err)
	}
	defer rows.Close()

	var messages []models.HomerRTCPMessage
	for rows.Next() {
		var createDate time.Time
		var protocolHeader, dataHeader, raw string

		if err := rows.Scan(&createDate, &protocolHeader, &dataHeader, &raw); err != nil {
			return nil, fmt.Errorf("erro lendo rtcp message %w", err)
		}

		var protocolHeaderStruct models.HomerMessageProtocolHeader
		if err := json.Unmarshal([]byte(protocolHeader), &protocolHeaderStruct); err != nil {
			return nil, fmt.Errorf("erro parseando rtcp protocol header %w", err)
		}

		var dataHeaderStruct models.HomerDataHeaderRTCP
		if err := json.Unmarshal([]byte(dataHeader), &dataHeaderStruct); err != nil {
			return nil, fmt.Errorf("erro parseando rtcp data header %w", err)
		}

		var rawStruct models.RTCPRaw
		if err := json.Unmarshal([]byte(raw), &rawStruct); err != nil {
			return nil, fmt.Errorf("erro parseando rtcp raw %w", err)
		}

		messages = append(messages, models.HomerRTCPMessage{
			CreateDate:     createDate,
			ProtocolHeader: protocolHeaderStruct,
			DataHeader:     dataHeaderStruct,
			Raw:            rawStruct,
			Type:           "rtcp",
		})
	}
	return messages, nil
}
