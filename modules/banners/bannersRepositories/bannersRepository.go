package bannersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/banners"
	"github.com/yporn/sirarom-backend/modules/banners/bannersPatterns"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
)

type IBannersRepository interface {
	FindOneBanner(bannerId string) (*banners.Banner, error)
	FindBanner(req *banners.BannerFilter) ([]*banners.Banner, int)
	InsertBanner(req *banners.Banner) (*banners.Banner, error)
	UpdateBanner(req *banners.Banner) (*banners.Banner, error) 
	DeleteBanner(bannerId string) error
}

type bannersRepository struct {
	db            *sqlx.DB
	cfg           config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func BannersRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IBannersRepository {
	return &bannersRepository{
		db:            db,
		cfg:           cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *bannersRepository) FindOneBanner(bannerId string) (*banners.Banner, error) {
	query := `
		SELECT
			to_jsonb("t")
		FROM (
			SELECT
				"b"."id",
				"b"."index",
				"b"."delay",
				"b"."display",
				(
					SELECT
						COALESCE(array_to_json(array_agg("it")), '[]'::json)
					FROM (
						SELECT
							"i"."id",
							"i"."filename",
							"i"."url"
						FROM "banner_images" "i"
						WHERE "i"."banner_id" = "b"."id"
					) AS "it"
				) AS "images"
				FROM "banners" "b"
		WHERE "b"."id" = $1
		LIMIT 1
		) AS "t";
	`
	bannerBytes := make([]byte, 0)
	banner := &banners.Banner{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&bannerBytes, query, bannerId); err != nil {
		return nil, fmt.Errorf("get banner failed: %v", err)
	}

	if err := json.Unmarshal(bannerBytes, &banner); err != nil {
		return nil, fmt.Errorf("unmarshal banner failed: %v", err)
	}
	return banner, nil
}

func (r *bannersRepository) FindBanner(req *banners.BannerFilter) ([]*banners.Banner, int) {
	builder := bannersPatterns.FindBannerBuilder(r.db, req)
	engineer := bannersPatterns.FindBannerEngineer(builder)

	result := engineer.FindBanner().Result()
	count := engineer.CountBanner().Count()
	return result, count
}

func (r *bannersRepository) InsertBanner(req *banners.Banner) (*banners.Banner, error) {
	builder := bannersPatterns.InsertBannerBuilder(r.db, req)
	bannerId, err := bannersPatterns.InsertBannerEngineer(builder).InsertBanner()
	if err != nil {
		return nil, err
	}

	banner, err := r.FindOneBanner(bannerId)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (r *bannersRepository) UpdateBanner(req *banners.Banner) (*banners.Banner, error) {
	builder := bannersPatterns.UpdateBannerBuilder(r.db, req, r.filesUsecase)
	engineer := bannersPatterns.UpdateBannerEngineer(builder)

	if err := engineer.UpdateBanner(); err != nil {
		return nil, err
	}

	banner, err := r.FindOneBanner(strconv.Itoa(req.Id))
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (r *bannersRepository) DeleteBanner(bannerId string) error {
	query := `DELETE FROM "banners" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, bannerId); err != nil {
		return fmt.Errorf("delete banner failed: %v", err)
	}
	return nil
}