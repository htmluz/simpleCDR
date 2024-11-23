--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: call_records; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.call_records (
    id integer NOT NULL,
    user_name character varying(255),
    acct_session_id character varying(255),
    calling_station_id character varying(255),
    called_station_id character varying(255),
    nas_identifier character varying(255),
    cisco_nas_port character varying(255),
    h323_call_origin character varying(255),
    release_source character varying(255),
    h323_call_type character varying(255),
    call_id character varying(255),
    acct_session_time character varying(255),
    h323_disconnect_cause character varying(255),
    nas_ip_address character varying(255),
    acct_status_type character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    h323_setup_time timestamp without time zone,
    h323_connect_time timestamp without time zone,
    h323_disconnect_time timestamp without time zone,
    protocol character varying(255),
    codec character varying(255),
    remote_rtp_ip inet,
    remote_rtp_port integer,
    remote_sip_ip inet,
    remote_sip_port integer,
    local_rtp_ip inet,
    local_rtp_port integer,
    local_sip_ip inet,
    local_sip_port integer,
    ring_start timestamp without time zone,
    mos_ingress character varying(255),
    mos_egress character varying(255)
);


ALTER TABLE public.call_records OWNER TO postgres;

--
-- Name: call_records_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.call_records_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.call_records_id_seq OWNER TO postgres;

--
-- Name: call_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.call_records_id_seq OWNED BY public.call_records.id;


--
-- Name: call_records id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.call_records ALTER COLUMN id SET DEFAULT nextval('public.call_records_id_seq'::regclass);


--
-- Name: call_records call_records_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.call_records
    ADD CONSTRAINT call_records_pkey PRIMARY KEY (id);


--
-- Name: call_records unique_call_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.call_records
    ADD CONSTRAINT unique_call_id UNIQUE (call_id);


--
-- Name: idx_call_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_call_id ON public.call_records USING btree (call_id);


--
-- Name: idx_calling_station; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_calling_station ON public.call_records USING btree (calling_station_id);


--
-- PostgreSQL database dump complete
--

