package jobsHandlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/jobs"
	"github.com/yporn/sirarom-backend/modules/jobs/jobsUsecases"
)

type jobsHandlersErrCode string

const (
	findOneJobErr jobsHandlersErrCode = "jobs-001"
	findJobErr    jobsHandlersErrCode = "jobs-002"
	insertJobErr  jobsHandlersErrCode = "jobs-003"
	updateJobErr  jobsHandlersErrCode = "jobs-004"
	deleteJobErr  jobsHandlersErrCode = "jobs-005"
)

type IJobsHandler interface {
	FindOneJob(c *fiber.Ctx) error
	FindJob(c *fiber.Ctx) error
	AddJob(c *fiber.Ctx) error
	UpdateJob(c *fiber.Ctx) error
	DeleteJob(c *fiber.Ctx) error
}

type jobsHandler struct {
	cfg         config.IConfig
	jobsUsecase jobsUsecases.IJobsUsecase
}

func JobsHandler(cfg config.IConfig, jobsUsecase jobsUsecases.IJobsUsecase) IJobsHandler {
	return &jobsHandler{
		cfg:         cfg,
		jobsUsecase: jobsUsecase,
	}
}

func (h *jobsHandler) FindOneJob(c *fiber.Ctx) error {
	jobId := strings.Trim(c.Params("job_id"), " ")

	job, err := h.jobsUsecase.FindOneJob(jobId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneJobErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, job).Res()
}

func (h *jobsHandler) FindJob(c *fiber.Ctx) error {
	req := &jobs.JobFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findJobErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}
	
	if req.Sort == "" {
		req.Sort = "asc"
	}

	jobs := h.jobsUsecase.FindJob(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, jobs).Res()
}

func (h *jobsHandler) AddJob(c *fiber.Ctx) error {
	req := &jobs.Job{}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertJobErr),
			err.Error(),
		).Res()
	}

	job, err := h.jobsUsecase.AddJob(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertJobErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, job).Res()
}

func (h *jobsHandler) UpdateJob(c *fiber.Ctx) error {
	jobIdStr := strings.Trim(c.Params("job_id"), " ")
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		// Handle the error if jobIdStr cannot be converted to an integer
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateJobErr),
			"Invalid job ID",
		).Res()
	}

	req := &jobs.Job{}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateJobErr),
			err.Error(),
		).Res()
	}

	req.Id = jobId

	job, err := h.jobsUsecase.UpdateJob(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateJobErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, job).Res()
}

func (h *jobsHandler) DeleteJob(c *fiber.Ctx) error {
	jobId := strings.Trim(c.Params("job_id"), " ")

	_, err := h.jobsUsecase.FindOneJob(jobId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteJobErr),
			err.Error(),
		).Res()
	}

	if err := h.jobsUsecase.DeleteJob(jobId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteJobErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "deleted").Res()
}
