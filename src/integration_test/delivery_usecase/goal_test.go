package deliveryusecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	integrationutils "timetracker/integration_test/integration_utils"
	"timetracker/internal/testutils"
	"timetracker/models"
	"timetracker/models/dto"

	goalDelivery "timetracker/internal/Goal/delivery"
	goalRepoMock "timetracker/internal/Goal/repository/mocks"
	goalUC "timetracker/internal/Goal/usecase"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type GoalTestIntegrationSuite struct {
	suite.Suite
	uc             goalUC.UsecaseI
	router         *echo.Echo
	delivery       *goalDelivery.Delivery
	goalRepoMock   *goalRepoMock.RepositoryI
	goalDtoBuilder *integrationutils.ReqCreateUpdateGoalBuilder
	goalBuilder    *testutils.GoalBuilder
}

func getJsonDataGoal(t provider.T, goalDto dto.ReqCreateUpdateGoal) (goalDtoJson []byte) {
	goalDtoJson, err := json.Marshal(goalDto)
	t.Require().NoError(err)

	return
}

func TestGoalTestIntegrationSuite(t *testing.T) {
	suite.RunSuite(t, new(GoalTestIntegrationSuite))
}

func (s *GoalTestIntegrationSuite) BeforeEach(t provider.T) {
	s.goalRepoMock = goalRepoMock.NewRepositoryI(t)
	s.goalRepoMock = goalRepoMock.NewRepositoryI(t)
	s.uc = goalUC.New(s.goalRepoMock)
	s.goalDtoBuilder = integrationutils.NewReqCreateUpdateGoalBuilder()
	s.goalBuilder = testutils.NewGoalBuilder()

	s.router = echo.New()
	s.delivery = &goalDelivery.Delivery{
		GoalUC: s.uc,
	}

	goalDelivery.NewDelivery(s.router, s.uc, nil)
}

func (s *GoalTestIntegrationSuite) TestCreateGoal(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	goalDto := s.goalDtoBuilder.
		WithName("name").
		WithDescription("desc").
		WithHoursCount(20).
		WithTimeStart(startTime).
		WithTimeEnd(startTime).
		WithProjectID(1).
		Build()

	goalDtoJson := getJsonDataGoal(t, goalDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(goalDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":          0,
			"user_id":     1,
			"name":        goalDto.Name,
			"project_id":  *goalDto.ProjectID,
			"description": goalDto.Description,
			"hours_count": goalDto.HoursCount,
			"time_start":  goalDto.TimeStart,
			"time_end":    goalDto.TimeEnd,
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	s.goalRepoMock.On("CreateGoal", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.CreateGoal(ctx))

	t.Assert().Equal(http.StatusCreated, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *GoalTestIntegrationSuite) TestUpdateGoal(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	goalDto := s.goalDtoBuilder.
		WithName("name").
		WithDescription("desc").
		WithHoursCount(20).
		WithTimeStart(startTime).
		WithTimeEnd(startTime).
		WithProjectID(1).
		Build()

	goalModel := goalDto.ToModelGoal()
	goalModel.UserID = &integrationutils.DefaultUser.ID
	goalDtoJson := getJsonDataGoal(t, goalDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(goalDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":          0,
			"user_id":     1,
			"name":        goalDto.Name,
			"project_id":  *goalDto.ProjectID,
			"description": goalDto.Description,
			"hours_count": goalDto.HoursCount,
			"time_start":  goalDto.TimeStart,
			"time_end":    goalDto.TimeEnd,
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.goalRepoMock.On("GetGoal", goalModel.ID).Return(goalModel, nil)
	s.goalRepoMock.On("UpdateGoal", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.UpdateGoal(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *GoalTestIntegrationSuite) TestDeleteGoal(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	goalDto := s.goalDtoBuilder.
		WithName("name").
		WithDescription("desc").
		WithHoursCount(20).
		WithTimeStart(startTime).
		WithTimeEnd(startTime).
		WithProjectID(1).
		Build()

	goalModel := goalDto.ToModelGoal()
	goalModel.UserID = &integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/goal/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(goalModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	s.goalRepoMock.On("GetGoal", goalModel.ID).Return(goalModel, nil)
	s.goalRepoMock.On("DeleteGoal", goalModel.ID).Return(nil)

	t.Require().NoError(s.delivery.DeleteGoal(ctx))
	t.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *GoalTestIntegrationSuite) TestGetGoal(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	goalDto := s.goalDtoBuilder.
		WithName("name").
		WithDescription("desc").
		WithHoursCount(20).
		WithTimeStart(startTime).
		WithTimeEnd(startTime).
		WithProjectID(1).
		Build()

	goalModel := goalDto.ToModelGoal()
	goalModel.UserID = &integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/goal/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(goalModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":          0,
			"user_id":     1,
			"name":        goalDto.Name,
			"project_id":  *goalDto.ProjectID,
			"description": goalDto.Description,
			"hours_count": goalDto.HoursCount,
			"time_start":  goalDto.TimeStart,
			"time_end":    goalDto.TimeEnd,
		},
	}

	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.goalRepoMock.On("GetGoal", goalModel.ID).Return(goalModel, nil)

	t.Require().NoError(s.delivery.GetGoal(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *GoalTestIntegrationSuite) TestGetMyGoal(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	goals := make([]*models.Goal, 0, 10)
	err := faker.FakeData(&goals)
	t.Assert().NoError(err)

	for idx := range goals {
		goals[idx].UserID = &integrationutils.DefaultUser.ID
	}

	s.goalRepoMock.On("GetUserGoals", *goals[0].UserID).Return(goals, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range goals {
		new := map[string]interface{}{
			"id":          goals[idx].ID,
			"user_id":     *goals[idx].UserID,
			"name":        goals[idx].Name,
			"project_id":  *goals[idx].ProjectID,
			"description": goals[idx].Description,
			"hours_count": goals[idx].HoursCount,
			"time_start":  goals[idx].TimeStart,
			"time_end":    goals[idx].TimeEnd,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetMyGoals(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *GoalTestIntegrationSuite) TestUserGoals(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/user/:user_id/goals")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(fmt.Sprint(integrationutils.DefaultUser.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	goals := make([]*models.Goal, 0, 10)
	err := faker.FakeData(&goals)
	t.Assert().NoError(err)

	for idx := range goals {
		goals[idx].UserID = &integrationutils.DefaultUser.ID
	}

	s.goalRepoMock.On("GetUserGoals", *goals[0].UserID).Return(goals, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range goals {
		new := map[string]interface{}{
			"id":          goals[idx].ID,
			"user_id":     *goals[idx].UserID,
			"name":        goals[idx].Name,
			"project_id":  *goals[idx].ProjectID,
			"description": goals[idx].Description,
			"hours_count": goals[idx].HoursCount,
			"time_start":  goals[idx].TimeStart,
			"time_end":    goals[idx].TimeEnd,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetUserGoals(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}
