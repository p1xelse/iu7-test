package models

type Project struct {
	ID        uint64
	UserID    *uint64
	Name      string
	About     string
	Color     string
	IsPrivate bool
	TotalCountHours float64
}
