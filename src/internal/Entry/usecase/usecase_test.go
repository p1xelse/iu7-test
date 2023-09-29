package usecase

import (
	"testing"

	entryRepoMock "timetracker/internal/Entry/repository/mocks"
	tagRepoMock "timetracker/internal/Tag/repository/mocks"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"github.com/bxcodec/faker"
)

type EntryTestSuite struct {
	suite.Suite
	uc            UsecaseI
	entryRepoMock *entryRepoMock.RepositoryI
	tagRepoMock   *tagRepoMock.RepositoryI
	entryBuilder  *testutils.EntryBuilder
	tagBuilder    *testutils.TagBuilder
}

func TestEntryTestSuite(t *testing.T) {
	suite.RunSuite(t, new(EntryTestSuite))
}

func (s *EntryTestSuite) BeforeEach(t provider.T) {
	s.entryRepoMock = entryRepoMock.NewRepositoryI(t)
	s.tagRepoMock = tagRepoMock.NewRepositoryI(t)
	s.uc = New(s.entryRepoMock, s.tagRepoMock)
	s.entryBuilder = testutils.NewEntryBuilder()
	s.tagBuilder = testutils.NewTagBuilder()
}

func (s *EntryTestSuite) TestCreateEntry(t provider.T) {
	entry := s.entryBuilder.WithID(1).WithDescription("entry").Build()

	s.entryRepoMock.On("CreateEntry", &entry).Return(nil)
	err := s.uc.CreateEntry(&entry)

	t.Assert().NoError(err)
}

func (s *EntryTestSuite) TestUpdateEntry(t provider.T) {
	entry := s.entryBuilder.
		WithID(1).
		WithDescription("entry").
		WithUserID(1).
		WithTagList([]models.Tag{s.tagBuilder.WithName("hello").Build()}).
		Build()
	notFoundEntry := s.entryBuilder.WithID(0).Build()

	s.entryRepoMock.On("GetEntry", entry.ID).Return(&entry, nil)
	s.entryRepoMock.On("UpdateEntry", &entry).Return(nil)
	s.tagRepoMock.On("UpdateEntryTags", entry.ID, entry.TagList).Return(nil)
	s.entryRepoMock.On("GetEntry", notFoundEntry.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		ArgData *models.Entry
		Error   error
	}{
		"success": {
			ArgData: &entry,
			Error:   nil,
		},
		"entry not found": {
			ArgData: &notFoundEntry,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.UpdateEntry(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *EntryTestSuite) TestGetEntry(t provider.T) {
	entry := s.entryBuilder.WithID(1).WithDescription("entry").Build()
	tags := make([]*models.Tag, 0, 10)
	err := faker.FakeData(&tags)
	t.Assert().NoError(err)
	for _, tag := range tags {
		entry.TagList = append(entry.TagList, *tag)
	}

	s.entryRepoMock.On("GetEntry", entry.ID).Return(&entry, nil)
	s.tagRepoMock.On("GetEntryTags", entry.ID).Return(tags, nil)
	result, err := s.uc.GetEntry(entry.ID)

	t.Require().NoError(err)
	t.Assert().Equal(entry, *result)
}

func (s *EntryTestSuite) TestDeleteEntry(t provider.T) {
	entry := s.entryBuilder.WithID(1).WithDescription("entry").WithUserID(2).Build()
	notFoundEntry := s.entryBuilder.WithID(0).WithUserID(3).Build()

	s.entryRepoMock.On("GetEntry", entry.ID).Return(&entry, nil)
	s.entryRepoMock.On("DeleteEntry", entry.ID).Return(nil)
	s.entryRepoMock.On("GetEntry", notFoundEntry.ID).Return(nil, models.ErrNotFound)

	s.tagRepoMock.On("DeleteEntryTags", entry.ID).Return(nil)

	cases := map[string]struct {
		EntryID uint64
		UserID  uint64
		Error   error
	}{
		"success": {
			EntryID: entry.ID,
			UserID:  *entry.UserID,
			Error:   nil,
		},
		"entry not found": {
			EntryID: notFoundEntry.ID,
			UserID:  *notFoundEntry.UserID,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.DeleteEntry(test.EntryID, test.UserID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *EntryTestSuite) TestUsecaseGetUserEntrys(t provider.T) {
	entries := make([]*models.Entry, 0, 10)
	err := faker.FakeData(&entries)
	tagList := []models.Tag{s.tagBuilder.Build(), s.tagBuilder.Build(), s.tagBuilder.Build()}
	t.Assert().NoError(err)

	for idx := range entries {
		entries[idx].UserID = entries[0].UserID
		entries[idx].TagList = tagList
	}
	s.entryRepoMock.On("GetUserEntries", *entries[0].UserID).Return(entries, nil)
	for _, entry := range entries {
		s.tagRepoMock.On("GetEntryTags", entry.ID).Return(testutils.MakePointerSlice(entry.TagList), nil)
	}

	cases := map[string]struct {
		UserID uint64
		Entrys []*models.Entry
		Error  error
	}{
		"success": {
			UserID: *entries[0].UserID,
			Entrys: entries,
			Error:  nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resEntrys, err := s.uc.GetUserEntries(test.UserID)
			t.Require().ErrorIs(err, test.Error)
			t.Assert().Equal(entries, resEntrys)
		})
	}
}
