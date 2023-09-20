package usecase

import (
	userRep "timetracker/internal/User/repository"
	"timetracker/models"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UsecaseI interface {
	UpdateUser(e *models.User) error
	GetUser(id uint64) (*models.User, error)
	GetUsers() ([]*models.User, error)
}

type usecase struct {
	userRepository userRep.RepositoryI
}

func (u *usecase) UpdateUser(user *models.User) error {
	_, err := u.userRepository.GetUser(user.ID)
	if err != nil {
		return errors.Wrap(err, "user repository error")
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			return errors.Wrap(err, "Error in func user.Usecase.UpdateUser bcrypt")
		}

		user.Password = string(hashedPassword)
	}

	err = u.userRepository.UpdateUser(user)
	if err != nil {
		return errors.Wrap(err, "Error in func user.Usecase.UpdateUser")
	}

	return nil
}

func (u *usecase) GetUser(id uint64) (*models.User, error) {
	user, err := u.userRepository.GetUser(id)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func user.Usecase.GetUser")
	}

	return user, nil
}

func (u *usecase) GetUsers() ([]*models.User, error) {
	users, err := u.userRepository.GetUsers()
	if err != nil {
		return nil, errors.Wrap(err, "Error in func user.Usecase.GetUsers")
	}

	return users, nil
}

func New(uRep userRep.RepositoryI) UsecaseI {
	return &usecase{
		userRepository: uRep,
	}
}
