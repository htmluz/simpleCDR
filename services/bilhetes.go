package services

import (
	"database/sql"
	"fmt"
	"log"
	"radiusgo/models"
	"radiusgo/utils"
)

func InsertBid(db *sql.DB, bilh models.BilheteFull) error {
	var callIDa, callIDb interface{}
	if bilh.LegA != nil && bilh.LegA.CallID != "" {
		callIDa = bilh.LegA.CallID
	}
	if bilh.LegB != nil && bilh.LegB.CallID != "" {
		callIDb = bilh.LegB.CallID
	}

	_, err := db.Exec(
		`INSERT INTO bilhetes (bid, lega, legb) VALUES ($1, $2, $3)
		ON CONFLICT (bid) 
		DO UPDATE SET lega = COALESCE(EXCLUDED.lega, bilhetes.lega),
		legb = COALESCE(EXCLUDED.legb, bilhetes.legb)
		`, bilh.Bid, callIDa, callIDb)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error inserting the ticket %v", err)
	}
	return nil
}

func BidExists(db *sql.DB, bid string) bool {
	var exists bool
	err := db.QueryRow(`SELECT COUNT(*) > 0 FROM bilhetes WHERE bid = $1`, bid).Scan(&exists)
	if err != nil {
		fmt.Println("erro query exists ", err)
		return false
	}
	return exists
}

func InsertBilhete(db *sql.DB, bilhete *models.Bilhete) error {
	h323SetupTime, e := utils.ConvertToTimestamp(bilhete.H323SetupTime)
	if e != nil {
		return e
	}
	h323ConnectTime, e := utils.ConvertToTimestamp(bilhete.H323ConnectTime)
	if e != nil {
		return e
	}
	h323DisconnectTime, e := utils.ConvertToTimestamp(bilhete.H323DisconnectTime)
	if e != nil {
		return e
	}
	fRingStart, e := utils.ConvertToTimestamp(bilhete.RingStart)
	if e != nil {
		return e
	}
	remoteRTPPort, e := utils.StrToInt(bilhete.RemoteRTPPort)
	if e != nil {
		return e
	}
	remoteSIPPort, e := utils.StrToInt(bilhete.RemoteSIPPort)
	if e != nil {
		return e
	}
	localRTPPort, e := utils.StrToInt(bilhete.LocalRTPPort)
	if e != nil {
		return e
	}
	localSIPPort, e := utils.StrToInt(bilhete.LocalSIPPort)
	if e != nil {
		return e
	}
	remoteRTPIP, e := utils.HasIP(bilhete.RemoteRTPIp)
	if e != nil {
		return e
	}
	remoteSIPIP, e := utils.HasIP(bilhete.RemoteSIPIp)
	if e != nil {
		return e
	}
	localRTPIP, e := utils.HasIP(bilhete.LocalRTPIp)
	if e != nil {
		return e
	}
	localSIPIP, e := utils.HasIP(bilhete.LocalSIPIp)
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
		return fmt.Errorf("Error inserting the ticket %v", err)
	}
	return nil
}
