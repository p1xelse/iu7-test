package usecase_test

import (
	"github.com/bxcodec/faker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	goalMocks "timetracker/internal/Project/repository/mocks"
	"timetracker/internal/Project/usecase"
	"timetracker/models"
)

type TestCaseGetProject struct {
	ArgData     uint64
	ExpectedRes *models.Project
	Error       error
}

type TestCaseDeleteProject struct {
	ArgData []*uint64
	Error   error
}

type TestCaseCreateUpdateProject struct {
	ArgData *models.Project
	Error   error
}

type TestCaseGetUserProjects struct {
	ArgData     *uint64
	ExpectedRes []*models.Project
	Error       error
}

func TestUsecaseGetProject(t *testing.T) {
	var mockProjectRes models.Project
	err := faker.FakeData(&mockProjectRes)
	assert.NoError(t, err)

	mockExpectedProject := mockProjectRes

	mockProjectRepo := goalMocks.NewRepositoryI(t)

	mockProjectRepo.On("GetProject", mockProjectRes.ID).Return(&mockProjectRes, nil)

	useCase := usecase.New(mockProjectRepo, nil)

	cases := map[string]TestCaseGetProject{
		"success": {
			ArgData:     mockProjectRes.ID,
			ExpectedRes: &mockExpectedProject,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetProject(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestUsecaseUpdateProject(t *testing.T) {
	var mockProject, invalidMockProject models.Project
	err := faker.FakeData(&mockProject)
	assert.NoError(t, err)

	invalidMockProject.ID += mockProject.ID + 1

	mockProjectRepo := goalMocks.NewRepositoryI(t)

	mockProjectRepo.On("GetProject", mockProject.ID).Return(&mockProject, nil)
	mockProjectRepo.On("UpdateProject", &mockProject).Return(nil)

	mockProjectRepo.On("GetProject", invalidMockProject.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockProjectRepo, nil)

	cases := map[string]TestCaseCreateUpdateProject{
		"success": {
			ArgData: &mockProject,
			Error:   nil,
		},
		"Project not found": {
			ArgData: &invalidMockProject,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.UpdateProject(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestUsecaseCreateProject(t *testing.T) {
	var mockProject models.Project
	err := faker.FakeData(&mockProject)
	assert.NoError(t, err)

	mockProjectRepo := goalMocks.NewRepositoryI(t)

	mockProjectRepo.On("CreateProject", &mockProject).Return(nil)

	useCase := usecase.New(mockProjectRepo, nil)

	cases := map[string]TestCaseCreateUpdateProject{
		"success": {
			ArgData: &mockProject,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.CreateProject(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestUsecaseDeleteProject(t *testing.T) {
	var mockProject, invalidMockProject models.Project
	err := faker.FakeData(&mockProject)
	assert.NoError(t, err)

	invalidMockProject.ID += mockProject.ID + 1
	*invalidMockProject.UserID += *mockProject.UserID + 1

	mockProjectRepo := goalMocks.NewRepositoryI(t)

	mockProjectRepo.On("GetProject", mockProject.ID).Return(&mockProject, nil)
	mockProjectRepo.On("DeleteProject", mockProject.ID).Return(nil)

	mockProjectRepo.On("GetProject", invalidMockProject.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockProjectRepo, nil)

	cases := map[string]TestCaseDeleteProject{
		"success": {
			ArgData: []*uint64{&mockProject.ID, mockProject.UserID},
			Error:   nil,
		},
		"Project not found": {
			ArgData: []*uint64{&invalidMockProject.ID, invalidMockProject.UserID},
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.DeleteProject(*test.ArgData[0], *test.ArgData[1])
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestUsecaseGetUserProjects(t *testing.T) {
	mockProjectRes := make([]*models.Project, 0, 10)
	err := faker.FakeData(&mockProjectRes)

	for idx := range mockProjectRes {
		mockProjectRes[idx].UserID = mockProjectRes[0].UserID
	}
	assert.NoError(t, err)

	mockExpectedProject := mockProjectRes

	mockProjectRepo := goalMocks.NewRepositoryI(t)

	mockProjectRepo.On("GetUserProjects", mockProjectRes[0].UserID).Return(mockProjectRes, nil)

	useCase := usecase.New(mockProjectRepo, nil)

	cases := map[string]TestCaseGetUserProjects{
		"success": {
			ArgData:     mockProjectRes[0].UserID,
			ExpectedRes: mockExpectedProject,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetUserProjects(*test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockProjectRepo.AssertExpectations(t)
}
