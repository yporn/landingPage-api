package houseModelsRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/houseModels"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsPatterns"
)

type IHouseModelsRepository interface {
	FindOneHouseModel(houseId string) (*houseModels.HouseModel, error)
	FindAllHouseModels() ([]houseModels.HouseModelName, error)
	InsertHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error)
	FindHouseModel(projectId string, req *houseModels.HouseModelFilter) ([]*houseModels.HouseModel, int)
	UpdateHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error)
	DeleteHouseModel(houseId string) error
}


type houseModelsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func HouseModelsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IHouseModelsRepository {
	return &houseModelsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *houseModelsRepository) FindAllHouseModels() ([]houseModels.HouseModelName, error) {
    query := `
        SELECT 
			"hm"."id",
			"hm"."name"
        FROM "house_models" "hm";
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("get house models failed: %v", err)
    }
    defer rows.Close()

    var houseModelNames []houseModels.HouseModelName
    for rows.Next() {
        var houseModelName houseModels.HouseModelName
        if err := rows.Scan(&houseModelName.Id, &houseModelName.Name); err != nil {
            return nil, fmt.Errorf("scan house model name failed: %v", err)
        }

        houseModelNames = append(houseModelNames, houseModelName)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("row error: %v", err)
    }

    return houseModelNames, nil
}


func (r *houseModelsRepository) FindOneHouseModel(houseId string) (*houseModels.HouseModel, error) {
	query := `
	SELECT to_jsonb("t")
	FROM (
		SELECT
			"hm".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("hmi")), '[]'::json)
				FROM (
					SELECT
						"hmi".*
					FROM "house_model_type_items" "hmi"
					WHERE "hmi"."house_model_id" = "hm"."id"
				) AS "hmi"
			) AS "type_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("ihm")), '[]'::json)
				FROM (
					SELECT
						"ihm"."id",
						"ihm"."filename",
						"ihm"."url"
					FROM "house_model_images" "ihm"
					WHERE "ihm"."house_model_id" = "hm"."id"
				) AS "ihm"
			) AS "house_images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("hmp")), '[]'::json)
				FROM (
					SELECT
						"hmp".*,
						(
							SELECT
								COALESCE(array_to_json(array_agg("hmpi")), '[]'::json)
							FROM (
								SELECT
									"hmpi".*
								FROM "house_model_plan_items" "hmpi"
								WHERE "hmpi"."house_model_plan_id" = "hmp"."id"
							) AS "hmpi"
						) AS "plan_items",
						(
							SELECT
								COALESCE(array_to_json(array_agg("ihmp")), '[]'::json)
							FROM (
								SELECT
									"ihmp"."id",
									"ihmp"."filename",
									"ihmp"."url"
								FROM "house_model_plan_images" "ihmp"
								WHERE "ihmp"."house_model_plan_id" = "hmp"."id"
							) AS "ihmp"
						) AS "plan_images"
					FROM "house_model_plans" "hmp"
					WHERE "hmp"."house_model_id" = "hm"."id"
				) AS "hmp"
			) AS "house_plan"
			FROM "house_models" "hm"
		WHERE "hm"."id" = $1
	) AS "t";
	`
	var houseModelJSON []byte

	err := r.db.QueryRow(query, houseId).Scan(&houseModelJSON)
	if err != nil {
		return nil, fmt.Errorf("get house model failed: %v", err)
	}

	var houseModel houseModels.HouseModel
	err = json.Unmarshal(houseModelJSON, &houseModel)
	if err != nil {
		return nil, fmt.Errorf("unmarshal house model failed: %v", err)
	}

	return &houseModel, nil
}

func (r *houseModelsRepository) FindHouseModel(projectId string, req *houseModels.HouseModelFilter) ([]*houseModels.HouseModel, int) {
	builder := houseModelsPatterns.FindHouseModelBuilder(r.db, projectId, req)
	engineer := houseModelsPatterns.FindHouseModelEngineer(builder)
	return engineer.FindHouseModel().Result(), engineer.CountHouseModel().Count()
}

func (r *houseModelsRepository) InsertHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error) {
	planItems := []*houseModels.HouseModelPlanItem{}
	builder := houseModelsPatterns.InsertHouseModelBuilder(r.db, req, planItems)
	projectId, err := houseModelsPatterns.InsertProjectEngineer(builder).InsertHouseModel()
	if err != nil {
		return nil, err
	}

	houseModel, err := r.FindOneHouseModel(projectId)
	if err != nil {
		return nil, err
	}
	return houseModel, nil
}

func (r *houseModelsRepository) UpdateHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error) {
	builder := houseModelsPatterns.UpdateHouseModelBuilder(r.db, req, r.filesUsecase)
	engineer := houseModelsPatterns.UpdateHouseModelEngineer(builder)

	if err := engineer.UpdateHouseModel(); err != nil {
		return nil, err
	}

	houseModel, err := r.FindOneHouseModel(strconv.Itoa(req.Id))
	if err != nil {
		return nil, err
	}
	return houseModel, nil
}

func (r *houseModelsRepository) DeleteHouseModel(houseId string) error {
	query := `DELETE FROM "house_models" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, houseId); err != nil {
		return fmt.Errorf("delete house_models failed: %v", err)
	}
	return nil
}
