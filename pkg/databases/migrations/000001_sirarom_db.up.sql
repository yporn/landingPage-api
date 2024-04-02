BEGIN;

--set timezone
SET TIME ZONE 'UTC';

--Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--Create sequence
CREATE SEQUENCE users_id_seq START WITH 1 INCREMENT BY 1;

--Auto update
CREATE OR REPLACE FUNCTION set_updated_at_column() 
RETURNS TRIGGER AS 
$$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$
language
'plpgsql'; 

--Create enum
CREATE TYPE "display" AS ENUM('published', 'unpublished');

CREATE TYPE "status_project" AS ENUM('ready', 'new');

CREATE TYPE "type_project" AS ENUM('present', 'future');

CREATE TYPE "status_career" AS ENUM('opening', 'closed');

CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY, 
    "username" VARCHAR UNIQUE NOT NULL, 
    "password" VARCHAR, 
    "name" VARCHAR, 
    "tel" VARCHAR, 
    "email" VARCHAR UNIQUE, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "user_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "user_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "oauth" (
    "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4 (), 
    "user_id" INTEGER NOT NULL, 
    "access_token" VARCHAR NOT NULL, 
    "refresh_token" VARCHAR NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "roles" (
    "id" SERIAL PRIMARY KEY, 
    "title" VARCHAR NOT NULL UNIQUE
);

CREATE TABLE "user_roles" (
    "id" SERIAL PRIMARY KEY, 
    "user_id" INTEGER,
    "role_id" INTEGER
);

CREATE TABLE "banners" (
    "id" SERIAL PRIMARY KEY, 
    "index" INTEGER, 
    "delay" INTEGER, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "banner_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "banner_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "logos" (
    "id" SERIAL PRIMARY KEY, 
    "name" VARCHAR, 
    "index" INTEGER, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "logo_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "logo_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "data_settings" (
    "id" SERIAL PRIMARY KEY, 
    "tel" VARCHAR, 
    "email" VARCHAR, 
    "link_facebook" VARCHAR, 
    "link_instagram" VARCHAR, 
    "link_twitter" VARCHAR, 
    "link_tiktok" VARCHAR, 
    "link_line" VARCHAR, 
    "link_website" VARCHAR 
);

CREATE TABLE "data_setting_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "data_setting_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "projects" (
    "id" SERIAL PRIMARY KEY, 
    "name" VARCHAR, 
    "index" INTEGER, 
    "heading" VARCHAR, 
    "text" VARCHAR, 
    "location" TEXT, 
    "price" FLOAT, 
    "status_project" status_project NOT NULL, 
    "type_project" type_project NOT NULL, 
    "description" TEXT, 
    "name_facebook" VARCHAR, 
    "link_facebook" VARCHAR, 
    "tel" VARCHAR, 
    "address" TEXT, 
    "link_location" VARCHAR, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "project_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "project_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "project_house_type_items" (
    "id" SERIAL PRIMARY KEY, 
    "project_id" INTEGER, 
    "name" VARCHAR, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "project_desc_area_items" (
    "id" SERIAL PRIMARY KEY, 
    "project_id" INTEGER, 
    "item" VARCHAR, 
    "amount" INTEGER, 
    "unit" VARCHAR, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "project_facility_items" (
    "id" SERIAL PRIMARY KEY, 
    "project_id" INTEGER, 
    "item" VARCHAR, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "house_models" (
    "id" SERIAL PRIMARY KEY, 
    "project_id" INTEGER, 
    "name" VARCHAR, 
    "description" TEXT, 
    "link_video" VARCHAR, 
    "link_virtual_tour" VARCHAR, 
    "display" display NOT NULL,
    "index" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "house_model_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "house_model_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "house_model_type_items" (
    "id" SERIAL PRIMARY KEY, 
    "house_model_id" INTEGER, 
    "room_type" VARCHAR, 
    "amount" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "house_model_plans" (
    "id" SERIAL PRIMARY KEY, 
    "house_model_id" INTEGER, 
    "floor" INTEGER, 
    "size" VARCHAR, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "house_model_plan_items" (
    "id" SERIAL PRIMARY KEY, 
    "house_model_plan_id" INTEGER, 
    "room_type" VARCHAR, 
    "amount" INTEGER,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "house_model_plan_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "house_model_plan_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "interests" (
    "id" SERIAL PRIMARY KEY, 
    "bank_name" VARCHAR, 
    "interest_rate" FLOAT, 
    "note" VARCHAR, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "interest_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "interest_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "careers" (
    "id" SERIAL PRIMARY KEY, 
    "position" VARCHAR, 
    "amount" INTEGER, 
    "location" VARCHAR, 
    "description" TEXT, 
    "qualification" TEXT, 
    "start_date" DATE, 
    "end_date" DATE, 
    "status" status_career NOT NULL, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "promotions" (
    "id" SERIAL PRIMARY KEY, 
    "index" INTEGER, 
    "heading" VARCHAR, 
    "description" VARCHAR, 
    "start_date" VARCHAR, 
    "end_date" VARCHAR, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "promotion_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "promotion_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "promotion_free_items" (
    "id" SERIAL PRIMARY KEY, 
    "promotion_id" INTEGER, 
    "description" VARCHAR, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "promotion_house_models" (
    "id" SERIAL PRIMARY KEY, 
    "promotion_id" INTEGER, 
    "house_model_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "activities" (
    "id" SERIAL PRIMARY KEY, 
    "index" INTEGER, 
    "heading" VARCHAR, 
    "description" TEXT, 
    "start_date" DATE, 
    "end_date" DATE, 
    "video_link" VARCHAR, 
    "display" display NOT NULL, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "activities_images" (
    "id" SERIAL PRIMARY KEY, 
    "filename" VARCHAR, 
    "url" VARCHAR, 
    "activity_id" INTEGER, 
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "activity_logs" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER,
    "action" VARCHAR,
    "details" VARCHAR,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(), 
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);


CREATE TABLE "seo" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR,
    "description" TEXT,
    "keywords" VARCHAR,
    "robot" VARCHAR,
    "google_bot" VARCHAR,
);

ALTER TABLE "user_images"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "user_roles"
ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
ALTER TABLE "user_roles"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "oauth"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "banner_images"
ADD FOREIGN KEY ("banner_id") REFERENCES "banners" ("id") ON DELETE CASCADE;

ALTER TABLE "logo_images"
ADD FOREIGN KEY ("logo_id") REFERENCES "logos" ("id") ON DELETE CASCADE;

ALTER TABLE "data_setting_images"
ADD FOREIGN KEY ("data_setting_id") REFERENCES "data_settings" ("id") ON DELETE CASCADE;

ALTER TABLE "project_images"
ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE;

ALTER TABLE "project_house_type_items"
ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE;

ALTER TABLE "project_desc_area_items"
ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE;

ALTER TABLE "project_facility_items"
ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE;

ALTER TABLE "house_models"
ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE;

ALTER TABLE "house_model_images"
ADD FOREIGN KEY ("house_model_id") REFERENCES "house_models" ("id") ON DELETE CASCADE;

ALTER TABLE "house_model_type_items"
ADD FOREIGN KEY ("house_model_id") REFERENCES "house_models" ("id") ON DELETE CASCADE;

ALTER TABLE "house_model_plans"
ADD FOREIGN KEY ("house_model_id") REFERENCES "house_models" ("id") ON DELETE CASCADE;

ALTER TABLE "house_model_plan_images"
ADD FOREIGN KEY ("house_model_plan_id") REFERENCES "house_model_plans" ("id") ON DELETE CASCADE;

ALTER TABLE "house_model_plan_items"
ADD FOREIGN KEY ("house_model_plan_id") REFERENCES "house_model_plans" ("id") ON DELETE CASCADE;

ALTER TABLE "interest_images"
ADD FOREIGN KEY ("interest_id") REFERENCES "interests" ("id") ON DELETE CASCADE;

ALTER TABLE "promotion_images"
ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("id") ON DELETE CASCADE;

ALTER TABLE "promotion_free_items"
ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("id") ON DELETE CASCADE;

ALTER TABLE "promotion_house_models"
ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("id") ON DELETE CASCADE;

ALTER TABLE "promotion_house_models"
ADD FOREIGN KEY ("house_model_id") REFERENCES "house_models" ("id") ON DELETE CASCADE;

ALTER TABLE "activities_images"
ADD FOREIGN KEY ("activity_id") REFERENCES "activities" ("id") ON DELETE CASCADE;
ALTER TABLE "activity_logs"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE
UPDATE ON "users" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();
CREATE TRIGGER set_updated_at_timestamp_user_images_table BEFORE
UPDATE ON "user_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();
CREATE TRIGGER set_updated_at_timestamp_oauth_table BEFORE
UPDATE ON "oauth" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_banners_table BEFORE
UPDATE ON "banners" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_banner_images_table BEFORE
UPDATE ON "banner_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_logos_table BEFORE
UPDATE ON "logos" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_projects_table BEFORE
UPDATE ON "projects" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_project_images_table BEFORE
UPDATE ON "project_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_project_house_type_items_table BEFORE
UPDATE ON "project_house_type_items" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_project_desc_area_items_table BEFORE
UPDATE ON "project_desc_area_items" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_project_facility_items_table BEFORE
UPDATE ON "project_facility_items" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_house_models_table BEFORE
UPDATE ON "house_models" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_house_model_images_table BEFORE
UPDATE ON "house_model_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_house_model_type_items_table BEFORE
UPDATE ON "house_model_type_items" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_house_model_plans_table BEFORE
UPDATE ON "house_model_plans" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_house_model_plan_items_table BEFORE
UPDATE ON "house_model_plan_items" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_house_model_plan_images_table BEFORE
UPDATE ON "house_model_plan_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_interests_table BEFORE
UPDATE ON "interests" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_interest_images_table BEFORE
UPDATE ON "interest_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_careers_table BEFORE
UPDATE ON "careers" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_promotions_table BEFORE
UPDATE ON "promotions" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_promotion_images_table BEFORE
UPDATE ON "promotion_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_promotion_free_items_table BEFORE
UPDATE ON "promotion_free_items" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_promotion_house_models_table BEFORE
UPDATE ON "promotion_house_models" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_activities_table BEFORE
UPDATE ON "activities" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();

CREATE TRIGGER set_updated_at_timestamp_activities_images_table BEFORE
UPDATE ON "activities_images" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();
CREATE TRIGGER set_updated_at_timestamp_activity_logs_table BEFORE
UPDATE ON "activity_logs" FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column ();
COMMIT;