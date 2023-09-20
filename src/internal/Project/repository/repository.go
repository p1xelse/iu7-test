package repository

import (
	"timetracker/models"
)

type RepositoryI interface {
	CreateProject(e *models.Project) error
	UpdateProject(e *models.Project) error
	GetProject(id uint64) (*models.Project, error)
	DeleteProject(id uint64) error
	GetUserProjects(userID uint64) ([]*models.Project, error)
}
