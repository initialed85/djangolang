--
-- init
--
CREATE SCHEMA IF NOT EXISTS public;

ALTER SCHEMA public OWNER TO postgres;

COMMENT ON SCHEMA public IS 'standard public schema';

SET
    default_tablespace = '';

SET
    default_table_access_method = heap;

CREATE EXTENSION IF NOT EXISTS postgis SCHEMA public;

CREATE EXTENSION IF NOT EXISTS postgis_raster SCHEMA public;

SET
    postgis.gdal_enabled_drivers = 'ENABLE_ALL';

CREATE EXTENSION IF NOT EXISTS hstore SCHEMA public;

ALTER ROLE postgres
SET
    search_path TO public,
    postgis,
    hstore;

SET
    search_path TO public,
    postgis,
    hstore;

--
-- physical_things
--
DROP TABLE IF EXISTS public.physical_things CASCADE;

CREATE TABLE
    public.physical_things (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT NOW(),
        updated_at timestamptz NOT NULL DEFAULT NOW(),
        deleted_at timestamptz NULL DEFAULT NULL,
        external_id text NULL UNIQUE CHECK (trim(external_id) != ''),
        name text NOT NULL UNIQUE CHECK (trim(name) != ''),
        type text NOT NULL CHECK (
            trim(
                type
            ) != ''
        ),
        tags text[] NOT NULL DEFAULT '{}',
        metadata hstore NOT NULL DEFAULT ''::hstore,
        raw_data jsonb NULL
    );

ALTER TABLE public.physical_things OWNER TO postgres;

--
-- logical_things
--
DROP TABLE IF EXISTS public.logical_things CASCADE;

CREATE TABLE
    public.logical_things (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT NOW(),
        updated_at timestamptz NOT NULL DEFAULT NOW(),
        deleted_at timestamptz NULL DEFAULT NULL,
        external_id text NULL UNIQUE CHECK (trim(external_id) != ''),
        name text NOT NULL UNIQUE CHECK (trim(name) != ''),
        type text NOT NULL CHECK (
            trim(
                type
            ) != ''
        ),
        tags text[] NOT NULL DEFAULT '{}',
        metadata hstore NOT NULL DEFAULT ''::hstore,
        raw_data jsonb NULL,
        parent_physical_thing_id uuid NULL REFERENCES public.physical_things (id),
        parent_logical_thing_id uuid NULL REFERENCES public.logical_things (id),
        CONSTRAINT is_not_own_parent CHECK (parent_logical_thing_id != id)
    );

ALTER TABLE public.logical_things OWNER TO postgres;

--
-- location_history
--
DROP TABLE IF EXISTS public.location_history CASCADE;

CREATE TABLE
    public.location_history (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT NOW(),
        updated_at timestamptz NOT NULL DEFAULT NOW(),
        deleted_at timestamptz NULL DEFAULT NULL,
        timestamp timestamptz NOT NULL,
        point point NULL,
        polygon polygon NULL,
        parent_physical_thing_id uuid NULL REFERENCES public.physical_things (id),
        CONSTRAINT has_point_or_polygon_but_not_neither_and_not_both CHECK (
            (
                point IS NOT null
                AND polygon IS null
            )
            OR (
                point IS null
                AND polygon IS NOT null
            )
        )
    );

ALTER TABLE public.location_history OWNER TO postgres;

--
-- fuzz
--
DROP TABLE IF EXISTS public.fuzz CASCADE;

CREATE TABLE
    public.fuzz (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        column1 timestamp without time zone NULL,
        column2 timestamp with time zone NULL,
        column3 json NULL,
        column4 jsonb NULL,
        column5 character varying[] NULL,
        column6 text[] NULL,
        column7 character varying NULL,
        column8 text NULL,
        column9 smallint[] NULL,
        column10 integer[] NULL,
        column11 bigint[] NULL,
        column12 smallint NULL,
        column13 integer NULL,
        column14 bigint NULL,
        column15 real[] NULL,
        column16 float[] NULL,
        column17 numeric[] NULL,
        column18 double precision[] NULL,
        column19 float NULL,
        column20 real NULL,
        column21 numeric NULL,
        column22 double precision NULL,
        column23 boolean[] NULL,
        column24 boolean NULL,
        column25 tsvector NULL,
        column26 uuid NULL,
        column27 hstore NULL,
        column28 point NULL,
        column29 polygon NULL,
        column30 geometry NULL,
        column31 geometry (PointZ) NULL,
        column32 inet NULL,
        column33 bytea NULL
    );

ALTER TABLE public.fuzz OWNER TO postgres;
