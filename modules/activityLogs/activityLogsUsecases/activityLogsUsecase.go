package activityLogsUsecases

import (
	"github.com/yporn/sirarom-backend/modules/activityLogs"
	"github.com/yporn/sirarom-backend/modules/activityLogs/activityLogsRepositories"
)

type IActivityLogsUsecase interface {
	FindAllActivityLogs() ([]activityLogs.ActivityLog, error)
}

type activityLogsUsecase struct {
	activityLogsRepository activityLogsRepositories.IActivityLogsRepository
}

func ActivityLogsUsecases(activityLogsRepository activityLogsRepositories.IActivityLogsRepository) IActivityLogsUsecase {
	return &activityLogsUsecase{
		activityLogsRepository: activityLogsRepository,
	}
}

func (u *activityLogsUsecase) FindAllActivityLogs() ([]activityLogs.ActivityLog, error) {
	activityLog, err := u.activityLogsRepository.FindActivityLog()
	if err != nil {
		return nil, err
	}
	return activityLog, nil
}
