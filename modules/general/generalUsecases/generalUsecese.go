package generalUsecases

import (
	"github.com/yporn/sirarom-backend/modules/general"
	"github.com/yporn/sirarom-backend/modules/general/generalRepositories"
)

type IGeneralUsecase interface {
	FindOneGeneral(generalId string) (*general.General, error)
	UpdateGeneral(req *general.General) (*general.General, error)
}

type generalUsecase struct {
	generalRepository generalRepositories.IGeneralRepository
}

func GeneralUsecase(generalRepository generalRepositories.IGeneralRepository) IGeneralUsecase {
	return &generalUsecase{
		generalRepository: generalRepository,
	}
}

func (u *generalUsecase) FindOneGeneral(generalId string) (*general.General, error) {
	general, err := u.generalRepository.FindOneGeneral(generalId)
	if err != nil {
		return nil, err
	}
	return general, nil
}

func (u *generalUsecase) UpdateGeneral(req *general.General) (*general.General, error) {
	general, err := u.generalRepository.UpdateGeneral(req)
	if err != nil {
		return nil, err
	}
	return general, nil
}