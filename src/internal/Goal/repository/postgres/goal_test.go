package postgres

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	goalRepo "timetracker/internal/Goal/repository"
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

type GoalRepoTestSuite struct {
	suite.Suite
	db          *sql.DB
	gormDB      *gorm.DB
	mock        sqlmock.Sqlmock
	repo        goalRepo.RepositoryI
	goalBuilder *testutils.GoalBuilder
}

func TestGoalRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(GoalRepoTestSuite))
}

func (s *GoalRepoTestSuite) BeforeEach(t provider.T) {
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

	s.repo = NewGoalRepository(s.gormDB)
	s.goalBuilder = testutils.NewGoalBuilder()
}

func (s *GoalRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *GoalRepoTestSuite) TestCreateGoal(t provider.T) {
	goal := s.goalBuilder.
		WithUserID(1).
		WithProjectID(1).
		WithName("goal").
		WithDescription("goal").
		WithTimeStart(time.Now()).
		WithTimeEnd(time.Now()).
		Build()
	goalPostgres := toPostgresGoal(&goal)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "goal" ("user_id","name","project_id","description","time_start","time_end","hours_count") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(
			*goalPostgres.UserID,
			goalPostgres.Name,
			*goalPostgres.ProjectID,
			goalPostgres.Description,
			goalPostgres.TimeStart,
			goalPostgres.TimeEnd,
			goalPostgres.HoursCount,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(startID))

	s.mock.ExpectCommit()

	err := s.repo.CreateGoal(&goal)
	t.Assert().NoError(err)
	t.Assert().Equal(startID, goal.ID)
}

func (s *GoalRepoTestSuite) TestUpdateGoal(t provider.T) {
	goal := s.goalBuilder.
		WithID(1).
		WithUserID(1).
		WithProjectID(1).
		WithDescription("goal").
		WithName("goal").
		Build()
	goalPostgres := toPostgresGoal(&goal)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "goal" SET "user_id"=$1,"name"=$2,"project_id"=$3,"description"=$4 WHERE "id" = $5`)).
		WithArgs(*goalPostgres.UserID, goalPostgres.Name, *goalPostgres.ProjectID, goalPostgres.Description, goalPostgres.ID).
		WillReturnResult(sqlmock.NewResult(int64(startID), 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateGoal(&goal)
	t.Assert().NoError(err)
}

func (s *GoalRepoTestSuite) TestGetGoal(t provider.T) {
	goal := s.goalBuilder.
		WithID(1).
		WithUserID(1).
		WithProjectID(1).
		WithDescription("goal").
		Build()
	goalPostgres := toPostgresGoal(&goal)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "project_id", "description", "time_start", "time_end", "hours_count",
	}).
		AddRow(
			goalPostgres.ID,
			*goalPostgres.UserID,
			goalPostgres.Name,
			*goalPostgres.ProjectID,
			goalPostgres.Description,
			goalPostgres.TimeStart,
			goalPostgres.TimeEnd,
			goalPostgres.HoursCount,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "goal" WHERE id = $1 LIMIT 1`)).
		WithArgs(goalPostgres.ID).
		WillReturnRows(rows)

	resGoal, err := s.repo.GetGoal(goal.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(&goal, resGoal)
}

func (s *GoalRepoTestSuite) TestDeleteGoal(t provider.T) {
	goal := s.goalBuilder.
		WithID(1).
		WithUserID(1).
		WithProjectID(1).
		WithDescription("goal").
		Build()
	goalPostgres := toPostgresGoal(&goal)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "goal" WHERE "goal"."id" = $1`)).
		WithArgs(goalPostgres.ID).WillReturnResult(sqlmock.NewResult(int64(goalPostgres.ID), 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteGoal(goalPostgres.ID)
	t.Assert().NoError(err)
}

func (s *GoalRepoTestSuite) TestGetUserGoals(t provider.T) {
	goalsPostgres := make([]*Goal, 10)
	err := faker.FakeData(&goalsPostgres)
	t.Assert().NoError(err)

	for idx := range goalsPostgres {
		goalsPostgres[idx].UserID = goalsPostgres[0].UserID
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "project_id", "description", "time_start", "time_end", "hours_count",
	})

	for idx := range goalsPostgres {
		rows.AddRow(
			goalsPostgres[idx].ID,
			*goalsPostgres[idx].UserID,
			goalsPostgres[idx].Name,
			*goalsPostgres[idx].ProjectID,
			goalsPostgres[idx].Description,
			goalsPostgres[idx].TimeStart,
			goalsPostgres[idx].TimeEnd,
			goalsPostgres[idx].HoursCount,
		)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "goal" WHERE "goal"."user_id" = $1`)).
		WithArgs(*goalsPostgres[0].UserID).
		WillReturnRows(rows)

	resGoal, err := s.repo.GetUserGoals(*goalsPostgres[0].UserID)
	t.Assert().NoError(err)
	t.Assert().Equal(toModelGoals(goalsPostgres), resGoal)
}
