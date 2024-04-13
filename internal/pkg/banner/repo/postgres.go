package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
)

const (
	addItem              = "INSERT INTO banner (title, feature_id, banner_data, url,  create_time, update_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	getAll               = "SELECT id, title, feature_id, banner_data, url,  create_time, update_time, is_active FROM banner LIMIT $1 OFFSET $2; "
	getTagsForBanner     = "SELECT tag_id from banner_tag WHERE banner_id=$1;"
	getBannersWithTagIds = `SELECT b.id, b.title, b.feature_id, b.banner_data, b.url,  b.create_time, b.update_time, b.is_active , ARRAY_AGG(bt.tag_id)
						FROM banner b
						JOIN banner_tag bt ON b.id = bt.banner_id
						Where b.feature_id=%s AND bt.tag_id=%s
						GROUP BY b.id
						LIMIT $1 OFFSET $2; ;
					`
	addBT          = "INSERT INTO banner_tag (banner_id, tag_id, feature_id) VALUES ($1, $2, $3);"
	getAllFiltered = `SELECT id, title, feature_id, banner_data, url,  create_time, update_time, is_active 
					FROM banner 
					WHERE feature_id=%s 
					LIMIT $1 OFFSET $2; `
	getContent = `SELECT b.title, b.banner_data, b.url FROM banner b 
					JOIN banner_tag bt ON b.id = bt.banner_id 
					WHERE bt.tag_id = $1 AND bt.feature_id = $2 AND b.is_active='true'; `
	getById = `SELECT  b.title, b.feature_id, b.banner_data, b.url,  b.is_active,
					coalesce(array_agg(bt.tag_id) filter (where bt.tag_id is not null), '{}')
					FROM banner b
					LEFT JOIN banner_tag bt ON b.id = bt.banner_id
					WHERE b.id=$1
					GROUP BY b.id;`
	updateBanner = `UPDATE banner SET title=$1, feature_id=$2, banner_data=$3, url=$4, is_active=$5, update_time=$6
							WHERE id=$7;  `
	deleteTagsByBannerId = `DELETE FROM banner_tag WHERE banner_id=$1;`
	deleteBannerById     = `DELETE FROM banner WHERE id = $1;`
)

type BannerRepo struct {
	db   pgxtype.Querier
	conn pgx.Conn
}

func NewBannerRepo(db pgxtype.Querier, conn pgx.Conn) *BannerRepo {
	return &BannerRepo{
		db:   db,
		conn: conn,
	}
}

func (repo *BannerRepo) AddItem(ctx context.Context, item models.Banner) error {
	tx, err := repo.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			fmt.Printf("ERROR: %v", err)
		}
	}()

	row := tx.QueryRow(ctx, addItem, item.Content.Title, item.FeatureId, item.Content.Data, item.Content.Url, item.CreateTime, item.UpdateTime)
	err = row.Scan(&item.Id)
	if err != nil {
		return err
	}
	for _, tag := range item.TagIds {
		_, err = tx.Exec(ctx, addBT, item.Id, tag, item.FeatureId)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil

}

func (repo *BannerRepo) GetAllFiltered(ctx context.Context, count int64, offset int64, featureId int64, tagId int64) ([]models.Banner, error) {
	result := make([]models.Banner, 0, count)
	var query string
	query = fmt.Sprintf(getBannersWithTagIds, strconv.FormatInt(featureId, 10), strconv.FormatInt(tagId, 10))
	if featureId == 0 && tagId == 0 {
		query = fmt.Sprintf(getBannersWithTagIds, "b.feature_id", "bt.tag_id")
	} else if featureId == 0 {
		query = fmt.Sprintf(getBannersWithTagIds, "b.feature_id", strconv.FormatInt(tagId, 10))

	} else if tagId == 0 {
		query = fmt.Sprintf(getBannersWithTagIds, strconv.FormatInt(featureId, 10), "bt.tag_id")

	}
	rows, err := repo.db.Query(ctx, query, count, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item models.Banner
		if err := rows.Scan(&item.Id, &item.Content.Title, &item.FeatureId, &item.Content.Data, &item.Content.Url, &item.CreateTime, &item.UpdateTime, &item.IsActive, &item.TagIds); err != nil {
			return nil, fmt.Errorf("error occured while scanning items:%w", err)
		}

		result = append(result, item)
	}

	return result, nil
}

func (repo *BannerRepo) GetOne(ctx context.Context, featureId int64, tagId int64) (models.BannerContent, error) {
	var result models.BannerContent
	err := repo.db.QueryRow(ctx, getContent, tagId, featureId).Scan(
		&result.Title,
		&result.Data,
		&result.Url,
	)
	if err != nil {
		return models.BannerContent{}, err
	}
	return result, nil
}

func (repo *BannerRepo) GetById(ctx context.Context, id int64) (models.BannerForm, error) {
	var result models.BannerForm

	err := repo.db.QueryRow(ctx, getById, id).Scan(
		&result.Content.Title,
		&result.FeatureId,
		&result.Content.Data,
		&result.Content.Url,
		&result.IsActive,
		&result.TagIds,
	)
	if err != nil {
		return models.BannerForm{}, err
	}

	return result, nil
}

func (repo *BannerRepo) UpdateBanner(ctx context.Context, banner models.BannerForm, id int64) error {
	tx, err := repo.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			fmt.Printf("ERROR: %v", err)
		}
	}()
	//обновляем баннер
	_, err = tx.Exec(ctx, updateBanner, banner.Content.Title, banner.FeatureId, banner.Content.Data, banner.Content.Url, banner.IsActive, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	//чистим связанные с баннером теги
	_, err = tx.Exec(ctx, deleteTagsByBannerId, id)
	if err != nil {
		return err
	}
	//записываем новый список тегов
	for _, tag := range banner.TagIds {
		_, err = tx.Exec(ctx, addBT, id, tag, banner.FeatureId)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (repo *BannerRepo) DeleteBanner(ctx context.Context, id int64) error {
	tx, err := repo.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			fmt.Printf("ERROR: %v", err)
		}
	}()
	//чистим связанные с баннером теги
	_, err = tx.Exec(ctx, deleteTagsByBannerId, id)
	if err != nil {
		return err
	}
	//удаляем баннер
	_, err = tx.Exec(ctx, deleteBannerById, id)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil

}
