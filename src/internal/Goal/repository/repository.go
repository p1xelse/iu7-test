package repository

import (
	"timetracker/models"
)

type RepositoryI interface {
	CreateGoal(g *models.Goal) error
	UpdateGoal(g *models.Goal) error
	GetGoal(id uint64) (*models.Goal, error)
	DeleteGoal(id uint64) error
	GetUserGoals(userID uint64) ([]*models.Goal, error)
}
