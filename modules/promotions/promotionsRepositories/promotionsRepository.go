package promotionsRepositories

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/promotions"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsPatterns"
)

type IPromotionsRepository interface {
	FindOnePromotion(promotionId string) (*promotions.Promotion, error)
	FindPromotion(req *promotions.PromotionFilter) ([]*promotions.Promotion, int)
	InsertPromotion(req *promotions.Promotion) (*promotions.Promotion, error)
	UpdatePromotion(req *promotions.Promotion) (*promotions.Promotion, error) 
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
			) AS "promotion_images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("pfi")), '[]'::json)
				FROM (
					SELECT
						"pfi".*
					FROM "promotion_free_items" "pfi"
					WHERE "pfi"."promotion_id" = "p"."id"
				) AS "pfi"
			) AS "free_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("phm")), '[]'::json)
				FROM (
					SELECT
						"phm"."id",
						"phm"."house_model_id",
						(
							SELECT
								COALESCE(array_to_json(array_agg("hm")), '[]'::json)
							FROM (
								SELECT
									"hm"."id",
									"hm"."name",
									(
										SELECT
											COALESCE(array_to_json(array_agg("hmt")), '[]'::json)
										FROM (
											SELECT
												"hmt"."id",
												"hmt"."room_type",
												"hmt"."amount"
											FROM "house_model_type_items" "hmt"
											WHERE "hmt"."house_model_id" = "hm"."id"
										) AS "hmt"
									) AS "house_type"
								FROM "house_models" "hm"
								WHERE "hm"."id" = "phm"."house_model_id"
							) AS "hm"
						) AS "house_model_name"
					FROM "promotion_house_models" "phm"
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

func (r *promotionsRepository) FindPromotion(req *promotions.PromotionFilter) ([]*promotions.Promotion, int) {
	builder := promotionsPatterns.FindPromotionBuilder(r.db, req)
	engineer := promotionsPatterns.FindPromotionEngineer(builder)

	result := engineer.FindPromotion().Result()
	count := engineer.CountPromotion().Count()
	return result, count
}

func (r *promotionsRepository) InsertPromotion(req *promotions.Promotion) (*promotions.Promotion, error) {
	builder := promotionsPatterns.InsertPromotionBuilder(r.db, req)
	promotionId, err := promotionsPatterns.InsertPromotionEngineer(builder).InsertPromotion()
	if err != nil {
		return nil, err
	}

	promotion, err := r.FindOnePromotion(promotionId)
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (r *promotionsRepository) UpdatePromotion(req *promotions.Promotion) (*promotions.Promotion, error) {
    // Initialize the update builder
    builder := promotionsPatterns.UpdatePromotionBuilder(r.db, req, r.filesUsecase)
    engineer := promotionsPatterns.UpdatePromotionEngineer(builder)

    // Perform the update
    if err := engineer.UpdatePromotion(); err != nil {
        return nil, err
    }

    // Retrieve the updated promotion
    promotion, err := r.FindOnePromotion(strconv.Itoa(req.Id))
    if err != nil {
        return nil, err
    }
    return promotion, nil
}
