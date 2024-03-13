package bannersUsecases

import (
	"math"

	"github.com/yporn/sirarom-backend/modules/banners"
	"github.com/yporn/sirarom-backend/modules/banners/bannersRepositories"
	"github.com/yporn/sirarom-backend/modules/entities"
)

type IBannersUsecase interface {
	FindOneBanner(bannerId string) (*banners.Banner, error) 
	FindBanner(req *banners.BannerFilter) *entities.PaginateRes
	AddBanner(req *banners.Banner) (*banners.Banner, error)
	UpdateBanner(req *banners.Banner) (*banners.Banner, error)
	DeleteBanner(bannerId string) error
}

type bannersUsecase struct {
	bannersRepository bannersRepositories.IBannersRepository
}

func BannersUsecase(bannersRepository bannersRepositories.IBannersRepository) IBannersUsecase {
	return &bannersUsecase{
		bannersRepository: bannersRepository,
	}
}

func (u *bannersUsecase) FindOneBanner(bannerId string) (*banners.Banner, error) {
	banner, err := u.bannersRepository.FindOneBanner(bannerId)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (u *bannersUsecase) FindBanner(req *banners.BannerFilter) *entities.PaginateRes {
	banners, count := u.bannersRepository.FindBanner(req)

	return &entities.PaginateRes{
		Data:      banners,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *bannersUsecase) AddBanner(req *banners.Banner) (*banners.Banner, error) {
	banner, err := u.bannersRepository.InsertBanner(req)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (u *bannersUsecase) UpdateBanner(req *banners.Banner) (*banners.Banner, error) {
	banner, err := u.bannersRepository.UpdateBanner(req)
	if err != nil {
		return nil, err
	}
	return banner, nil
}


func (u *bannersUsecase) DeleteBanner(bannerId string) error {
	if err := u.bannersRepository.DeleteBanner(bannerId); err != nil {
		return err
	}
	return nil
}
