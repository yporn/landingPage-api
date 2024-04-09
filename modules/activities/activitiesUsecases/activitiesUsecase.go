package activitiesUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/activities"
	"github.com/yporn/sirarom-backend/modules/activities/activitiesRepositories"
	"github.com/yporn/sirarom-backend/modules/entities"
)

type IActivitiesUsecase interface {
	FindOneActivity(activityId string) (*activities.Activity, error) 
	FindActivity(req *activities.ActivityFilter) *entities.PaginateRes
	AddActivity(req *activities.Activity) (*activities.Activity, error)
	UpdateActivity(req *activities.Activity) (*activities.Activity, error)
	DeleteActivity(activityId string) error
}

type activitiesUsecase struct {
	activitiesRepository activitiesRepositories.IActivitiesRepository
}

func ActivitiesUsecase(activitiesRepository activitiesRepositories.IActivitiesRepository) IActivitiesUsecase {
	return &activitiesUsecase{
		activitiesRepository: activitiesRepository,
	}
}

func (u *activitiesUsecase) FindOneActivity(activityId string) (*activities.Activity, error) {
	activity, err := u.activitiesRepository.FindOneActivity(activityId)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (u *activitiesUsecase) FindActivity(req *activities.ActivityFilter) *entities.PaginateRes {
	activities, count := u.activitiesRepository.FindActivity(req)

	return &entities.PaginateRes{
		Data:      activities,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *activitiesUsecase) AddActivity(req *activities.Activity) (*activities.Activity, error) {
	activity, err := u.activitiesRepository.InsertActivity(req)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (u *activitiesUsecase) UpdateActivity(req *activities.Activity) (*activities.Activity, error) {
	activity, err := u.activitiesRepository.UpdateActivity(req)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (u *activitiesUsecase) DeleteActivity(activityId string) error {
	if err := u.activitiesRepository.DeleteActivity(activityId); err != nil {
		return err
	}
	return nil
}
