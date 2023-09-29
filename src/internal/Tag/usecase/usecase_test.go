package usecase

import (
	"testing"

	tagRepoMock "timetracker/internal/Tag/repository/mocks"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type TagTestSuite struct {
	suite.Suite
	uc          UsecaseI
	tagRepoMock *tagRepoMock.RepositoryI
	tagBuilder  *testutils.TagBuilder
}

func TestTagTestSuite(t *testing.T) {
	suite.RunSuite(t, new(TagTestSuite))
}

func (s *TagTestSuite) BeforeEach(t provider.T) {
	s.tagRepoMock = tagRepoMock.NewRepositoryI(t)
	s.uc = New(s.tagRepoMock)
	s.tagBuilder = testutils.NewTagBuilder()
}

func (s *TagTestSuite) TestCreateTag(t provider.T) {
	tag := s.tagBuilder.WithID(1).WithName("tag").Build()

	s.tagRepoMock.On("CreateTag", &tag).Return(nil)
	err := s.uc.CreateTag(&tag)

	t.Assert().NoError(err)
}

func (s *TagTestSuite) TestUpdateTag(t provider.T) {
	tag := s.tagBuilder.WithID(1).WithName("tag").Build()
	notFoundTag := s.tagBuilder.WithID(0).Build()

	s.tagRepoMock.On("GetTag", tag.ID).Return(&tag, nil)
	s.tagRepoMock.On("UpdateTag", &tag).Return(nil)
	s.tagRepoMock.On("GetTag", notFoundTag.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		ArgData *models.Tag
		Error   error
	}{
		"success": {
			ArgData: &tag,
			Error:   nil,
		},
		"tag not found": {
			ArgData: &notFoundTag,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.UpdateTag(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *TagTestSuite) TestGetTag(t provider.T) {
	tag := s.tagBuilder.WithID(1).WithName("tag").Build()

	s.tagRepoMock.On("GetTag", tag.ID).Return(&tag, nil)
	result, err := s.uc.GetTag(tag.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(tag, *result)
}

func (s *TagTestSuite) TestDeleteTag(t provider.T) {
	tag := s.tagBuilder.WithID(1).WithName("tag").WithUserID(2).Build()
	notFoundTag := s.tagBuilder.WithID(0).WithUserID(3).Build()

	s.tagRepoMock.On("GetTag", tag.ID).Return(&tag, nil)
	s.tagRepoMock.On("DeleteTag", tag.ID).Return(nil)
	s.tagRepoMock.On("GetTag", notFoundTag.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		TagID  uint64
		UserID uint64
		Error  error
	}{
		"success": {
			TagID:  tag.ID,
			UserID: tag.UserID,
			Error:  nil,
		},
		"tag not found": {
			TagID:  notFoundTag.ID,
			UserID: notFoundTag.UserID,
			Error:  models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.DeleteTag(test.TagID, test.UserID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *TagTestSuite) TestUsecaseGetUserTags(t provider.T) {
	tags := make([]*models.Tag, 0, 10)
	err := faker.FakeData(&tags)

	for idx := range tags {
		tags[idx].UserID = tags[0].UserID
	}
	t.Assert().NoError(err)
	s.tagRepoMock.On("GetUserTags", tags[0].UserID).Return(tags, nil)

	cases := map[string]struct {
		UserID uint64
		Tags   []*models.Tag
		Error  error
	}{
		"success": {
			UserID: tags[0].UserID,
			Tags:   tags,
			Error:  nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resTags, err := s.uc.GetUserTags(test.UserID)
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(tags, resTags)
		})
	}
}
