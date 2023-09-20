package usecase_test

import (
	"github.com/bxcodec/faker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	tagMocks "timetracker/internal/Tag/repository/mocks"
	"timetracker/internal/Tag/usecase"
	"timetracker/models"
)

type TestCaseGetTag struct {
	ArgData     uint64
	ExpectedRes *models.Tag
	Error       error
}

type TestCaseDeleteTag struct {
	ArgData []uint64
	Error   error
}

type TestCaseCreateUpdateTag struct {
	ArgData *models.Tag
	Error   error
}

type TestCaseGetUserTags struct {
	ArgData     uint64
	ExpectedRes []*models.Tag
	Error       error
}

func TestUsecaseGetTag(t *testing.T) {
	var mockTagRes models.Tag
	err := faker.FakeData(&mockTagRes)
	assert.NoError(t, err)

	mockExpectedTag := mockTagRes

	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockTagRepo.On("GetTag", mockTagRes.ID).Return(&mockTagRes, nil)

	useCase := usecase.New(mockTagRepo)

	cases := map[string]TestCaseGetTag{
		"success": {
			ArgData:     mockTagRes.ID,
			ExpectedRes: &mockExpectedTag,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetTag(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseUpdateTag(t *testing.T) {
	var mockTag, invalidMockTag models.Tag
	err := faker.FakeData(&mockTag)
	assert.NoError(t, err)

	invalidMockTag.ID += mockTag.ID + 1

	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockTagRepo.On("GetTag", mockTag.ID).Return(&mockTag, nil)
	mockTagRepo.On("UpdateTag", &mockTag).Return(nil)

	mockTagRepo.On("GetTag", invalidMockTag.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockTagRepo)

	cases := map[string]TestCaseCreateUpdateTag{
		"success": {
			ArgData: &mockTag,
			Error:   nil,
		},
		"Tag not found": {
			ArgData: &invalidMockTag,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.UpdateTag(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseCreateTag(t *testing.T) {
	var mockTag models.Tag
	err := faker.FakeData(&mockTag)
	assert.NoError(t, err)

	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockTagRepo.On("CreateTag", &mockTag).Return(nil)

	useCase := usecase.New(mockTagRepo)

	cases := map[string]TestCaseCreateUpdateTag{
		"success": {
			ArgData: &mockTag,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.CreateTag(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseDeleteTag(t *testing.T) {
	var mockTag, invalidMockTag models.Tag
	err := faker.FakeData(&mockTag)
	assert.NoError(t, err)

	invalidMockTag.ID += mockTag.ID + 1
	invalidMockTag.UserID += mockTag.UserID + 1

	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockTagRepo.On("GetTag", mockTag.ID).Return(&mockTag, nil)
	mockTagRepo.On("DeleteTag", mockTag.ID).Return(nil)

	mockTagRepo.On("GetTag", invalidMockTag.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockTagRepo)

	cases := map[string]TestCaseDeleteTag{
		"success": {
			ArgData: []uint64{mockTag.ID, mockTag.UserID},
			Error:   nil,
		},
		"Tag not found": {
			ArgData: []uint64{invalidMockTag.ID, invalidMockTag.UserID},
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.DeleteTag(test.ArgData[0], test.ArgData[1])
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseGetUserTags(t *testing.T) {
	mockTagRes := make([]*models.Tag, 0, 10)
	err := faker.FakeData(&mockTagRes)

	for idx := range mockTagRes {
		mockTagRes[idx].UserID = mockTagRes[0].UserID
	}
	assert.NoError(t, err)

	mockExpectedTag := mockTagRes

	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockTagRepo.On("GetUserTags", mockTagRes[0].UserID).Return(mockTagRes, nil)

	useCase := usecase.New(mockTagRepo)

	cases := map[string]TestCaseGetUserTags{
		"success": {
			ArgData:     mockTagRes[0].UserID,
			ExpectedRes: mockExpectedTag,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetUserTags(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockTagRepo.AssertExpectations(t)
}
