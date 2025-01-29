package models

import "time"

type BriefHomerCall struct {
	Sid       string             `json:"sid"`
	CallInfo  BriefHomerCallInfo `json:"brief_call_info"`
	StartTime string             `json:"start_time"`
	EndTime   string             `json:"end_time"`
}

type BriefHomerCallInfo struct {
	RuriUser   string `json:"ruri_user"`
	RuriDomain string `json:"ruri_domain"`
	FromUser   string `json:"from_user"`
	FromTag    string `json:"from_tag"`
	ToUser     string `json:"to_user"`
	CallID     string `json:"callid"`
	Cseq       string `json:"cseq"`
	Method     string `json:"method"`
	UserAgent  string `json:"user_agent"`
}

type HomerFilterParams struct {
	RuriDomain string `query:"ruri_domain"`
	RuriUser   string `query:"ruri_user"`
	FromUser   string `query:"from_user"`
	CallID     string `query:"call_id"`
	StartTime  string `query:"start_time"`
	EndTime    string `query:"end_time"`
}

type HomerMessageProtocolHeader struct {
	ProtocolFamily int    `json:"protocolFamily"`
	Protocol       int    `json:"protocol"`
	SrcIP          string `json:"srcIp"`
	DstIP          string `json:"dstIp"`
	SrcPort        int    `json:"srcPort"`
	DstPort        int    `json:"dstPort"`
	TimeSeconds    int    `json:"timeSeconds"`
	TimeUSeconds   int    `json:"timeUseconds"`
	PayloadType    int    `json:"payloadType"`
	CaptureID      string `json:"captureId"`
	CorrelationID  string `json:"correlation_id"`
}

type HomerResponse struct {
	CallID   string         `json:"call_id"`
	Messages []HomerMessage `json:"messages"`
}

type HomerDataHeaderRTCP struct {
	Node  string `json:"node"`
	Proto string `json:"proto"`
}

type HomerSIPMessage struct {
	CreateDate     time.Time                  `json:"create_date"`
	ProtocolHeader HomerMessageProtocolHeader `json:"protocol_header"`
	DataHeader     BriefHomerCallInfo         `json:"data_header"`
	Raw            string                     `json:"raw"`
	Type           string                     `json:"type"`
}

type HomerRTCPMessage struct {
	CreateDate     time.Time                  `json:"create_date"`
	ProtocolHeader HomerMessageProtocolHeader `json:"protocol_header"`
	DataHeader     HomerDataHeaderRTCP        `json:"data_header"`
	Raw            RTCPRaw                    `json:"raw"`
	Type           string                     `json:"type"`
}

type RTCPFlow struct {
	Type     string             `json:"type"`
	SrcIP    string             `json:"src_ip"`
	DstIP    string             `json:"dst_ip"`
	Messages []HomerRTCPMessage `json:"messages"`
}

type RTCPRaw struct {
	SenderInformation RTCPSenderInformation `json:"sender_information"`
	Ssrc              int                   `json:"ssrc"`
	Type              int                   `json:"type"`
	ReportCount       int                   `json:"report_count"`
	ReportBlocks      []RTCPReportBlocks    `json:"report_blocks"`
	ReportBlocksXr    RTCPReportBlocksXr    `json:"report_blocks_xr"`
	SdesSsrc          int                   `json:"sdes_ssrc"`
}

type RTCPSenderInformation struct {
	NtpTimestampSec  int `json:"ntp_timestamp_sec"`
	NtpTimestampUSec int `json:"ntp_timestamp_usec"`
	RtpTimestamp     int `json:"rtp_timestamp"`
	Packets          int `json:"packets"`
	Octets           int `json:"octets"`
}

type RTCPReportBlocks struct {
	SourceSsrc   int `json:"source_ssrc"`
	FractionLost int `json:"fraction_lost"`
	PacketsLost  int `json:"packets_lost"`
	HighestSeqNo int `json:"highest_seq_no"`
	IaJitter     int `json:"ia_jitter"`
	Lsr          int `json:"lsr"`
	Dlsr         int `json:"dlsr"`
}

type RTCPReportBlocksXr struct {
	Type            int `json:"type"`
	Id              int `json:"id"`
	FractionLost    int `json:"fraction_lost"`
	FractionDiscard int `json:"fraction_discard"`
	BurstDensity    int `json:"burst_density"`
	GapDensity      int `json:"gap_density"`
	BurstDuration   int `json:"burst_duration"`
	GapDuration     int `json:"gap_duration"`
	RountTripDelay  int `json:"round_trip_delay"`
	EndSystemDelay  int `json:"end_system_delay"`
}

type HomerMessage interface {
	GetCreateDate() time.Time
	GetType() string
}

func (m HomerSIPMessage) GetCreateDate() time.Time {
	return m.CreateDate
}

func (m HomerSIPMessage) GetType() string {
	return m.Type
}

func (f RTCPFlow) GetCreateDate() time.Time {
	if len(f.Messages) > 0 {
		return f.Messages[0].CreateDate
	}
	return time.Time{}
}

func (f RTCPFlow) GetType() string {
	return "rtcp_flow"
}
