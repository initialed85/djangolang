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
    "assets" (
        "id" uuid,
        "name" varchar NOT NULL UNIQUE,
        "type" varchar NOT NULL,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    "devices" (
        "id" uuid,
        "name" varchar NOT NULL,
        "asset_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "asset_id" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id")
    );

CREATE TABLE
    "tags" (
        "id" uuid,
        "name" varchar NOT NULL UNIQUE,
        "type" varchar NOT NULL,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    "assets_tags" (
        "id" uuid,
        "asset_id" uuid NOT NULL,
        "tag_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "asset_id" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id"),
        CONSTRAINT "tag_id" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id")
    );

CREATE TABLE
    "devices_tags" (
        "id" uuid,
        "device_id" uuid NOT NULL,
        "tag_id" uuid NOT NULL,
        PRIMARY KEY ("id"),
        CONSTRAINT "device_id" FOREIGN KEY ("device_id") REFERENCES "devices" ("id"),
        CONSTRAINT "tag_id" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id")
    );
