package models

import (
	"time"
)

type Goal struct {
	ID          uint64
	UserID      *uint64
	Name        string
	ProjectID   *uint64
	HoursCount  float64
	Description string
	TimeStart   time.Time
	TimeEnd     time.Time
}
