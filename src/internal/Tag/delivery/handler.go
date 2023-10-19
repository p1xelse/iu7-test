package delivery

import (
	"net/http"
	"strconv"
	tagUsecase "timetracker/internal/Tag/usecase"
	"timetracker/internal/middleware"
	"timetracker/models"
	"timetracker/models/dto"
	"timetracker/pkg"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Delivery struct {
	TagUC tagUsecase.UsecaseI
}

func (del *Delivery) ownerOrAdminValidate(c echo.Context, tag *models.Tag) error {
	user, ok := c.Get("user").(*models.User)

	if !ok {
		c.Logger().Error("can't get user from context")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	if user.Role == models.Admin.String() || tag.UserID == user.ID {
		return nil
	}

	return models.ErrPermissionDenied
}

// CreateTag godoc
// @Summary      Create tag
// @Description  Create tag
// @Tags     	 tag
// @Accept	 application/json
// @Produce  application/json
// @Param    tag body dto.ReqCreateUpdateTag true "tag info"
// @Success  200 {object} pkg.Response{body=dto.RespTag} "success update tag"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 400 {object} echo.HTTPError "bad req"
// @Failure 403 {object} echo.HTTPError "invalid csrf or permission denied"
// @Router   /tag/create [post]
func (delivery *Delivery) CreateTag(c echo.Context) error {

	var reqTag dto.ReqCreateUpdateTag
	err := c.Bind(&reqTag)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqTag); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	tag := reqTag.ToModelTag()
	tag.UserID = userId
	err = delivery.TagUC.CreateTag(tag)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respTag := dto.GetResponseFromModelTag(tag)

	return c.JSON(http.StatusCreated, pkg.Response{Body: *respTag})
}

// GetTag godoc
// @Summary      Show a post
// @Description  Get tag by id
// @Tags     	 tag
// @Accept	 application/json
// @Produce  application/json
// @Param id  path int  true  "Tag ID"
// @Success  200 {object} pkg.Response{body=dto.RespTag} "success get tag"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /tag/{id} [get]
func (delivery *Delivery) GetTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	tag, err := delivery.TagUC.GetTag(id)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	err = delivery.ownerOrAdminValidate(c, tag)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respTag := dto.GetResponseFromModelTag(tag)
	return c.JSON(http.StatusOK, pkg.Response{Body: *respTag})
}

// UpdateTag godoc
// @Summary      Update an tag
// @Description  Update an tag. Acl: owner
// @Tags     	 tag
// @Accept	 application/json
// @Produce  application/json
// @Param    tag body dto.ReqCreateUpdateTag true "tag info"
// @Success  200 {object} pkg.Response{body=dto.RespTag} "success update tag"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf or permission denied"
// @Router   /tag/edit [post]
func (delivery *Delivery) UpdateTag(c echo.Context) error {

	var reqTag dto.ReqCreateUpdateTag
	err := c.Bind(&reqTag)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqTag); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	tag := reqTag.ToModelTag()
	tag.UserID = userId
	err = delivery.TagUC.UpdateTag(tag)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respTag := dto.GetResponseFromModelTag(tag)

	return c.JSON(http.StatusOK, pkg.Response{Body: *respTag})
}

// DeleteTag godoc
// @Summary      Delete an tag
// @Description  Delete an tag. Acl: owner
// @Tags     	 tag
// @Accept	 application/json
// @Param id path int  true  "Tag ID"
// @Success  204
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 404 {object} echo.HTTPError "can't find tag with such id"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Router   /tag/{id} [delete]
func (delivery *Delivery) DeleteTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	err = delivery.TagUC.DeleteTag(id, userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMyTags godoc
// @Summary      Get my tags
// @Description  Get my tags.
// @Tags     tag
// @Produce  application/json
// @Param        day    query     string  false  "day for events"
// @Success  200 {object} pkg.Response{body=[]dto.RespTag} "success get tags"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /me/tags [get]
func (delivery *Delivery) GetMyTags(c echo.Context) error {
	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	tags, err := delivery.TagUC.GetUserTags(userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respTags := dto.GetResponseFromModelTags(tags)

	return c.JSON(http.StatusOK, pkg.Response{Body: respTags})
}

// GetUserTags godoc
// @Summary      Get user tags
// @Description  Get user tags.
// @Tags     tag
// @Produce  application/json
// @Param        day    query     string  false  "day for events"
// @Success  200 {object} pkg.Response{body=[]dto.RespTag} "success get tags"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /user/{user_id}/tags [get]
func (delivery *Delivery) GetUserTags(c echo.Context) error {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	tags, err := delivery.TagUC.GetUserTags(userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respTags := dto.GetResponseFromModelTags(tags)

	return c.JSON(http.StatusOK, pkg.Response{Body: respTags})
}

func handleError(err error) *echo.HTTPError {
	causeErr := errors.Cause(err)
	switch {
	case errors.Is(causeErr, models.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, models.ErrNotFound.Error())
	case errors.Is(causeErr, models.ErrBadRequest):
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	case errors.Is(causeErr, models.ErrPermissionDenied):
		return echo.NewHTTPError(http.StatusForbidden, models.ErrPermissionDenied.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, causeErr.Error())
	}
}

func NewDelivery(e *echo.Echo, eu tagUsecase.UsecaseI, aclM *middleware.AclMiddleware) {
	handler := &Delivery{
		TagUC: eu,
	}

	e.POST("/tag/create", handler.CreateTag)
	e.POST("/tag/edit", handler.UpdateTag)  // acl: owner
	e.GET("/tag/:id", handler.GetTag)       // acl: owner, admin
	e.DELETE("/tag/:id", handler.DeleteTag) // acl: owner
	e.GET("/me/tags", handler.GetMyTags)
	e.GET("/user/:user_id/tags", handler.GetUserTags, aclM.FriendsOrAdminOnly)
}
