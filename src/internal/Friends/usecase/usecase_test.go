package usecase

import (
	"testing"

	friendsRepoMock "timetracker/internal/Friends/repository/mocks"
	userRepoMock "timetracker/internal/User/repository/mocks"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type FriendsTestSuite struct {
	suite.Suite
	uc          UsecaseI
	fRepoMock   *friendsRepoMock.RepositoryI
	uRepoMock   *userRepoMock.RepositoryI
	fRelBuilder *testutils.FriendRelationBuilder
	userBuilder *testutils.UserBuilder
}

func TestFriendsTestSuite(t *testing.T) {
	suite.RunSuite(t, new(FriendsTestSuite))
}

func (s *FriendsTestSuite) BeforeEach(t provider.T) {
	s.fRepoMock = friendsRepoMock.NewRepositoryI(t)
	s.uRepoMock = userRepoMock.NewRepositoryI(t)
	s.uc = New(s.fRepoMock, s.uRepoMock)
	s.fRelBuilder = testutils.NewFriendRelationBuilder()
	s.userBuilder = testutils.NewUserBuilder()
}

func (s *FriendsTestSuite) TestCreateFriendRelation(t provider.T) {
	friendRel := s.fRelBuilder.WithSubID(1).WithUserID(2).Build()
	friendRelBad := s.fRelBuilder.WithSubID(1).WithUserID(1).Build()

	s.fRepoMock.On("CheckFriends", &friendRel).Return(false, nil)
	s.fRepoMock.On("CreateFriendRelation", &friendRel).Return(nil)

	t.Assert().NoError(s.uc.CreateFriendRelation(&friendRel))
	t.Assert().ErrorIs(s.uc.CreateFriendRelation(&friendRelBad), models.ErrBadRequest)
}

func (s *FriendsTestSuite) TestDeleteFriendRelation(t provider.T) {
	friendRel := s.fRelBuilder.WithSubID(1).WithUserID(2).Build()
	friendRelBad := s.fRelBuilder.WithSubID(1).WithUserID(1).Build()

	s.fRepoMock.On("CheckFriends", &friendRel).Return(true, nil)
	s.fRepoMock.On("DeleteFriendRelation", &friendRel).Return(nil)

	t.Assert().NoError(s.uc.DeleteFriendRelation(&friendRel))
	t.Assert().ErrorIs(s.uc.DeleteFriendRelation(&friendRelBad), models.ErrBadRequest)
}

func (s *FriendsTestSuite) TestCheckIsFriends(t provider.T) {
	userID1 := uint64(1)
	userID2 := uint64(2)
	fRel1 := s.fRelBuilder.WithSubID(userID1).WithUserID(userID2).Build()
	fRel2 := s.fRelBuilder.WithSubID(userID2).WithUserID(userID1).Build()

	s.fRepoMock.On("CheckFriends", &fRel1).Return(true, nil)
	s.fRepoMock.On("CheckFriends", &fRel2).Return(true, nil)

	result1, err1 := s.uc.CheckIsFriends(userID1, userID2)
	_, err2 := s.uc.CheckIsFriends(userID1, userID1)

	t.Assert().NoError(err1)
	t.Assert().Equal(true, result1)
	t.Assert().ErrorIs(err2, models.ErrBadRequest)
}

func (s *FriendsTestSuite) TestGetUserSubs(t provider.T) {
	userID := uint64(1)
	subsIDs := []uint64{2, 3, 4, 5}
	subs := testutils.BuildUsersByIDs(s.userBuilder, subsIDs)

	s.fRepoMock.On("GetUserSubs", userID).Return(subsIDs, nil)
	s.uRepoMock.On("GetUsersByIDs", subsIDs).Return(subs, nil)

	result, err := s.uc.GetUserSubs(userID)

	t.Assert().NoError(err)
	t.Assert().Equal(subs, result)
}

func (s *FriendsTestSuite) TestGetUserFriends(t provider.T) {
	userID := uint64(1)
	friendIDs := []uint64{2, 3, 4, 5}
	friends := testutils.BuildUsersByIDs(s.userBuilder, friendIDs)

	s.fRepoMock.On("GetUserFriends", userID).Return(friendIDs, nil)
	s.uRepoMock.On("GetUsersByIDs", friendIDs).Return(friends, nil)

	result, err := s.uc.GetUserFriends(userID)

	t.Assert().NoError(err)
	t.Assert().Equal(friends, result)
}
