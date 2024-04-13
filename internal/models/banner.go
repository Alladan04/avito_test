package models

import (
	"encoding/json"
	"time"
)

type BannerContent struct {
	Title string `json:"title"`
	Data  string `json:"text"`
	Url   string `json:"url"`
}

func (i BannerContent) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(i)
	return bytes, err
}

type Banner struct {
	Id         int64         `json:"id"`
	Content    BannerContent `json:"content"`
	TagIds     []int64       `json:"tag_ids"`
	CreateTime time.Time     `json:"create_time"`
	UpdateTime time.Time     `json:"update_time"`
	FeatureId  int64         `json:"feature_id"`
	IsActive   bool          `json:"is_active"`
}

type UpdateContentForm struct {
	Title *string `json:"title,omitempty"`
	Data  *string `json:"text,omitempty"`
	Url   *string `json:"url,omitempty"`
}
type BannerUpdateForm struct {
	Content   *UpdateContentForm `json:"content,omitempty"`
	FeatureId *int64             `json:"feature_id,omitempty"`
	TagIds    []int64            `json:"tag_ids,omitempty"`
	IsActive  *bool              `json:"is_active,omitempty"`
}
type BannerForm struct {
	Content   BannerContent `json:"content"`
	FeatureId int64         `json:"feature_id"`
	TagIds    []int64       `json:"tag_ids"`
	IsActive  bool          `json:"is_active"`
}
