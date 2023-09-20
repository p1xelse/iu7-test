package repository

import (
	"time"
	"timetracker/models"
)

type RepositoryI interface {
	CreateEntry(e *models.Entry) error
	UpdateEntry(e *models.Entry) error
	GetEntry(id uint64) (*models.Entry, error)
	DeleteEntry(id uint64) error
	GetUserEntries(userID uint64) ([]*models.Entry, error)
	GetUserEntriesForDay(userID uint64, date time.Time) ([]*models.Entry, error)
}
