package models

type Tag struct {
	ID     uint64
	UserID uint64
	Name   string
	About  string
	Color  string
}
