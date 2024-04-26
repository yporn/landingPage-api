package seoRepositories

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/seo"
	"github.com/yporn/sirarom-backend/modules/seo/seoPatterns"
)

type ISeoRepository interface {
	FindOneSeo(seoId string) (*seo.Seo, error)
	UpdateSeo(req *seo.Seo) (*seo.Seo, error)
}

type seoRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
}

func SeoRepository(db *sqlx.DB, cfg config.IConfig) ISeoRepository {
	return &seoRepository{
		db:           db,
		cfg:          cfg,
	}
}

func (r *seoRepository) FindOneSeo(seoId string) (*seo.Seo, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"s".*
		FROM "seo" "s"
		WHERE "s"."id" = $1
		LIMIT 1
	) AS "t";`

	seoBytes := make([]byte, 0)
	seo := &seo.Seo{}

	if err := r.db.Get(&seoBytes, query, seoId); err != nil {
		return nil, fmt.Errorf("get seo failed: %v", err)
	}
	if err := json.Unmarshal(seoBytes, seo); err != nil {
		return nil, fmt.Errorf("unmarshal seo failed: %v", err)
	}
	return seo, nil
}

func (r *seoRepository) UpdateSeo(req *seo.Seo) (*seo.Seo, error) {
	builder := seoPatterns.UpdateSeoBuilder(r.db, req)
	engineer := seoPatterns.UpdateSeoEngineer(builder)

	if err := engineer.UpdateSeo(); err != nil {
		return nil, err
	}

	seoId := strconv.Itoa(req.Id)

	seo, err := r.FindOneSeo(seoId)
	if err != nil {
		return nil, err
	}

	return seo, nil
}
