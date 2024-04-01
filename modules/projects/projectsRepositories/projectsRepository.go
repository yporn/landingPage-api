package projectsRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/projects"
	"github.com/yporn/sirarom-backend/modules/projects/projectsPatterns"
)

type IProjectRepository interface {
	FindOneProject(projectId string) (*projects.Project, error)
	FindProject(req *projects.ProjectFilter) ([]*projects.Project, int)
	InsertProject(req *projects.Project) (*projects.Project, error)
	UpdateProject(req *projects.Project) (*projects.Project, error)
	DeleteProject(projectId string) error
	FindProjectHouseModel(projectID string) (*projects.Project, error)
}

type projectsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func ProjectsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IProjectRepository {
	return &projectsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *projectsRepository) FindOneProject(projectId string) (*projects.Project, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p".*,
			(
				SELECT
				COALESCE(array_to_json(array_agg("htit")), '[]'::json)
				FROM (
					SELECT
						"hti".*
					FROM "project_house_type_items" "hti"
					WHERE "hti"."project_id" = "p"."id"
					
				) AS "htit"
			) AS "house_type_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("dait")), '[]'::json)
				FROM (
					SELECT
						"dai".*
					FROM "project_desc_area_items" "dai"
					WHERE "dai"."project_id" = "p"."id"
					
				) AS "dait"
			) AS "area_items",
			(
				SELECT
				COALESCE(array_to_json(array_agg("cit")), '[]'::json)
				FROM (
					SELECT
						"ci".*
					FROM "project_facility_items" "ci"
					WHERE "ci"."project_id" = "p"."id"
					
				) AS "cit"
			) AS "facilities_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "project_images" "i"
					WHERE "i"."project_id" = "p"."id"
				) AS "it"
			) AS "images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("hm")), '[]'::json)
				FROM (
					SELECT
						"hm".*,
						(
							SELECT
								COALESCE(array_to_json(array_agg("hmti")), '[]'::json)
							FROM (
								SELECT
									"hmti".*
								FROM "house_model_type_items" "hmti"
								WHERE "hmti"."house_model_id" = "hm"."id"
							) AS "hmti"
						) AS "type_items",
						(
							SELECT
								COALESCE(array_to_json(array_agg("ihm")), '[]'::json)
							FROM (
								SELECT
									"ihm".*
								FROM "house_model_images" "ihm"
								WHERE "ihm"."house_model_id" = "hm"."id"
							) AS "ihm"
						) AS "house_images"
					FROM "house_models" "hm"
					WHERE "hm"."project_id" = "p"."id"
				) AS "hm"
			) AS "house_models"
			FROM "projects" "p"
		WHERE "p"."id" = $1
	) AS "t";
	`
	projectBytes := make([]byte, 0)
	project := &projects.Project{
		Images:        make([]*entities.Image, 0),
		HouseTypeItem: make([]*projects.ProjectHouseTypeItem, 0),
		DescAreaItem:  make([]*projects.ProjectDescAreaItem, 0),
		FacilityItem:  make([]*projects.ProjectFacilityItem, 0),
	}

	if err := r.db.Get(&projectBytes, query, projectId); err != nil {
		return nil, fmt.Errorf("get project failed: %v", err)
	}
	if err := json.Unmarshal(projectBytes, &project); err != nil {
		return nil, fmt.Errorf("unmarshal project failed: %v", err)
	}
	return project, nil
}

func (r *projectsRepository) FindProjectHouseModel(projectId string) (*projects.Project, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("hm")), '[]'::json)
				FROM (
					SELECT
						"hm".*,
						(
							SELECT
								COALESCE(array_to_json(array_agg("hmti")), '[]'::json)
							FROM (
								SELECT
									"hmti".*
								FROM "house_model_type_items" "hmti"
								WHERE "hmti"."house_model_id" = "hm"."id"
								AND ("hmti"."room_type" = 'ห้องนอน' OR "hmti"."room_type" = 'ห้องน้ำ' OR "hmti"."room_type" = 'ที่จอดรถ')
							) AS "hmti"
						) AS "type_items",
						(
							SELECT
								COALESCE(array_to_json(array_agg("ihm")), '[]'::json)
							FROM (
								SELECT
									"ihm".*
								FROM "house_model_images" "ihm"
								WHERE "ihm"."house_model_id" = "hm"."id"
							) AS "ihm"
						) AS "house_images"
					FROM "house_models" "hm"
					WHERE "hm"."project_id" = "p"."id"
				) AS "hm"
			) AS "house_models"
			FROM "projects" "p"
		WHERE "p"."id" = $1
	) AS "t";
	`
	projectBytes := make([]byte, 0)
	project := &projects.Project{
		Images:        make([]*entities.Image, 0),
		HouseTypeItem: make([]*projects.ProjectHouseTypeItem, 0),
		DescAreaItem:  make([]*projects.ProjectDescAreaItem, 0),
		FacilityItem:  make([]*projects.ProjectFacilityItem, 0),
	}

	if err := r.db.Get(&projectBytes, query, projectId); err != nil {
		return nil, fmt.Errorf("get project failed: %v", err)
	}
	if err := json.Unmarshal(projectBytes, &project); err != nil {
		return nil, fmt.Errorf("unmarshal project failed: %v", err)
	}
	return project, nil
}

func (r *projectsRepository) FindProject(req *projects.ProjectFilter) ([]*projects.Project, int) {
	builder := projectsPatterns.FindProjectBuilder(r.db, req)
	engineer := projectsPatterns.FindProjectEngineer(builder)
	return engineer.FindProject(), engineer.CountOrder()
}

func (r *projectsRepository) InsertProject(req *projects.Project) (*projects.Project, error) {
	builder := projectsPatterns.InsertProjectBuilder(r.db, req)
	projectId, err := projectsPatterns.InsertProjectEngineer(builder).InsertProject()
	if err != nil {
		return nil, err
	}

	project, err := r.FindOneProject(projectId)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectsRepository) UpdateProject(req *projects.Project) (*projects.Project, error) {
	builder := projectsPatterns.UpdateProjectBuilder(r.db, req, r.filesUsecase)
	engineer := projectsPatterns.UpdateProjectEngineer(builder)

	if err := engineer.UpdateProject(); err != nil {
		return nil, err
	}

	project, err := r.FindOneProject(strconv.Itoa(req.Id))
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectsRepository) DeleteProject(projectId string) error {
	query := `DELETE FROM "projects" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, projectId); err != nil {
		return fmt.Errorf("delete project failed: %v", err)
	}
	return nil
}
