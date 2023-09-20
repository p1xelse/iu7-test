package redis

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"timetracker/internal/Auth/repository"
	"timetracker/models"
)

type authRepository struct {
	db  *redis.Client
	ctx context.Context
}

func (ar authRepository) CreateCookie(cookie *models.Cookie) error {
	err := ar.db.Set(ar.ctx, cookie.SessionToken, cookie.UserID, cookie.MaxAge).Err()

	if err != nil {
		return errors.Wrap(err, "redis error")
	}

	return nil
}

func (ar authRepository) GetUserByCookie(value string) (string, error) {
	userIdStr, err := ar.db.Get(ar.ctx, value).Result()

	if errors.Is(err, redis.Nil) {
		return "", models.ErrNotFound
	} else if err != nil {
		return "", errors.Wrap(err, "redis error")
	}

	return userIdStr, nil
}

func (ar authRepository) DeleteCookie(value string) error {
	err := ar.db.Del(ar.ctx, value).Err()

	if err != nil {
		return errors.Wrap(err, "redis error")
	}

	return nil
}

func NewAuthRepository(db *redis.Client) repository.RepositoryI {
	return &authRepository{
		db:  db,
		ctx: context.Background(),
	}
}
