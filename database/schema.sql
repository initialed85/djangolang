ALTER SYSTEM
SET
    log_statement = 'all';

ALTER DATABASE some_db
SET
    log_statement = 'all';

SELECT
    pg_reload_conf();

--
-- PostgreSQL database dump
--
-- Dumped from database version 14.11 (Debian 14.11-1.pgdg110+2)
-- Dumped by pg_dump version 15.5 (Homebrew)
SET
    statement_timeout = 0;

SET
    lock_timeout = 0;

SET
    idle_in_transaction_session_timeout = 0;

SET
    client_encoding = 'UTF8';

SET
    standard_conforming_strings = on;

SELECT
    pg_catalog.set_config ('search_path', '', false);

SET
    check_function_bodies = false;

SET
    xmloption = content;

SET
    client_min_messages = warning;

SET
    row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--
CREATE SCHEMA IF NOT EXISTS public;

ALTER SCHEMA public OWNER TO postgres;

SET
    search_path = 'public';

--
-- Name: aggregate_detection(); Type: FUNCTION; Schema: public; Owner: postgres
--
CREATE FUNCTION public.aggregate_detection () RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.status = 'needs tracking' THEN
        -- id bigint NOT NULL PRIMARY KEY,
        -- start_timestamp timestamp with time zone NOT NULL,
        -- end_timestamp timestamp with time zone NOT NULL,
        -- class_id bigint NOT NULL,
        -- class_name text NOT NULL,
        -- score float NOT NULL,
        -- count bigint NOT NULL,
        -- weighted_score float NOT NULL,
        -- event_id bigint NULL

        WITH cte1 AS (
            SELECT
                d.class_id AS class_id,
                d.class_name AS class_name,
                avg(d.score) AS score,
                count(d.*) AS count,
                avg(d.score) * count(d.*) AS weighted_score,
                NEW.id AS event_id
            FROM detections d
            WHERE d.event_id = NEW.id
            GROUP BY (d.class_id, d.class_name)
        ),
        cte2 AS (
            SELECT NEW.start_timestamp AS start_timestamp,
                NEW.end_timestamp AS end_timestamp,
                cte1.*
            FROM cte1
        )
        INSERT INTO aggregated_detection (
            start_timestamp,
            end_timestamp,
            class_id,
            class_name,
            score,
            count,
            weighted_score,
            event_id
        )
        SELECT * FROM cte2;

    END IF;

    RETURN NEW;
END;
$$;

ALTER FUNCTION public.aggregate_detection () OWNER TO postgres;

SET
    default_tablespace = '';

SET
    default_table_access_method = heap;

--
-- Name: aggregated_detection; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.aggregated_detection (
        id bigint NOT NULL,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        class_id bigint NOT NULL,
        class_name text NOT NULL,
        score double precision NOT NULL,
        count bigint NOT NULL,
        weighted_score double precision NOT NULL,
        event_id bigint
    );

ALTER TABLE public.aggregated_detection OWNER TO postgres;

--
-- Name: aggregated_detection_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.aggregated_detection_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.aggregated_detection_id_seq OWNER TO postgres;

--
-- Name: aggregated_detection_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.aggregated_detection_id_seq OWNED BY public.aggregated_detection.id;

--
-- Name: camera; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.camera (id bigint NOT NULL, name text NOT NULL, stream_url text NOT NULL);

ALTER TABLE public.camera OWNER TO postgres;

--
-- Name: camera_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.camera_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.camera_id_seq OWNER TO postgres;

--
-- Name: camera_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.camera_id_seq OWNED BY public.camera.id;

--
-- Name: cameras; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.cameras AS
SELECT
    camera.id,
    camera.name,
    camera.stream_url
FROM
    public.camera;

ALTER TABLE public.cameras OWNER TO postgres;

--
-- Name: detection; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.detection (
        id bigint NOT NULL,
        "timestamp" timestamp with time zone NOT NULL,
        class_id bigint NOT NULL,
        class_name text NOT NULL,
        score double precision NOT NULL,
        centroid point NOT NULL,
        bounding_box polygon NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint,
        object_id bigint,
        colour public.geometry (PointZ)
    );

ALTER TABLE public.detection OWNER TO postgres;

--
-- Name: detection_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.detection_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.detection_id_seq OWNER TO postgres;

--
-- Name: detection_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.detection_id_seq OWNED BY public.detection.id;

--
-- Name: detections; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.detections AS
SELECT
    detection.id,
    detection."timestamp",
    detection.class_id,
    detection.class_name,
    detection.score,
    detection.centroid,
    detection.bounding_box,
    detection.camera_id,
    detection.event_id,
    detection.object_id,
    detection.colour
FROM
    public.detection;

ALTER TABLE public.detections OWNER TO postgres;

--
-- Name: event; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.event (
        id bigint NOT NULL,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        duration interval DEFAULT '00:00:00'::interval NOT NULL,
        original_video_id bigint NOT NULL,
        thumbnail_image_id bigint NOT NULL,
        processed_video_id bigint,
        source_camera_id bigint NOT NULL,
        status text DEFAULT true NOT NULL,
        CONSTRAINT event_status_check CHECK (
            (
                status = ANY (
                    ARRAY['needs detection'::text, 'detection underway'::text, 'needs tracking'::text, 'tracking underway'::text, 'done'::text]
                )
            )
        )
    );

ALTER TABLE public.event OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.event_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.event_id_seq OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.event_id_seq OWNED BY public.event.id;

--
-- Name: events; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.events AS
SELECT
    event.id,
    event.start_timestamp,
    event.end_timestamp,
    event.duration,
    event.original_video_id,
    event.thumbnail_image_id,
    event.processed_video_id,
    event.source_camera_id,
    event.status
FROM
    public.event;

ALTER TABLE public.events OWNER TO postgres;

--
-- Name: event_with_detection; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.event_with_detection AS
WITH
    detections_1 AS (
        SELECT
            detections.id,
            detections."timestamp",
            detections.class_id,
            detections.class_name,
            detections.score,
            detections.centroid,
            detections.bounding_box,
            detections.camera_id,
            detections.event_id,
            detections.object_id,
            detections.colour
        FROM
            public.detections
        ORDER BY
            detections."timestamp" DESC
    ),
    detection_2 AS (
        SELECT
            d.event_id,
            d.class_id,
            d.class_name,
            avg(d.score) AS score,
            count(d.id) AS count,
            (avg(d.score) * (count(d.id))::double precision) AS weighted_score
        FROM
            public.detection d
        GROUP BY
            d.event_id,
            d.class_id,
            d.class_name
    ),
    events_1 AS (
        SELECT
            e_1.id,
            e_1.start_timestamp,
            e_1.end_timestamp,
            e_1.duration,
            e_1.original_video_id,
            e_1.thumbnail_image_id,
            e_1.processed_video_id,
            e_1.source_camera_id,
            e_1.status,
            d.event_id,
            d.class_id,
            d.class_name,
            d.score,
            d.count,
            d.weighted_score
        FROM
            (
                public.events e_1
                LEFT JOIN detection_2 d ON ((d.event_id = e_1.id))
            )
        ORDER BY
            e_1.start_timestamp DESC
    )
SELECT
    e.id,
    e.start_timestamp,
    e.end_timestamp,
    e.duration,
    e.original_video_id,
    e.thumbnail_image_id,
    e.processed_video_id,
    e.source_camera_id,
    e.status,
    e.event_id,
    e.class_id,
    e.class_name,
    e.score,
    e.count,
    e.weighted_score
FROM
    events_1 e
ORDER BY
    e.start_timestamp DESC;

ALTER TABLE public.event_with_detection OWNER TO postgres;

--
-- Name: image; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.image (
        id bigint NOT NULL,
        "timestamp" timestamp with time zone NOT NULL,
        size double precision DEFAULT 0 NOT NULL,
        file_path text NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint
    );

ALTER TABLE public.image OWNER TO postgres;

--
-- Name: image_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.image_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.image_id_seq OWNER TO postgres;

--
-- Name: image_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.image_id_seq OWNED BY public.image.id;

--
-- Name: images; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.images AS
SELECT
    image.id,
    image."timestamp",
    image.size,
    image.file_path,
    image.camera_id,
    image.event_id
FROM
    public.image;

ALTER TABLE public.images OWNER TO postgres;

--
-- Name: object; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.object (
        id bigint NOT NULL,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        class_id bigint NOT NULL,
        class_name text NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint
    );

ALTER TABLE public.object OWNER TO postgres;

--
-- Name: object_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.object_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.object_id_seq OWNER TO postgres;

--
-- Name: object_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.object_id_seq OWNED BY public.object.id;

--
-- Name: objects; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.objects AS
SELECT
    object.id,
    object.start_timestamp,
    object.end_timestamp,
    object.class_id,
    object.class_name,
    object.camera_id,
    object.event_id
FROM
    public.object;

ALTER TABLE public.objects OWNER TO postgres;

--
-- Name: v_introspect_table_oids; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.v_introspect_table_oids AS
WITH
    table_oids AS (
        SELECT
            c.relname,
            c.oid,
            n.nspname,
            c.reltuples,
            c.relkind,
            c.relam,
            n.nspacl,
            c.relacl,
            c.reltype,
            c.relowner,
            c.relhasindex
        FROM
            (
                pg_class c
                LEFT JOIN pg_namespace n ON ((n.oid = c.relnamespace))
            )
        WHERE
            pg_table_is_visible(c.oid)
    )
SELECT
    table_oids.relname,
    table_oids.oid,
    table_oids.nspname,
    table_oids.reltuples,
    table_oids.relkind,
    table_oids.relam,
    table_oids.nspacl,
    table_oids.relacl,
    table_oids.reltype,
    table_oids.relowner,
    table_oids.relhasindex
FROM
    table_oids
ORDER BY
    table_oids.nspname,
    table_oids.relname;

ALTER TABLE public.v_introspect_table_oids OWNER TO postgres;

--
-- Name: v_introspect_columns; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.v_introspect_columns AS
WITH
    foreign_keys AS (
        SELECT DISTINCT
            tc.table_schema AS schema_name,
            kcu.table_name,
            kcu.column_name,
            ccu.table_name AS foreign_table_name,
            ccu.column_name AS foreign_column_name
        FROM
            (
                (
                    information_schema.table_constraints tc
                    JOIN information_schema.key_column_usage kcu ON (
                        (
                            ((tc.constraint_name)::name = (kcu.constraint_name)::name)
                            AND ((tc.table_schema)::name = (kcu.table_schema)::name)
                        )
                    )
                )
                JOIN information_schema.constraint_column_usage ccu ON (
                    (
                        ((ccu.constraint_name)::name = (tc.constraint_name)::name)
                        AND ((ccu.table_schema)::name = (tc.table_schema)::name)
                    )
                )
            )
        WHERE
            ((tc.constraint_type)::text = 'FOREIGN KEY'::text)
    ),
    primary_keys AS (
        SELECT DISTINCT
            tc.table_schema AS schema_name,
            kcu.table_name,
            kcu.column_name
        FROM
            (
                information_schema.table_constraints tc
                JOIN information_schema.key_column_usage kcu ON (
                    (
                        ((tc.constraint_name)::name = (kcu.constraint_name)::name)
                        AND ((tc.table_schema)::name = (kcu.table_schema)::name)
                    )
                )
            )
        WHERE
            ((tc.constraint_type)::text = 'PRIMARY KEY'::text)
    )
SELECT
    toids.relname AS tablename,
    a.attname AS "column",
    format_type(a.atttypid, a.atttypmod) AS datatype,
    toids.nspname AS schema,
    toids.oid,
    row_number() OVER (
        ORDER BY
            toids.nspname DESC,
            toids.relname,
            a.attnum
    ) AS column_id,
    a.attnum AS pos,
    a.atttypid AS typeid,
    a.attlen AS typelen,
    a.atttypmod AS typemod,
    a.attndims AS ndims,
    a.attnotnull AS "notnull",
    a.atthasdef AS hasdefault,
    a.atthasmissing AS hasmissing,
    (
        (primary_keys.column_name IS NOT NULL)
        AND (foreign_keys.foreign_column_name IS NULL)
    ) AS ispkey,
    foreign_keys.foreign_table_name AS ftable,
    foreign_keys.foreign_column_name AS fcolumn
FROM
    pg_attribute a,
    (
        (
            public.v_introspect_table_oids toids
            LEFT JOIN LATERAL (
                SELECT
                    primary_keys_1.schema_name,
                    primary_keys_1.table_name,
                    primary_keys_1.column_name
                FROM
                    primary_keys primary_keys_1
                WHERE
                    (
                        ((primary_keys_1.schema_name)::name = toids.nspname)
                        AND ((primary_keys_1.table_name)::name = toids.relname)
                        AND ((primary_keys_1.column_name)::name = a.attname)
                    )
            ) primary_keys ON (true)
        )
        LEFT JOIN LATERAL (
            SELECT
                foreign_keys_1.schema_name,
                foreign_keys_1.table_name,
                foreign_keys_1.column_name,
                foreign_keys_1.foreign_table_name,
                foreign_keys_1.foreign_column_name
            FROM
                foreign_keys foreign_keys_1
            WHERE
                (
                    ((foreign_keys_1.schema_name)::name = toids.nspname)
                    AND ((foreign_keys_1.table_name)::name = toids.relname)
                    AND ((foreign_keys_1.column_name)::name = a.attname)
                )
        ) foreign_keys ON (true)
    )
WHERE
    (
        (a.attnum > 0)
        AND (NOT a.attisdropped)
        AND (a.attrelid = toids.oid)
    )
ORDER BY
    toids.nspname DESC,
    toids.relname,
    a.attnum;

ALTER TABLE public.v_introspect_columns OWNER TO postgres;

--
-- Name: v_introspect_tables; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.v_introspect_tables AS
SELECT
    col.tablename,
    col.oid,
    col.schema,
    ts.reltuples,
    ts.relkind,
    ts.relam,
    ts.relacl,
    ts.reltype,
    ts.relowner,
    ts.relhasindex,
    json_agg(
        json_build_object(
            'column',
            col."column",
            'datatype',
            col.datatype,
            'table',
            col.tablename,
            'pos',
            col.pos,
            'typeid',
            col.typeid,
            'typelen',
            col.typelen,
            'typemod',
            col.typemod,
            'notnull',
            col."notnull",
            'hasdefault',
            col.hasdefault,
            'hasmissing',
            col.hasmissing,
            'ispkey',
            col.ispkey,
            'ftable',
            col.ftable,
            'fcolumn',
            col.fcolumn,
            'parent_id',
            col.oid
        )
    ) AS columns
FROM
    public.v_introspect_columns col,
    public.v_introspect_table_oids ts
WHERE
    (ts.oid = col.oid)
GROUP BY
    col.tablename,
    col.oid,
    col.schema,
    ts.reltuples,
    ts.relkind,
    ts.relam,
    ts.nspacl,
    ts.relacl,
    ts.reltype,
    ts.relowner,
    ts.relhasindex
ORDER BY
    col.schema DESC,
    col.tablename,
    ts.relkind;

ALTER TABLE public.v_introspect_tables OWNER TO postgres;

--
-- Name: v_introspect_schemas; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.v_introspect_schemas AS
SELECT
    v_introspect_tables.schema,
    row_number() OVER (
        ORDER BY
            v_introspect_tables.schema
    ) AS id,
    json_agg(
        json_build_object(
            'id',
            v_introspect_tables.oid,
            'tablename',
            v_introspect_tables.tablename,
            'columns',
            v_introspect_tables.columns,
            'schema',
            v_introspect_tables.schema
        )
    ) AS tables
FROM
    public.v_introspect_tables
GROUP BY
    v_introspect_tables.schema
ORDER BY
    v_introspect_tables.schema;

ALTER TABLE public.v_introspect_schemas OWNER TO postgres;

--
-- Name: video; Type: TABLE; Schema: public; Owner: postgres
--
CREATE TABLE
    public.video (
        id bigint NOT NULL,
        start_timestamp timestamp with time zone NOT NULL,
        end_timestamp timestamp with time zone NOT NULL,
        duration interval DEFAULT '00:00:00'::interval NOT NULL,
        size double precision DEFAULT 0 NOT NULL,
        file_path text NOT NULL,
        camera_id bigint NOT NULL,
        event_id bigint
    );

ALTER TABLE public.video OWNER TO postgres;

--
-- Name: video_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.video_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE public.video_id_seq OWNER TO postgres;

--
-- Name: video_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--
ALTER SEQUENCE public.video_id_seq OWNED BY public.video.id;

--
-- Name: videos; Type: VIEW; Schema: public; Owner: postgres
--
CREATE VIEW
    public.videos AS
SELECT
    video.id,
    video.start_timestamp,
    video.end_timestamp,
    video.duration,
    video.size,
    video.file_path,
    video.camera_id,
    video.event_id
FROM
    public.video;

ALTER TABLE public.videos OWNER TO postgres;

--
-- Name: aggregated_detection id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.aggregated_detection
ALTER COLUMN id
SET DEFAULT nextval('public.aggregated_detection_id_seq'::regclass);

--
-- Name: camera id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.camera
ALTER COLUMN id
SET DEFAULT nextval('public.camera_id_seq'::regclass);

--
-- Name: detection id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.detection
ALTER COLUMN id
SET DEFAULT nextval('public.detection_id_seq'::regclass);

--
-- Name: event id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.event
ALTER COLUMN id
SET DEFAULT nextval('public.event_id_seq'::regclass);

--
-- Name: image id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.image
ALTER COLUMN id
SET DEFAULT nextval('public.image_id_seq'::regclass);

--
-- Name: object id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.object
ALTER COLUMN id
SET DEFAULT nextval('public.object_id_seq'::regclass);

--
-- Name: video id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.video
ALTER COLUMN id
SET DEFAULT nextval('public.video_id_seq'::regclass);

--
-- Name: aggregated_detection aggregated_detection_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.aggregated_detection
ADD CONSTRAINT aggregated_detection_pkey PRIMARY KEY (id);

--
-- Name: camera camera_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.camera
ADD CONSTRAINT camera_name_key UNIQUE (name);

--
-- Name: camera camera_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.camera
ADD CONSTRAINT camera_pkey PRIMARY KEY (id);

--
-- Name: detection detection_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_pkey PRIMARY KEY (id);

--
-- Name: event event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.event
ADD CONSTRAINT event_pkey PRIMARY KEY (id);

--
-- Name: image image_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.image
ADD CONSTRAINT image_pkey PRIMARY KEY (id);

--
-- Name: object object_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.object
ADD CONSTRAINT object_pkey PRIMARY KEY (id);

--
-- Name: video video_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.video
ADD CONSTRAINT video_pkey PRIMARY KEY (id);

--
-- Name: detection_class_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX detection_class_id_idx ON public.detection USING btree (class_id);

--
-- Name: detection_class_name_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX detection_class_name_idx ON public.detection USING btree (class_name);

--
-- Name: detection_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX detection_timestamp_idx ON public.detection USING btree ("timestamp");

--
-- Name: event_end_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX event_end_timestamp_idx ON public.event USING btree (end_timestamp);

--
-- Name: event_end_timestamp_source_camera_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX event_end_timestamp_source_camera_id_idx ON public.event USING btree (end_timestamp, source_camera_id);

--
-- Name: event_start_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX event_start_timestamp_idx ON public.event USING btree (start_timestamp);

--
-- Name: event_start_timestamp_source_camera_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX event_start_timestamp_source_camera_id_idx ON public.event USING btree (start_timestamp, source_camera_id);

--
-- Name: image_event_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX image_event_id_idx ON public.image USING btree (event_id);

--
-- Name: image_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX image_timestamp_idx ON public.image USING btree ("timestamp");

--
-- Name: object_class_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX object_class_id_idx ON public.object USING btree (class_id);

--
-- Name: object_class_name_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX object_class_name_idx ON public.object USING btree (class_name);

--
-- Name: object_end_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX object_end_timestamp_idx ON public.object USING btree (end_timestamp);

--
-- Name: object_event_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX object_event_id_idx ON public.object USING btree (event_id);

--
-- Name: object_start_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX object_start_timestamp_idx ON public.object USING btree (start_timestamp);

--
-- Name: video_end_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX video_end_timestamp_idx ON public.video USING btree (end_timestamp);

--
-- Name: video_event_id_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX video_event_id_idx ON public.video USING btree (event_id);

--
-- Name: video_start_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--
CREATE INDEX video_start_timestamp_idx ON public.video USING btree (start_timestamp);

--
-- Name: event aggregate_detection_trigger; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER aggregate_detection_trigger
AFTER
UPDATE ON public.event FOR EACH ROW WHEN ((new.status = 'needs tracking'::text))
EXECUTE FUNCTION public.aggregate_detection ();

--
-- Name: aggregated_detection aggregated_detection_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.aggregated_detection
ADD CONSTRAINT aggregated_detection_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: detection detection_camera_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: detection detection_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: detection detection_object_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.detection
ADD CONSTRAINT detection_object_id_fkey FOREIGN KEY (object_id) REFERENCES public.object (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: event event_original_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.event
ADD CONSTRAINT event_original_video_id_fkey FOREIGN KEY (original_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: event event_processed_video_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.event
ADD CONSTRAINT event_processed_video_id_fkey FOREIGN KEY (processed_video_id) REFERENCES public.video (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: event event_source_camera_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.event
ADD CONSTRAINT event_source_camera_id_fkey FOREIGN KEY (source_camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: event event_thumbnail_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.event
ADD CONSTRAINT event_thumbnail_image_id_fkey FOREIGN KEY (thumbnail_image_id) REFERENCES public.image (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: image image_camera_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.image
ADD CONSTRAINT image_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: image image_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.image
ADD CONSTRAINT image_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: object objet_camera_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.object
ADD CONSTRAINT objet_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: object objet_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.object
ADD CONSTRAINT objet_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: video video_camera_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.video
ADD CONSTRAINT video_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: video video_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY public.video
ADD CONSTRAINT video_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.event (id) ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--
REVOKE USAGE ON SCHEMA public
FROM
    PUBLIC;

GRANT ALL ON SCHEMA public TO PUBLIC;

--
-- PostgreSQL database dump complete
--
