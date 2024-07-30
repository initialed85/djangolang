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
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        external_id text NULL CHECK (trim(external_id) != ''),
        name text NOT NULL CHECK (trim(name) != ''),
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

CREATE UNIQUE INDEX physical_things_unique_external_id_not_deleted ON public.physical_things (external_id)
WHERE
    deleted_at IS null;

CREATE UNIQUE INDEX physical_things_unique_external_id_deleted ON public.physical_things (external_id, deleted_at)
WHERE
    deleted_at IS NOT null;

CREATE UNIQUE INDEX physical_things_unique_name_not_deleted ON public.physical_things (name)
WHERE
    deleted_at IS null;

CREATE UNIQUE INDEX physical_things_unique_name_deleted ON public.physical_things (name, deleted_at)
WHERE
    deleted_at IS NOT null;

CREATE
OR REPLACE FUNCTION create_physical_things () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_physical_things BEFORE INSERT ON physical_things FOR EACH ROW
EXECUTE PROCEDURE create_physical_things ();

CREATE
OR REPLACE FUNCTION update_physical_things () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--
-- logical_things
--
DROP TABLE IF EXISTS public.logical_things CASCADE;

CREATE TABLE
    public.logical_things (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        external_id text NULL CHECK (trim(external_id) != ''),
        name text NOT NULL CHECK (trim(name) != ''),
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

CREATE UNIQUE INDEX logical_things_unique_external_id_not_deleted ON public.logical_things (external_id)
WHERE
    deleted_at IS null;

CREATE UNIQUE INDEX logical_things_unique_external_id_deleted ON public.logical_things (external_id, deleted_at)
WHERE
    deleted_at IS NOT null;

CREATE UNIQUE INDEX logical_things_unique_name_not_deleted ON public.logical_things (name)
WHERE
    deleted_at IS null;

CREATE UNIQUE INDEX logical_things_unique_name_deleted ON public.logical_things (name, deleted_at)
WHERE
    deleted_at IS NOT null;

--
-- location_history
--
DROP TABLE IF EXISTS public.location_history CASCADE;

CREATE TABLE
    public.location_history (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
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

--
-- triggers for physical_things
--
CREATE TRIGGER update_physical_things BEFORE
UPDATE ON physical_things FOR EACH ROW
EXECUTE PROCEDURE update_physical_things ();

CREATE RULE "delete_physical_things" AS ON DELETE TO "physical_things"
DO INSTEAD (
    UPDATE physical_things
    SET
        created_at = old.created_at,
        updated_at = now(),
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);

CREATE RULE "delete_physical_things_cascade_to_logical_things" AS ON DELETE TO "physical_things"
DO ALSO (
    DELETE FROM logical_things
    WHERE
        parent_physical_thing_id = old.id
        AND deleted_at IS null
);

CREATE RULE "delete_physical_things_cascade_to_location_history" AS ON DELETE TO "physical_things"
DO ALSO (
    DELETE FROM location_history
    WHERE
        parent_physical_thing_id = old.id
        AND deleted_at IS null
);

--
-- triggers for logical_things
--
CREATE
OR REPLACE FUNCTION create_logical_things () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_logical_things BEFORE INSERT ON logical_things FOR EACH ROW
EXECUTE PROCEDURE create_logical_things ();

CREATE
OR REPLACE FUNCTION update_logical_things () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_logical_things BEFORE
UPDATE ON logical_things FOR EACH ROW
EXECUTE PROCEDURE update_logical_things ();

CREATE RULE "delete_logical_things" AS ON DELETE TO "logical_things"
DO INSTEAD (
    UPDATE logical_things
    SET
        created_at = old.created_at,
        updated_at = old.updated_at,
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);

-- TODO
-- CREATE RULE "delete_logical_things_cascade_to_logical_things" AS ON DELETE TO "logical_things"
-- DO ALSO (
--     DELETE FROM logical_things
--     WHERE
--         parent_logical_thing_id = old.id
--         AND deleted_at IS null
--         AND id != old.id
--         AND pg_trigger_depth() < 1
-- );
--
-- triggers for location_history
--
CREATE
OR REPLACE FUNCTION create_location_history () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_location_history BEFORE INSERT ON location_history FOR EACH ROW
EXECUTE PROCEDURE create_location_history ();

CREATE
OR REPLACE FUNCTION update_location_history () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_location_history BEFORE
UPDATE ON location_history FOR EACH ROW
EXECUTE PROCEDURE update_location_history ();

CREATE RULE "delete_location_history" AS ON DELETE TO "location_history"
DO INSTEAD (
    UPDATE location_history
    SET
        created_at = old.created_at,
        updated_at = now(),
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);
