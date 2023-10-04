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

	entryDelivery "timetracker/internal/Entry/delivery"
	entryRepoMock "timetracker/internal/Entry/repository/mocks"
	entryUC "timetracker/internal/Entry/usecase"
	tagRepoMock "timetracker/internal/Tag/repository/mocks"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type EntryTestIntegrationSuite struct {
	suite.Suite
	uc              entryUC.UsecaseI
	router          *echo.Echo
	delivery        *entryDelivery.Delivery
	entryRepoMock   *entryRepoMock.RepositoryI
	tagRepoMock     *tagRepoMock.RepositoryI
	entryDtoBuilder *integrationutils.ReqCreateUpdateEntryBuilder
	tagBuilder      *testutils.TagBuilder
}

func getJsonData(t provider.T, entryDto dto.ReqCreateUpdateEntry) (entryDtoJson []byte) {
	entryDtoJson, err := json.Marshal(entryDto)
	t.Require().NoError(err)

	return
}

func TestEntryTestIntegrationSuite(t *testing.T) {
	suite.RunSuite(t, new(EntryTestIntegrationSuite))
}

func (s *EntryTestIntegrationSuite) BeforeEach(t provider.T) {
	s.entryRepoMock = entryRepoMock.NewRepositoryI(t)
	s.tagRepoMock = tagRepoMock.NewRepositoryI(t)
	s.uc = entryUC.New(s.entryRepoMock, s.tagRepoMock)
	s.entryDtoBuilder = integrationutils.NewReqCreateUpdateEntryBuilder()
	s.tagBuilder = testutils.NewTagBuilder()

	s.router = echo.New()
	s.delivery = &entryDelivery.Delivery{
		EntryUC: s.uc,
	}

	entryDelivery.NewDelivery(s.router, s.uc, nil)
}

func (s *EntryTestIntegrationSuite) TestCreateEntry(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entryDto := s.entryDtoBuilder.
		WithDescription("entry").
		WithTimeStart(startTime).
		WithTimeEnd(startTime.Add(2 * time.Hour)).
		WithTagList([]uint64{}).
		Build()

	entryModel := entryDto.ToModelEntry()
	entryModel.UserID = &integrationutils.DefaultUser.ID
	entryDtoJson := getJsonData(t, entryDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(entryDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	tagList := []models.Tag{}
	expData := map[string]map[string]interface{}{
		"body": {
			"id":          0,
			"user_id":     1,
			"project_id":  nil,
			"description": "entry",
			"tag_list":    tagList,
			"time_start":  startTime.Format("2006-01-02T15:04:05-07:00"),
			"time_end":    startTime.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00"),
			"duration":    "2h0m0s",
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.entryRepoMock.On("CreateEntry", mock.Anything).Return(nil)

	t.Require().NoError(s.delivery.CreateEntry(ctx))

	t.Assert().Equal(http.StatusCreated, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *EntryTestIntegrationSuite) TestUpdateEntry(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entryDto := s.entryDtoBuilder.
		WithDescription("entry").
		WithTimeStart(startTime).
		WithTimeEnd(startTime.Add(2 * time.Hour)).
		WithTagList([]uint64{}).
		WithID(1).
		Build()

	entryModel := entryDto.ToModelEntry()
	entryModel.UserID = &integrationutils.DefaultUser.ID
	entryDtoJson := getJsonData(t, entryDto)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(entryDtoJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	tagList := []models.Tag{}
	expData := map[string]map[string]interface{}{
		"body": {
			"id":          1,
			"user_id":     1,
			"project_id":  nil,
			"description": "entry",
			"tag_list":    tagList,
			"time_start":  startTime.Format("2006-01-02T15:04:05-07:00"),
			"time_end":    startTime.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00"),
			"duration":    "2h0m0s",
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.entryRepoMock.On("GetEntry", entryModel.ID).Return(entryModel, nil)
	s.entryRepoMock.On("UpdateEntry", mock.Anything).Return(nil)
	s.tagRepoMock.On("UpdateEntryTags", entryModel.ID, entryModel.TagList).Return(nil)

	t.Require().NoError(s.delivery.UpdateEntry(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *EntryTestIntegrationSuite) TestDeleteEntry(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entryDto := s.entryDtoBuilder.
		WithDescription("entry").
		WithTimeStart(startTime).
		WithTimeEnd(startTime.Add(2 * time.Hour)).
		WithTagList([]uint64{}).
		WithID(1).
		Build()

	entryModel := entryDto.ToModelEntry()
	entryModel.UserID = &integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/entry/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(entryModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.entryRepoMock.On("GetEntry", entryModel.ID).Return(entryModel, nil)
	s.entryRepoMock.On("DeleteEntry", entryModel.ID).Return(nil)
	s.tagRepoMock.On("DeleteEntryTags", entryModel.ID).Return(nil)

	t.Require().NoError(s.delivery.DeleteEntry(ctx))
	t.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *EntryTestIntegrationSuite) TestGetEntry(t provider.T) {
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entryDto := s.entryDtoBuilder.
		WithDescription("entry").
		WithTimeStart(startTime).
		WithTimeEnd(startTime.Add(2 * time.Hour)).
		WithTagList([]uint64{}).
		WithID(1).
		Build()

	entryModel := entryDto.ToModelEntry()
	entryModel.UserID = &integrationutils.DefaultUser.ID

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/entry/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(fmt.Sprint(entryModel.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	tags := make([]*models.Tag, 0, 10)
	err := faker.FakeData(&tags)
	t.Assert().NoError(err)
	for _, tag := range tags {
		entryModel.TagList = append(entryModel.TagList, *tag)
	}

	expData := map[string]map[string]interface{}{
		"body": {
			"id":          1,
			"user_id":     1,
			"project_id":  nil,
			"description": "entry",
			"tag_list":    entryModel.TagList,
			"time_start":  startTime.Format("2006-01-02T15:04:05-07:00"),
			"time_end":    startTime.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00"),
			"duration":    "2h0m0s",
		},
	}
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	// mock.Anything так как в функции берется указатель на user_id, в тесте этот указатель мы взять не можем
	s.entryRepoMock.On("GetEntry", entryModel.ID).Return(entryModel, nil)
	s.tagRepoMock.On("GetEntryTags", entryModel.ID).Return(tags, nil)

	t.Require().NoError(s.delivery.GetEntry(ctx))

	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *EntryTestIntegrationSuite) TestGetMyEntry(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entries := make([]*models.Entry, 0, 10)
	err := faker.FakeData(&entries)
	t.Assert().NoError(err)
	entries = entries[:1]
	tagList := []models.Tag{s.tagBuilder.Build(), s.tagBuilder.Build(), s.tagBuilder.Build()}

	for idx := range entries {
		entries[idx].UserID = &integrationutils.DefaultUser.ID
		entries[idx].TagList = tagList
		entries[idx].TimeStart = startTime
		entries[idx].TimeEnd = startTime.Add(2 * time.Hour)
	}

	s.entryRepoMock.On("GetUserEntries", *entries[0].UserID).Return(entries, nil)
	for _, entry := range entries {
		s.tagRepoMock.On("GetEntryTags", entry.ID).Return(testutils.MakePointerSlice(entry.TagList), nil)
	}

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range entries {
		new := map[string]interface{}{
			"id":          entries[idx].ID,
			"user_id":     *entries[0].UserID,
			"project_id":  *entries[idx].ProjectID,
			"description": entries[idx].Description,
			"tag_list":    entries[idx].TagList,
			"time_start":  startTime.Format("2006-01-02T15:04:05-07:00"),
			"time_end":    startTime.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00"),
			"duration":    "2h0m0s",
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetMyEntries(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}

func (s *EntryTestIntegrationSuite) TestUserEntries(t provider.T) {
	// http prepare
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := s.router.NewContext(req, rec)
	ctx.SetPath("/user/:user_id/entries")
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(fmt.Sprint(integrationutils.DefaultUser.ID))

	// имитируем вызов мидлвары
	integrationutils.AuthMiddlewareAction(ctx)

	//model prepare
	startTime := time.Date(2021, 8, 15, 14, 30, 45, 0, time.Local)
	entries := make([]*models.Entry, 0, 10)
	err := faker.FakeData(&entries)
	t.Assert().NoError(err)
	entries = entries[:1]
	tagList := []models.Tag{s.tagBuilder.Build(), s.tagBuilder.Build(), s.tagBuilder.Build()}

	for idx := range entries {
		entries[idx].UserID = &integrationutils.DefaultUser.ID
		entries[idx].TagList = tagList
		entries[idx].TimeStart = startTime
		entries[idx].TimeEnd = startTime.Add(2 * time.Hour)
	}

	s.entryRepoMock.On("GetUserEntries", *entries[0].UserID).Return(entries, nil)
	for _, entry := range entries {
		s.tagRepoMock.On("GetEntryTags", entry.ID).Return(testutils.MakePointerSlice(entry.TagList), nil)
	}

	expData := map[string][]map[string]interface{}{
		"body": {},
	}
	bodyArr := []map[string]interface{}{}
	for idx := range entries {
		new := map[string]interface{}{
			"id":          entries[idx].ID,
			"user_id":     *entries[0].UserID,
			"project_id":  *entries[idx].ProjectID,
			"description": entries[idx].Description,
			"tag_list":    entries[idx].TagList,
			"time_start":  startTime.Format("2006-01-02T15:04:05-07:00"),
			"time_end":    startTime.Add(2 * time.Hour).Format("2006-01-02T15:04:05-07:00"),
			"duration":    "2h0m0s",
		}
		bodyArr = append(bodyArr, new)
	}
	expData["body"] = bodyArr
	expDataJson, err := json.Marshal(expData)
	t.Require().NoError(err, "error while prepare expDataJson")

	t.Require().NoError(s.delivery.GetUserEntries(ctx))
	t.Assert().Equal(http.StatusOK, rec.Code)
	t.Assert().JSONEq(string(expDataJson), rec.Body.String())
}
