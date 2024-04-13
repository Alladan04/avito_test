package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	db redis.Client
}

const DefaultExpTime = time.Minute * 10

func NewCacheRepo(db redis.Client) *CacheRepo {
	return &CacheRepo{
		db: db,
	}
}

func (repo *CacheRepo) GetBanner(ctx context.Context, featureId int64, tagId int64) (models.BannerContent, error) {
	var result models.BannerContent
	key := fmt.Sprintf("%d:%d", featureId, tagId)
	banner, err := repo.db.Get(ctx, key).Result()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(banner), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repo *CacheRepo) AddBanner(ctx context.Context, featureId int64, tagId int64, banner models.BannerContent) error {
	key := fmt.Sprintf("%d:%d", featureId, tagId)
	err := repo.db.Set(ctx, key, banner, DefaultExpTime).Err()
	if err != nil {
		return err
	}
	return nil
}
