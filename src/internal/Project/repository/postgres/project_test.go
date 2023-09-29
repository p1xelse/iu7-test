package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	projectRepo "timetracker/internal/Project/repository"
	"timetracker/internal/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	startID = uint64(1)
)

var (
	columns = []string{"id", "user_id", "name", "about", "color", "is_private", "total_count_hours"}
)

type ProjectRepoTestSuite struct {
	suite.Suite
	db             *sql.DB
	gormDB         *gorm.DB
	mock           sqlmock.Sqlmock
	repo           projectRepo.RepositoryI
	projectBuilder *testutils.ProjectBuilder
}

func TestProjectRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(ProjectRepoTestSuite))
}

func (s *ProjectRepoTestSuite) BeforeEach(t provider.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error while creating sql mock")
	}
	s.db = db

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		t.Fatal("error while creating gorm db")
	}

	s.gormDB = gormDB
	s.mock = mock

	s.repo = NewProjectRepository(s.gormDB)
	s.projectBuilder = testutils.NewProjectBuilder()
}

func (s *ProjectRepoTestSuite) TearDownTest(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *ProjectRepoTestSuite) TestCreateProject(t provider.T) {
	project := s.projectBuilder.
		WithUserID(1).
		WithName("project").
		Build()
	projectPostgres := toPostgresProject(&project)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "project" ("user_id","name","about","color","is_private","total_count_hours") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			*projectPostgres.UserID,
			projectPostgres.Name,
			projectPostgres.About,
			projectPostgres.Color,
			projectPostgres.IsPrivate,
			projectPostgres.TotalCountHours,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(startID))

	s.mock.ExpectCommit()

	err := s.repo.CreateProject(&project)
	t.Assert().NoError(err)
	t.Assert().Equal(startID, project.ID)
}

func (s *ProjectRepoTestSuite) TestUpdateProject(t provider.T) {
	project := s.projectBuilder.
		WithID(1).
		WithUserID(1).
		WithName("project").
		WithAbout("about").
		Build()
	projectPostgres := toPostgresProject(&project)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "project" SET "user_id"=$1,"name"=$2,"about"=$3 WHERE "id" = $4`)).
		WithArgs(
			*projectPostgres.UserID,
			projectPostgres.Name,
			projectPostgres.About,
			projectPostgres.ID).
		WillReturnResult(sqlmock.NewResult(int64(startID), 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateProject(&project)
	t.Assert().NoError(err)
}

func (s *ProjectRepoTestSuite) TestGetProject(t provider.T) {
	project := s.projectBuilder.
		WithID(1).
		WithUserID(1).
		WithName("project").
		Build()
	projectPostgres := toPostgresProject(&project)

	rows := sqlmock.NewRows(columns).
		AddRow(
			projectPostgres.ID,
			*projectPostgres.UserID,
			projectPostgres.Name,
			projectPostgres.About,
			projectPostgres.Color,
			projectPostgres.IsPrivate,
			projectPostgres.TotalCountHours,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "project" WHERE id = $1 LIMIT 1`)).
		WithArgs(projectPostgres.ID).
		WillReturnRows(rows)

	resProject, err := s.repo.GetProject(project.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(&project, resProject)
}

func (s *ProjectRepoTestSuite) TestDeleteProject(t provider.T) {
	project := s.projectBuilder.
		WithID(1).
		WithUserID(1).
		WithName("project").
		Build()
	projectPostgres := toPostgresProject(&project)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "project" WHERE "project"."id" = $1`)).
		WithArgs(projectPostgres.ID).WillReturnResult(sqlmock.NewResult(int64(projectPostgres.ID), 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteProject(projectPostgres.ID)
	t.Assert().NoError(err)
}

func (s *ProjectRepoTestSuite) TestGetUserProjects(t provider.T) {
	projectsPostgres := make([]*Project, 10)
	err := faker.FakeData(&projectsPostgres)
	t.Assert().NoError(err)

	for idx := range projectsPostgres {
		projectsPostgres[idx].UserID = projectsPostgres[0].UserID
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "about", "color", "is_private", "total_count_hours",
	})

	for idx := range projectsPostgres {
		rows.AddRow(
			projectsPostgres[idx].ID,
			*projectsPostgres[idx].UserID,
			projectsPostgres[idx].Name,
			projectsPostgres[idx].About,
			projectsPostgres[idx].Color,
			projectsPostgres[idx].IsPrivate,
			projectsPostgres[idx].TotalCountHours,
		)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "project" WHERE "project"."user_id" = $1`)).
		WithArgs(*projectsPostgres[0].UserID).
		WillReturnRows(rows)

	resProject, err := s.repo.GetUserProjects(*projectsPostgres[0].UserID)
	t.Assert().NoError(err)
	t.Assert().Equal(toModelProjects(projectsPostgres), resProject)
}
