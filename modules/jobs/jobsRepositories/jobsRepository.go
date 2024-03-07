package jobsRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/jobs"
	"github.com/yporn/sirarom-backend/modules/jobs/jobsPatterns"
)

type IJobRepository interface {
	FindOneJob(jobId string) (*jobs.Job, error)
	FindJob(req *jobs.JobFilter) ([]*jobs.Job, int)
	InsertJob(req *jobs.Job) (*jobs.Job, error)
	UpdateJob(req *jobs.Job) (*jobs.Job, error)
	DeleteJob(jobId string) error
}

type jobsRepository struct {
	db  *sqlx.DB
	cfg config.IConfig
}

func JobsRepository(db *sqlx.DB, cfg config.IConfig) IJobRepository {
	return &jobsRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *jobsRepository) FindOneJob(jobId string) (*jobs.Job, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
		*
		FROM "careers"  
		WHERE "id" = $1
		LIMIT 1
	) AS "t";`

	jobBytes := make([]byte, 0)
	job := &jobs.Job{}

	
	if err := r.db.Get(&jobBytes, query, jobId); err != nil {
		return nil, fmt.Errorf("get job failed: %v", err)
	}
	if err := json.Unmarshal(jobBytes, job); err != nil {
		return nil, fmt.Errorf("unmarshal job failed: %v", err)
	}
	return job, nil
}

func (r *jobsRepository) FindJob(req *jobs.JobFilter) ([]*jobs.Job, int) {
	builder := jobsPatterns.FindJobBuilder(r.db, req)
	engineer := jobsPatterns.FindJobEngineer(builder)

	result := engineer.FindJob().Result()
	count := engineer.CountJob().Count()

	return result, count
}

func (r *jobsRepository) InsertJob(req *jobs.Job) (*jobs.Job, error) {
	builder := jobsPatterns.InsertJobBuilder(r.db, req)
	jobId, err := jobsPatterns.InsertJobEngineer(builder).InsertJob()
	if err != nil {
		return nil, err
	}

	job, err := r.FindOneJob(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobsRepository) UpdateJob(req *jobs.Job) (*jobs.Job, error) {
	builder := jobsPatterns.UpdateJobBuilder(r.db, req)
	engineer := jobsPatterns.UpdateJobEngineer(builder)

	if err := engineer.UpdateJob(); err != nil {
		return nil, err
	}

	jobId := strconv.Itoa(req.Id)
	
	job, err := r.FindOneJob(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobsRepository) DeleteJob(jobId string) error {
	query := `DELETE FROM "careers" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, jobId); err != nil {
		return fmt.Errorf("delete job failed: %v", err)
	}
	return nil
}