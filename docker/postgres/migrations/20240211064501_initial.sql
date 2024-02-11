-- Create "firsts" table
CREATE TABLE "firsts" (
  "id" uuid NOT NULL,
  "name" character varying NOT NULL,
  "type" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "firsts_name_key" to table: "firsts"
CREATE UNIQUE INDEX "firsts_name_key" ON "firsts" ("name");
-- Create "thirds" table
CREATE TABLE "thirds" (
  "id" uuid NOT NULL,
  "name" character varying NOT NULL,
  "type" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "thirds_name_key" to table: "thirds"
CREATE UNIQUE INDEX "thirds_name_key" ON "thirds" ("name");
-- Create "firsts_thirds" table
CREATE TABLE "firsts_thirds" (
  "id" uuid NOT NULL,
  "first_id" uuid NOT NULL,
  "third_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "first_id" FOREIGN KEY ("first_id") REFERENCES "firsts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "third_id" FOREIGN KEY ("third_id") REFERENCES "thirds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "seconds" table
CREATE TABLE "seconds" (
  "id" uuid NOT NULL,
  "name" character varying NOT NULL,
  "first_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "first_id" FOREIGN KEY ("first_id") REFERENCES "firsts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "seconds_thirds" table
CREATE TABLE "seconds_thirds" (
  "id" uuid NOT NULL,
  "second_id" uuid NOT NULL,
  "third_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "second_id" FOREIGN KEY ("second_id") REFERENCES "seconds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "third_id" FOREIGN KEY ("third_id") REFERENCES "thirds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
