package dto

import (
	"time"
	"timetracker/models"
)

type ReqCreateUpdateGoal struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	ProjectID   *uint64   `json:"project_id" validate:"required"`
	HoursCount  float64   `json:"hours_count" validate:"required"`
	Description string    `json:"description"`
	TimeStart   time.Time `json:"time_start" validate:"required"`
	TimeEnd     time.Time `json:"time_end" validate:"required"`
}

func (req *ReqCreateUpdateGoal) ToModelGoal() *models.Goal {
	return &models.Goal{
		ID:          req.ID,
		Name:        req.Name,
		ProjectID:   req.ProjectID,
		HoursCount:  req.HoursCount,
		Description: req.Description,
		TimeEnd:     req.TimeEnd,
		TimeStart:   req.TimeStart,
	}
}

type RespGoal struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	UserID      *uint64   `json:"user_id"`
	ProjectID   *uint64   `json:"project_id"`
	HoursCount  float64   `json:"hours_count"`
	Description string    `json:"description"`
	TimeStart   time.Time `json:"time_start"`
	TimeEnd     time.Time `json:"time_end"`
}

func GetResponseFromModelGoal(goal *models.Goal) *RespGoal {
	return &RespGoal{
		ID:          goal.ID,
		Name:        goal.Name,
		UserID:      goal.UserID,
		ProjectID:   goal.ProjectID,
		HoursCount:  goal.HoursCount,
		Description: goal.Description,
		TimeEnd:     goal.TimeEnd,
		TimeStart:   goal.TimeStart,
	}
}

func GetResponseFromModelGoals(goals []*models.Goal) []*RespGoal {
	result := make([]*RespGoal, 0, 10)
	for _, goal := range goals {
		result = append(result, GetResponseFromModelGoal(goal))
	}

	return result
}
