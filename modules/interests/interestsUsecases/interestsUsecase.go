package interestsUsecases

import (
	"github.com/yporn/sirarom-backend/modules/interests"
	"github.com/yporn/sirarom-backend/modules/interests/interestsRepositories"
)

type IInterestsUsecase interface {
	FindOneInterest(interestId string) (*interests.Interest, error)
	AddInterest(req *interests.Interest) (*interests.Interest, error)
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

func (u *interestsUsecase) AddInterest(req *interests.Interest) (*interests.Interest, error) {
	interest, err := u.interestsRepository.InsertInterest(req)
	if err != nil {
		return nil, err
	}
	return interest, nil
}

