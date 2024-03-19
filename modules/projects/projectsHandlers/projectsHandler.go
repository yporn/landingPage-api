package projectsHandlers

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/projects"
	"github.com/yporn/sirarom-backend/modules/projects/projectsUsecases"
)

type projectsHandlersErrCode string

const (
	findOneProjectErr projectsHandlersErrCode = "projects-001"
	findProjectErr    projectsHandlersErrCode = "projects-002"
	insertProjectErr  projectsHandlersErrCode = "projects-003"
	deleteProjectErr  projectsHandlersErrCode = "projects-004"
	updateProjectErr  projectsHandlersErrCode = "projects-005"
)

type IProjectsHandler interface {
	FindOneProject(c *fiber.Ctx) error
	FindProject(c *fiber.Ctx) error
	AddProject(c *fiber.Ctx) error
	UpdateProject(c *fiber.Ctx) error
	DeleteProject(c *fiber.Ctx) error
	FindProjectHouseModel(c *fiber.Ctx) error
}

type projectsHandler struct {
	cfg              config.IConfig
	projectsUsecases projectsUsecases.IProjectsUsecase
	filesUsecase     filesUsecases.IFilesUsecase
}

func ProjectsHandler(cfg config.IConfig, projectsUsecase projectsUsecases.IProjectsUsecase, filesUsecase filesUsecases.IFilesUsecase) IProjectsHandler {
	return &projectsHandler{
		cfg:              cfg,
		projectsUsecases: projectsUsecase,
		filesUsecase:     filesUsecase,
	}
}

func (h *projectsHandler) FindOneProject(c *fiber.Ctx) error {
	projectId := strings.Trim(c.Params("project_id"), " ")

	project, err := h.projectsUsecases.FindOneProject(projectId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProjectErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, project).Res()
}

func (h *projectsHandler) FindProjectHouseModel(c *fiber.Ctx) error {
	projectId := strings.Trim(c.Params("project_id"), " ")

	project, err := h.projectsUsecases.FindProjectHouseModel(projectId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProjectErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, project).Res()
}


func (h *projectsHandler) FindProject(c *fiber.Ctx) error {
	req := &projects.ProjectFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProjectErr),
			err.Error(),
		).Res()
	}

	// Paginate
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Sort
	orderByMap := map[string]string{
		"id":         `"p"."id"`,
		"created_at": `"p"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		h.projectsUsecases.FindProject(req),
	).Res()
}

func (h *projectsHandler) AddProject(c *fiber.Ctx) error {
	req := &projects.Project{
		HouseTypeItem:   make([]*projects.ProjectHouseTypeItem, 0),
		DescAreaItem:    make([]*projects.ProjectDescAreaItem, 0),
		ComfortableItem: make([]*projects.ProjectComfortableItem, 0),
		Images:          make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProjectErr),
			err.Error(),
		).Res()
	}

	if len(req.HouseTypeItem) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProjectErr),
			"house type item id is invalid",
		).Res()
	}

	if len(req.DescAreaItem) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProjectErr),
			"area item id is invalid",
		).Res()
	}

	// facilities
	if len(req.ComfortableItem) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProjectErr),
			"house type item id is invalid",
		).Res()
	}

	project, err := h.projectsUsecases.AddProject(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProjectErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, project).Res()
}

func (h *projectsHandler) UpdateProject(c *fiber.Ctx) error {
	projectIdStr := strings.Trim(c.Params("project_id"), " ")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProjectErr),
			err.Error(),
		).Res()
	}

	req := &projects.Project{
		HouseTypeItem:   make([]*projects.ProjectHouseTypeItem, 0),
		DescAreaItem:    make([]*projects.ProjectDescAreaItem, 0),
		ComfortableItem: make([]*projects.ProjectComfortableItem, 0),
		Images:          make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProjectErr),
			err.Error(),
		).Res()
	}
	req.Id = projectId

	project, err := h.projectsUsecases.UpdateProject(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProjectErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, project).Res()
}

func (h *projectsHandler) DeleteProject(c *fiber.Ctx) error {
	projectId := strings.Trim(c.Params("project_id"), " ")

	project, err := h.projectsUsecases.FindOneProject(projectId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProjectErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range project.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("projects/%s", path.Base(p.Url)),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProjectErr),
			err.Error(),
		).Res()
	}

	if err := h.projectsUsecases.DeleteProject(projectId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProjectErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
