BEGIN;

TRUNCATE TABLE "users" CASCADE;
TRUNCATE TABLE "oauth" CASCADE;
TRUNCATE TABLE "roles" CASCADE;
TRUNCATE TABLE "user_roles" CASCADE;
TRUNCATE TABLE "banners" CASCADE;
TRUNCATE TABLE "banner_images" CASCADE;
TRUNCATE TABLE "brands" CASCADE;
TRUNCATE TABLE "data_settings" CASCADE;
TRUNCATE TABLE "projects" CASCADE;
TRUNCATE TABLE "project_images" CASCADE;
TRUNCATE TABLE "projects_house_type_items" CASCADE;
TRUNCATE TABLE "project_desc_area_items" CASCADE;
TRUNCATE TABLE "project_facility_items" CASCADE;
TRUNCATE TABLE "house_models" CASCADE;
TRUNCATE TABLE "house_model_images" CASCADE;
TRUNCATE TABLE "house_model_type_items" CASCADE;
TRUNCATE TABLE "house_model_plans" CASCADE;
TRUNCATE TABLE "house_model_plans_items" CASCADE;
TRUNCATE TABLE "house_model_plan_images" CASCADE;
TRUNCATE TABLE "interests" CASCADE;
TRUNCATE TABLE "interest_images" CASCADE;
TRUNCATE TABLE "careers" CASCADE;
TRUNCATE TABLE "promotions" CASCADE;
TRUNCATE TABLE "promotion_images" CASCADE;
TRUNCATE TABLE "promotion_free_items" CASCADE;
TRUNCATE TABLE "promotion_house_models" CASCADE;
TRUNCATE TABLE "activities" CASCADE;
TRUNCATE TABLE "activities_images" CASCADE;
TRUNCATE TABLE "activity_logs" CASCADE;
COMMIT;