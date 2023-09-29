package usecase

import (
	"testing"

	projectRepoMock "timetracker/internal/Project/repository/mocks"
	"timetracker/internal/testutils"
	"timetracker/models"

	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ProjectTestSuite struct {
	suite.Suite
	uc             UsecaseI
	prRepoMock     *projectRepoMock.RepositoryI
	projectBuilder *testutils.ProjectBuilder
}

func TestProjectTestSuite(t *testing.T) {
	suite.RunSuite(t, new(ProjectTestSuite))
}

func (s *ProjectTestSuite) BeforeEach(t provider.T) {
	s.prRepoMock = projectRepoMock.NewRepositoryI(t)
	s.uc = New(s.prRepoMock, nil)
	s.projectBuilder = testutils.NewProjectBuilder()
}

func (s *ProjectTestSuite) TestCreateProject(t provider.T) {
	project := s.projectBuilder.WithID(1).WithName("project").Build()

	s.prRepoMock.On("CreateProject", &project).Return(nil)
	err := s.uc.CreateProject(&project)

	t.Assert().NoError(err)
}

func (s *ProjectTestSuite) TestUpdateProject(t provider.T) {
	project := s.projectBuilder.WithID(1).WithName("project").Build()
	notFoundProject := s.projectBuilder.WithID(0).Build()

	s.prRepoMock.On("GetProject", project.ID).Return(&project, nil)
	s.prRepoMock.On("UpdateProject", &project).Return(nil)
	s.prRepoMock.On("GetProject", notFoundProject.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		ArgData *models.Project
		Error   error
	}{
		"success": {
			ArgData: &project,
			Error:   nil,
		},
		"project not found": {
			ArgData: &notFoundProject,
			Error:   models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.UpdateProject(test.ArgData)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *ProjectTestSuite) TestGetProject(t provider.T) {
	project := s.projectBuilder.WithID(1).WithName("project").Build()

	s.prRepoMock.On("GetProject", project.ID).Return(&project, nil)
	result, err := s.uc.GetProject(project.ID)

	t.Assert().NoError(err)
	t.Assert().Equal(project, *result)
}

func (s *ProjectTestSuite) TestDeleteProject(t provider.T) {
	project := s.projectBuilder.WithID(1).WithName("project").WithUserID(2).Build()
	notFoundProject := s.projectBuilder.WithID(0).WithUserID(3).Build()

	s.prRepoMock.On("GetProject", project.ID).Return(&project, nil)
	s.prRepoMock.On("DeleteProject", project.ID).Return(nil)
	s.prRepoMock.On("GetProject", notFoundProject.ID).Return(nil, models.ErrNotFound)

	cases := map[string]struct {
		ProjectID uint64
		UserID    uint64
		Error     error
	}{
		"success": {
			ProjectID: project.ID,
			UserID:    *project.UserID,
			Error:     nil,
		},
		"project not found": {
			ProjectID: notFoundProject.ID,
			UserID:    *notFoundProject.UserID,
			Error:     models.ErrNotFound,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			err := s.uc.DeleteProject(test.ProjectID, test.UserID)
			t.Assert().ErrorIs(err, test.Error)
		})
	}
}

func (s *ProjectTestSuite) TestUsecaseGetUserProjects(t provider.T) {
	projects := make([]*models.Project, 0, 10)
	err := faker.FakeData(&projects)

	for idx := range projects {
		projects[idx].UserID = projects[0].UserID
	}
	t.Assert().NoError(err)
	s.prRepoMock.On("GetUserProjects", *projects[0].UserID).Return(projects, nil)

	cases := map[string]struct {
		UserID   uint64
		Projects []*models.Project
		Error    error
	}{
		"success": {
			UserID:   *projects[0].UserID,
			Projects: projects,
			Error:    nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t provider.T) {
			resProjects, err := s.uc.GetUserProjects(test.UserID)
			t.Assert().ErrorIs(err, test.Error)
			t.Assert().Equal(projects, resProjects)
		})
	}
}
