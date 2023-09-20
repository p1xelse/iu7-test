package dto

import (
	"timetracker/models"
)

type ReqCreateUpdateTag struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name" validate:"required"`
	About string `json:"about"`
	Color string `json:"color"`
}

func (req *ReqCreateUpdateTag) ToModelTag() *models.Tag {
	return &models.Tag{
		ID:    req.ID,
		Name:  req.Name,
		About: req.About,
		Color: req.Color,
	}
}

type RespTag struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
	About  string `json:"about"`
	Color  string `json:"color"`
}

func GetResponseFromModelTag(tag *models.Tag) *RespTag {
	return &RespTag{
		ID:     tag.ID,
		UserID: tag.UserID,
		Name:   tag.Name,
		About:  tag.About,
		Color:  tag.Color,
	}
}

func GetResponseFromModelTags(tags []*models.Tag) []*RespTag {
	result := make([]*RespTag, 0, 10)
	for _, tag := range tags {
		result = append(result, GetResponseFromModelTag(tag))
	}

	return result
}
