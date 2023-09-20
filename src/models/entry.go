package models

import (
	"time"
	"timetracker/pkg"
)

type Entry struct {
	ID          uint64    `json:"id"`
	UserID      *uint64   `json:"user_id"`
	ProjectID   *uint64   `json:"project_id"`
	Description string    `json:"description"`
	TagList     []Tag     `json:"tag_list"`
	TimeStart   time.Time `json:"time_start"`
	TimeEnd     time.Time `json:"time_end"`
	Duration    string    `json:"-"`
}

func (e *Entry) CalcDuration() {
	e.Duration = pkg.GetPrettyDuration(e.TimeStart, e.TimeEnd)
}
