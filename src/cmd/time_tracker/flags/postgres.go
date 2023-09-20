package flags

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresFlags struct {
	ConnectionDSN      string        `toml:"dsn"`
	MaxOpenConnections int           `toml:"max-open-connections"`
	ConnectionLifetime time.Duration `toml:"conn-lifetime"`
}

func (f PostgresFlags) Init() (*gorm.DB, error) {
	cfg := postgres.Config{DSN: f.ConnectionDSN}
	db, err := gorm.Open(postgres.New(cfg),
		&gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(f.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(f.ConnectionLifetime)

	return db, nil
}
