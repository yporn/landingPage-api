package logosUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/logos"
	"github.com/yporn/sirarom-backend/modules/logos/logosRepositories"
)

type ILogosUsecase interface {
	FindOneLogo(logoId string) (*logos.Logo, error)
	FindLogo(req *logos.LogoFilter) *entities.PaginateRes
	AddLogo(req *logos.Logo) (*logos.Logo, error) 
	UpdateLogo(req *logos.Logo) (*logos.Logo, error)
	DeleteLogo(logoId string) error
}

type logosUsecase struct {
	logosRepository logosRepositories.ILogosRepository
}

func LogosUsecase(logosRepository logosRepositories.ILogosRepository) ILogosUsecase {
	return &logosUsecase{
		logosRepository: logosRepository,
	}
}

func (u *logosUsecase) FindOneLogo(logoId string) (*logos.Logo, error) {
	logo, err := u.logosRepository.FindOneLogo(logoId)
	if err != nil {
		return nil, err
	}
	return logo, nil
}

func (u *logosUsecase) FindLogo(req *logos.LogoFilter) *entities.PaginateRes {
	logos, count := u.logosRepository.FindLogo(req)

	return &entities.PaginateRes{
		Data:      logos,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *logosUsecase) AddLogo(req *logos.Logo) (*logos.Logo, error) {
	logo, err := u.logosRepository.InsertLogo(req)
	if err != nil {
		return nil, err
	}
	return logo, nil
}


func (u *logosUsecase) UpdateLogo(req *logos.Logo) (*logos.Logo, error) {
	logo, err := u.logosRepository.UpdateLogo(req)
	if err != nil {
		return nil, err
	}
	return logo, nil
}

func (u *logosUsecase) DeleteLogo(logoId string) error {
	if err := u.logosRepository.DeleteLogo(logoId); err != nil {
		return err
	}
	return nil
}
