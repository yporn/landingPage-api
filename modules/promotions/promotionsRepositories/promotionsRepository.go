package promotionsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/promotions"
)

type IPromotionsRepository interface {
	FindOnePromotion(promotionId string) (*promotions.Promotion, error)
}

type promotionsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func PromotionsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IPromotionsRepository {
	return &promotionsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *promotionsRepository) FindOnePromotion(promotionId string) (*promotions.Promotion, error) {
	query := `
	SELECT to_jsonb("t")
	FROM (
		SELECT
			"p".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("i")), '[]'::json)
				FROM (
					SELECT
						"i".*
					FROM "promotion_images" "i"
					WHERE "i"."promotion_id" = "p"."id"
				) AS "i"
			) AS "house_images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("phm")), '[]'::json)
				FROM (
					SELECT
						"phm"."id",
						"hm"."name",
						"hmt"."room_type",
						"hmt"."amount"
					FROM "promotion_house_models" "phm"
					LEFT JOIN "house_models" "hm" ON "phm"."house_model_id" = "hm"."id"
					LEFT JOIN "house_model_type_items" "hmt" ON "phm"."house_model_id" = "hmt"."house_model_id"
					WHERE "phm"."promotion_id" = "p"."id"
				) AS "phm"
			) AS "house_models"
			FROM "promotions" "p"
		WHERE "p"."id" = $1
	) AS "t";
	`
	var promotionJSON []byte

	err := r.db.QueryRow(query, promotionId).Scan(&promotionJSON)
	if err != nil {
		return nil, fmt.Errorf("get promotion failed: %v", err)
	}

	var promotion promotions.Promotion
	err = json.Unmarshal(promotionJSON, &promotion)
	if err != nil {
		return nil, fmt.Errorf("unmarshal promotion failed: %v", err)
	}

	return &promotion, nil
}
