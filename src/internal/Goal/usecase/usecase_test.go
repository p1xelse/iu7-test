package usecase_test

import (
	"github.com/bxcodec/faker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	goalMocks "timetracker/internal/Goal/repository/mocks"
	"timetracker/internal/Goal/usecase"
	"timetracker/models"
)

type TestCaseGetGoal struct {
	ArgData     uint64
	ExpectedRes *models.Goal
	Error       error
}

type TestCaseDeleteGoal struct {
	ArgData []uint64
	Error   error
}

type TestCaseCreateUpdateGoal struct {
	ArgData *models.Goal
	Error   error
}

type TestCaseGetUserGoals struct {
	ArgData     uint64
	ExpectedRes []*models.Goal
	Error       error
}

func TestUsecaseGetGoal(t *testing.T) {
	var mockGoalRes models.Goal
	err := faker.FakeData(&mockGoalRes)
	assert.NoError(t, err)

	mockExpectedGoal := mockGoalRes

	mockGoalRepo := goalMocks.NewRepositoryI(t)

	mockGoalRepo.On("GetGoal", mockGoalRes.ID).Return(&mockGoalRes, nil)

	useCase := usecase.New(mockGoalRepo)

	cases := map[string]TestCaseGetGoal{
		"success": {
			ArgData:     mockGoalRes.ID,
			ExpectedRes: &mockExpectedGoal,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetGoal(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockGoalRepo.AssertExpectations(t)
}

func TestUsecaseUpdateGoal(t *testing.T) {
	var mockGoal, invalidMockGoal models.Goal
	err := faker.FakeData(&mockGoal)
	assert.NoError(t, err)

	invalidMockGoal.ID += mockGoal.ID + 1

	mockGoalRepo := goalMocks.NewRepositoryI(t)

	mockGoalRepo.On("GetGoal", mockGoal.ID).Return(&mockGoal, nil)
	mockGoalRepo.On("UpdateGoal", &mockGoal).Return(nil)

	mockGoalRepo.On("GetGoal", invalidMockGoal.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockGoalRepo)

	cases := map[string]TestCaseCreateUpdateGoal{
		"success": {
			ArgData: &mockGoal,
			Error:   nil,
		},
		"Goal not found": {
			ArgData: &invalidMockGoal,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.UpdateGoal(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockGoalRepo.AssertExpectations(t)
}

func TestUsecaseCreateGoal(t *testing.T) {
	var mockGoal models.Goal
	err := faker.FakeData(&mockGoal)
	assert.NoError(t, err)

	mockGoalRepo := goalMocks.NewRepositoryI(t)

	mockGoalRepo.On("CreateGoal", &mockGoal).Return(nil)

	useCase := usecase.New(mockGoalRepo)

	cases := map[string]TestCaseCreateUpdateGoal{
		"success": {
			ArgData: &mockGoal,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.CreateGoal(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockGoalRepo.AssertExpectations(t)
}

func TestUsecaseDeleteGoal(t *testing.T) {
	var mockGoal, invalidMockGoal models.Goal
	err := faker.FakeData(&mockGoal)
	assert.NoError(t, err)

	invalidMockGoal.ID += mockGoal.ID + 1
	*invalidMockGoal.UserID += *mockGoal.UserID + 1

	mockGoalRepo := goalMocks.NewRepositoryI(t)

	mockGoalRepo.On("GetGoal", mockGoal.ID).Return(&mockGoal, nil)
	mockGoalRepo.On("DeleteGoal", mockGoal.ID).Return(nil)

	mockGoalRepo.On("GetGoal", invalidMockGoal.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockGoalRepo)

	cases := map[string]TestCaseDeleteGoal{
		"success": {
			ArgData: []uint64{mockGoal.ID, *mockGoal.UserID},
			Error:   nil,
		},
		"Goal not found": {
			ArgData: []uint64{invalidMockGoal.ID, *invalidMockGoal.UserID},
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.DeleteGoal(test.ArgData[0], test.ArgData[1])
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockGoalRepo.AssertExpectations(t)
}

func TestUsecaseGetUserGoals(t *testing.T) {
	mockGoalRes := make([]*models.Goal, 0, 10)
	err := faker.FakeData(&mockGoalRes)

	for idx := range mockGoalRes {
		mockGoalRes[idx].UserID = mockGoalRes[0].UserID
	}
	assert.NoError(t, err)

	mockExpectedGoal := mockGoalRes

	mockGoalRepo := goalMocks.NewRepositoryI(t)

	mockGoalRepo.On("GetUserGoals", mockGoalRes[0].UserID).Return(mockGoalRes, nil)

	useCase := usecase.New(mockGoalRepo)

	cases := map[string]TestCaseGetUserGoals{
		"success": {
			ArgData:     *mockGoalRes[0].UserID,
			ExpectedRes: mockExpectedGoal,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetUserGoals(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockGoalRepo.AssertExpectations(t)
}
