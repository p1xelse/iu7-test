package deliveryusecase

import (
	"encoding/json"
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
	userBuilder    *testutils.UserBuilder
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
	s.userBuilder = testutils.NewUserBuilder()

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

func (s *FriendTestIntegrationSuite) TestUnsubscribe(t provider.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/friends/unsubscribe/:user_id")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues("2")

	friendRel := s.fRelBuilder.WithSubID(1).WithUserID(2).Build()
	s.friendRepoMock.On("CheckFriends", mock.MatchedBy(func(arg *models.FriendRelation) bool {
		return *arg.SubscriberID == *friendRel.SubscriberID && *arg.UserID == *friendRel.UserID
	})).Return(true, nil)

	s.friendRepoMock.On("DeleteFriendRelation", mock.MatchedBy(func(arg *models.FriendRelation) bool {
		return *arg.SubscriberID == *friendRel.SubscriberID && *arg.UserID == *friendRel.UserID
	})).Return(nil)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)
	t.Require().NoError(s.delivery.Unsubscribe(ctx))
	t.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *FriendTestIntegrationSuite) TestGetUserSubs(t provider.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/user/:user_id/subs")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues("1")

	userID := uint64(1)
	subsIDs := []uint64{2, 3, 4, 5}
	subs := testutils.BuildUsersByIDs(s.userBuilder, subsIDs)

	s.friendRepoMock.On("GetUserSubs", userID).Return(subsIDs, nil)
	s.uRepoMock.On("GetUsersByIDs", subsIDs).Return(subs, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range subs {
		new := map[string]interface{}{
			"id":    subs[idx].ID,
			"name":  subs[idx].Name,
			"email": subs[idx].Email,
			"about": subs[idx].About,
			"role":  subs[idx].Role,
		}
		bodyArr = append(bodyArr, new)
	}

	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)
	t.Require().NoError(s.delivery.GetUserSubs(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *FriendTestIntegrationSuite) TestGetUserFriends(t provider.T) {

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/user/:user_id/friends")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues("1")

	userID := uint64(1)
	friendIDs := []uint64{2, 3, 4, 5}
	friends := testutils.BuildUsersByIDs(s.userBuilder, friendIDs)

	s.friendRepoMock.On("GetUserFriends", userID).Return(friendIDs, nil)
	s.uRepoMock.On("GetUsersByIDs", friendIDs).Return(friends, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range friends {
		new := map[string]interface{}{
			"id":    friends[idx].ID,
			"name":  friends[idx].Name,
			"email": friends[idx].Email,
			"about": friends[idx].About,
			"role":  friends[idx].Role,
		}
		bodyArr = append(bodyArr, new)
	}

	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)
	t.Require().NoError(s.delivery.GetUserFriends(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}
