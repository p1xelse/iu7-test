package postgres

import (
	"strconv"
	"time"
	"timetracker/internal/Auth/repository"
	"timetracker/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type authRepositoryPostgres struct {
	db *gorm.DB
}

type Cookie struct {
	UserID       *uint64   `gorm:"column:user_id"`
	SessionToken string    `gorm:"column:session_token"`
	ExpireTime   time.Time `gorm:"column:expire_time"`
}

func (Cookie) TableName() string {
	return "cookie"
}

func toPostgresCookie(e *models.Cookie) *Cookie {
	return &Cookie{
		UserID:       &e.UserID,
		SessionToken: e.SessionToken,
		ExpireTime:   time.Now().Add(e.MaxAge),
	}
}

func toModelCookie(e *Cookie) *models.Cookie {
	return &models.Cookie{
		UserID:       *e.UserID,
		SessionToken: e.SessionToken,
	}
}

func NewAuthRepositoryPostgres(db *gorm.DB) repository.RepositoryI {
	return &authRepositoryPostgres{
		db: db,
	}
}

func (ar authRepositoryPostgres) CreateCookie(cookie *models.Cookie) error {
	postgresCookie := toPostgresCookie(cookie)

	tx := ar.db.Create(postgresCookie)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table cookie)")
	}

	return nil
}

func (ar authRepositoryPostgres) GetUserByCookie(value string) (string, error) {
	var postgresCookie Cookie
	tx := ar.db.Where(&Cookie{SessionToken: value}).Take(&postgresCookie)

	if tx.Error != nil {
		return "", errors.Wrap(tx.Error, "database error (table cookie)")
	}

	if postgresCookie.ExpireTime.Before(time.Now()) {
		err := ar.DeleteCookie(value)

		if err != nil {
			return "", errors.Wrap(tx.Error, "database error (table cookie)")
		}

		return "", models.ErrNotFound
	}


	return strconv.Itoa(int(*postgresCookie.UserID)), nil
}

func (ar authRepositoryPostgres) DeleteCookie(value string) error {
	tx := ar.db.Delete(&Cookie{SessionToken: value})

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table cookie)")
	}

	return nil
}
