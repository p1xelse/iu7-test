package postgres

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	entryRepo "timetracker/internal/Entry/repository"
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
	columns = []string{"id", "user_id", "project_id", "description", "time_start", "time_end"}
)

type EntryRepoTestSuite struct {
	suite.Suite
	db           *sql.DB
	gormDB       *gorm.DB
	mock         sqlmock.Sqlmock
	repo         entryRepo.RepositoryI
	entryBuilder *testutils.EntryBuilder
}

func TestEntryRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(EntryRepoTestSuite))
}

func (s *EntryRepoTestSuite) BeforeEach(t provider.T) {
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

	s.repo = NewEntryRepository(s.gormDB)
	s.entryBuilder = testutils.NewEntryBuilder()
}

func (s *EntryRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *EntryRepoTestSuite) TestCreateEntry(t provider.T) {
	entry := s.entryBuilder.
		WithUserID(1).
		WithProjectID(1).
		WithDescription("entrt").
		Build()
	entryPostgres := toPostgresEntry(&entry)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "entry" ("user_id","description","time_start","time_end","project_id") VALUES ($1,$2,$3,$4,$5)`)).
		WithArgs(entryPostgres.UserID, entryPostgres.Description, entryPostgres.TimeStart, entryPostgres.TimeEnd, entryPostgres.ProjectID).
		WillReturnRows(sqlmock.NewRows([]string{"project_id", "id"}).AddRow(*entryPostgres.ProjectID, startID))

	s.mock.ExpectCommit()

	err := s.repo.CreateEntry(&entry)
	t.Assert().NoError(err)
	t.Assert().Equal(startID, entry.ID)
}

func (s *EntryRepoTestSuite) TestUpdateEntry(t provider.T) {
	entry := s.entryBuilder.
		WithID(1).
		WithUserID(1).
		WithProjectID(1).
		WithDescription("entrt").
		Build()
	entryPostgres := toPostgresEntry(&entry)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "entry" SET "user_id"=$1,"project_id"=$2,"description"=$3 WHERE "id" = $4`)).
		WithArgs(*entryPostgres.UserID, *entryPostgres.ProjectID, entryPostgres.Description, entryPostgres.ID).WillReturnResult(sqlmock.NewResult(int64(startID), 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateEntry(&entry)
	t.Assert().NoError(err)
}

func (s *EntryRepoTestSuite) TestGetEntry(t provider.T) {
	entry := s.entryBuilder.
		WithID(1).
		WithUserID(1).
		WithProjectID(1).
		WithDescription("entrt").
		Build()
	entryPostgres := toPostgresEntry(&entry)

	rows := sqlmock.NewRows(columns).
		AddRow(
			entryPostgres.ID,
			*entryPostgres.UserID,
			*entryPostgres.ProjectID,
			entryPostgres.Description,
			entryPostgres.TimeStart,
			entryPostgres.TimeEnd,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "entry" WHERE id = $1 LIMIT 1`)).
		WithArgs(entryPostgres.ID).
		WillReturnRows(rows)

	resEntry, err := s.repo.GetEntry(entry.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(&entry, resEntry)
}

func (s *EntryRepoTestSuite) TestDeleteEntry(t provider.T) {
	entry := s.entryBuilder.
		WithID(1).
		WithUserID(1).
		WithProjectID(1).
		WithDescription("entrt").
		Build()
	entryPostgres := toPostgresEntry(&entry)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "entry" WHERE "entry"."id" = $1`)).
		WithArgs(entryPostgres.ID).WillReturnResult(sqlmock.NewResult(int64(entryPostgres.ID), 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteEntry(entryPostgres.ID)
	t.Assert().NoError(err)
}

func (s *EntryRepoTestSuite) TestGetUserEntries(t provider.T) {
	entriesPostgres := make([]*Entry, 10)
	err := faker.FakeData(&entriesPostgres)
	t.Assert().NoError(err)

	for idx := range entriesPostgres {
		entriesPostgres[idx].UserID = entriesPostgres[0].UserID
	}

	rows := sqlmock.NewRows(columns)

	for idx := range entriesPostgres {
		rows.AddRow(
			entriesPostgres[idx].ID,
			*entriesPostgres[idx].UserID,
			*entriesPostgres[idx].ProjectID,
			entriesPostgres[idx].Description,
			entriesPostgres[idx].TimeStart,
			entriesPostgres[idx].TimeEnd)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "entry" WHERE "entry"."user_id" = $1`)).
		WithArgs(*entriesPostgres[0].UserID).
		WillReturnRows(rows)

	resEntry, err := s.repo.GetUserEntries(*entriesPostgres[0].UserID)
	t.Assert().NoError(err)
	t.Assert().Equal(toModelEntries(entriesPostgres), resEntry)
}

func (s *EntryRepoTestSuite) TestGetUserEntriesForDay(t provider.T) {
	entriesPostgres := make([]*Entry, 10)
	err := faker.FakeData(&entriesPostgres)
	t.Assert().NoError(err)

	today := time.Now()
	for idx := range entriesPostgres {
		entriesPostgres[idx].UserID = entriesPostgres[0].UserID
		if idx < 5 {
			entriesPostgres[idx].TimeStart = today
			entriesPostgres[idx].TimeEnd = today
		} else {
			entriesPostgres[idx].TimeStart = today.Add(time.Hour * 24)
			entriesPostgres[idx].TimeEnd = today.Add(time.Hour * 24)
		}
	}

	rows := sqlmock.NewRows(columns)

	for idx := range entriesPostgres[:5] {
		rows.AddRow(
			entriesPostgres[idx].ID,
			*entriesPostgres[idx].UserID,
			*entriesPostgres[idx].ProjectID,
			entriesPostgres[idx].Description,
			entriesPostgres[idx].TimeStart,
			entriesPostgres[idx].TimeEnd)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "entry" WHERE "entry"."user_id" = $1 AND (time_start BETWEEN $2 AND $3)`)).
		WithArgs(*entriesPostgres[0].UserID,
			time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
			time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, today.Location()),
		).
		WillReturnRows(rows)

	resEntry, err := s.repo.GetUserEntriesForDay(*entriesPostgres[0].UserID, today)
	t.Assert().NoError(err)
	t.Assert().Equal(toModelEntries(entriesPostgres[:5]), resEntry)
}
