package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func RemoveSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func ConvertToTimestamp(value string) (string, error) {
	if value == "" {
		cTime := time.Now()
		z := cTime.AddDate(-2000, 0, 0)
		return z.Format("2006-01-02 15:04:05.000-07:00"), nil
	}
	layout := "15:04:05.000 -0700 Mon Jan 02 2006"
	parsed_time, err := time.Parse(layout, value)
	if err != nil {
		return "", fmt.Errorf("Erro convertendo tempo %s, %v", value, err)
	}
	return parsed_time.Format("2006-01-02 15:04:05.000-07:00"), nil
}

func StrToInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func HasIP(s string) (string, error) {
	if s == "" {
		return "0.0.0.0", nil
	}
	return s, nil
}
