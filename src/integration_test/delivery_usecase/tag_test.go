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

	tagDelivery "timetracker/internal/Tag/delivery"
	tagRepoMock "timetracker/internal/Tag/repository/mocks"
	tagUC "timetracker/internal/Tag/usecase"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type TagTestIntegrationSuite struct {
	suite.Suite
	uc            tagUC.UsecaseI
	router        *echo.Echo
	delivery      *tagDelivery.Delivery
	tagRepoMock   *tagRepoMock.RepositoryI
	tagDtoBuilder *integrationutils.ReqCreateUpdateTagBuilder
	tagBuilder    *testutils.TagBuilder
}

func getJsonDataTag(t provider.T, tagDto dto.ReqCreateUpdateTag) (tagDtoJson []byte) {
	tagDtoJson, err := json.Marshal(tagDto)
	t.Require().NoError(err)

	return
}

func TestTagTestIntegrationSuite(t *testing.T) {
	suite.RunSuite(t, new(TagTestIntegrationSuite))
}

func (s *TagTestIntegrationSuite) BeforeEach(t provider.T) {
	s.tagRepoMock = tagRepoMock.NewRepositoryI(t)
	s.tagRepoMock = tagRepoMock.NewRepositoryI(t)
	s.uc = tagUC.New(s.tagRepoMock)
	s.tagDtoBuilder = integrationutils.NewReqCreateUpdateTagBuilder()
	s.tagBuilder = testutils.NewTagBuilder()

	s.router = echo.New()
	s.delivery = &tagDelivery.Delivery{
		TagUC: s.uc,
	}

	tagDelivery.NewDelivery(s.router, s.uc, nil)
}

func (s *TagTestIntegrationSuite) TestCreateTag(t provider.T) {
	tagDto := s.tagDtoBuilder.
		WithName("name").
		WithAbout("about").
		WithColor("green").
		Build()

	tagDtoJson := getJsonDataTag(t, tagDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tagDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":      0,
			"user_id": 1,
			"name":    tagDto.Name,
			"about":   tagDto.About,
			"color":   tagDto.Color,
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	s.tagRepoMock.On("CreateTag", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.CreateTag(ctx))

	t.Assert().Equal(http.StatusCreated, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *TagTestIntegrationSuite) TestUpdateTag(t provider.T) {
	tagDto := s.tagDtoBuilder.
		WithID(1).
		WithName("name").
		WithAbout("about").
		WithColor("green").
		Build()

	tagModel := tagDto.ToModelTag()
	tagModel.UserID = integrationutils.DefaultUser.ID
	tagDtoJson := getJsonDataTag(t, tagDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tagDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":      tagDto.ID,
			"user_id": 1,
			"name":    tagDto.Name,
			"about":   tagDto.About,
			"color":   tagDto.Color,
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.tagRepoMock.On("GetTag", tagModel.ID).Return(tagModel, nil)
	s.tagRepoMock.On("UpdateTag", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.UpdateTag(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *TagTestIntegrationSuite) TestDeleteTag(t provider.T) {
	tagDto := s.tagDtoBuilder.
		WithID(1).
		WithName("name").
		WithAbout("about").
		WithColor("green").
		Build()

	tagModel := tagDto.ToModelTag()
	tagModel.UserID = integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/tag/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(tagModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	s.tagRepoMock.On("GetTag", tagModel.ID).Return(tagModel, nil)
	s.tagRepoMock.On("DeleteTag", tagModel.ID).Return(nil)

	t.Require().NoError(s.delivery.DeleteTag(ctx))
	t.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *TagTestIntegrationSuite) TestGetTag(t provider.T) {
	tagDto := s.tagDtoBuilder.
		WithID(1).
		WithName("name").
		WithAbout("about").
		WithColor("green").
		Build()

	tagModel := tagDto.ToModelTag()
	tagModel.UserID = integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/tag/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(tagModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	expData := map[string]map[string]interface{}{
		"body": {
			"id":      tagDto.ID,
			"user_id": tagModel.UserID,
			"name":    tagDto.Name,
			"about":   tagDto.About,
			"color":   tagDto.Color,
		},
	}

	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.tagRepoMock.On("GetTag", tagModel.ID).Return(tagModel, nil)

	t.Require().NoError(s.delivery.GetTag(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *TagTestIntegrationSuite) TestGetMyTag(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	tags := make([]*models.Tag, 0, 10)
	err := faker.FakeData(&tags)
	t.Assert().NoError(err)

	for idx := range tags {
		tags[idx].UserID = integrationutils.DefaultUser.ID
	}

	s.tagRepoMock.On("GetUserTags", tags[0].UserID).Return(tags, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range tags {
		new := map[string]interface{}{
			"id":      tags[idx].ID,
			"user_id": tags[idx].UserID,
			"name":    tags[idx].Name,
			"about":   tags[idx].About,
			"color":   tags[idx].Color,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetMyTags(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *TagTestIntegrationSuite) TestUserTags(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/user/:user_id/tags")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(fmt.Sprint(integrationutils.DefaultUser.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	tags := make([]*models.Tag, 0, 10)
	err := faker.FakeData(&tags)
	t.Assert().NoError(err)

	for idx := range tags {
		tags[idx].UserID = integrationutils.DefaultUser.ID
	}

	s.tagRepoMock.On("GetUserTags", tags[0].UserID).Return(tags, nil)

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range tags {
		new := map[string]interface{}{
			"id":      tags[idx].ID,
			"user_id": tags[idx].UserID,
			"name":    tags[idx].Name,
			"about":   tags[idx].About,
			"color":   tags[idx].Color,
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetUserTags(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}
