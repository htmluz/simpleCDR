package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"radiusgo/models"
	"radiusgo/services"
	"radiusgo/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandlePostBilhete(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := c.Body()
		bilhete := new(models.Bilhete)

		if err := json.Unmarshal(data, bilhete); err != nil {
			log.Fatal("Erro fazendo parsing do json")
		}

		err := services.InsertBilhete(db, bilhete)
		if err != nil {
			log.Fatal("Erro inserindo bilhete")
		}

		return c.SendStatus(201)
	}
}

// TODO moodularizar
func HandleGetBilhetes(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filters := new(models.FilterParams)
		if err := c.QueryParser(filters); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":  "Invalid parameters",
				"detail": err.Error(),
			})
		}
		if filters.Page < 1 {
			filters.Page = 1
		}
		if filters.PerPage < 1 || filters.PerPage > 500 {
			filters.PerPage = 25
		}

		q := strings.Builder{}
		q.WriteString(`
    SELECT 
      cr.user_name, cr.acct_session_id, cr.calling_station_id, cr.called_station_id,
      cr.h323_setup_time, cr.h323_connect_time, cr.h323_disconnect_time,
      cr.nas_identifier, cr.cisco_nas_port, cr.h323_call_origin,
      cr.release_source, cr.h323_call_type, cr.call_id,
      cr.acct_session_time, cr.h323_disconnect_cause,
      cr.nas_ip_address, cr.acct_status_type, cr.protocol,
      cr.codec, cr.remote_rtp_ip, cr.remote_rtp_port,
      cr.remote_sip_ip, cr.remote_sip_port,
      cr.local_rtp_ip, cr.local_rtp_port,
      cr.local_sip_ip, cr.local_sip_port,
      cr.ring_start, cr.mos_ingress, cr.mos_egress,
      COALESCE(g.name, '') AS gwname
    FROM call_records cr
    LEFT JOIN gateways g
    ON cr.nas_ip_address::inet = g.ip
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
			args = append(args, "%"+(filters.CalledPhone)+"%")
			argPosition++
		}
		if filters.CallingPhone != "" {
			q.WriteString(fmt.Sprintf(" AND calling_station_id ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.CallingPhone)+"%")
			argPosition++
		}
		if filters.AnyPhone != "" {
			q.WriteString(fmt.Sprintf(" AND (calling_station_id ILIKE $%d OR called_station_id ILIKE $%d)",
				argPosition, argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.AnyPhone)+"%")
			argPosition++
		}
		if filters.NapA != "" {
			q.WriteString(fmt.Sprintf(" AND user_name ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.NapA)+"%")
			argPosition++
		}
		if filters.NapB != "" {
			q.WriteString(fmt.Sprintf(" AND cisco_nas_port ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.NapB)+"%")
			argPosition++
		}
		if filters.DisconnCause != "" {
			q.WriteString(fmt.Sprintf(" AND h323_disconnect_cause ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.DisconnCause)+"%")
			argPosition++
		}
		if filters.CallID != "" {
			q.WriteString(fmt.Sprintf(" AND call_id ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.CallID)+"%")
			argPosition++
		}
		if filters.GatewayIP != "" {
			q.WriteString(fmt.Sprintf(" AND nas_ip_address ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.GatewayIP)+"%")
			argPosition++
		}
		if filters.Codec != "" {
			q.WriteString(fmt.Sprintf(" AND codec ILIKE $%d", argPosition))
			args = append(args, "%"+utils.RemoveSpaces(filters.Codec)+"%")
			argPosition++
		}

		countQ := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", q.String())
		var total int
		err := db.QueryRow(countQ, args...).Scan(&total)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Error counting records",
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
				"error":   "Error searching for records",
				"details": err.Error(),
			})
		}
		defer rows.Close()

		bilhetes := []models.Bilhete{}
		for rows.Next() {
			var b models.Bilhete
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
				&b.GwName,
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error":   "Error reading records",
					"details": err.Error(),
				})
			}
			bilhetes = append(bilhetes, b)
		}

		totalPages := (total + filters.PerPage - 1) / filters.PerPage
		// TODO voltar pagina anterior se maior que total

		r := models.BilhetesResponse{
			Data:        bilhetes,
			Total:       total,
			CurrentPage: filters.Page,
			PerPage:     filters.PerPage,
			TotalPages:  totalPages,
		}
		return c.JSON(r)
	}
}
