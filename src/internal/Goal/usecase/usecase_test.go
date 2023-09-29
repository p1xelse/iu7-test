package usecase

import (
	"testing"

	goalRepoMock "timetracker/internal/Goal/repository/mocks"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type GoalTestSuite struct {
	suite.Suite
	uc           UsecaseI
	goalRepoMock *goalRepoMock.RepositoryI
	goalBuilder  *testutils.GoalBuilder
}

func TestGoalTestSuite(t *testing.T) {
	suite.RunSuite(t, new(GoalTestSuite))
}

func (s *GoalTestSuite) BeforeEach(t provider.T) {
	s.goalRepoMock = goalRepoMock.NewRepositoryI(t)
	s.uc = New(s.goalRepoMock)
	s.goalBuilder = testutils.NewGoalBuilder()
}

func (s *GoalTestSuite) TestCreateGoal(t provider.T) {
	goal := s.goalBuilder.WithID(1).WithName("goal").Build()

	s.goalRepoMock.On("CreateGoal", &goal).Return(nil)
	err := s.uc.CreateGoal(&goal)

	t.Assert().NoError(err)
}

func (s *GoalTestSuite) TestUpdateGoal(t provider.T) {
	goal := s.goalBuilder.WithID(1).WithName("goal").Build()
	notFoundGoal := s.goalBuilder.WithID(0).Build()

	s.goalRepoMock.On("GetGoal", goal.ID).Return(&goal, nil)
	s.goalRepoMock.On("UpdateGoal", &goal).Return(nil)
	s.goalRepoMock.On("GetGoal", notFoundGoal.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		ArgData *models.Goal
		Error   error
	}{
		"success": {
			ArgData: &goal,
			Error:   nil,
		},
		"goal not found": {
			ArgData: &notFoundGoal,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.UpdateGoal(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *GoalTestSuite) TestGetGoal(t provider.T) {
	goal := s.goalBuilder.WithID(1).WithName("goal").Build()

	s.goalRepoMock.On("GetGoal", goal.ID).Return(&goal, nil)
	result, err := s.uc.GetGoal(goal.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(goal, *result)
}

func (s *GoalTestSuite) TestDeleteGoal(t provider.T) {
	goal := s.goalBuilder.WithID(1).WithName("goal").WithUserID(2).Build()
	notFoundGoal := s.goalBuilder.WithID(0).WithUserID(3).Build()

	s.goalRepoMock.On("GetGoal", goal.ID).Return(&goal, nil)
	s.goalRepoMock.On("DeleteGoal", goal.ID).Return(nil)
	s.goalRepoMock.On("GetGoal", notFoundGoal.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		GoalID uint64
		UserID uint64
		Error  error
	}{
		"success": {
			GoalID: goal.ID,
			UserID: *goal.UserID,
			Error:  nil,
		},
		"goal not found": {
			GoalID: notFoundGoal.ID,
			UserID: *notFoundGoal.UserID,
			Error:  models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.DeleteGoal(test.GoalID, test.UserID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *GoalTestSuite) TestUsecaseGetUserGoals(t provider.T) {
	goals := make([]*models.Goal, 0, 10)
	err := faker.FakeData(&goals)

	for idx := range goals {
		goals[idx].UserID = goals[0].UserID
	}
	t.Assert().NoError(err)
	s.goalRepoMock.On("GetUserGoals", *goals[0].UserID).Return(goals, nil)

	cases := map[string]struct {
		UserID uint64
		Goals  []*models.Goal
		Error  error
	}{
		"success": {
			UserID: *goals[0].UserID,
			Goals:  goals,
			Error:  nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resGoals, err := s.uc.GetUserGoals(test.UserID)
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(goals, resGoals)
		})
	}
}
