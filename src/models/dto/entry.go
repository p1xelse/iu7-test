package dto

import (
	"time"
	"timetracker/models"
)

type ReqCreateUpdateEntry struct {
	ID          uint64    `json:"id"`
	ProjectID   *uint64   `json:"project_id"`
	Description string    `json:"description"`
	TagList     []uint64  `json:"tag_list"`
	TimeStart   time.Time `json:"time_start" validate:"required"`
	TimeEnd     time.Time `json:"time_end" validate:"required"`
}

func (req *ReqCreateUpdateEntry) ToModelEntry() *models.Entry {
	tagListModel := make([]models.Tag, len(req.TagList))

	for idx := range req.TagList {
		tagListModel[idx] = models.Tag{ID: req.TagList[idx]}
	}

	return &models.Entry{
		ID:          req.ID,
		ProjectID:   req.ProjectID,
		Description: req.Description,
		TagList:     tagListModel,
		TimeEnd:     req.TimeEnd,
		TimeStart:   req.TimeStart,
	}
}

type RespEntry struct {
	ID          uint64       `json:"id"`
	UserID      *uint64      `json:"user_id"`
	ProjectID   *uint64      `json:"project_id"`
	Description string       `json:"description"`
	TagList     []models.Tag `json:"tag_list"`
	TimeStart   time.Time    `json:"time_start"`
	TimeEnd     time.Time    `json:"time_end"`
	Duration    string       `json:"duration"`
}

func GetResponseFromModelEntry(entry *models.Entry) *RespEntry {
	entry.CalcDuration()

	return &RespEntry{
		ID:          entry.ID,
		UserID:      entry.UserID,
		ProjectID:   entry.ProjectID,
		Description: entry.Description,
		TagList:     entry.TagList,
		TimeEnd:     entry.TimeEnd,
		TimeStart:   entry.TimeStart,
		Duration:    entry.Duration,
	}
}

func GetResponseFromModelEntries(entries []*models.Entry) []*RespEntry {
	result := make([]*RespEntry, 0, 10)
	for _, entry := range entries {
		result = append(result, GetResponseFromModelEntry(entry))
	}

	return result
}
