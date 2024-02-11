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

CREATE TABLE
    "thirds" (
        "id" uuid,
        "name" varchar NOT NULL UNIQUE,
        "type" varchar NOT NULL,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    "firsts_thirds" (
        "id" uuid,
        "first_id" uuid NOT NULL,
        "third_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "first_id" FOREIGN KEY ("first_id") REFERENCES "firsts" ("id"),
        CONSTRAINT "third_id" FOREIGN KEY ("third_id") REFERENCES "thirds" ("id")
    );

CREATE TABLE
    "seconds_thirds" (
        "id" uuid,
        "second_id" uuid NOT NULL,
        "third_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "second_id" FOREIGN KEY ("second_id") REFERENCES "seconds" ("id"),
        CONSTRAINT "third_id" FOREIGN KEY ("third_id") REFERENCES "thirds" ("id")
    );
