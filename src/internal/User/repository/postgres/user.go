package postgres

import (
	"fmt"
	"timetracker/internal/User/repository"
	"timetracker/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	About    string `gorm:"column:about"`
	Role     string `gorm:"column:role"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "users"
}

func toPostgresUser(u *models.User) *User {
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		About:    u.About,
		Role:     u.Role,
		Password: u.Password,
	}
}

func toModelUser(u *User) *models.User {
	return &models.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		About:    u.About,
		Role:     u.Role,
		Password: u.Password,
	}
}

func toModelUsers(entries []*User) []*models.User {
	out := make([]*models.User, len(entries))

	for i, b := range entries {
		out[i] = toModelUser(b)
	}

	return out
}

type userRepository struct {
	db *gorm.DB
}

func (ur userRepository) CreateUser(user *models.User) error {
	postgresUser := toPostgresUser(user)

	tx := ur.db.Create(postgresUser)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table user)")
	}

	user.ID = postgresUser.ID
	return nil
}

func (ur userRepository) UpdateUser(user *models.User) error {
	postgresUser := toPostgresUser(user)

	tx := ur.db.Omit("id").Updates(postgresUser)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table user)")
	}

	return nil
}

func (ur userRepository) GetUser(id uint64) (*models.User, error) {
	var user User

	tx := ur.db.Where(&User{ID: id}).Take(&user)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table user)")
	}

	return toModelUser(&user), nil
}

func (ur userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user User
	fmt.Println("email", email)

	tx := ur.db.Where(&User{Email: email}).Take(&user)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table user)")
	}

	return toModelUser(&user), nil
}

func (ur userRepository) GetUsers() ([]*models.User, error) {
	users := make([]*User, 0, 10)
	tx := ur.db.Omit("password").Find(&users)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return toModelUsers(users), nil
}

func (ur userRepository) GetUsersByIDs(userIDs []uint64) ([]*models.User, error) {
	users := make([]*User, 0, 10)
	tx := ur.db.Omit("password").Find(&users, userIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return toModelUsers(users), nil
}

func NewUserRepository(db *gorm.DB) repository.RepositoryI {
	return &userRepository{
		db: db,
	}
}
