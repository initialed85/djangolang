SET
    search_path TO public;

DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;

CREATE EXTENSION "uuid-ossp" SCHEMA public;

CREATE EXTENSION postgis;

CREATE EXTENSION postgis_raster;

CREATE EXTENSION postgis_sfcgal;

CREATE EXTENSION fuzzystrmatch;

CREATE EXTENSION address_standardizer;

CREATE EXTENSION address_standardizer_data_us;

CREATE EXTENSION postgis_tiger_geocoder;

CREATE EXTENSION postgis_topology;

CREATE TABLE
    public.firsts (
        "id" uuid DEFAULT public.uuid_generate_v4 (),
        "name" varchar NOT NULL UNIQUE,
        "type" varchar NOT NULL,
        "raw_data_jsonb" jsonb,
        "raw_data_json" json,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    public.seconds (
        "id" uuid DEFAULT public.uuid_generate_v4 (),
        "name" varchar NOT NULL,
        "type" varchar NOT NULL,
        "raw_data_jsonb" jsonb,
        "raw_data_json" json,
        "first_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "first_id" FOREIGN KEY ("first_id") REFERENCES "firsts" ("id")
    );

INSERT INTO
    public.firsts (
        name,
        type,
        raw_data_jsonb,
        raw_data_json
    )
VALUES
    ('First 1', 'Type 1', null, null);

INSERT INTO
    public.firsts (
        name,
        type,
        raw_data_jsonb,
        raw_data_json
    )
VALUES
    ('First 2', 'Type 1', 'null'::jsonb, 'null'::json);

INSERT INTO
    public.firsts (
        name,
        type,
        raw_data_jsonb,
        raw_data_json
    )
VALUES
    ('First 3', 'Type 2', '[]'::jsonb, '[]'::json);

INSERT INTO
    firsts (
        name,
        type,
        raw_data_jsonb,
        raw_data_json
    )
VALUES
    ('First 4', 'Type 3', '{}'::jsonb, '{}'::json);
