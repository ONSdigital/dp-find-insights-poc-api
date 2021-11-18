--
-- PostgreSQL database dump
--

-- Dumped from database version 13.5 (Debian 13.5-0+deb11u1)
-- Dumped by pg_dump version 13.5 (Debian 13.5-0+deb11u1)

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
-- Name: data_vers; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.data_vers (
    id integer NOT NULL,
    year integer,
    ver_string text,
    notes text,
    public boolean
);


ALTER TABLE public.data_vers OWNER TO steve;

--
-- Name: data_vers_id_seq; Type: SEQUENCE; Schema: public; Owner: steve
--

CREATE SEQUENCE public.data_vers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.data_vers_id_seq OWNER TO steve;

--
-- Name: data_vers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: steve
--

ALTER SEQUENCE public.data_vers_id_seq OWNED BY public.data_vers.id;


--
-- Name: geo; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.geo (
    id integer NOT NULL,
    geo_type_id integer,
    geo_code text,
    geo_name text
);


ALTER TABLE public.geo OWNER TO steve;

--
-- Name: geo_metric; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.geo_metric (
    id integer NOT NULL,
    geo_id integer,
    category_id integer,
    metric numeric,
    data_ver_id integer,
    year integer
);


ALTER TABLE public.geo_metric OWNER TO steve;

--
-- Name: geo_metric_id_seq; Type: SEQUENCE; Schema: public; Owner: steve
--

CREATE SEQUENCE public.geo_metric_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.geo_metric_id_seq OWNER TO steve;

--
-- Name: geo_metric_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: steve
--

ALTER SEQUENCE public.geo_metric_id_seq OWNED BY public.geo_metric.id;


--
-- Name: geo_type; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.geo_type (
    id integer NOT NULL,
    geo_type_name text
);


ALTER TABLE public.geo_type OWNER TO steve;

--
-- Name: lsoa2011_lad2020_lookup; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.lsoa2011_lad2020_lookup (
    id integer NOT NULL,
    lsoa2011code text,
    lad2020code text
);


ALTER TABLE public.lsoa2011_lad2020_lookup OWNER TO steve;

--
-- Name: lsoa2011_lad2020_lookup_id_seq; Type: SEQUENCE; Schema: public; Owner: steve
--

CREATE SEQUENCE public.lsoa2011_lad2020_lookup_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.lsoa2011_lad2020_lookup_id_seq OWNER TO steve;

--
-- Name: lsoa2011_lad2020_lookup_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: steve
--

ALTER SEQUENCE public.lsoa2011_lad2020_lookup_id_seq OWNED BY public.lsoa2011_lad2020_lookup.id;


--
-- Name: nomis_category; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.nomis_category (
    id integer NOT NULL,
    nomis_desc_id integer NOT NULL,
    category_name text,
    measurement_unit text,
    stat_unit text,
    long_nomis_code text,
    year integer
);


ALTER TABLE public.nomis_category OWNER TO steve;

--
-- Name: nomis_category_id_seq; Type: SEQUENCE; Schema: public; Owner: steve
--

CREATE SEQUENCE public.nomis_category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nomis_category_id_seq OWNER TO steve;

--
-- Name: nomis_category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: steve
--

ALTER SEQUENCE public.nomis_category_id_seq OWNED BY public.nomis_category.id;


--
-- Name: nomis_desc; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.nomis_desc (
    id integer NOT NULL,
    long_desc text,
    short_desc text,
    short_nomis_code text,
    year integer
);


ALTER TABLE public.nomis_desc OWNER TO steve;

--
-- Name: nomis_desc_id_seq; Type: SEQUENCE; Schema: public; Owner: steve
--

CREATE SEQUENCE public.nomis_desc_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nomis_desc_id_seq OWNER TO steve;

--
-- Name: nomis_desc_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: steve
--

ALTER SEQUENCE public.nomis_desc_id_seq OWNED BY public.nomis_desc.id;


--
-- Name: schema_ver; Type: TABLE; Schema: public; Owner: steve
--

CREATE TABLE public.schema_ver (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    build_time text,
    git_commit text,
    version text
);


ALTER TABLE public.schema_ver OWNER TO steve;

--
-- Name: schema_ver_id_seq; Type: SEQUENCE; Schema: public; Owner: steve
--

CREATE SEQUENCE public.schema_ver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.schema_ver_id_seq OWNER TO steve;

--
-- Name: schema_ver_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: steve
--

ALTER SEQUENCE public.schema_ver_id_seq OWNED BY public.schema_ver.id;


--
-- Name: data_vers id; Type: DEFAULT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.data_vers ALTER COLUMN id SET DEFAULT nextval('public.data_vers_id_seq'::regclass);


--
-- Name: geo_metric id; Type: DEFAULT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo_metric ALTER COLUMN id SET DEFAULT nextval('public.geo_metric_id_seq'::regclass);


--
-- Name: lsoa2011_lad2020_lookup id; Type: DEFAULT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.lsoa2011_lad2020_lookup ALTER COLUMN id SET DEFAULT nextval('public.lsoa2011_lad2020_lookup_id_seq'::regclass);


--
-- Name: nomis_category id; Type: DEFAULT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.nomis_category ALTER COLUMN id SET DEFAULT nextval('public.nomis_category_id_seq'::regclass);


--
-- Name: nomis_desc id; Type: DEFAULT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.nomis_desc ALTER COLUMN id SET DEFAULT nextval('public.nomis_desc_id_seq'::regclass);


--
-- Name: schema_ver id; Type: DEFAULT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.schema_ver ALTER COLUMN id SET DEFAULT nextval('public.schema_ver_id_seq'::regclass);


--
-- Name: data_vers data_vers_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.data_vers
    ADD CONSTRAINT data_vers_pkey PRIMARY KEY (id);


--
-- Name: geo_metric geo_metric_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT geo_metric_pkey PRIMARY KEY (id);


--
-- Name: geo geo_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo
    ADD CONSTRAINT geo_pkey PRIMARY KEY (id);


--
-- Name: geo_type geo_type_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo_type
    ADD CONSTRAINT geo_type_pkey PRIMARY KEY (id);


--
-- Name: lsoa2011_lad2020_lookup lsoa2011_lad2020_lookup_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.lsoa2011_lad2020_lookup
    ADD CONSTRAINT lsoa2011_lad2020_lookup_pkey PRIMARY KEY (id);


--
-- Name: nomis_category nomis_category_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.nomis_category
    ADD CONSTRAINT nomis_category_pkey PRIMARY KEY (id, nomis_desc_id);


--
-- Name: nomis_desc nomis_desc_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.nomis_desc
    ADD CONSTRAINT nomis_desc_pkey PRIMARY KEY (id);


--
-- Name: schema_ver schema_ver_pkey; Type: CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.schema_ver
    ADD CONSTRAINT schema_ver_pkey PRIMARY KEY (id);


--
-- Name: idx_geo_metric_geo_id; Type: INDEX; Schema: public; Owner: steve
--

CREATE INDEX idx_geo_metric_geo_id ON public.geo_metric USING btree (geo_id);


--
-- Name: idx_nomis_category_id; Type: INDEX; Schema: public; Owner: steve
--

CREATE UNIQUE INDEX idx_nomis_category_id ON public.nomis_category USING btree (id);


--
-- Name: idx_schema_ver_deleted_at; Type: INDEX; Schema: public; Owner: steve
--

CREATE INDEX idx_schema_ver_deleted_at ON public.schema_ver USING btree (deleted_at);


--
-- Name: unique; Type: INDEX; Schema: public; Owner: steve
--

CREATE UNIQUE INDEX "unique" ON public.geo USING btree (geo_code);


--
-- Name: geo_metric fk_data_vers_go_metrics; Type: FK CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT fk_data_vers_go_metrics FOREIGN KEY (data_ver_id) REFERENCES public.data_vers(id);


--
-- Name: geo_metric fk_geo_go_metrics; Type: FK CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT fk_geo_go_metrics FOREIGN KEY (geo_id) REFERENCES public.geo(id);


--
-- Name: geo fk_geo_type_geos; Type: FK CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo
    ADD CONSTRAINT fk_geo_type_geos FOREIGN KEY (geo_type_id) REFERENCES public.geo_type(id);


--
-- Name: geo_metric fk_nomis_category_go_metrics; Type: FK CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT fk_nomis_category_go_metrics FOREIGN KEY (category_id) REFERENCES public.nomis_category(id);


--
-- Name: nomis_category fk_nomis_desc_nomis_categories; Type: FK CONSTRAINT; Schema: public; Owner: steve
--

ALTER TABLE ONLY public.nomis_category
    ADD CONSTRAINT fk_nomis_desc_nomis_categories FOREIGN KEY (nomis_desc_id) REFERENCES public.nomis_desc(id);


--
-- PostgreSQL database dump complete
--
