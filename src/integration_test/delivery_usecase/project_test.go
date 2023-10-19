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

	projectDelivery "timetracker/internal/Project/delivery"
	projectRepoMock "timetracker/internal/Project/repository/mocks"
	projectUC "timetracker/internal/Project/usecase"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type ProjectTestIntegrationSuite struct {
	suite.Suite
	uc                projectUC.UsecaseI
	router            *echo.Echo
	delivery          *projectDelivery.Delivery
	projectRepoMock   *projectRepoMock.RepositoryI
	projectDtoBuilder *integrationutils.ReqCreateUpdateProjectBuilder
	projectBuilder    *testutils.ProjectBuilder
}

func getJsonDataProject(t provider.T, projectDto dto.ReqCreateUpdateProject) (projectDtoJson []byte) {
	projectDtoJson, err := json.Marshal(projectDto)
	t.Require().NoError(err)

	return
}

func TestProjectTestIntegrationSuite(t *testing.T) {
	suite.RunSuite(t, new(ProjectTestIntegrationSuite))
}

func (s *ProjectTestIntegrationSuite) BeforeEach(t provider.T) {
	s.projectRepoMock = projectRepoMock.NewRepositoryI(t)
	s.projectRepoMock = projectRepoMock.NewRepositoryI(t)
	s.uc = projectUC.New(s.projectRepoMock, nil)
	s.projectDtoBuilder = integrationutils.NewReqCreateUpdateProjectBuilder()
	s.projectBuilder = testutils.NewProjectBuilder()

	s.router = echo.New()
	s.delivery = &projectDelivery.Delivery{
		ProjectUC: s.uc,
	}

	projectDelivery.NewDelivery(s.router, s.uc, nil)
}

func (s *ProjectTestIntegrationSuite) TestCreateProject(t provider.T) {
	projectDto := s.projectDtoBuilder.
		WithName("name").
		WithAbout("about").
		WithColor("green").
		WithIsPrivate(true).
		Build()

	projectDtoJson := getJsonDataProject(t, projectDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(projectDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":                0,
			"user_id":           1,
			"name":              projectDto.Name,
			"about":             projectDto.About,
			"color":             projectDto.Color,
			"is_private":        projectDto.IsPrivate,
			"total_count_hours": 0,
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	s.projectRepoMock.On("CreateProject", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.CreateProject(ctx))

	t.Assert().Equal(http.StatusCreated, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *ProjectTestIntegrationSuite) TestUpdateProject(t provider.T) {
	projectDto := s.projectDtoBuilder.
		WithID(1).
		WithName("name").
		WithAbout("about").
		WithColor("green").
		WithIsPrivate(true).
		Build()

	projectModel := projectDto.ToModelProject()
	projectModel.UserID = &integrationutils.DefaultUser.ID
	projectDtoJson := getJsonDataProject(t, projectDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(projectDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":                projectDto.ID,
			"user_id":           1,
			"name":              projectDto.Name,
			"about":             projectDto.About,
			"color":             projectDto.Color,
			"is_private":        projectDto.IsPrivate,
			"total_count_hours": projectModel.TotalCountHours,
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.projectRepoMock.On("GetProject", projectModel.ID).Return(projectModel, nil)
	s.projectRepoMock.On("UpdateProject", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.UpdateProject(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *ProjectTestIntegrationSuite) TestDeleteProject(t provider.T) {
	projectDto := s.projectDtoBuilder.
		WithID(1).
		WithName("name").
		WithAbout("about").
		WithColor("green").
		WithIsPrivate(true).
		Build()

	projectModel := projectDto.ToModelProject()
	projectModel.UserID = &integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/project/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(projectModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	s.projectRepoMock.On("GetProject", projectModel.ID).Return(projectModel, nil)
	s.projectRepoMock.On("DeleteProject", projectModel.ID).Return(nil)

	t.Require().NoError(s.delivery.DeleteProject(ctx))
	t.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *ProjectTestIntegrationSuite) TestGetProject(t provider.T) {
	projectDto := s.projectDtoBuilder.
		WithID(1).
		WithName("name").
		WithAbout("about").
		WithColor("green").
		WithIsPrivate(true).
		Build()

	projectModel := projectDto.ToModelProject()
	projectModel.UserID = &integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/project/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(projectModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":                projectDto.ID,
			"user_id":           projectModel.UserID,
			"name":              projectDto.Name,
			"about":             projectDto.About,
			"color":             projectDto.Color,
			"is_private":        projectDto.IsPrivate,
			"total_count_hours": projectModel.TotalCountHours,
		},
	}

	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.projectRepoMock.On("GetProject", projectModel.ID).Return(projectModel, nil)

	t.Require().NoError(s.delivery.GetProject(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *ProjectTestIntegrationSuite) TestGetMyProject(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	projects := make([]*models.Project, 0, 10)
	err := faker.FakeData(&projects)
	t.Assert().NoError(err)

	for idx := range projects {
		projects[idx].UserID = &integrationutils.DefaultUser.ID
	}

	s.projectRepoMock.On("GetUserProjects", *projects[0].UserID).Return(projects, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range projects {
		new := map[string]interface{}{
			"id":                projects[idx].ID,
			"user_id":           projects[idx].UserID,
			"name":              projects[idx].Name,
			"about":             projects[idx].About,
			"color":             projects[idx].Color,
			"is_private":        projects[idx].IsPrivate,
			"total_count_hours": projects[idx].TotalCountHours,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetMyProjects(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *ProjectTestIntegrationSuite) TestUserProjects(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/user/:user_id/projects")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(fmt.Sprint(integrationutils.DefaultUser.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	projects := make([]*models.Project, 0, 10)
	err := faker.FakeData(&projects)
	t.Assert().NoError(err)

	for idx := range projects {
		projects[idx].UserID = &integrationutils.DefaultUser.ID
	}

	s.projectRepoMock.On("GetUserProjects", *projects[0].UserID).Return(projects, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range projects {
		new := map[string]interface{}{
			"id":                projects[idx].ID,
			"user_id":           projects[idx].UserID,
			"name":              projects[idx].Name,
			"about":             projects[idx].About,
			"color":             projects[idx].Color,
			"is_private":        projects[idx].IsPrivate,
			"total_count_hours": projects[idx].TotalCountHours,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetUserProjects(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}
