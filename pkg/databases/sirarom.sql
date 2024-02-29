CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "username" varchar UNIQUE,
  "password" varchar,
  "email" varchar UNIQUE,
  "role_id" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "oauth" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "access_token" varchar,
  "refresh_token" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "roles" (
  "id" int PRIMARY KEY,
  "title" varchar
);

CREATE TABLE "banners" (
  "id" bigInt PRIMARY KEY,
  "number" int,
  "delay" int,
  "display" boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "banners_images" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "banner_id" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "logos" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "display" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "data_settings" (
  "id" varchar PRIMARY KEY,
  "tel" varchar,
  "email" varchar,
  "link_facebook" varchar,
  "link_instagram" varchar,
  "link_twitter" varchar,
  "link_tiktok" varchar,
  "link_line" varchar,
  "link_website" varchar
);

CREATE TABLE "projects" (
  "id" varchar PRIMARY KEY,
  "name" varchar,
  "sort" varchar,
  "heading" varchar,
  "text" varchar,
  "location" text,
  "price" int,
  "status_project" varchar,
  "type_house_estate" varchar,
  "description" text,
  "name_facebook" varchar,
  "link_facebook" varchar,
  "tel" varchar,
  "address" varchar,
  "link_location" varchar,
  "display" boolean,
  "status" varchar
);

CREATE TABLE "projects_images" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "house_estate_id" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "projects_house_type_items" (
  "id" varchar PRIMARY KEY,
  "project_id" varchar,
  "type" varchar
);

CREATE TABLE "projects_desc_area_items" (
  "id" varchar PRIMARY KEY,
  "project_id" varchar,
  "item" varchar,
  "number" varchar,
  "unit" varchar
);

CREATE TABLE "projects_comfortable_items" (
  "id" varchar PRIMARY KEY,
  "project_id" varchar,
  "item" varchar
);

CREATE TABLE "house_plans" (
  "id" varchar PRIMARY KEY,
  "project_id" varchar,
  "name" varchar,
  "description" text,
  "link_video" varchar,
  "link_virtual_tour" varchar,
  "display" boolean,
  "sort" boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "house_plans_images" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "house_collection_plan_id" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "interests" (
  "id" varchar PRIMARY KEY,
  "bank_name" varchar,
  "interest_rate" varchar,
  "start_date" date,
  "end_date" date,
  "note" varchar,
  "display" boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "careers" (
  "id" varchar PRIMARY KEY,
  "position" varchar,
  "quantity" int,
  "description" text,
  "qualification" text,
  "start_date" date,
  "end_date" date,
  "status" varchar,
  "display" boolean,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "promotions" (
  "id" varchar PRIMARY KEY,
  "heading" varchar,
  "description" text,
  "start_date" date,
  "end_date" date,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "promotions_free_items" (
  "id" varchar PRIMARY KEY,
  "promotion_id" varchar,
  "description" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "promotions_house_plans" (
  "id" varchar PRIMARY KEY,
  "promotion_id" varchar,
  "house_collection_id" varchar
);

CREATE TABLE "activities" (
  "id" varchar PRIMARY KEY,
  "heading" varchar,
  "description" text,
  "start_date" date,
  "end_date" date,
  "video_link" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "activities_images" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "activity_id" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "banners_images" ADD FOREIGN KEY ("banner_id") REFERENCES "banners" ("id");

ALTER TABLE "projects_images" ADD FOREIGN KEY ("house_estate_id") REFERENCES "projects" ("id");

ALTER TABLE "projects_house_type_items" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "projects_desc_area_items" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "projects_comfortable_items" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "house_plans" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "house_plans_images" ADD FOREIGN KEY ("house_collection_plan_id") REFERENCES "house_plans" ("id");

ALTER TABLE "promotions_free_items" ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("id");

ALTER TABLE "promotions_house_plans" ADD FOREIGN KEY ("promotion_id") REFERENCES "promotions" ("id");

ALTER TABLE "promotions_house_plans" ADD FOREIGN KEY ("house_collection_id") REFERENCES "house_plans" ("id");

ALTER TABLE "activities_images" ADD FOREIGN KEY ("activity_id") REFERENCES "activities" ("id");
