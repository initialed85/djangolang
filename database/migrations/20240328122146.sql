-- Create "aggregated_detection" table
CREATE TABLE "aggregated_detection" (
  "id" bigserial NOT NULL,
  "start_timestamp" timestamptz NOT NULL,
  "end_timestamp" timestamptz NOT NULL,
  "class_id" bigint NOT NULL,
  "class_name" text NOT NULL,
  "score" double precision NOT NULL,
  "count" bigint NOT NULL,
  "weighted_score" double precision NOT NULL,
  "event_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create "camera" table
CREATE TABLE "camera" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "stream_url" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "camera_name_key" to table: "camera"
CREATE UNIQUE INDEX "camera_name_key" ON "camera" ("name");
-- Create "detection" table
CREATE TABLE "detection" (
  "id" bigserial NOT NULL,
  "timestamp" timestamptz NOT NULL,
  "class_id" bigint NOT NULL,
  "class_name" text NOT NULL,
  "score" double precision NOT NULL,
  "centroid" point NOT NULL,
  "bounding_box" polygon NOT NULL,
  "camera_id" bigint NOT NULL,
  "event_id" bigint NULL,
  "object_id" bigint NULL,
  "colour" geometry(pointz) NULL,
  PRIMARY KEY ("id")
);
-- Create index "detection_class_id_idx" to table: "detection"
CREATE INDEX "detection_class_id_idx" ON "detection" ("class_id");
-- Create index "detection_class_name_idx" to table: "detection"
CREATE INDEX "detection_class_name_idx" ON "detection" ("class_name");
-- Create index "detection_timestamp_idx" to table: "detection"
CREATE INDEX "detection_timestamp_idx" ON "detection" ("timestamp");
-- Create "event" table
CREATE TABLE "event" (
  "id" bigserial NOT NULL,
  "start_timestamp" timestamptz NOT NULL,
  "end_timestamp" timestamptz NOT NULL,
  "duration" interval NOT NULL DEFAULT '00:00:00'::interval,
  "original_video_id" bigint NOT NULL,
  "thumbnail_image_id" bigint NOT NULL,
  "processed_video_id" bigint NULL,
  "source_camera_id" bigint NOT NULL,
  "status" text NOT NULL DEFAULT 'true',
  PRIMARY KEY ("id"),
  CONSTRAINT "event_status_check" CHECK (status = ANY (ARRAY['needs detection'::text, 'detection underway'::text, 'needs tracking'::text, 'tracking underway'::text, 'done'::text]))
);
-- Create index "event_end_timestamp_idx" to table: "event"
CREATE INDEX "event_end_timestamp_idx" ON "event" ("end_timestamp");
-- Create index "event_end_timestamp_source_camera_id_idx" to table: "event"
CREATE INDEX "event_end_timestamp_source_camera_id_idx" ON "event" ("end_timestamp", "source_camera_id");
-- Create index "event_start_timestamp_idx" to table: "event"
CREATE INDEX "event_start_timestamp_idx" ON "event" ("start_timestamp");
-- Create index "event_start_timestamp_source_camera_id_idx" to table: "event"
CREATE INDEX "event_start_timestamp_source_camera_id_idx" ON "event" ("start_timestamp", "source_camera_id");
-- Create "image" table
CREATE TABLE "image" (
  "id" bigserial NOT NULL,
  "timestamp" timestamptz NOT NULL,
  "size" double precision NOT NULL DEFAULT 0,
  "file_path" text NOT NULL,
  "camera_id" bigint NOT NULL,
  "event_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "image_event_id_idx" to table: "image"
CREATE INDEX "image_event_id_idx" ON "image" ("event_id");
-- Create index "image_timestamp_idx" to table: "image"
CREATE INDEX "image_timestamp_idx" ON "image" ("timestamp");
-- Create "object" table
CREATE TABLE "object" (
  "id" bigserial NOT NULL,
  "start_timestamp" timestamptz NOT NULL,
  "end_timestamp" timestamptz NOT NULL,
  "class_id" bigint NOT NULL,
  "class_name" text NOT NULL,
  "camera_id" bigint NOT NULL,
  "event_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "object_class_id_idx" to table: "object"
CREATE INDEX "object_class_id_idx" ON "object" ("class_id");
-- Create index "object_class_name_idx" to table: "object"
CREATE INDEX "object_class_name_idx" ON "object" ("class_name");
-- Create index "object_end_timestamp_idx" to table: "object"
CREATE INDEX "object_end_timestamp_idx" ON "object" ("end_timestamp");
-- Create index "object_event_id_idx" to table: "object"
CREATE INDEX "object_event_id_idx" ON "object" ("event_id");
-- Create index "object_start_timestamp_idx" to table: "object"
CREATE INDEX "object_start_timestamp_idx" ON "object" ("start_timestamp");
-- Create "video" table
CREATE TABLE "video" (
  "id" bigserial NOT NULL,
  "start_timestamp" timestamptz NOT NULL,
  "end_timestamp" timestamptz NOT NULL,
  "duration" interval NOT NULL DEFAULT '00:00:00'::interval,
  "size" double precision NOT NULL DEFAULT 0,
  "file_path" text NOT NULL,
  "camera_id" bigint NOT NULL,
  "event_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "video_end_timestamp_idx" to table: "video"
CREATE INDEX "video_end_timestamp_idx" ON "video" ("end_timestamp");
-- Create index "video_event_id_idx" to table: "video"
CREATE INDEX "video_event_id_idx" ON "video" ("event_id");
-- Create index "video_start_timestamp_idx" to table: "video"
CREATE INDEX "video_start_timestamp_idx" ON "video" ("start_timestamp");
-- Modify "aggregated_detection" table
ALTER TABLE "aggregated_detection" ADD
 CONSTRAINT "aggregated_detection_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "detection" table
ALTER TABLE "detection" ADD
 CONSTRAINT "detection_camera_id_fkey" FOREIGN KEY ("camera_id") REFERENCES "camera" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "detection_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "detection_object_id_fkey" FOREIGN KEY ("object_id") REFERENCES "object" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "event" table
ALTER TABLE "event" ADD
 CONSTRAINT "event_original_video_id_fkey" FOREIGN KEY ("original_video_id") REFERENCES "video" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "event_processed_video_id_fkey" FOREIGN KEY ("processed_video_id") REFERENCES "video" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "event_source_camera_id_fkey" FOREIGN KEY ("source_camera_id") REFERENCES "camera" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "event_thumbnail_image_id_fkey" FOREIGN KEY ("thumbnail_image_id") REFERENCES "image" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "image" table
ALTER TABLE "image" ADD
 CONSTRAINT "image_camera_id_fkey" FOREIGN KEY ("camera_id") REFERENCES "camera" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "image_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "object" table
ALTER TABLE "object" ADD
 CONSTRAINT "objet_camera_id_fkey" FOREIGN KEY ("camera_id") REFERENCES "camera" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "objet_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
-- Modify "video" table
ALTER TABLE "video" ADD
 CONSTRAINT "video_camera_id_fkey" FOREIGN KEY ("camera_id") REFERENCES "camera" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT, ADD
 CONSTRAINT "video_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "event" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
