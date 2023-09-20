package usecase_test

import (
	"github.com/bxcodec/faker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"testing"
	authMocks "timetracker/internal/Auth/repository/mocks"
	authUsecase "timetracker/internal/Auth/usecase"
	userMocks "timetracker/internal/User/repository/mocks"
	"timetracker/models"
)

type TestCaseSignUp struct {
	ArgData     *models.User
	ExpectedRes uint64
	Error       error
}

type TestCaseSignIn struct {
	ArgData           *models.User
	ExpectedResUser   *models.User
	ExpectedResCookie uint64
	Error             error
}

type TestCaseDeleteCookie struct {
	ArgData string
	Error   error
}

type TestCaseAuth struct {
	ArgData  string
	Expected *models.User
	Error    error
}

func TestUsecaseSignUp(t *testing.T) {
	var mockUserSuccess models.User
	err := faker.FakeData(&mockUserSuccess)
	assert.NoError(t, err)

	var mockUserConflictEmail models.User
	err = faker.FakeData(&mockUserConflictEmail)
	assert.NoError(t, err)

	mockAuthRepo := authMocks.NewRepositoryI(t)
	mockUserRepo := userMocks.NewRepositoryI(t)

	mockUserRepo.On("GetUserByEmail", mockUserSuccess.Email).Return(&mockUserSuccess, models.ErrNotFound)
	mockUserRepo.On("CreateUser", &mockUserSuccess).Return(nil)
	mockAuthRepo.On("CreateCookie", mock.AnythingOfType("*models.Cookie")).Return(nil)
	mockUserRepo.On("GetUserByEmail", mockUserConflictEmail.Email).Return(&mockUserConflictEmail, models.ErrConflictEmail)

	useCase := authUsecase.New(mockUserRepo, mockAuthRepo)

	cases := map[string]TestCaseSignUp{
		"success": {
			ArgData:     &mockUserSuccess,
			ExpectedRes: mockUserSuccess.ID,
			Error:       nil,
		},
		"conflict_email": {
			ArgData: &mockUserConflictEmail,
			Error:   models.ErrConflictEmail,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			cookie, err := useCase.SignUp(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))

			if err == nil {
				assert.Equal(t, test.ExpectedRes, cookie.UserID)
			}
		})
	}
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestUsecaseSignIn(t *testing.T) {
	var mockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	var mockUserSignIn models.User
	mockUserSignIn.Email = mockUser.Email
	mockUserSignIn.Password = mockUser.Password

	var mockUserSignInInvalidPassword models.User
	err = faker.FakeData(&mockUserSignInInvalidPassword.Email)
	assert.NoError(t, err)
	mockUserSignInInvalidPassword.Password = "dfvdf"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(mockUser.Password), 8)
	assert.NoError(t, err)

	mockUser.Password = string(hashedPassword)

	mockAuthRepo := authMocks.NewRepositoryI(t)
	mockUserRepo := userMocks.NewRepositoryI(t)

	mockUserFail := mockUser
	mockUserRepo.On("GetUserByEmail", mockUserSignInInvalidPassword.Email).Return(&mockUserFail, nil)

	mockUserRepo.On("GetUserByEmail", mockUserSignIn.Email).Return(&mockUser, nil)
	mockAuthRepo.On("CreateCookie", mock.AnythingOfType("*models.Cookie")).Return(nil)

	useCase := authUsecase.New(mockUserRepo, mockAuthRepo)

	expectedUser := mockUser
	expectedUser.Password = ""
	cases := map[string]TestCaseSignIn{
		"success": {
			ArgData:           &mockUserSignIn,
			ExpectedResUser:   &expectedUser,
			ExpectedResCookie: mockUser.ID,
			Error:             nil,
		},
		"invalid_password": {
			ArgData: &mockUserSignInInvalidPassword,
			Error:   models.ErrInvalidPassword,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, _, err := useCase.SignIn(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedResUser, user)
			}
		})
	}
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestUsecaseDeleteCookie(t *testing.T) {
	var cookie models.Cookie
	err := faker.FakeData(&cookie)
	assert.NoError(t, err)

	var cookieGetFail models.Cookie
	err = faker.FakeData(&cookieGetFail)
	assert.NoError(t, err)

	var cookieDeleteFail models.Cookie
	err = faker.FakeData(&cookieDeleteFail)
	assert.NoError(t, err)

	mockAuthRepo := authMocks.NewRepositoryI(t)
	mockUserRepo := userMocks.NewRepositoryI(t)

	mockAuthRepo.On("GetUserByCookie", cookie.SessionToken).Return(strconv.Itoa(int(cookie.UserID)), nil)
	mockAuthRepo.On("DeleteCookie", cookie.SessionToken).Return(nil)

	mockAuthRepo.On("GetUserByCookie", cookieGetFail.SessionToken).Return("", models.ErrNotFound)

	mockAuthRepo.On("GetUserByCookie", cookieDeleteFail.SessionToken).Return(strconv.Itoa(int(cookieDeleteFail.UserID)), nil)
	mockAuthRepo.On("DeleteCookie", cookieDeleteFail.SessionToken).Return(models.ErrInternalServerError)

	useCase := authUsecase.New(mockUserRepo, mockAuthRepo)

	cases := map[string]TestCaseDeleteCookie{
		"success": {
			ArgData: cookie.SessionToken,
			Error:   nil,
		},
		"fail_get": {
			ArgData: cookieGetFail.SessionToken,
			Error:   models.ErrNotFound,
		},
		"fail_delete": {
			ArgData: cookieDeleteFail.SessionToken,
			Error:   models.ErrInternalServerError,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.DeleteCookie(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockAuthRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestUsecaseAuth(t *testing.T) {
	var cookie, invalidCookie models.Cookie
	err := faker.FakeData(&cookie)
	assert.NoError(t, err)
	invalidCookie.SessionToken += "lol"

	mockAuthRepo := authMocks.NewRepositoryI(t)
	mockUserRepo := userMocks.NewRepositoryI(t)

	var user models.User
	err = faker.FakeData(&user)
	assert.NoError(t, err)

	user.ID = cookie.UserID

	mockAuthRepo.On("GetUserByCookie", cookie.SessionToken).Return(strconv.Itoa(int(cookie.UserID)), nil)
	mockAuthRepo.On("GetUserByCookie", invalidCookie.SessionToken).Return("", models.ErrNotFound)
	mockUserRepo.On("GetUser", cookie.UserID).Return(&user, nil)

	user.Password = ""

	useCase := authUsecase.New(mockUserRepo, mockAuthRepo)

	cases := map[string]TestCaseAuth{
		"success": {
			ArgData:  cookie.SessionToken,
			Expected: &user,
			Error:    nil,
		},
		"not_found": {
			ArgData:  invalidCookie.SessionToken,
			Expected: nil,
			Error:    models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			gotUser, err := useCase.Auth(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
			if err == nil {
				assert.Equal(t, test.Expected, gotUser)
			}
		})
	}
	mockAuthRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
