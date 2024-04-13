package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/Alladan04/avito_test/internal/pkg/banner"
)

const (
	pageElementsCount = 10
)

type BannerUsecase struct {
	repo  banner.BannerRepo
	cache banner.CacheRepo
}

func NewBannerUsecase(repo banner.BannerRepo, cache banner.CacheRepo) *BannerUsecase {
	return &BannerUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *BannerUsecase) AddItem(ctx context.Context, data models.BannerForm) (models.Banner, error) {

	item := models.Banner{
		Content: models.BannerContent{
			Title: data.Content.Title,
			Data:  data.Content.Data,
			Url:   data.Content.Url,
		},
		FeatureId:  data.FeatureId,
		TagIds:     data.TagIds,
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
	}
	err := uc.repo.AddItem(ctx, item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (uc *BannerUsecase) GetAll(ctx context.Context, count int64, offset int64, featureId int64, tagId int64) ([]models.Banner, error) {
	//get data from db
	var data []models.Banner
	var err error
	if count == 0 {
		count = pageElementsCount
	}

	data, err = uc.repo.GetAllFiltered(ctx, count, offset, featureId, tagId)

	if err != nil {
		return nil, err
	}
	return data, nil

}
func (uc *BannerUsecase) GetOne(ctx context.Context, featureId int64, tagId int64, showLastRevision bool) (models.BannerContent, error) {
	result, err := uc.cache.GetBanner(ctx, featureId, tagId)
	if err == nil {
		return result, nil
	}
	result, err = uc.repo.GetOne(ctx, featureId, tagId)
	if err != nil {
		return result, err
	}
	err = uc.cache.AddBanner(ctx, featureId, tagId, result)
	if err != nil {
		fmt.Printf("error while trying to cache:%s", err.Error())
	}
	return result, nil
}

func (uc *BannerUsecase) UpdateBanner(ctx context.Context, payload models.BannerUpdateForm, id int64) error {
	banner, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return errors.New("not found")
	}
	if payload.Content != nil {
		if payload.Content.Title != nil {
			banner.Content.Title = *payload.Content.Title
		}
		if payload.Content.Data != nil {
			banner.Content.Url = *payload.Content.Url
		}
		if payload.Content.Url != nil {
			banner.Content.Url = *payload.Content.Url
		}
	}
	if payload.FeatureId != nil {
		banner.FeatureId = *payload.FeatureId
	}
	if payload.TagIds != nil {
		banner.TagIds = payload.TagIds
	}
	if payload.IsActive != nil {
		banner.IsActive = *payload.IsActive
	}
	err = uc.repo.UpdateBanner(ctx, banner, id)
	if err != nil {
		return errors.New("internal")
	}
	return nil

}

func (uc *BannerUsecase) DeleteBanner(ctx context.Context, id int64) error {
	err := uc.repo.DeleteBanner(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
