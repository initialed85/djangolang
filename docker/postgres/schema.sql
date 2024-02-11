CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION postgis;

CREATE EXTENSION postgis_raster;

CREATE EXTENSION postgis_sfcgal;

CREATE EXTENSION fuzzystrmatch;

CREATE EXTENSION address_standardizer;

CREATE EXTENSION address_standardizer_data_us;

CREATE EXTENSION postgis_tiger_geocoder;

CREATE EXTENSION postgis_topology;

CREATE TABLE
    "firsts" (
        "id" uuid,
        "name" varchar NOT NULL UNIQUE,
        "type" varchar NOT NULL,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    "seconds" (
        "id" uuid,
        "name" varchar NOT NULL,
        "first_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "first_id" FOREIGN KEY ("first_id") REFERENCES "firsts" ("id")
    );
