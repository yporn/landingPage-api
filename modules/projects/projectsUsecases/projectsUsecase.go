package projectsUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/projects"
	"github.com/yporn/sirarom-backend/modules/projects/projectsRepositories"
)

type IProjectsUsecase interface {
	FindOneProject(projectId string) (*projects.Project, error)
	FindProject(req *projects.ProjectFilter) *entities.PaginateRes
	AddProject(req *projects.Project) (*projects.Project, error)
	UpdateProject(req *projects.Project) (*projects.Project, error)
	DeleteProject(projectId string) error
}

type projectsUsecase struct {
	projectsRepository projectsRepositories.IProjectRepository
}

func ProjectsUsecase(projectsRepository projectsRepositories.IProjectRepository) IProjectsUsecase {
	return &projectsUsecase{
		projectsRepository: projectsRepository,
	}
}

func (u *projectsUsecase) FindOneProject(projectId string) (*projects.Project, error) {
	project, err := u.projectsRepository.FindOneProject(projectId)
	if err != nil {
		return nil, err
	}
	return project, nil
}


func (u *projectsUsecase) AddProject(req *projects.Project) (*projects.Project, error) {
	project, err := u.projectsRepository.InsertProject(req)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (u *projectsUsecase) FindProject(req *projects.ProjectFilter) *entities.PaginateRes {
	projects, count := u.projectsRepository.FindProject(req)
	return &entities.PaginateRes{
		Data:      projects,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *projectsUsecase) UpdateProject(req *projects.Project) (*projects.Project, error) {
	project, err := u.projectsRepository.UpdateProject(req)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (u *projectsUsecase) DeleteProject(projectId string) error {
	if err := u.projectsRepository.DeleteProject(projectId); err != nil {
		return err
	}
	return nil
}