package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"timetracker/models"

	"github.com/redis/go-redis/v9"
)

const (
	ttl = time.Duration(time.Hour)
)

type CacheStorageI interface {
	Set(key string, data interface{}) error
	Get(key string) ([]byte, error)
}

type StorageRedis struct {
	db *redis.Client

	ttl time.Duration
	ctx context.Context
}

func NewStorageRedis(m *redis.Client) *StorageRedis {
	storage := StorageRedis{
		db:  m,
		ttl: ttl,
		ctx: context.Background(),
	}

	return &storage
}
func (sr *StorageRedis) Set(key string, data interface{}) error {
	rawValue, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("fail to marshal data for cache: %w", err)
	}

	err = sr.db.Set(sr.ctx, key, rawValue, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (sr *StorageRedis) Get(key string) ([]byte, error) {
	resp, err := sr.db.Get(sr.ctx, key).Bytes()
	if err != nil && err == redis.Nil{
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return resp, nil
}
