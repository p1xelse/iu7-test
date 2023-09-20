package usecase_test

import (
	"github.com/bxcodec/faker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	entryMocks "timetracker/internal/Entry/repository/mocks"
	"timetracker/internal/Entry/usecase"
	tagMocks "timetracker/internal/Tag/repository/mocks"
	"timetracker/models"
)

type TestCaseGetEntry struct {
	ArgData     uint64
	ExpectedRes *models.Entry
	Error       error
}

type TestCaseDeleteEntry struct {
	ArgData []*uint64
	Error   error
}

type TestCaseCreateUpdateEntry struct {
	ArgData *models.Entry
	Error   error
}

type TestCaseGetUserEntries struct {
	ArgData     *uint64
	ExpectedRes []*models.Entry
	Error       error
}

func TestUsecaseGetEntry(t *testing.T) {
	var mockEntryRes models.Entry
	err := faker.FakeData(&mockEntryRes)
	assert.NoError(t, err)

	mockExpectedEntry := mockEntryRes
	mockExpectedEntry.TagList = nil

	mockEntryRepo := entryMocks.NewRepositoryI(t)
	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockTags := make([]*models.Tag, 0, 10)
	err = faker.FakeData(&mockTags)
	assert.NoError(t, err)
	for _, tag := range mockTags {
		mockExpectedEntry.TagList = append(mockExpectedEntry.TagList, *tag)
	}

	mockEntryRepo.On("GetEntry", mockEntryRes.ID).Return(&mockEntryRes, nil)
	mockTagRepo.On("GetEntryTags", mockEntryRes.ID).Return(mockTags, nil)

	useCase := usecase.New(mockEntryRepo, mockTagRepo, nil)

	cases := map[string]TestCaseGetEntry{
		"success": {
			ArgData:     mockEntryRes.ID,
			ExpectedRes: &mockExpectedEntry,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetEntry(test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockEntryRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseCreateEntry(t *testing.T) {
	var mockEntry models.Entry
	err := faker.FakeData(&mockEntry)
	assert.NoError(t, err)

	mockEntryRepo := entryMocks.NewRepositoryI(t)
	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockEntryRepo.On("CreateEntry", &mockEntry).Return(nil)
	mockTagRepo.On("CreateEntryTags", mockEntry.ID, mockEntry.TagList).Return(nil)

	useCase := usecase.New(mockEntryRepo, mockTagRepo, nil)

	cases := map[string]TestCaseCreateUpdateEntry{
		"success": {
			ArgData: &mockEntry,
			Error:   nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.CreateEntry(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockEntryRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseUpdateEntry(t *testing.T) {
	var mockEntry, invalidMockEntry models.Entry
	err := faker.FakeData(&mockEntry)
	assert.NoError(t, err)

	invalidMockEntry.ID = mockEntry.ID + 1

	mockEntryRepo := entryMocks.NewRepositoryI(t)
	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockEntryRepo.On("UpdateEntry", &mockEntry).Return(nil)
	mockTagRepo.On("UpdateEntryTags", mockEntry.ID, mockEntry.TagList).Return(nil)

	mockEntryRepo.On("GetEntry", mockEntry.ID).Return(&mockEntry, nil)
	mockEntryRepo.On("GetEntry", invalidMockEntry.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockEntryRepo, mockTagRepo, nil)

	cases := map[string]TestCaseCreateUpdateEntry{
		"success": {
			ArgData: &mockEntry,
			Error:   nil,
		},
		"Entry not found": {
			ArgData: &invalidMockEntry,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.UpdateEntry(test.ArgData)
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockEntryRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseDeleteEntry(t *testing.T) {
	var mockEntry, invalidMockEntry models.Entry
	err := faker.FakeData(&mockEntry)
	assert.NoError(t, err)

	invalidMockEntry.ID += mockEntry.ID + 1
	*invalidMockEntry.UserID += *mockEntry.UserID + 1

	mockEntryRepo := entryMocks.NewRepositoryI(t)
	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockEntryRepo.On("DeleteEntry", mockEntry.ID).Return(nil)
	mockTagRepo.On("DeleteEntryTags", mockEntry.ID).Return(nil)

	mockEntryRepo.On("GetEntry", mockEntry.ID).Return(&mockEntry, nil)
	mockEntryRepo.On("GetEntry", invalidMockEntry.ID).Return(nil, models.ErrNotFound)

	useCase := usecase.New(mockEntryRepo, mockTagRepo, nil)

	cases := map[string]TestCaseDeleteEntry{
		"success": {
			ArgData: []*uint64{&mockEntry.ID, mockEntry.UserID},
			Error:   nil,
		},
		"Entry not found": {
			ArgData: []*uint64{&invalidMockEntry.ID, invalidMockEntry.UserID},
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := useCase.DeleteEntry(*test.ArgData[0], *test.ArgData[1])
			require.Equal(t, test.Error, errors.Cause(err))
		})
	}
	mockEntryRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}

func TestUsecaseGetUserEntries(t *testing.T) {
	mockEntryRes := make([]*models.Entry, 0, 10)
	err := faker.FakeData(&mockEntryRes)
	assert.NoError(t, err)

	mockExpectedEntry := mockEntryRes
	mockExpectedEntry[0].TagList = nil

	mockTags := make([]*models.Tag, 0, 10)
	err = faker.FakeData(&mockTags)
	assert.NoError(t, err)

	for _, tag := range mockTags {
		mockExpectedEntry[0].TagList = append(mockExpectedEntry[0].TagList, *tag)
	}

	for idx := range mockExpectedEntry {
		mockExpectedEntry[idx].UserID = mockExpectedEntry[0].UserID
		mockExpectedEntry[idx].TagList = mockExpectedEntry[0].TagList
	}

	mockEntryRepo := entryMocks.NewRepositoryI(t)
	mockTagRepo := tagMocks.NewRepositoryI(t)

	mockEntryRepo.On("GetUserEntries", mockExpectedEntry[0].UserID).Return(mockExpectedEntry, nil)

	for idx := range mockExpectedEntry {
		mockTagRepo.On("GetEntryTags", mockExpectedEntry[idx].ID).Return(mockTags, nil)
	}

	useCase := usecase.New(mockEntryRepo, mockTagRepo, nil)

	cases := map[string]TestCaseGetUserEntries{
		"success": {
			ArgData:     mockEntryRes[0].UserID,
			ExpectedRes: mockExpectedEntry,
			Error:       nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			user, err := useCase.GetUserEntries(*test.ArgData)
			require.Equal(t, test.Error, err)

			if err == nil {
				assert.Equal(t, test.ExpectedRes, user)
			}
		})
	}
	mockEntryRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
}
