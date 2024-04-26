-- Create "physical_things" table
CREATE TABLE "physical_things" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "external_id" text NULL,
  "name" text NOT NULL,
  "type" text NOT NULL,
  "tags" text[] NOT NULL DEFAULT '{}',
  "metadata" hstore NOT NULL DEFAULT ''::hstore,
  "raw_data" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "physical_things_external_id_key" UNIQUE ("external_id"),
  CONSTRAINT "physical_things_name_key" UNIQUE ("name"),
  CONSTRAINT "physical_things_external_id_check" CHECK (TRIM(BOTH FROM external_id) <> ''::text),
  CONSTRAINT "physical_things_name_check" CHECK (TRIM(BOTH FROM name) <> ''::text),
  CONSTRAINT "physical_things_type_check" CHECK (TRIM(BOTH FROM type) <> ''::text)
);
-- Create "location_history" table
CREATE TABLE "location_history" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "timestamp" timestamptz NOT NULL,
  "point" point NULL,
  "polygon" polygon NULL,
  "parent_physical_thing_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "location_history_parent_physical_thing_id_fkey" FOREIGN KEY ("parent_physical_thing_id") REFERENCES "physical_things" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "has_point_or_polygon_but_not_neither_and_not_both" CHECK (((point IS NOT NULL) AND (polygon IS NULL)) OR ((point IS NULL) AND (polygon IS NOT NULL)))
);
-- Create "logical_things" table
CREATE TABLE "logical_things" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "external_id" text NULL,
  "name" text NOT NULL,
  "type" text NOT NULL,
  "tags" text[] NOT NULL DEFAULT '{}',
  "metadata" hstore NOT NULL DEFAULT ''::hstore,
  "raw_data" jsonb NULL,
  "parent_physical_thing_id" uuid NULL,
  "parent_logical_thing_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "logical_things_external_id_key" UNIQUE ("external_id"),
  CONSTRAINT "logical_things_name_key" UNIQUE ("name"),
  CONSTRAINT "logical_things_parent_logical_thing_id_fkey" FOREIGN KEY ("parent_logical_thing_id") REFERENCES "logical_things" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "logical_things_parent_physical_thing_id_fkey" FOREIGN KEY ("parent_physical_thing_id") REFERENCES "physical_things" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "is_not_own_parent" CHECK (parent_logical_thing_id <> id),
  CONSTRAINT "logical_things_external_id_check" CHECK (TRIM(BOTH FROM external_id) <> ''::text),
  CONSTRAINT "logical_things_name_check" CHECK (TRIM(BOTH FROM name) <> ''::text),
  CONSTRAINT "logical_things_type_check" CHECK (TRIM(BOTH FROM type) <> ''::text)
);
