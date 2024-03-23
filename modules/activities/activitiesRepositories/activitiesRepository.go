package activitiesRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"

	"github.com/yporn/sirarom-backend/modules/activities"
	"github.com/yporn/sirarom-backend/modules/activities/activitiesPatterns"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
)

type IActivitiesRepository interface {
	FindOneActivity(activityId string) (*activities.Activity, error)
	FindActivity(req *activities.ActivityFilter) ([]*activities.Activity, int)
	InsertActivity(req *activities.Activity) (*activities.Activity, error)
	DeleteActivity(activityId string) error
	UpdateActivity(req *activities.Activity) (*activities.Activity, error)
}

type activitiesRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func ActivitiesRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IActivitiesRepository {
	return &activitiesRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *activitiesRepository) FindOneActivity(activityId string) (*activities.Activity, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"a"."id",
			"a"."index",
			"a"."heading",
			"a"."description",
			"a"."start_date",
			"a"."end_date",
			"a"."video_link",
			"a"."display",
			"a"."created_at",
			"a"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "activities_images" "i"
					WHERE "i"."activity_id" = "a"."id"
				) AS "it"
			) AS "images"
		FROM "activities" "a"
		WHERE "a"."id" = $1
		LIMIT 1
	) AS "t";`

	activityBytes := make([]byte, 0)
	activity := &activities.Activity{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&activityBytes, query, activityId); err != nil {
		return nil, fmt.Errorf("get activity failed: %v", err)
	}
	if err := json.Unmarshal(activityBytes, &activity); err != nil {
		return nil, fmt.Errorf("unmarshal activity failed: %v", err)
	}
	return activity, nil
}

func (r *activitiesRepository) FindActivity(req *activities.ActivityFilter) ([]*activities.Activity, int) {
	builder := activitiesPatterns.FindActivityBuilder(r.db, req)
	engineer := activitiesPatterns.FindActivityEngineer(builder)

	result := engineer.FindActivity().Result()
	count := engineer.CountActivity().Count()
	return result, count
}

func (r *activitiesRepository) InsertActivity(req *activities.Activity) (*activities.Activity, error) {
	builder := activitiesPatterns.InsertActivityBuilder(r.db, req)
	activityId, err := activitiesPatterns.InsertActivityEngineer(builder).InsertActivity()
	if err != nil {
		return nil, err
	}

	activity, err := r.FindOneActivity(activityId)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activitiesRepository) UpdateActivity(req *activities.Activity) (*activities.Activity, error) {
	builder := activitiesPatterns.UpdateActivityBuilder(r.db, req, r.filesUsecase)
	engineer := activitiesPatterns.UpdateActivityEngineer(builder)

	if err := engineer.UpdateActivity(); err != nil {
		return nil, err
	}

	activity, err := r.FindOneActivity(strconv.Itoa(req.Id))
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activitiesRepository) DeleteActivity(activityId string) error {
	query := `DELETE FROM "activities" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, activityId); err != nil {
		return fmt.Errorf("delete activity failed: %v", err)
	}
	return nil
}