package usecase

import (
	"strconv"
	"time"
	authRep "timetracker/internal/Auth/repository"
	userRep "timetracker/internal/User/repository"
	"timetracker/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UsecaseI interface {
	Auth(cookie string) (*models.User, error)
	SignIn(user *models.User) (*models.User, *models.Cookie, error)
	SignUp(user *models.User) (*models.Cookie, error)
	DeleteCookie(value string) error
}

type usecase struct {
	authRepository authRep.RepositoryI
	userRepository userRep.RepositoryI
}

func (u usecase) Auth(cookie string) (*models.User, error) {
	userIdStr, err := u.authRepository.GetUserByCookie(cookie)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.Auth")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.Auth")
	}

	gotUser, err := u.userRepository.GetUser(userId)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.Auth")
	}
	gotUser.Password = ""

	return gotUser, nil
}

func (u usecase) SignIn(user *models.User) (*models.User, *models.Cookie, error) {
	repUsr, err := u.userRepository.GetUserByEmail(user.Email)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Error in func auth.Usecase.SignIn")
	}

	err = bcrypt.CompareHashAndPassword([]byte(repUsr.Password), []byte(user.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, nil, models.ErrInvalidPassword
	} else if err != nil {
		return nil, nil, errors.Wrap(err, "Error in func auth.Usecase.SignIn bcrypt error")
	}

	repUsr.Password = ""

	cookie := models.Cookie{
		UserID:       repUsr.ID,
		SessionToken: uuid.NewString(),
		MaxAge:       (3600 * 24 * 365) * time.Second,
	}

	err = u.authRepository.CreateCookie(&cookie)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error in func auth.Usecase.SignIn")
	}

	return repUsr, &cookie, nil
}

func (u usecase) SignUp(user *models.User) (*models.Cookie, error) {
	_, err := u.userRepository.GetUserByEmail(user.Email)

	if err != models.ErrNotFound && err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.SignUp")
	} else if err == nil {
		return nil, models.ErrConflictEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.SignUp bcrypt error")
	}

	user.Password = string(hashedPassword)

	err = u.userRepository.CreateUser(user)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.SignUp")
	}
	user.Password = ""

	cookie := models.Cookie{
		UserID:       user.ID,
		SessionToken: uuid.NewString(),
		MaxAge:       (3600 * 24 * 365) * time.Second,
	}

	err = u.authRepository.CreateCookie(&cookie)
	if err != nil {
		return nil, errors.Wrap(err, "Error in func auth.Usecase.SignUp")
	}

	return &cookie, nil
}

func (u usecase) DeleteCookie(value string) error {
	_, err := u.authRepository.GetUserByCookie(value)
	if err != nil {
		return errors.Wrap(err, "Error in func auth.Usecase.DeleteCookie")
	}

	err = u.authRepository.DeleteCookie(value)
	if err != nil {
		return errors.Wrap(err, "Error in func auth.Usecase.DeleteCookie")
	}

	return nil
}

func New(uRep userRep.RepositoryI, aRep authRep.RepositoryI) UsecaseI {
	return &usecase{
		userRepository: uRep,
		authRepository: aRep,
	}
}
