package models

type Bilhete struct {
	UserName            string `json:"User-Name"`
	AcctSessionID       string `json:"Acct-Session-Id"`
	CallingStationID    string `json:"Calling-Station-Id"`
	CalledStationID     string `json:"Called-Station-Id"`
	H323SetupTime       string `json:"h323-setup-time"`
	H323ConnectTime     string `json:"h323-connect-time"`
	H323DisconnectTime  string `json:"h323-disconnect-time"`
	NASIdentifier       string `json:"NAS-Identifier"`
	CiscoNASPort        string `json:"Cisco-NAS-Port"`
	H323CallOrigin      string `json:"h323-call-origin"`
	ReleaseSource       string `json:"release-source"`
	H323CallType        string `json:"h323-call-type"`
	CallID              string `json:"call-id"`
	AcctSessionTime     string `json:"Acct-Session-Time"`
	H323DisconnectCause string `json:"h323-disconnect-cause"`
	NASIPAddress        string `json:"NAS-IP-Address"`
	AcctStatusType      string `json:"Acct-Status-Type"`
	Protocol            string `json:"Protocol"`
	Codec               string `json:"Codec"`
	RemoteRTPPort       string `json:"Remote-RTP-Port"`
	RemoteRTPIp         string `json:"Remote-RTP-IP"`
	RemoteSIPPort       string `json:"Remote-SIP-Port"`
	RemoteSIPIp         string `json:"Remote-SIP-IP"`
	LocalRTPPort        string `json:"Local-RTP-Port"`
	LocalRTPIp          string `json:"Local-RTP-IP"`
	LocalSIPPort        string `json:"Local-SIP-Port"`
	LocalSIPIp          string `json:"Local-SIP-IP"`
	RingStart           string `json:"Ring-Start"`
	MosIngress          string `json:"MOS-Ingres"`
	MosEgress           string `json:"MOS-Egress"`
	GwName              string `json:"Gw-Name"`
}

type BilhetesResponse struct {
	Data        []Bilhete `json:"data"`
	Total       int       `json:"total"`
	CurrentPage int       `json:"currentPage"`
	PerPage     int       `json:"perPage"`
	TotalPages  int       `json:"totalPages"`
}

type FilterParams struct {
	StartDate    string `query:"startDate"`
	EndDate      string `query:"endDate"`
	CalledPhone  string `query:"calledPhone"`
	CallingPhone string `query:"callingPhone"`
	AnyPhone     string `query:"anyPhone"`
	NapA         string `query:"napA"`
	NapB         string `query:"napB"`
	DisconnCause string `query:"disconnCause"`
	CallID       string `query:"callId"`
	GatewayIP    string `query:"gatewayIp"`
	Codec        string `query:"codec"`
	Page         int    `query:"page"`
	PerPage      int    `query:"perPage"`
}
