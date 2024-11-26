-- Name: cleanup_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cleanup_config (
    id integer NOT NULL,
    cleanup_days integer DEFAULT 90 NOT NULL,
    updated_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.cleanup_config OWNER TO postgres;

--
-- Name: cleanup_config_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cleanup_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cleanup_config_id_seq OWNER TO postgres;

--
-- Name: cleanup_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cleanup_config_id_seq OWNED BY public.cleanup_config.id;


--
-- Name: cleanup_config id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cleanup_config ALTER COLUMN id SET DEFAULT nextval('public.cleanup_config_id_seq'::regclass);


--
-- Name: cleanup_config cleanup_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cleanup_config
    ADD CONSTRAINT cleanup_config_pkey PRIMARY KEY (id);

--
-- PostgreSQL database dump complete
--

