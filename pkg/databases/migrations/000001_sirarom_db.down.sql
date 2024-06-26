BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "users"; 
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauth_table ON "oauth";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_banners_table ON "banners";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_banner_images_table ON "banner_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_brands_table ON "brands";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_projects_table ON "projects";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_project_images_table ON "project_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_projects_house_type_items_table ON "projects_house_type_items";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_project_desc_area_items_table ON "project_desc_area_items";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_project_facility_items_table ON "project_facility_items";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_house_models_table ON "house_models";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_house_model_images_table ON "house_model_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_house_model_type_items_table ON "house_model_type_items";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_house_model_plans_table ON "house_model_plans";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_house_model_plans_items_table ON "house_model_plan_items";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_house_model_plans_items_table ON "house_model_plan_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_interests_table ON "interests";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_interest_images_table ON "interest_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_careers_table ON "careers";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_promotions_table ON "promotions";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_promotion_images_table ON "promotion_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_promotions_free_items_table ON "promotion_free_items";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_promotions_house_models_table ON "promotion_house_models";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_activities_table ON "activities";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_activities_images_table ON "activities_images";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_activity_logs_table ON "activity_logs";
DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "oauth" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;
DROP TABLE IF EXISTS "user_roles" CASCADE;
DROP TABLE IF EXISTS "banners" CASCADE;
DROP TABLE IF EXISTS "banner_images" CASCADE;
DROP TABLE IF EXISTS "logos" CASCADE;
DROP TABLE IF EXISTS "logo_images" CASCADE;
DROP TABLE IF EXISTS "data_settings" CASCADE;
DROP TABLE IF EXISTS "projects" CASCADE;
DROP TABLE IF EXISTS "project_images" CASCADE;
DROP TABLE IF EXISTS "project_house_type_items" CASCADE;
DROP TABLE IF EXISTS "project_desc_area_items" CASCADE;
DROP TABLE IF EXISTS "project_facility_items" CASCADE;
DROP TABLE IF EXISTS "house_models" CASCADE;
DROP TABLE IF EXISTS "house_model_images" CASCADE;
DROP TABLE IF EXISTS "house_model_type_items" CASCADE;
DROP TABLE IF EXISTS "house_model_plans" CASCADE;
DROP TABLE IF EXISTS "house_model_plans_items" CASCADE;
DROP TABLE IF EXISTS "interests" CASCADE;
DROP TABLE IF EXISTS "interest_images" CASCADE;
DROP TABLE IF EXISTS "careers" CASCADE;
DROP TABLE IF EXISTS "promotions" CASCADE;
DROP TABLE IF EXISTS "promotion_images" CASCADE;
DROP TABLE IF EXISTS "promotion_free_items" CASCADE;
DROP TABLE IF EXISTS "promotion_house_models" CASCADE;
DROP TABLE IF EXISTS "activities" CASCADE;
DROP TABLE IF EXISTS "activities_images" CASCADE;
DROP TABLE IF EXISTS "activity_logs" CASCADE;


DROP TYPE IF EXISTS "display";
DROP TYPE IF EXISTS "status_project";
DROP TYPE IF EXISTS "type_project";
DROP TYPE IF EXISTS "status_career";

COMMIT;