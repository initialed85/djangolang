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
        age interval NOT NULL DEFAULT interval '0 seconds',
        optional_age interval NULL,
        count int NOT NULL,
        optional_count int NULL,
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
-- not_null_fuzz
--
DROP TABLE IF EXISTS public.not_null_fuzz CASCADE;

CREATE TABLE
    public.not_null_fuzz (
        id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid (),
        some_bigint bigint NOT NULL DEFAULT 1,
        some_bigint_array bigint[] NOT NULL DEFAULT '{1}',
        some_boolean boolean NOT NULL DEFAULT true,
        some_boolean_array boolean[] NOT NULL default '{true}',
        some_bytea bytea NOT NULL DEFAULT '\x65'::bytea,
        some_character_varying character varying NOT NULL DEFAULT 'A',
        some_character_varying_array character varying[] NOT NULL DEFAULT '{A}',
        some_double_precision double precision NOT NULL DEFAULT 1.0,
        some_double_precision_array double precision[] NOT NULL DEFAULT '{1.0}',
        some_float float NOT NULL DEFAULT 1.0,
        some_float_array float[] NOT NULL DEFAULT '{1.0}',
        -- some_geometry geometry NOT NULL,
        some_geometry_point_z geometry (PointZ) NOT NULL DEFAULT ST_PointZ (1.337, 69.420, 800.8135),
        some_hstore hstore NOT NULL DEFAULT 'A=>1'::hstore,
        some_inet inet NOT NULL DEFAULT '192.168.137.222/24'::inet,
        some_integer integer NOT NULL DEFAULT 1,
        some_integer_array integer[] NOT NULL DEFAULT '{1}',
        some_interval interval NOT NULL DEFAULT interval '1337 milliseconds',
        some_json json NOT NULL DEFAULT '{"some": "data"}'::json,
        some_jsonb jsonb NOT NULL DEFAULT '{"some": "data"}'::jsonb,
        some_numeric numeric NOT NULL DEFAULT 1.0,
        some_numeric_array numeric[] NOT NULL DEFAULT '{1.0}',
        some_point point NOT NULL DEFAULT ST_MakePoint (1.337, 69.420)::point,
        some_polygon polygon NOT NULL DEFAULT ST_MakePolygon (ST_GeomFromText ('LINESTRING(75 29,77 29,77 29, 75 29)'))::polygon,
        some_real real NOT NULL DEFAULT 1.0,
        some_real_array real[] NOT NULL DEFAULT '{1.0}',
        some_smallint smallint NOT NULL DEFAULT 1,
        some_smallint_array smallint[] NOT NULL DEFAULT '{1}',
        some_text text NOT NULL DEFAULT 'A',
        some_text_array text[] NOT NULL DEFAULT '{A}',
        some_timestamptz timestamp with time zone NOT NULL DEFAULT '2024-07-19T11:45:00+08:00',
        some_timestamp timestamp without time zone NOT NULL DEFAULT '2020-03-27T08:30:00',
        some_tsvector tsvector NOT NULL DEFAULT 'a'::tsvector,
        some_uuid uuid NOT NULL DEFAULT '11111111-1111-1111-1111-111111111111'::uuid
    );

ALTER TABLE public.not_null_fuzz OWNER TO postgres;

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

--
-- camry stuff
--
--
-- camera
--
DROP TABLE IF EXISTS public.camera CASCADE;

CREATE TABLE
    public.camera (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        name text NOT NULL CHECK (trim(name) != ''),
        stream_url text NOT NULL CHECK (trim(stream_url) != ''),
        last_seen timestamptz NULL DEFAULT NULL
    );

ALTER TABLE public.camera OWNER TO postgres;

CREATE UNIQUE INDEX camera_unique_name_not_deleted ON public.camera (name)
WHERE
    deleted_at IS null;

CREATE UNIQUE INDEX camera_unique_name_deleted ON public.camera (name, deleted_at)
WHERE
    deleted_at IS NOT null;

--
-- video
--
DROP TABLE IF EXISTS public.video CASCADE;

CREATE TABLE
    public.video (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        file_name text NOT NULL CHECK (trim(file_name) != ''),
        started_at timestamptz NOT NULL,
        ended_at timestamptz NULL,
        duration interval NULL,
        file_size float NULL,
        thumbnail_name text NULL,
        status text NULL,
        camera_id uuid NOT NULL REFERENCES public.camera (id)
    );

ALTER TABLE public.video OWNER TO postgres;

CREATE INDEX video_started_at ON public.video (started_at);

CREATE INDEX video_camera_id_started_at ON public.video (camera_id, started_at);

CREATE INDEX video_ended_at ON public.video (ended_at);

CREATE INDEX video_camera_id_ended_at ON public.video (camera_id, ended_at);

CREATE INDEX video_file_name ON public.video (file_name);

CREATE INDEX video_camera_id_file_name ON public.video (camera_id, file_name);

CREATE INDEX video_status ON public.video (status);

CREATE INDEX video_camera_id_status ON public.video (camera_id, status);

--
-- detection
--
DROP TABLE IF EXISTS public.detection CASCADE;

CREATE TABLE
    public.detection (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        seen_at timestamptz NOT NULL,
        class_id bigint NOT NULL,
        class_name text NOT NULL,
        score float NOT NULL,
        centroid Point NOT NULL,
        bounding_box Polygon NOT NULL,
        video_id uuid NOT NULL REFERENCES public.video (id),
        camera_id uuid NOT NULL REFERENCES public.camera (id)
    );

ALTER TABLE public.detection OWNER TO postgres;

CREATE INDEX detection_seen_at ON public.detection (seen_at);

CREATE INDEX detection_class_id_seen_at ON public.detection (class_id, seen_at);

CREATE INDEX detection_class_name_seen_at ON public.detection (class_name, seen_at);

CREATE INDEX detection_video_id_seen_at ON public.detection (video_id, seen_at);

CREATE INDEX detection_video_id_class_id_seen_at ON public.detection (video_id, class_id, seen_at);

CREATE INDEX detection_video_id_class_name_seen_at ON public.detection (video_id, class_name, seen_at);

--
-- triggers for camera
--
CREATE
OR REPLACE FUNCTION create_camera () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_camera BEFORE INSERT ON camera FOR EACH ROW
EXECUTE PROCEDURE create_camera ();

CREATE
OR REPLACE FUNCTION update_camera () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_camera BEFORE
UPDATE ON camera FOR EACH ROW
EXECUTE PROCEDURE update_camera ();

CREATE RULE "delete_camera" AS ON DELETE TO "camera"
DO INSTEAD (
    UPDATE camera
    SET
        created_at = old.created_at,
        updated_at = now(),
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);

CREATE RULE "delete_camera_cascade_to_video" AS ON DELETE TO "camera"
DO ALSO (
    DELETE FROM video
    WHERE
        camera_id = old.id
        AND deleted_at IS null
);

CREATE RULE "delete_camera_cascade_to_detection" AS ON DELETE TO "camera"
DO ALSO (
    DELETE FROM detection
    WHERE
        camera_id = old.id
        AND deleted_at IS null
);

--
-- triggers for video
--
CREATE
OR REPLACE FUNCTION create_video () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_video BEFORE INSERT ON video FOR EACH ROW
EXECUTE PROCEDURE create_video ();

CREATE
OR REPLACE FUNCTION update_video () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_video BEFORE
UPDATE ON video FOR EACH ROW
EXECUTE PROCEDURE update_video ();

CREATE RULE "delete_video" AS ON DELETE TO "video"
DO INSTEAD (
    UPDATE video
    SET
        created_at = old.created_at,
        updated_at = old.updated_at,
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);

--
-- triggers for detection
--
CREATE
OR REPLACE FUNCTION create_detection () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_detection BEFORE INSERT ON detection FOR EACH ROW
EXECUTE PROCEDURE create_detection ();

CREATE
OR REPLACE FUNCTION update_detection () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_detection BEFORE
UPDATE ON detection FOR EACH ROW
EXECUTE PROCEDURE update_detection ();

CREATE RULE "delete_detection" AS ON DELETE TO "detection"
DO INSTEAD (
    UPDATE detection
    SET
        created_at = old.created_at,
        updated_at = old.updated_at,
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);
