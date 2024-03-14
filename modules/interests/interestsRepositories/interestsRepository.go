package interestsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/interests"
	"github.com/yporn/sirarom-backend/modules/interests/interestsPatterns"
)

type IInterestRepository interface {
	FindOneInterest(interestId string) (*interests.Interest, error)
	// FindInterest(req *products.ProductFilter) ([]*products.Product, int)
	InsertInterest(req *interests.Interest) (*interests.Interest, error)
	// UpdateInterest(req *interests.Interest) (*interests.Interest, error)
	// DeleteInterest(interestId string) error
}

type interestsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func InterestsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IInterestRepository {
	return &interestsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *interestsRepository) FindOneInterest(interestId string) (*interests.Interest, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"bi".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "interest_images" "i"
					WHERE "i"."interest_id" = "bi"."id"
				) AS "it"
			) AS "images"
		FROM "interests" "bi"
		WHERE "id" = $1
		LIMIT 1
	) AS "t";`

	interestBytes := make([]byte, 0)
	interest := &interests.Interest{}

	if err := r.db.Get(&interestBytes, query, interestId); err != nil {
		return nil, fmt.Errorf("get interest failed: %v", err)
	}
	if err := json.Unmarshal(interestBytes, &interest); err != nil {
		return nil, fmt.Errorf("unmarshal interest failed: %v", err)
	}
	return interest, nil
}

func (r *interestsRepository) InsertInterest(req *interests.Interest) (*interests.Interest, error) {

	builder := interestsPatterns.InsertInterestBuilder(r.db, req)
	interestId, err := interestsPatterns.InsertInterestEngineer(builder).InsertInterest()
	if err != nil {
		return nil, err
	}

	interest, err := r.FindOneInterest(interestId)
	if err != nil {
		return nil, err
	}

	return interest, nil
}
