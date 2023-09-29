package usecase

import (
	"testing"

	userRepoMock "timetracker/internal/User/repository/mocks"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type UserTestSuite struct {
	suite.Suite
	uc           UsecaseI
	userRepoMock *userRepoMock.RepositoryI
	userBuilder  *testutils.UserBuilder
}

func TestUserTestSuite(t *testing.T) {
	suite.RunSuite(t, new(UserTestSuite))
}

func (s *UserTestSuite) BeforeEach(t provider.T) {
	s.userRepoMock = userRepoMock.NewRepositoryI(t)
	s.uc = New(s.userRepoMock)
	s.userBuilder = testutils.NewUserBuilder()
}

func (s *UserTestSuite) TestUpdateUser(t provider.T) {
	user := s.userBuilder.WithID(1).WithName("user").Build()
	notFoundUser := s.userBuilder.WithID(0).Build()

	s.userRepoMock.On("GetUser", user.ID).Return(&user, nil)
	s.userRepoMock.On("UpdateUser", &user).Return(nil)
	s.userRepoMock.On("GetUser", notFoundUser.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		ArgData *models.User
		Error   error
	}{
		"success": {
			ArgData: &user,
			Error:   nil,
		},
		"user not found": {
			ArgData: &notFoundUser,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.UpdateUser(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *UserTestSuite) TestGetUser(t provider.T) {
	user := s.userBuilder.WithID(1).WithName("user").Build()

	s.userRepoMock.On("GetUser", user.ID).Return(&user, nil)
	result, err := s.uc.GetUser(user.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(user, *result)
}

func (s *UserTestSuite) TestUsecaseGetUsers(t provider.T) {
	users := make([]*models.User, 0, 10)
	err := faker.FakeData(&users)
	t.Assert().NoError(err)

	s.userRepoMock.On("GetUsers").Return(users, nil)

	cases := map[string]struct {
		Users []*models.User
		Error error
	}{
		"success": {
			Users: users,
			Error: nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resUsers, err := s.uc.GetUsers()
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(test.Users, resUsers)
		})
	}
}
