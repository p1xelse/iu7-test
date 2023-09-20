package models

type RoleType int64

const (
	DefaultUser RoleType = iota
	Admin
)

func (s RoleType) String() string {
	switch s {
	case DefaultUser:
		return "user"
	case Admin:
		return "admin"
	}
	return "unknown"
}

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	About    string `json:"about"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
