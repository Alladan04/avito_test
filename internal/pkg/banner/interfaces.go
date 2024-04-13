package banner

import (
	"context"

	"github.com/Alladan04/avito_test/internal/models"
)

type BannerRepo interface {
	AddItem(ctx context.Context, item models.Banner) error
	GetById(ctx context.Context, id int64) (models.BannerForm, error)
	UpdateBanner(ctx context.Context, banner models.BannerForm, id int64) error
	//GetAll(ctx context.Context, count int64, offset int64) ([]models.Banner, error)
	GetOne(ctx context.Context, featureId int64, tagId int64) (models.BannerContent, error)
	GetAllFiltered(ctx context.Context, count int64, offset int64, featureId int64, tagId int64) ([]models.Banner, error)
	DeleteBanner(ctx context.Context, id int64) error
}

type BannerUsecase interface {
	AddItem(ctx context.Context, data models.BannerForm) (models.Banner, error)
	GetOne(ctx context.Context, featureId int64, tagId int64, showLastRevision bool) (models.BannerContent, error)
	GetAll(ctx context.Context, count int64, offset int64, featureId int64, tagId int64) ([]models.Banner, error)
	UpdateBanner(ctx context.Context, payload models.BannerUpdateForm, id int64) error
	DeleteBanner(ctx context.Context, id int64) error
}

type CacheRepo interface {
	GetBanner(ctx context.Context, featureId int64, tagId int64) (models.BannerContent, error)
	AddBanner(ctx context.Context, featureId int64, tagId int64, banner models.BannerContent) error
}
