package promotionsUsecases

import (
	"github.com/yporn/sirarom-backend/modules/promotions"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsRepositories"
)

type IPromotionsUsecase interface {
	FindOnePromotion(promotionId string) (*promotions.Promotion, error)
}

type promotionsUsecase struct {
	promotionsRepository promotionsRepositories.IPromotionsRepository
}

func PromotionsUsecase (promotionsRepository promotionsRepositories.IPromotionsRepository) IPromotionsUsecase {
	return &promotionsUsecase{
		promotionsRepository: promotionsRepository,
	}
}

func (u *promotionsUsecase) FindOnePromotion(promotionId string) (*promotions.Promotion, error) {
	promotion, err := u.promotionsRepository.FindOnePromotion(promotionId)
	if err != nil {
		return nil, err
	}
	return promotion, nil
}
