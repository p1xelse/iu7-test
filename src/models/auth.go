package models

import "time"

type Cookie struct {
	SessionToken string
	UserID       uint64
	MaxAge       time.Duration
}
