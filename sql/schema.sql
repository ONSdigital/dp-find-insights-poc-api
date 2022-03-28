--
-- PostgreSQL database dump
--

-- Dumped from database version 13.4
-- Dumped by pg_dump version 13.4

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

--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: data_ver; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.data_ver (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    census_year integer,
    ver_string text,
    source text,
    notes text,
    public boolean
);


ALTER TABLE public.data_ver OWNER TO insights;

--
-- Name: geo; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.geo (
    id integer NOT NULL,
    type_id integer,
    code text,
    name text,
    lat numeric,
    long numeric,
    valid boolean DEFAULT true,
    wkb_geometry public.geometry(Geometry,4326),
    wkb_long_lat_geom public.geometry(Geometry,4326)
);


ALTER TABLE public.geo OWNER TO insights;

--
-- Name: geo_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.geo_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.geo_id_seq OWNER TO insights;

--
-- Name: geo_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.geo_id_seq OWNED BY public.geo.id;


--
-- Name: geo_metric; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.geo_metric (
    id integer NOT NULL,
    geo_id integer,
    category_id integer,
    metric numeric,
    data_ver_id integer
);


ALTER TABLE public.geo_metric OWNER TO insights;

--
-- Name: geo_metric_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.geo_metric_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.geo_metric_id_seq OWNER TO insights;

--
-- Name: geo_metric_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.geo_metric_id_seq OWNED BY public.geo_metric.id;


--
-- Name: geo_type; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.geo_type (
    id integer NOT NULL,
    name text
);


ALTER TABLE public.geo_type OWNER TO insights;

--
-- Name: lsoa2011_lad2020_lookup; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.lsoa2011_lad2020_lookup (
    id integer NOT NULL,
    lsoa2011code text,
    lad2020code text
);


ALTER TABLE public.lsoa2011_lad2020_lookup OWNER TO insights;

--
-- Name: lsoa2011_lad2020_lookup_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.lsoa2011_lad2020_lookup_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.lsoa2011_lad2020_lookup_id_seq OWNER TO insights;

--
-- Name: lsoa2011_lad2020_lookup_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.lsoa2011_lad2020_lookup_id_seq OWNED BY public.lsoa2011_lad2020_lookup.id;


--
-- Name: nomis_category; Type: TABLE; Schema: public; Owner: insights
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


ALTER TABLE public.nomis_category OWNER TO insights;

--
-- Name: nomis_category_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.nomis_category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nomis_category_id_seq OWNER TO insights;

--
-- Name: nomis_category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.nomis_category_id_seq OWNED BY public.nomis_category.id;


--
-- Name: nomis_desc; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.nomis_desc (
    id integer NOT NULL,
    nomis_topic_id integer NOT NULL,
    name text,
    pop_stat text,
    short_nomis_code text,
    year integer
);


ALTER TABLE public.nomis_desc OWNER TO insights;

--
-- Name: nomis_desc_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.nomis_desc_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nomis_desc_id_seq OWNER TO insights;

--
-- Name: nomis_desc_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.nomis_desc_id_seq OWNED BY public.nomis_desc.id;


--
-- Name: nomis_topic; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.nomis_topic (
    id integer NOT NULL,
    top_nomis_code text,
    name text
);


ALTER TABLE public.nomis_topic OWNER TO insights;

--
-- Name: nomis_topic_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.nomis_topic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nomis_topic_id_seq OWNER TO insights;

--
-- Name: nomis_topic_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.nomis_topic_id_seq OWNED BY public.nomis_topic.id;


--
-- Name: postcode; Type: TABLE; Schema: public; Owner: insights
--

CREATE TABLE public.postcode (
    id integer NOT NULL,
    geo_id integer,
    pcds text
);


ALTER TABLE public.postcode OWNER TO insights;

--
-- Name: postcode_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.postcode_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.postcode_id_seq OWNER TO insights;

--
-- Name: postcode_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.postcode_id_seq OWNED BY public.postcode.id;


--
-- Name: schema_ver; Type: TABLE; Schema: public; Owner: insights
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


ALTER TABLE public.schema_ver OWNER TO insights;

--
-- Name: schema_ver_id_seq; Type: SEQUENCE; Schema: public; Owner: insights
--

CREATE SEQUENCE public.schema_ver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.schema_ver_id_seq OWNER TO insights;

--
-- Name: schema_ver_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: insights
--

ALTER SEQUENCE public.schema_ver_id_seq OWNED BY public.schema_ver.id;


--
-- Name: geo id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo ALTER COLUMN id SET DEFAULT nextval('public.geo_id_seq'::regclass);


--
-- Name: geo_metric id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo_metric ALTER COLUMN id SET DEFAULT nextval('public.geo_metric_id_seq'::regclass);


--
-- Name: lsoa2011_lad2020_lookup id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.lsoa2011_lad2020_lookup ALTER COLUMN id SET DEFAULT nextval('public.lsoa2011_lad2020_lookup_id_seq'::regclass);


--
-- Name: nomis_category id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_category ALTER COLUMN id SET DEFAULT nextval('public.nomis_category_id_seq'::regclass);


--
-- Name: nomis_desc id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_desc ALTER COLUMN id SET DEFAULT nextval('public.nomis_desc_id_seq'::regclass);


--
-- Name: nomis_topic id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_topic ALTER COLUMN id SET DEFAULT nextval('public.nomis_topic_id_seq'::regclass);


--
-- Name: postcode id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.postcode ALTER COLUMN id SET DEFAULT nextval('public.postcode_id_seq'::regclass);


--
-- Name: schema_ver id; Type: DEFAULT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.schema_ver ALTER COLUMN id SET DEFAULT nextval('public.schema_ver_id_seq'::regclass);


--
-- Name: data_ver data_ver_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.data_ver
    ADD CONSTRAINT data_ver_pkey PRIMARY KEY (id);


--
-- Name: geo_metric geo_metric_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT geo_metric_pkey PRIMARY KEY (id);


--
-- Name: geo geo_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo
    ADD CONSTRAINT geo_pkey PRIMARY KEY (id);


--
-- Name: geo_type geo_type_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo_type
    ADD CONSTRAINT geo_type_pkey PRIMARY KEY (id);


--
-- Name: lsoa2011_lad2020_lookup lsoa2011_lad2020_lookup_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.lsoa2011_lad2020_lookup
    ADD CONSTRAINT lsoa2011_lad2020_lookup_pkey PRIMARY KEY (id);


--
-- Name: nomis_category nomis_category_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_category
    ADD CONSTRAINT nomis_category_pkey PRIMARY KEY (id, nomis_desc_id);


--
-- Name: nomis_desc nomis_desc_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_desc
    ADD CONSTRAINT nomis_desc_pkey PRIMARY KEY (id, nomis_topic_id);


--
-- Name: nomis_topic nomis_topic_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_topic
    ADD CONSTRAINT nomis_topic_pkey PRIMARY KEY (id);


--
-- Name: postcode postcode_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT postcode_pkey PRIMARY KEY (id);


--
-- Name: schema_ver schema_ver_pkey; Type: CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.schema_ver
    ADD CONSTRAINT schema_ver_pkey PRIMARY KEY (id);


--
-- Name: geo_long_lat_geom_idx; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX geo_long_lat_geom_idx ON public.geo USING gist (wkb_long_lat_geom);


--
-- Name: geo_wkb_geometry_geom_idx; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX geo_wkb_geometry_geom_idx ON public.geo USING gist (wkb_geometry);


--
-- Name: idx_data_ver_deleted_at; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX idx_data_ver_deleted_at ON public.data_ver USING btree (deleted_at);


--
-- Name: idx_geo_metric_category_id; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX idx_geo_metric_category_id ON public.geo_metric USING btree (category_id);


--
-- Name: idx_geo_metric_geo_id; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX idx_geo_metric_geo_id ON public.geo_metric USING btree (geo_id);


--
-- Name: idx_nomis_category_id; Type: INDEX; Schema: public; Owner: insights
--

CREATE UNIQUE INDEX idx_nomis_category_id ON public.nomis_category USING btree (id);


--
-- Name: idx_nomis_category_long_nomis_code; Type: INDEX; Schema: public; Owner: insights
--

CREATE UNIQUE INDEX idx_nomis_category_long_nomis_code ON public.nomis_category USING btree (long_nomis_code);


--
-- Name: idx_nomis_desc_id; Type: INDEX; Schema: public; Owner: insights
--

CREATE UNIQUE INDEX idx_nomis_desc_id ON public.nomis_desc USING btree (id);


--
-- Name: idx_nomis_desc_short_nomis_code; Type: INDEX; Schema: public; Owner: insights
--

CREATE UNIQUE INDEX idx_nomis_desc_short_nomis_code ON public.nomis_desc USING btree (short_nomis_code);


--
-- Name: idx_postcode_geo_id; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX idx_postcode_geo_id ON public.postcode USING btree (geo_id);


--
-- Name: idx_schema_ver_deleted_at; Type: INDEX; Schema: public; Owner: insights
--

CREATE INDEX idx_schema_ver_deleted_at ON public.schema_ver USING btree (deleted_at);


--
-- Name: unique; Type: INDEX; Schema: public; Owner: insights
--

CREATE UNIQUE INDEX "unique" ON public.geo USING btree (code);


--
-- Name: geo_metric fk_data_ver_go_metrics; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT fk_data_ver_go_metrics FOREIGN KEY (data_ver_id) REFERENCES public.data_ver(id);


--
-- Name: geo_metric fk_geo_go_metrics; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT fk_geo_go_metrics FOREIGN KEY (geo_id) REFERENCES public.geo(id);


--
-- Name: postcode fk_geo_post_codes; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.postcode
    ADD CONSTRAINT fk_geo_post_codes FOREIGN KEY (geo_id) REFERENCES public.geo(id);


--
-- Name: geo fk_geo_type_geos; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo
    ADD CONSTRAINT fk_geo_type_geos FOREIGN KEY (type_id) REFERENCES public.geo_type(id);


--
-- Name: geo_metric fk_nomis_category_go_metrics; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.geo_metric
    ADD CONSTRAINT fk_nomis_category_go_metrics FOREIGN KEY (category_id) REFERENCES public.nomis_category(id);


--
-- Name: nomis_category fk_nomis_desc_nomis_categories; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_category
    ADD CONSTRAINT fk_nomis_desc_nomis_categories FOREIGN KEY (nomis_desc_id) REFERENCES public.nomis_desc(id);


--
-- Name: nomis_desc fk_nomis_topic_nomis_descs; Type: FK CONSTRAINT; Schema: public; Owner: insights
--

ALTER TABLE ONLY public.nomis_desc
    ADD CONSTRAINT fk_nomis_topic_nomis_descs FOREIGN KEY (nomis_topic_id) REFERENCES public.nomis_topic(id);


--
-- PostgreSQL database dump complete
--

