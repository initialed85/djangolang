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
