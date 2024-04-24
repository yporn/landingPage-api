package seoUsecases

import (
	"github.com/yporn/sirarom-backend/modules/seo"
	"github.com/yporn/sirarom-backend/modules/seo/seoRepositories"
)

type ISeoUsecase interface {
	FindOneSeo(seoId string) (*seo.Seo, error)
	UpdateSeo(req *seo.Seo) (*seo.Seo, error)
}

type seoUsecase struct {
	seoRepository seoRepositories.ISeoRepository
}

func SeoUsecase(seoRepository seoRepositories.ISeoRepository) ISeoUsecase {
	return &seoUsecase{
		seoRepository: seoRepository,
	}
}

func (u *seoUsecase) FindOneSeo(seoId string) (*seo.Seo, error) {
	seo, err := u.seoRepository.FindOneSeo(seoId)
	if err != nil {
		return nil, err
	}
	return seo, nil
}

func (u *seoUsecase) UpdateSeo(req *seo.Seo) (*seo.Seo, error) {
	seo, err := u.seoRepository.UpdateSeo(req)
	if err != nil {
		return nil, err
	}
	return seo, nil
}