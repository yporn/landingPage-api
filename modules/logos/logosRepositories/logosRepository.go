package logosRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/logos"
	"github.com/yporn/sirarom-backend/modules/logos/logosPatterns"
)

type ILogosRepository interface {
	FindOneLogo(logoId string) (*logos.Logo, error)
	FindLogo(req *logos.LogoFilter) ([]*logos.Logo, int)
	InsertLogo(req *logos.Logo) (*logos.Logo, error)
	UpdateLogo(req *logos.Logo) (*logos.Logo, error)
	DeleteLogo(logoId string) error
}

type logosRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func LogosRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) ILogosRepository {
	return &logosRepository{
		db: db,
		cfg: cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *logosRepository) FindOneLogo(logoId string) (*logos.Logo, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"l".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "logo_images" "i"
					WHERE "i"."logo_id" = "l"."id"
				) AS "it"
			) AS "images"
		FROM "logos" "l"
		WHERE "l"."id" = $1
		LIMIT 1
	) AS "t";`

	logoBytes := make([]byte, 0)
	logo := &logos.Logo{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&logoBytes, query, logoId); err != nil {
		return nil, fmt.Errorf("get logo failed: %v", err)
	}
	if err := json.Unmarshal(logoBytes, &logo); err != nil {
		return nil, fmt.Errorf("unmarshal logo failed: %v", err)
	}
	return logo, nil
}


func (r *logosRepository) FindLogo(req *logos.LogoFilter) ([]*logos.Logo, int) {
	builder := logosPatterns.FindLogoBuilder(r.db, req)
	engineer := logosPatterns.FindLogoEngineer(builder)

	result := engineer.FindLogo().Result()
	count := engineer.CountLogo().Count()
	return result, count
}


func (r *logosRepository) InsertLogo(req *logos.Logo) (*logos.Logo, error) {
	builder := logosPatterns.InsertLogoBuilder(r.db, req)
	logoId, err := logosPatterns.InsertLogoEngineer(builder).InsertLogo()
	if err != nil {
		return nil, err
	}

	logo, err := r.FindOneLogo(logoId)
	if err != nil {
		return nil, err
	}
	return logo, nil
}

func (r *logosRepository) UpdateLogo(req *logos.Logo) (*logos.Logo, error) {
	builder := logosPatterns.UpdateLogoBuilder(r.db, req, r.filesUsecase)
	engineer := logosPatterns.UpdateLogoEngineer(builder)

	if err := engineer.UpdateLogo(); err != nil {
		return nil, err
	}

	logo, err := r.FindOneLogo(strconv.Itoa(req.Id))
	if err != nil {
		return nil, err
	}
	return logo, nil
}

func (r *logosRepository) DeleteLogo(logoId string) error {
	query := `DELETE FROM "logos" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, logoId); err != nil {
		return fmt.Errorf("delete logo failed: %v", err)
	}
	return nil
}