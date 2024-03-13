package jobsUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/jobs"
	"github.com/yporn/sirarom-backend/modules/jobs/jobsRepositories"
)

type IJobsUsecase interface {
	FindOneJob(jobId string) (*jobs.Job, error)
	FindJob(req *jobs.JobFilter) *entities.PaginateRes
	AddJob(req *jobs.Job) (*jobs.Job, error)
	UpdateJob(req *jobs.Job) (*jobs.Job, error)
	DeleteJob(jobId string) error
}

type jobsUsecase struct {
	jobsRepository jobsRepositories.IJobRepository
}

func JobsUsecase(jobsRepository jobsRepositories.IJobRepository) IJobsUsecase {
	return &jobsUsecase{
		jobsRepository: jobsRepository,
	}
}

func (u *jobsUsecase) FindOneJob(jobId string) (*jobs.Job, error) {
	job, err := u.jobsRepository.FindOneJob(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (u *jobsUsecase) FindJob(req *jobs.JobFilter) *entities.PaginateRes {
	jobs, count := u.jobsRepository.FindJob(req)

	return &entities.PaginateRes{
		Data: jobs,
		Page: req.Page,
		Limit: req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *jobsUsecase) AddJob(req *jobs.Job) (*jobs.Job, error) {
	job, err := u.jobsRepository.InsertJob(req)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (u *jobsUsecase) UpdateJob(req *jobs.Job) (*jobs.Job, error) {
	job, err := u.jobsRepository.UpdateJob(req)
	if err != nil {
		return nil, err
	}
	return job, nil
}


func (u *jobsUsecase) DeleteJob(jobId string) error {
	if err := u.jobsRepository.DeleteJob(jobId); err != nil {
		return err
	}
	return nil
}