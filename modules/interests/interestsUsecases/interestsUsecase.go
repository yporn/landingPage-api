package interestsUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/interests"
	"github.com/yporn/sirarom-backend/modules/interests/interestsRepositories"
)

type IInterestsUsecase interface {
	FindOneInterest(interestId string) (*interests.Interest, error)
	FindInterest(req *interests.InterestFilter) *entities.PaginateRes
	AddInterest(req *interests.Interest) (*interests.Interest, error)
	UpdateInterest(req *interests.Interest) (*interests.Interest, error)
	DeleteInterest(interestId string) error
}

type interestsUsecase struct {
	interestsRepository interestsRepositories.IInterestRepository
}

func InterestsUsecase(interestsRepository interestsRepositories.IInterestRepository) IInterestsUsecase {
	return &interestsUsecase{
		interestsRepository: interestsRepository,
	}
}

func (u *interestsUsecase) FindOneInterest(interestId string) (*interests.Interest, error) {
	interest, err := u.interestsRepository.FindOneInterest(interestId)
	if err != nil {
		return nil, err
	}
	return interest, nil
}

func (u *interestsUsecase) FindInterest(req *interests.InterestFilter) *entities.PaginateRes {
	activities, count := u.interestsRepository.FindInterest(req)

	return &entities.PaginateRes{
		Data:      activities,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *interestsUsecase) AddInterest(req *interests.Interest) (*interests.Interest, error) {
	interest, err := u.interestsRepository.InsertInterest(req)
	if err != nil {
		return nil, err
	}
	return interest, nil
}

func (u *interestsUsecase) UpdateInterest(req *interests.Interest) (*interests.Interest, error) {
	activity, err := u.interestsRepository.UpdateInterest(req)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (u *interestsUsecase) DeleteInterest(interestId string) error {
	if err := u.interestsRepository.DeleteInterest(interestId); err != nil {
		return err
	}
	return nil
}
