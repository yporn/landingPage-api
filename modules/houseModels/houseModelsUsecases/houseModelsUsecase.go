package houseModelsUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/houseModels"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsRepositories"
)

type IHouseModelsUsecase interface {
	FindOneHouseModel(houseId string) (*houseModels.HouseModel, error)
	FindAllHouseModels() ([]houseModels.HouseModelName, error) 
	FindHouseModel(projectId string, req *houseModels.HouseModelFilter) *entities.PaginateRes
	AddHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error)
	UpdateHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error)
	DeleteHouseModel(houseId string) error
}

type houseModelsUsecase struct {
	houseModelsRepository houseModelsRepositories.IHouseModelsRepository
}

func HouseModelsUsecases(houseModelsRepository houseModelsRepositories.IHouseModelsRepository) IHouseModelsUsecase {
	return &houseModelsUsecase{
		houseModelsRepository: houseModelsRepository,
	}
}

func (u *houseModelsUsecase) FindOneHouseModel(houseId string) (*houseModels.HouseModel, error) {
	houseModel, err := u.houseModelsRepository.FindOneHouseModel(houseId)
	if err != nil {
		return nil, err
	}
	return houseModel, nil
}

func (u *houseModelsUsecase) FindAllHouseModels() ([]houseModels.HouseModelName, error) {
	houseModel, err := u.houseModelsRepository.FindAllHouseModels()
	if err != nil {
		return nil, err
	}
	return houseModel, nil
}

func (u *houseModelsUsecase) FindHouseModel(projectId string, req *houseModels.HouseModelFilter) *entities.PaginateRes {
	houseModels, count := u.houseModelsRepository.FindHouseModel(projectId ,req)
	return &entities.PaginateRes{
		Data:      houseModels,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *houseModelsUsecase) AddHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error) {
	houseModel, err := u.houseModelsRepository.InsertHouseModel(req)
	if err != nil {
		return nil, err
	}
	return houseModel, nil
}

func (u *houseModelsUsecase) UpdateHouseModel(req *houseModels.HouseModel) (*houseModels.HouseModel, error) {
	project, err := u.houseModelsRepository.UpdateHouseModel(req)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (u *houseModelsUsecase) DeleteHouseModel(houseId string) error {
	if err := u.houseModelsRepository.DeleteHouseModel(houseId); err != nil {
		return err
	}
	return nil
}