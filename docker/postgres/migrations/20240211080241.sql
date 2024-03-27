CREATE EXTENSION postgis;

CREATE EXTENSION postgis_raster;

CREATE EXTENSION postgis_sfcgal;

CREATE EXTENSION fuzzystrmatch;

CREATE EXTENSION address_standardizer;

CREATE EXTENSION address_standardizer_data_us;

CREATE EXTENSION postgis_tiger_geocoder;

CREATE EXTENSION postgis_topology;

SET
  search_path TO public;

DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;

CREATE EXTENSION "uuid-ossp";

CREATE TABLE
  firsts (
    "id" uuid DEFAULT uuid_generate_v4 (),
    "name" varchar NOT NULL UNIQUE,
    "type" varchar NOT NULL,
    "raw_data_jsonb" jsonb,
    "raw_data_json" json,
    "geom" geometry,
    PRIMARY KEY ("id")
  );

CREATE TABLE
  seconds (
    "id" uuid DEFAULT uuid_generate_v4 (),
    "name" varchar NOT NULL,
    "type" varchar NOT NULL,
    "raw_data_jsonb" jsonb,
    "raw_data_json" json,
    "first_id" uuid NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "first_id" FOREIGN KEY ("first_id") REFERENCES "firsts" ("id")
  );

INSERT INTO
  firsts (
    name,
    type,
    raw_data_jsonb,
    raw_data_json,
    geom
  )
VALUES
  ('First 1', 'Type 1', null, null, 'POINT(0 0)');

INSERT INTO
  firsts (
    name,
    type,
    raw_data_jsonb,
    raw_data_json,
    geom
  )
VALUES
  (
    'First 2',
    'Type 1',
    'null'::jsonb,
    'null'::json,
    'LINESTRING (0 0, 1 1, 2 1, 2 2)'
  );

INSERT INTO
  firsts (
    name,
    type,
    raw_data_jsonb,
    raw_data_json,
    geom
  )
VALUES
  (
    'First 3',
    'Type 2',
    '[]'::jsonb,
    '[]'::json,
    'POLYGON((0 0, 1 0, 1 1, 0 1, 0 0))'
  );

INSERT INTO
  firsts (
    name,
    type,
    raw_data_jsonb,
    raw_data_json,
    geom
  )
VALUES
  (
    'First 4',
    'Type 2',
    '{}'::jsonb,
    '{}'::json,
    'POLYGON((0 0, 10 0, 10 10, 0 10, 0 0), (1 1, 1 2, 2 2, 2 1, 1 1))'
  );

INSERT INTO
  firsts (
    name,
    type,
    raw_data_jsonb,
    raw_data_json,
    geom
  )
VALUES
  (
    'First 5',
    'Type 3',
    '{}'::jsonb,
    '{}'::json,
    'GEOMETRYCOLLECTION (POINT(2 0), POLYGON((0 0, 1 0, 1 1, 0 1, 0 0)))'
  );

INSERT INTO
  seconds (
    name,
    type,
    first_id
  )
VALUES
  (
    'Second A',
    'Type A',
    (
      SELECT
        id
      FROM
        firsts
      WHERE
        name = 'First 1'
      LIMIT
        1
    )
  );

INSERT INTO
  seconds (
    name,
    type,
    first_id
  )
VALUES
  (
    'Second B',
    'Type A',
    (
      SELECT
        id
      FROM
        firsts
      WHERE
        name = 'First 1'
      LIMIT
        1
    )
  );

INSERT INTO
  seconds (
    name,
    type,
    first_id
  )
VALUES
  (
    'Second C',
    'Type B',
    (
      SELECT
        id
      FROM
        firsts
      WHERE
        name = 'First 2'
      LIMIT
        1
    )
  );

INSERT INTO
  seconds (
    name,
    type,
    first_id
  )
VALUES
  (
    'Second D',
    'Type B',
    (
      SELECT
        id
      FROM
        firsts
      WHERE
        name = 'First 2'
      LIMIT
        1
    )
  );
