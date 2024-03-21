package promotionsUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/promotions"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsRepositories"
)

type IPromotionsUsecase interface {
	FindOnePromotion(promotionId string) (*promotions.Promotion, error)
	FindPromotion(req *promotions.PromotionFilter) *entities.PaginateRes 
	AddPromotion(req *promotions.Promotion) (*promotions.Promotion, error)
	UpdatePromotion(req *promotions.Promotion) (*promotions.Promotion, error)
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

func (u *promotionsUsecase) FindPromotion(req *promotions.PromotionFilter) *entities.PaginateRes {
	promotions, count := u.promotionsRepository.FindPromotion(req)

	return &entities.PaginateRes{
		Data:      promotions,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *promotionsUsecase) AddPromotion(req *promotions.Promotion) (*promotions.Promotion, error) {
	promotion, err := u.promotionsRepository.InsertPromotion(req)
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (u *promotionsUsecase) UpdatePromotion(req *promotions.Promotion) (*promotions.Promotion, error) {
	promotion, err := u.promotionsRepository.UpdatePromotion(req)
	if err != nil {
		return nil, err
	}
	return promotion, nil
}
