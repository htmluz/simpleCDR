package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

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

func convertToTimestamp(value string) (string, error) {
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

func strToInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func hasIP(s string) (string, error) {
	if s == "" {
		return "0.0.0.0", nil
	}
	return s, nil
}

func insertBilhete(bilhete *Bilhete) error {
	h323SetupTime, e := convertToTimestamp(bilhete.H323SetupTime)
	if e != nil {
		return e
	}
	h323ConnectTime, e := convertToTimestamp(bilhete.H323ConnectTime)
	if e != nil {
		return e
	}
	h323DisconnectTime, e := convertToTimestamp(bilhete.H323DisconnectTime)
	if e != nil {
		return e
	}
	fRingStart, e := convertToTimestamp(bilhete.RingStart)
	if e != nil {
		return e
	}
	remoteRTPPort, e := strToInt(bilhete.RemoteRTPPort)
	if e != nil {
		return e
	}
	remoteSIPPort, e := strToInt(bilhete.RemoteSIPPort)
	if e != nil {
		return e
	}
	localRTPPort, e := strToInt(bilhete.LocalRTPPort)
	if e != nil {
		return e
	}
	localSIPPort, e := strToInt(bilhete.LocalSIPPort)
	if e != nil {
		return e
	}
	remoteRTPIP, e := hasIP(bilhete.RemoteRTPIp)
	if e != nil {
		return e
	}
	remoteSIPIP, e := hasIP(bilhete.RemoteSIPIp)
	if e != nil {
		return e
	}
	localRTPIP, e := hasIP(bilhete.LocalRTPIp)
	if e != nil {
		return e
	}
	localSIPIP, e := hasIP(bilhete.LocalSIPIp)
	if e != nil {
		return e
	}

	query := `
        INSERT INTO call_records (
            user_name, acct_session_id, calling_station_id, called_station_id,
            h323_setup_time, h323_connect_time, h323_disconnect_time,
            nas_identifier, cisco_nas_port, h323_call_origin,
            release_source, h323_call_type, call_id,
            acct_session_time, h323_disconnect_cause,
            nas_ip_address, acct_status_type,
            protocol, codec, remote_rtp_ip,
            remote_rtp_port, remote_sip_ip,
            remote_sip_port, local_rtp_ip,
						local_rtp_port, local_sip_ip,
	          local_sip_port, ring_start,
						mos_ingress, mos_egress
        ) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
					$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
				)
				ON CONFLICT (call_id)
				DO UPDATE SET
					user_name = EXCLUDED.user_name,
					acct_session_id = EXCLUDED.acct_session_id,
					calling_station_id = EXCLUDED.calling_station_id,
					called_station_id = EXCLUDED.called_station_id,
					h323_setup_time = EXCLUDED.h323_setup_time,
					h323_connect_time = EXCLUDED.h323_connect_time,
					h323_disconnect_time = EXCLUDED.h323_disconnect_time,
					nas_identifier = EXCLUDED.nas_identifier,
					cisco_nas_port = EXCLUDED.cisco_nas_port,
					h323_call_origin = EXCLUDED.h323_call_origin,
					release_source = EXCLUDED.release_source,
					h323_call_type = EXCLUDED.h323_call_type,
					call_id = EXCLUDED.call_id,
					acct_session_time = EXCLUDED.acct_session_time,
					h323_disconnect_cause = EXCLUDED.h323_disconnect_cause,
					nas_ip_address = EXCLUDED.nas_ip_address,
					acct_status_type = EXCLUDED.acct_status_type,
					protocol = EXCLUDED.protocol,
					codec = EXCLUDED.codec,
					remote_rtp_ip = EXCLUDED.remote_rtp_ip,
					remote_rtp_port = EXCLUDED.remote_rtp_port,
					remote_sip_ip = EXCLUDED.remote_sip_ip,
					remote_sip_port = EXCLUDED.remote_sip_port,
					local_rtp_ip = EXCLUDED.local_rtp_ip,
					local_rtp_port = EXCLUDED.local_rtp_port,
					local_sip_ip = EXCLUDED.local_sip_ip,
					local_sip_port = EXCLUDED.local_sip_port,
					ring_start = EXCLUDED.ring_start,
					mos_ingress = EXCLUDED.mos_ingress,
					mos_egress = EXCLUDED.mos_egress
        RETURNING id
  `
	var id int
	err := db.QueryRow(query,
		bilhete.UserName,
		bilhete.AcctSessionID,
		bilhete.CallingStationID,
		bilhete.CalledStationID,
		h323SetupTime,
		h323ConnectTime,
		h323DisconnectTime,
		bilhete.NASIdentifier,
		bilhete.CiscoNASPort,
		bilhete.H323CallOrigin,
		bilhete.ReleaseSource,
		bilhete.H323CallType,
		bilhete.CallID,
		bilhete.AcctSessionTime,
		bilhete.H323DisconnectCause,
		bilhete.NASIPAddress,
		bilhete.AcctStatusType,
		bilhete.Protocol,
		bilhete.Codec,
		remoteRTPIP,
		remoteRTPPort,
		remoteSIPIP,
		remoteSIPPort,
		localRTPIP,
		localRTPPort,
		localSIPIP,
		localSIPPort,
		fRingStart,
		bilhete.MosIngress,
		bilhete.MosEgress,
	).Scan(&id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Erro ao inserir bilhete %v", err)
	}
	return nil
}

func handlePostBilhete(c *fiber.Ctx) error {
	data := c.Body()

	bilhete := new(Bilhete)

	if err := json.Unmarshal(data, bilhete); err != nil {
		log.Fatal("Erro fazendo parsing do json")
	}

	err := insertBilhete(bilhete)
	if err != nil {
		log.Fatal("Erro ao inserir bilhete")
	}

	return c.SendStatus(200)
}

func debugFilters(filters *FilterParams, q string, args []interface{}) {
	log.Printf("\n=== Debug Filtros ===")
	log.Printf("StartDate: '%v'", filters.StartDate)
	log.Printf("EndDate: '%v'", filters.EndDate)
	log.Printf("PhoneNumber: '%v'", filters.AnyPhone)
	log.Printf("UserName: '%v'", filters.NapA)
	log.Printf("Page: %d", filters.Page)
	log.Printf("PerPage: %d", filters.PerPage)
	log.Printf("\n=== Query Construída ===")
	log.Printf("%s", q)
	log.Printf("\n=== Argumentos ===")
	for i, arg := range args {
		log.Printf("$%d: %v", i+1, arg)
	}
	log.Printf("\n==================")
}

func handleGetCleanupDays(c *fiber.Ctx) error {
	type GetDaysResponse struct {
		Days      int    `json:"days"`
		UpdatedAt string `json:"updated_at"`
	}
	var res GetDaysResponse
	q := `SELECT cleanup_days, updated_at FROM cleanup_config LIMIT 1;`
	if err := db.QueryRow(q).Scan(&res.Days, &res.UpdatedAt); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "Erro consultando os dias",
			"detail": err.Error(),
		})
	}
	return c.JSON(res)
}

func handleUpdateCleanupDays(c *fiber.Ctx) error {
	type UpdateDaysRequest struct {
		Days int `json:"days"`
	}
	var req UpdateDaysRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":  "Body inválido",
			"detail": err.Error(),
		})
	}
	if req.Days <= 30 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Valor deve ser maior que 30 dias.",
		})
	}
	q := `
		UPDATE cleanup_config
		SET cleanup_days = $1, updated_at = NOW()
		WHERE id = 1;
	`
	if _, err := db.Exec(q, req.Days); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":  "Erro ao atualizar cleanup_days",
			"detail": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "cleanup_days atualizado com sucesso",
		"days":    req.Days,
	})
}

func handleGetBilhetes(c *fiber.Ctx) error {
	filters := new(FilterParams)
	if err := c.QueryParser(filters); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":  "Parâmetros inválidos",
			"detail": err.Error(),
		})
	}
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 20
	}

	q := strings.Builder{}
	q.WriteString(`
		SELECT 
			user_name, acct_session_id, calling_station_id, called_station_id,
			h323_setup_time, h323_connect_time, h323_disconnect_time,
			nas_identifier, cisco_nas_port, h323_call_origin,
			release_source, h323_call_type, call_id,
			acct_session_time, h323_disconnect_cause,
			nas_ip_address, acct_status_type, protocol,
		  codec, remote_rtp_ip, remote_rtp_port,
		  remote_sip_ip, remote_sip_port,
		  local_rtp_ip, local_rtp_port,
		  local_sip_ip, local_sip_port,
		  ring_start, mos_ingress, mos_egress
		FROM call_records
		WHERE 1=1
	`)

	args := []interface{}{}
	argPosition := 1

	if filters.StartDate != "" {
		q.WriteString(fmt.Sprintf(" AND h323_setup_time >= $%d", argPosition))
		args = append(args, filters.StartDate)
		argPosition++
	}
	if filters.EndDate != "" {
		q.WriteString(fmt.Sprintf(" AND h323_setup_time <= $%d", argPosition))
		args = append(args, filters.EndDate)
		argPosition++
	}
	if filters.CalledPhone != "" {
		q.WriteString(fmt.Sprintf(" AND called_station_id ILIKE $%d", argPosition))
		args = append(args, "%"+filters.CalledPhone+"%")
		argPosition++
	}
	if filters.CallingPhone != "" {
		q.WriteString(fmt.Sprintf(" AND calling_station_id ILIKE $%d", argPosition))
		args = append(args, "%"+filters.CallingPhone+"%")
		argPosition++
	}
	if filters.AnyPhone != "" {
		q.WriteString(fmt.Sprintf(" AND (calling_station_id ILIKE $%d OR called_station_id ILIKE $%d)",
			argPosition, argPosition))
		args = append(args, "%"+filters.AnyPhone+"%")
		argPosition++
	}
	if filters.NapA != "" {
		q.WriteString(fmt.Sprintf(" AND user_name ILIKE $%d", argPosition))
		args = append(args, "%"+filters.NapA+"%")
		argPosition++
	}
	if filters.NapB != "" {
		q.WriteString(fmt.Sprintf(" AND cisco_nas_port ILIKE $%d", argPosition))
		args = append(args, "%"+filters.NapB+"%")
		argPosition++
	}
	if filters.DisconnCause != "" {
		q.WriteString(fmt.Sprintf(" AND h323_disconnect_cause ILIKE $%d", argPosition))
		args = append(args, "%"+filters.DisconnCause+"%")
		argPosition++
	}
	if filters.CallID != "" {
		q.WriteString(fmt.Sprintf(" AND call_id ILIKE $%d", argPosition))
		args = append(args, "%"+filters.CallID+"%")
		argPosition++
	}
	if filters.GatewayIP != "" {
		q.WriteString(fmt.Sprintf(" AND nas_ip_address ILIKE $%d", argPosition))
		args = append(args, "%"+filters.GatewayIP+"%")
		argPosition++
	}
	if filters.Codec != "" {
		q.WriteString(fmt.Sprintf(" AND codec ILIKE $%d", argPosition))
		args = append(args, "%"+filters.Codec+"%")
		argPosition++
	}

	countQ := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", q.String())
	var total int
	err := db.QueryRow(countQ, args...).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Erro ao contar os registros",
			"details": err.Error(),
		})
	}

	offset := (filters.Page - 1) * filters.PerPage
	q.WriteString(fmt.Sprintf(" ORDER BY h323_setup_time DESC LIMIT $%d OFFSET $%d",
		argPosition, argPosition+1))
	args = append(args, filters.PerPage, offset)

	rows, err := db.Query(q.String(), args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Erro ao buscar os registros",
			"details": err.Error(),
		})
	}
	defer rows.Close()

	bilhetes := []Bilhete{}
	for rows.Next() {
		var b Bilhete
		err := rows.Scan(
			&b.UserName,
			&b.AcctSessionID,
			&b.CallingStationID,
			&b.CalledStationID,
			&b.H323SetupTime,
			&b.H323ConnectTime,
			&b.H323DisconnectTime,
			&b.NASIdentifier,
			&b.CiscoNASPort,
			&b.H323CallOrigin,
			&b.ReleaseSource,
			&b.H323CallType,
			&b.CallID,
			&b.AcctSessionTime,
			&b.H323DisconnectCause,
			&b.NASIPAddress,
			&b.AcctStatusType,
			&b.Protocol,
			&b.Codec,
			&b.RemoteRTPIp,
			&b.RemoteRTPPort,
			&b.RemoteSIPIp,
			&b.RemoteSIPPort,
			&b.LocalRTPIp,
			&b.LocalRTPPort,
			&b.LocalSIPIp,
			&b.LocalSIPPort,
			&b.RingStart,
			&b.MosIngress,
			&b.MosEgress,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Erro lendo os registros",
				"details": err.Error(),
			})
		}
		bilhetes = append(bilhetes, b)
	}

	totalPages := (total + filters.PerPage - 1) / filters.PerPage
	// TODO voltar pagina anterior se maior que total

	r := BilhetesResponse{
		Data:        bilhetes,
		Total:       total,
		CurrentPage: filters.Page,
		PerPage:     filters.PerPage,
		TotalPages:  totalPages,
	}
	return c.JSON(r)
}

func autoClean(db *sql.DB, interval time.Duration) {
	go func() {
		for {
			var days int
			q := `SELECT cleanup_days FROM cleanup_config LIMIT 1;`
			if err := db.QueryRow(q).Scan(&days); err != nil {
				log.Printf("Erro ao obter o cleanup_days %v\n", err)
				time.Sleep(interval)
				continue
			}
			var count int
			countQ := fmt.Sprintf(`SELECT COUNT(*) FROM call_records WHERE created_at < now() - INTERVAL '%d day'`, days)
			if err := db.QueryRow(countQ).Scan(&count); err != nil {
				log.Printf("Erro ao contar a quantia de tickets a serem apagadas %v\n", err)
				time.Sleep(interval)
				continue
			}

			delQ := fmt.Sprintf(`DELETE FROM call_records WHERE created_at < NOW() - INTERVAL '%d day'`, days)
			if _, err := db.Exec(delQ); err != nil {
				log.Printf("Erro ao executar a limpeza %v\n", err)
			} else {
				log.Printf("Limpeza concluída, %d registros com mais de %d dias apagados", count, days)
			}
			time.Sleep(interval)
			// TODO gravar historico de limpezas
		}
	}()
}

var db *sql.DB

func main() {
	conn := "postgres://postgres:12345@localhost:5432/radius?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	autoClean(db, 24*time.Hour)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Post("/bilhetes", handlePostBilhete)
	app.Get("/bilhetes", handleGetBilhetes)
	app.Post("/rotinas/limpezadias", handleUpdateCleanupDays)
	app.Get("/rotinas/limpezadias", handleGetCleanupDays)
	log.Fatal(app.Listen(":5000"))
}
