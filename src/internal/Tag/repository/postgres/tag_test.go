package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	tagRepo "timetracker/internal/Tag/repository"
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
	columns = []string{"id", "user_id", "name", "about", "color"}
)

type TagRepoTestSuite struct {
	suite.Suite
	db         *sql.DB
	gormDB     *gorm.DB
	mock       sqlmock.Sqlmock
	repo       tagRepo.RepositoryI
	tagBuilder *testutils.TagBuilder
}

func TestTagRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(TagRepoTestSuite))
}

func (s *TagRepoTestSuite) BeforeEach(t provider.T) {
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

	s.repo = NewTagRepository(s.gormDB)
	s.tagBuilder = testutils.NewTagBuilder()
}

func (s *TagRepoTestSuite) TearDownTest(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *TagRepoTestSuite) TestCreateTag(t provider.T) {
	tag := s.tagBuilder.
		WithUserID(1).
		WithName("tag").
		WithAbout("about").
		WithColor("green").
		Build()
	tagPostgres := toPostgresTag(&tag)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "tag" ("user_id","name","about","color") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(
			tagPostgres.UserID,
			tagPostgres.Name,
			tagPostgres.About,
			tagPostgres.Color,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(startID))

	s.mock.ExpectCommit()

	err := s.repo.CreateTag(&tag)
	t.Assert().NoError(err)
	t.Assert().Equal(startID, tag.ID)
}

func (s *TagRepoTestSuite) TestUpdateTag(t provider.T) {
	tag := s.tagBuilder.
		WithID(1).
		WithUserID(1).
		WithName("tag").
		WithAbout("about").
		WithColor("green").
		Build()
	tagPostgres := toPostgresTag(&tag)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "tag" SET "user_id"=$1,"name"=$2,"about"=$3,"color"=$4 WHERE "id" = $5`)).
		WithArgs(
			tagPostgres.UserID,
			tagPostgres.Name,
			tagPostgres.About,
			tagPostgres.Color,
			tagPostgres.ID).
		WillReturnResult(sqlmock.NewResult(int64(tagPostgres.ID), 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateTag(&tag)
	t.Assert().NoError(err)
}

func (s *TagRepoTestSuite) TestGetTag(t provider.T) {
	tag := s.tagBuilder.
		WithID(1).
		WithUserID(1).
		WithName("tag").
		WithAbout("about").
		WithColor("green").
		Build()
	tagPostgres := toPostgresTag(&tag)

	rows := sqlmock.NewRows(columns).
		AddRow(
			tagPostgres.ID,
			tagPostgres.UserID,
			tagPostgres.Name,
			tagPostgres.About,
			tagPostgres.Color,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag" WHERE id = $1 LIMIT 1`)).
		WithArgs(tagPostgres.ID).
		WillReturnRows(rows)

	resTag, err := s.repo.GetTag(tag.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(&tag, resTag)
}

func (s *TagRepoTestSuite) TestDeleteTag(t provider.T) {
	tag := s.tagBuilder.
		WithID(1).
		WithUserID(1).
		WithName("tag").
		WithAbout("about").
		WithColor("green").
		Build()
	tagPostgres := toPostgresTag(&tag)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tag" WHERE "tag"."id" = $1`)).
		WithArgs(tagPostgres.ID).WillReturnResult(sqlmock.NewResult(int64(tagPostgres.ID), 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteTag(tagPostgres.ID)
	t.Assert().NoError(err)
}

func (s *TagRepoTestSuite) TestGetUserTags(t provider.T) {
	tagsPostgres := make([]*Tag, 10)
	err := faker.FakeData(&tagsPostgres)
	t.Assert().NoError(err)

	for idx := range tagsPostgres {
		tagsPostgres[idx].UserID = tagsPostgres[0].UserID
	}

	rows := sqlmock.NewRows(columns)

	for idx := range tagsPostgres {
		rows.AddRow(
			tagsPostgres[idx].ID,
			tagsPostgres[idx].UserID,
			tagsPostgres[idx].Name,
			tagsPostgres[idx].About,
			tagsPostgres[idx].Color,
		)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tag" WHERE "tag"."user_id" = $1`)).
		WithArgs(tagsPostgres[0].UserID).
		WillReturnRows(rows)

	resTag, err := s.repo.GetUserTags(tagsPostgres[0].UserID)
	t.Assert().NoError(err)
	t.Assert().Equal(toModelTags(tagsPostgres), resTag)
}
