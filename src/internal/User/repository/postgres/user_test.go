package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	userRepo "timetracker/internal/User/repository"
	"timetracker/internal/testutils"
	"timetracker/models"

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
	columns = []string{"id", "name", "email", "about", "role", "password"}
)

type UserRepoTestSuite struct {
	suite.Suite
	db          *sql.DB
	gormDB      *gorm.DB
	mock        sqlmock.Sqlmock
	repo        userRepo.RepositoryI
	userBuilder *testutils.UserBuilder
}

func TestUserRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(UserRepoTestSuite))
}

func (s *UserRepoTestSuite) BeforeEach(t provider.T) {
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

	s.repo = NewUserRepository(s.gormDB)
	s.userBuilder = testutils.NewUserBuilder()
}

func (s *UserRepoTestSuite) TearDownTest(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *UserRepoTestSuite) TestCreateUser(t provider.T) {
	user := s.userBuilder.
		WithName("name").
		WithEmail("email").
		WithAbout("about").
		WithRole("admin").
		WithPassword("password").
		Build()
	userPostgres := toPostgresUser(&user)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("name","email","about","role","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			userPostgres.Name,
			userPostgres.Email,
			userPostgres.About,
			userPostgres.Role,
			userPostgres.Password,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(startID))

	s.mock.ExpectCommit()

	err := s.repo.CreateUser(&user)
	t.Assert().NoError(err)
	t.Assert().Equal(startID, user.ID)
}

func (s *UserRepoTestSuite) TestUpdateUser(t provider.T) {
	user := s.userBuilder.
		WithID(1).
		WithName("name").
		WithEmail("email").
		WithAbout("about").
		WithRole("admin").
		WithPassword("password").
		Build()
	userPostgres := toPostgresUser(&user)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "name"=$1,"email"=$2,"about"=$3,"role"=$4,"password"=$5 WHERE "id" = $6`)).
		WithArgs(
			userPostgres.Name,
			userPostgres.Email,
			userPostgres.About,
			userPostgres.Role,
			userPostgres.Password,
			userPostgres.ID).
		WillReturnResult(sqlmock.NewResult(int64(userPostgres.ID), 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateUser(&user)
	t.Assert().NoError(err)
}

func (s *UserRepoTestSuite) TestGetUser(t provider.T) {
	user := s.userBuilder.
		WithID(1).
		WithName("name").
		WithEmail("email").
		WithAbout("about").
		WithRole("admin").
		WithPassword("password").
		Build()
	userPostgres := toPostgresUser(&user)

	rows := sqlmock.NewRows(columns).
		AddRow(
			userPostgres.ID,
			userPostgres.Name,
			userPostgres.Email,
			userPostgres.About,
			userPostgres.Role,
			userPostgres.Password,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1 LIMIT 1`)).
		WithArgs(userPostgres.ID).
		WillReturnRows(rows)

	resUser, err := s.repo.GetUser(user.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(&user, resUser)
}

func (s *UserRepoTestSuite) TestGetUserByEmail(t provider.T) {
	user := s.userBuilder.
		WithID(1).
		WithName("name").
		WithEmail("email").
		WithAbout("about").
		WithRole("admin").
		WithPassword("password").
		Build()
	userPostgres := toPostgresUser(&user)

	rows := sqlmock.NewRows(columns).
		AddRow(
			userPostgres.ID,
			userPostgres.Name,
			userPostgres.Email,
			userPostgres.About,
			userPostgres.Role,
			userPostgres.Password,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."email" = $1 LIMIT 1`)).
		WithArgs(userPostgres.Email).
		WillReturnRows(rows)

	resUser, err := s.repo.GetUserByEmail(user.Email)
	t.Assert().NoError(err)
	t.Assert().Equal(&user, resUser)
}

func (s *UserRepoTestSuite) equalSlicePointers(expected []*models.User, actual []*models.User, t provider.T) {
	t.Assert().Equal(len(expected), len(actual))
	for idx := range expected {
		t.Assert().EqualValues(expected[idx], actual[idx])
	}
}
func (s *UserRepoTestSuite) TestGetUsers(t provider.T) {
	usersPostgres := make([]*User, 10)
	err := faker.FakeData(&usersPostgres)
	t.Assert().NoError(err)

	for idx := range usersPostgres {
		usersPostgres[idx].Password = ""
	}

	//withput password
	rows := sqlmock.NewRows([]string{"id", "name", "email", "about", "role"})

	for idx := range usersPostgres {
		rows.AddRow(
			usersPostgres[idx].ID,
			usersPostgres[idx].Name,
			usersPostgres[idx].Email,
			usersPostgres[idx].About,
			usersPostgres[idx].Role,
		)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "users"."id","users"."name","users"."email","users"."about","users"."role" FROM "users"`)).
		WillReturnRows(rows)

	resUsers, err := s.repo.GetUsers()
	t.Assert().NoError(err)
	s.equalSlicePointers(resUsers, toModelUsers(usersPostgres), t)
}
func (s *UserRepoTestSuite) TestGetUsersByIDs(t provider.T) {
	usersPostgres := make([]*User, 10)
	err := faker.FakeData(&usersPostgres)
	usersPostgres = usersPostgres[:2]
	t.Assert().NoError(err)
	ids := make([]uint64, len(usersPostgres))

	for idx := range usersPostgres {
		usersPostgres[idx].Password = ""
		ids[idx] = usersPostgres[idx].ID
	}

	//without password
	rows := sqlmock.NewRows([]string{"id", "name", "email", "about", "role"})

	for idx := range usersPostgres {
		rows.AddRow(
			usersPostgres[idx].ID,
			usersPostgres[idx].Name,
			usersPostgres[idx].Email,
			usersPostgres[idx].About,
			usersPostgres[idx].Role,
		)
	}

	// without ids in `IN`
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "users"."id","users"."name","users"."email","users"."about","users"."role" FROM "users" WHERE "users"."id" IN`)).
		WillReturnRows(rows)

	resUsers, err := s.repo.GetUsersByIDs(ids)
	t.Assert().NoError(err)
	s.equalSlicePointers(resUsers, toModelUsers(usersPostgres), t)
}
