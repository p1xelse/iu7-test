package deliveryusecase

import (
	"net/http"
	"net/http/httptest"
	"testing"
	integrationutils "timetracker/integration_test/integration_utils"
	"timetracker/models"

	friendDelivery "timetracker/internal/Friends/delivery"
	friendRepoMock "timetracker/internal/Friends/repository/mocks"
	friendUC "timetracker/internal/Friends/usecase"
	userRepoMock "timetracker/internal/User/repository/mocks"
	"timetracker/internal/testutils"

	"github.com/labstack/echo/v4"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type FriendTestIntegrationSuite struct {
	suite.Suite
	uc             friendUC.UsecaseI
	router         *echo.Echo
	delivery       *friendDelivery.Delivery
	friendRepoMock *friendRepoMock.RepositoryI
	fRelBuilder    *testutils.FriendRelationBuilder
	uRepoMock      *userRepoMock.RepositoryI
}

func TestFriendTestIntegrationSuite(t *testing.T) {
	suite.RunSuite(t, new(FriendTestIntegrationSuite))
}

func (s *FriendTestIntegrationSuite) BeforeEach(t provider.T) {
	s.friendRepoMock = friendRepoMock.NewRepositoryI(t)
	s.uRepoMock = userRepoMock.NewRepositoryI(t)
	s.uc = friendUC.New(s.friendRepoMock, s.uRepoMock)
	s.fRelBuilder = testutils.NewFriendRelationBuilder()

	s.router = echo.New()
	s.delivery = &friendDelivery.Delivery{
		FriendsUC: s.uc,
	}

	friendDelivery.NewDelivery(s.router, s.uc, nil)
}

func (s *FriendTestIntegrationSuite) TestSubscribe(t provider.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/friends/subscribe/:user_id")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues("2")

	friendRel := s.fRelBuilder.WithSubID(1).WithUserID(2).Build()
	s.friendRepoMock.On("CheckFriends", mock.MatchedBy(func(arg *models.FriendRelation) bool {
		return *arg.SubscriberID == *friendRel.SubscriberID && *arg.UserID == *friendRel.UserID
	})).Return(false, nil)

	s.friendRepoMock.On("CreateFriendRelation", mock.MatchedBy(func(arg *models.FriendRelation) bool {
		return *arg.SubscriberID == *friendRel.SubscriberID && *arg.UserID == *friendRel.UserID
	})).Return(nil)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)
	t.Require().NoError(s.delivery.Subscribe(ctx))
	t.Assert().Equal(http.StatusCreated, rec.Code)
}
