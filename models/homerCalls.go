package models

type LokiResponse struct {
	Status string   `json:"status"`
	Data   LokiData `json:"data"`
}

type LokiData struct {
	ResultType string       `json:"resultType"`
	Result     []LokiResult `json:"result"`
}

type LokiResult struct {
	Stream LokiStream `json:"stream"`
	Values [][]string `json:"values"`
}

type LokiStream struct {
	CallID   string `json:"call_id"`
	From     string `json:"from"`
	To       string `json:"to"`
	DstIP    string `json:"dst_ip"`
	DstPort  string `json:"dst_port"`
	Hostname string `json:"hostname"`
	Job      string `json:"job"`
	Method   string `json:"method"`
	Node     string `json:"node"`
	Protocol string `json:"protocol"`
	Response string `json:"response"`
	SrcIP    string `json:"src_ip"`
	SrcPort  string `json:"src_port"`
	Type     string `json:"type"`
}

type HomerCall struct {
	CallID     string       `json:"call_id"`
	Messages   []LokiResult `json:"messages"`
	StartTime  string       `json:"start_time"`
	EndTime    string       `json:"end_time"`
	FromNumber string       `json:"from_number"`
	ToNumber   string       `json:"to_number"`
}

type FilterParamsHomer struct {
	StartDate    string `query:"startDate"`
	EndDate      string `query:"endDate"`
	CalledPhone  string `query:"calledPhone"`
	CallingPhone string `query:"callingPhone"`
	AnyPhone     string `query:"anyPhone"`
	Domain       string `query:"domain"`
	CallID       string `query:"callId"`
}
