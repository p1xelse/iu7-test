package repository

import "timetracker/models"

type RepositoryI interface {
	CreateCookie(cookie *models.Cookie) error
	GetUserByCookie(value string) (string, error)
	DeleteCookie(value string) error
}
