package generalRepositories

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/general"
	"github.com/yporn/sirarom-backend/modules/general/generalPatterns"
)

type IGeneralRepository interface {
	FindOneGeneral(generalId string) (*general.General, error)
	UpdateGeneral(req *general.General) (*general.General, error)
}


type generalRepository struct {
	db  *sqlx.DB
	cfg config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func GeneralRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IGeneralRepository {
	return &generalRepository{
		db:  db,
		cfg: cfg,
		filesUsecase: filesUsecase,
	}
}


func (r *generalRepository) FindOneGeneral(generalId string) (*general.General, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
		*
		FROM "data_settings"  
		WHERE "id" = $1
		LIMIT 1
	) AS "t";`

	generalBytes := make([]byte, 0)
	general := &general.General{}

	
	if err := r.db.Get(&generalBytes, query, generalId); err != nil {
		return nil, fmt.Errorf("get general failed: %v", err)
	}
	if err := json.Unmarshal(generalBytes, general); err != nil {
		return nil, fmt.Errorf("unmarshal general failed: %v", err)
	}
	return general, nil
}

func (r *generalRepository) UpdateGeneral(req *general.General) (*general.General, error) {
	builder := generalPatterns.UpdateGeneralBuilder(r.db, req, r.filesUsecase)
	engineer := generalPatterns.UpdateGeneralEngineer(builder)

	if err := engineer.UpdateGeneral(); err != nil {
		return nil, err
	}

	generalId := strconv.Itoa(req.Id)
	
	general, err := r.FindOneGeneral(generalId)
	if err != nil {
		return nil, err
	}

	return general, nil
}