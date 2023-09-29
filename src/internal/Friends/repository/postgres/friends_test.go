package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	fRelRepo "timetracker/internal/Friends/repository"
	"timetracker/internal/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	startID = uint64(1)
)

var (
	columns = []string{"subscriber_id", "user_id"}
)

type FriendRelationRepoTestSuite struct {
	suite.Suite
	db          *sql.DB
	gormDB      *gorm.DB
	mock        sqlmock.Sqlmock
	repo        fRelRepo.RepositoryI
	fRelBuilder *testutils.FriendRelationBuilder
}

func TestFriendRelationRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(FriendRelationRepoTestSuite))
}

func (s *FriendRelationRepoTestSuite) BeforeEach(t provider.T) {
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

	s.repo = NewFriendRepository(s.gormDB)
	s.fRelBuilder = testutils.NewFriendRelationBuilder()
}

func (s *FriendRelationRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *FriendRelationRepoTestSuite) TestCreateFriendRelation(t provider.T) {
	fRel := s.fRelBuilder.
		WithUserID(1).
		WithSubID(2).
		Build()
	fRelPostgres := toPostgresFriendRelation(&fRel)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "friend_relation" ("subscriber_id","user_id") VALUES ($1,$2)`)).
		WithArgs(
			*fRelPostgres.SubscriberID,
			*fRelPostgres.UserID,
		).WillReturnResult(sqlmock.NewResult(int64(startID), 1))

	s.mock.ExpectCommit()

	err := s.repo.CreateFriendRelation(&fRel)
	t.Assert().NoError(err)
}

func (s *FriendRelationRepoTestSuite) TestCheckFriends(t provider.T) {
	fRel := s.fRelBuilder.
		WithUserID(1).
		WithSubID(2).
		Build()
	fRelPostgres := toPostgresFriendRelation(&fRel)

	rows := sqlmock.NewRows(columns).
		AddRow(
			fRelPostgres.SubscriberID,
			fRelPostgres.UserID,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "friend_relation" WHERE "friend_relation"."subscriber_id" = $1 AND "friend_relation"."user_id" = $2 LIMIT 1`)).
		WithArgs(
			fRelPostgres.SubscriberID,
			fRelPostgres.UserID,
		).
		WillReturnRows(rows)

	isFriends, err := s.repo.CheckFriends(&fRel)
	t.Assert().NoError(err)
	t.Assert().Equal(true, isFriends)
}

func (s *FriendRelationRepoTestSuite) TestDeleteFriendRelation(t provider.T) {
	fRel := s.fRelBuilder.
		WithSubID(2).
		WithUserID(1).
		Build()
	fRelPostgres := toPostgresFriendRelation(&fRel)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "friend_relation" WHERE "friend_relation"."subscriber_id" = $1 AND "friend_relation"."user_id" = $2`)).
		WithArgs(*fRelPostgres.SubscriberID, *fRelPostgres.UserID).
		WillReturnResult(sqlmock.NewResult(int64(1), 1))
	s.mock.ExpectCommit()

	err := s.repo.DeleteFriendRelation(&fRel)
	t.Assert().NoError(err)
}

func (s *FriendRelationRepoTestSuite) TestGetUserSubs(t provider.T) {
	fRel := s.fRelBuilder.
		WithUserID(1).
		WithSubID(2).
		Build()
	fRelPostgres := toPostgresFriendRelation(&fRel)
	subs := []uint64{*fRel.SubscriberID}

	rows := sqlmock.NewRows([]string{"subscriber_id"}).
		AddRow(
			fRelPostgres.SubscriberID,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT f1.subscriber_id FROM friend_relation f1 left join friend_relation f2 on f2.user_id = f1.subscriber_id and f2.subscriber_id = f1.user_id WHERE f1.user_id = $1 and f2.user_id is null`)).
		WithArgs(*fRelPostgres.UserID).
		WillReturnRows(rows)

	resSubs, err := s.repo.GetUserSubs(*fRelPostgres.UserID)
	t.Assert().NoError(err)
	t.Assert().Equal(subs, resSubs)
}

func (s *FriendRelationRepoTestSuite) TestGetUserFriends(t provider.T) {
	fRel := s.fRelBuilder.
		WithUserID(1).
		WithSubID(2).
		Build()
	fRelPostgres := toPostgresFriendRelation(&fRel)
	subs := []uint64{*fRel.SubscriberID}

	rows := sqlmock.NewRows([]string{"subscriber_id"}).
		AddRow(
			fRelPostgres.SubscriberID,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT f1.subscriber_id FROM friend_relation f1 join friend_relation f2 on f2.user_id = f1.subscriber_id and f2.subscriber_id = f1.user_id WHERE f1.user_id = $1`)).
		WithArgs(*fRelPostgres.UserID).
		WillReturnRows(rows)

	resSubs, err := s.repo.GetUserFriends(*fRelPostgres.UserID)
	t.Assert().NoError(err)
	t.Assert().Equal(subs, resSubs)
}
