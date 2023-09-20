package repository

import (
	"timetracker/models"
)

type RepositoryI interface {
	CreateTag(e *models.Tag) error
	UpdateTag(e *models.Tag) error
	GetTag(id uint64) (*models.Tag, error)
	DeleteTag(id uint64) error
	GetUserTags(userID uint64) ([]*models.Tag, error)
	GetEntryTags(entryID uint64) ([]*models.Tag, error)
	CreateEntryTags(entryID uint64, tagList []models.Tag) error
	UpdateEntryTags(entryID uint64, tagList []models.Tag) error
	DeleteEntryTags(entryID uint64) error
}
