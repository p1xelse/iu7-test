package usecase_test

import (
	"github.com/bxcodec/faker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	userMocks "timetracker/internal/User/repository/mocks"
	"timetracker/internal/User/usecase"
	"timetracker/models"
)

type TestCaseGetUser struct {
	ArgData     uint64
	ExpectedRes *models.User
	Error       error
}

type TestCaseCreateUpdateUser struct {
	ArgData *models.User
	Error   error
}

func TestUsecaseGetUser(t *testing.T) {
	var mockUserRes models.User
	err := faker.FakeData(&mockUserRes)
	assert.NoError(t, err)
	mockUserRes.Password = ""

	mockExpectedUser := mockUserRes

	mockUserRepo := userMocks.NewRepositoryI(t)

	mockUserRepo.On("GetUser", mockUserRes.ID).Return(&mockUserRes, nil)

	useCase := usecase.New(mockUserRepo)

	cases := map[string]TestCaseGetUser{
		"success": {
			ArgData:     mockUserRes.ID,
			ExpectedRes: &mockExpectedUser,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetUser(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockUserRepo.AssertExpectations(t)
}

func TestUsecaseUpdateUser(t *testing.T) {
	var mockUser, invalidMockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	invalidMockUser.ID += mockUser.ID + 1

	mockUserRepo := userMocks.NewRepositoryI(t)

	mockUserRepo.On("GetUser", mockUser.ID).Return(&mockUser, nil)
	mockUserRepo.On("UpdateUser", &mockUser).Return(nil)

	mockUserRepo.On("GetUser", invalidMockUser.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockUserRepo)

	cases := map[string]TestCaseCreateUpdateUser{
		"success": {
			ArgData: &mockUser,
			Error:   nil,
		},
		"User not found": {
			ArgData: &invalidMockUser,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.UpdateUser(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockUserRepo.AssertExpectations(t)
}
