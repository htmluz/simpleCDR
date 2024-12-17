CREATE INDEX idx_call_records_setup_time ON call_records (h323_setup_time);
CREATE INDEX idx_call_records_disconnect_time ON call_records (h323_disconnect_time);
CREATE INDEX idx_call_records_nas_ip ON call_records (nas_ip_address);
CREATE INDEX idx_gateways_ip ON gateways (ip);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_calling_station_id ON call_records USING gin (calling_station_id gin_trgm_ops);
CREATE INDEX idx_called_station_id ON call_records USING gin (called_station_id gin_trgm_ops);
