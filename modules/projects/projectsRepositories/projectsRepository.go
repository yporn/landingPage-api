package projectsRepositories

import (
	"context"
	"encoding/json"
	"fmt"

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
	DeleteProject(projectId string) error
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
					FROM "project_comfortable_items" "ci"
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
			) AS "images"
			FROM "projects" "p"
		WHERE "p"."id" = $1
	) AS "t";
	`
	projectBytes := make([]byte, 0)
	project := &projects.Project{
		Images:          make([]*entities.Image, 0),
		HouseTypeItem:   make([]*projects.ProjectHouseTypeItem, 0),
		DescAreaItem:    make([]*projects.ProjectDescAreaItem, 0),
		ComfortableItem: make([]*projects.ProjectComfortableItem, 0),
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

func (r *projectsRepository) DeleteProject(projectId string) error {
	query := `DELETE FROM "projects" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, projectId); err != nil {
		return fmt.Errorf("delete project failed: %v", err)
	}
	return nil
}