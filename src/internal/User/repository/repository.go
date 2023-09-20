package repository

import (
	"timetracker/models"
)

type RepositoryI interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUser(id uint64) (*models.User, error)
	GetUsers() ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsersByIDs(userIDs []uint64) ([]*models.User, error)
}
