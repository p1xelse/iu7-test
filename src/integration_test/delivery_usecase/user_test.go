package deliveryusecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	integrationutils "timetracker/integration_test/integration_utils"
	"timetracker/internal/testutils"
	"timetracker/models"
	"timetracker/models/dto"

	userDelivery "timetracker/internal/User/delivery"
	userRepoMock "timetracker/internal/User/repository/mocks"
	userUC "timetracker/internal/User/usecase"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type UserTestIntegrationSuite struct {
	suite.Suite
	uc             userUC.UsecaseI
	router         *echo.Echo
	delivery       *userDelivery.Delivery
	userRepoMock   *userRepoMock.RepositoryI
	userDtoBuilder *integrationutils.ReqUpdateUserBuilder
	userBuilder    *testutils.UserBuilder
}

func getJsonDataUser(t provider.T, userDto dto.ReqUpdateUser) (userDtoJson []byte) {
	userDtoJson, err := json.Marshal(userDto)
	t.Require().NoError(err)

	return
}

func TestUserTestIntegrationSuite(t *testing.T) {
	suite.RunSuite(t, new(UserTestIntegrationSuite))
}

func (s *UserTestIntegrationSuite) BeforeEach(t provider.T) {
	s.userRepoMock = userRepoMock.NewRepositoryI(t)
	s.userRepoMock = userRepoMock.NewRepositoryI(t)
	s.uc = userUC.New(s.userRepoMock)
	s.userDtoBuilder = integrationutils.NewReqUpdateUserBuilder()
	s.userBuilder = testutils.NewUserBuilder()

	s.router = echo.New()
	s.delivery = &userDelivery.Delivery{
		UserUC: s.uc,
	}

	userDelivery.NewDelivery(s.router, s.uc, nil)
}

func (s *UserTestIntegrationSuite) TestUpdateUser(t provider.T) {
	userDto := s.userDtoBuilder.
		WithName("name").
		WithAbout("about").
		WithPassword("green").
		WithRole("admin").
		Build()

	userModel := userDto.ToModelUser()
	userModel.ID = integrationutils.DefaultUser.ID
	userDtoJson := getJsonDataUser(t, userDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(userDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.userRepoMock.On("GetUser", userModel.ID).Return(userModel, nil)
	s.userRepoMock.On("UpdateUser", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.UpdateUser(ctx))

	t.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *UserTestIntegrationSuite) TestGetUsers(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	users := make([]*models.User, 0, 10)
	err := faker.FakeData(&users)
	t.Assert().NoError(err)

	s.userRepoMock.On("GetUsers").Return(users, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}

	bodyArr := []map[string]interface{}{}
	for idx := range users {
		new := map[string]interface{}{
			"id":    users[idx].ID,
			"name":  users[idx].Name,
			"email": users[idx].Email,
			"about": users[idx].About,
			"role":  users[idx].Role,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetUsers(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *UserTestIntegrationSuite) TestGetUser(t provider.T) {
	userDto := s.userDtoBuilder.
		WithName("name").
		WithAbout("about").
		WithPassword("green").
		WithRole("admin").
		Build()

	userModel := userDto.ToModelUser()
	userModel.ID = integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/users/:user_id")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(fmt.Sprint(userModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":    userModel.ID,
			"name":  userDto.Name,
			"email": userDto.Email,
			"about": userDto.About,
			"role":  userDto.Role,
		},
	}

	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.userRepoMock.On("GetUser", userModel.ID).Return(userModel, nil)

	t.Require().NoError(s.delivery.GetUser(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *UserTestIntegrationSuite) TestGetMe(t provider.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/me")

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":    integrationutils.DefaultUser.ID,
			"name":  integrationutils.DefaultUser.Name,
			"email": integrationutils.DefaultUser.Email,
			"about": integrationutils.DefaultUser.About,
			"role":  integrationutils.DefaultUser.Role,
		},
	}

	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.userRepoMock.On("GetUser", integrationutils.DefaultUser.ID).Return(&integrationutils.DefaultUser, nil)

	t.Require().NoError(s.delivery.GetMe(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}
