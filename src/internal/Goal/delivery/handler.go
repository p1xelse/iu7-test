package delivery

import (
	"net/http"
	"strconv"
	goalUsecase "timetracker/internal/Goal/usecase"
	"timetracker/internal/middleware"
	"timetracker/models"
	"timetracker/models/dto"
	"timetracker/pkg"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Delivery struct {
	GoalUC goalUsecase.UsecaseI
}

func (del *Delivery) ownerOrAdminValidate(c echo.Context, goal *models.Goal) error {
	user, ok := c.Get("user").(*models.User)

	if !ok {
		c.Logger().Error("can't get user from context")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	if user.Role == models.Admin.String() || *goal.UserID == user.ID {
		return nil
	}

	return models.ErrPermissionDenied
}

// CreateGoal godoc
// @Summary      Create goal
// @Description  Create goal
// @Tags     	 goal
// @Accept	 application/json
// @Produce  application/json
// @Param    goal body dto.ReqCreateUpdateGoal true "goal info"
// @Success  200 {object} pkg.Response{body=dto.RespGoal} "success update goal"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 400 {object} echo.HTTPError "bad req"
// @Failure 403 {object} echo.HTTPError "invalid csrf or permission denied"
// @Router   /goal/create [post]
func (delivery *Delivery) CreateGoal(c echo.Context) error {

	var reqGoal dto.ReqCreateUpdateGoal
	err := c.Bind(&reqGoal)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqGoal); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	goal := reqGoal.ToModelGoal()
	goal.UserID = &userId
	err = delivery.GoalUC.CreateGoal(goal)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respGoal := dto.GetResponseFromModelGoal(goal)

	return c.JSON(http.StatusOK, pkg.Response{Body: *respGoal})
}

// GetGoal godoc
// @Summary      Show a post
// @Description  Get goal by id. Acl: admin, owner
// @Tags     	 goal
// @Accept	 application/json
// @Produce  application/json
// @Param id  path int  true  "Goal ID"
// @Success  200 {object} pkg.Response{body=dto.RespGoal} "success get goal"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /goal/{id} [get]
func (delivery *Delivery) GetGoal(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	goal, err := delivery.GoalUC.GetGoal(id)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	err = delivery.ownerOrAdminValidate(c, goal)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respGoal := dto.GetResponseFromModelGoal(goal)
	return c.JSON(http.StatusOK, pkg.Response{Body: *respGoal})
}

// UpdateGoal godoc
// @Summary      Update a goal
// @Description  Update a goal. Acl: owner only
// @Tags     	 goal
// @Accept	 application/json
// @Produce  application/json
// @Param    goal body dto.ReqCreateUpdateGoal true "goal info"
// @Success  200 {object} pkg.Response{body=dto.RespGoal} "success update goal"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf or permission denied"
// @Router   /goal/edit [post]
func (delivery *Delivery) UpdateGoal(c echo.Context) error {

	var reqGoal dto.ReqCreateUpdateGoal
	err := c.Bind(&reqGoal)

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqGoal); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	goal := reqGoal.ToModelGoal()
	goal.UserID = &userId
	err = delivery.GoalUC.UpdateGoal(goal)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respGoal := dto.GetResponseFromModelGoal(goal)

	return c.JSON(http.StatusOK, pkg.Response{Body: *respGoal})
}

// DeleteGoal godoc
// @Summary      Delete a goal. Acl: owner only
// @Description  Delete a goal
// @Tags     	 goal
// @Accept	 application/json
// @Param id path int  true  "Goal ID"
// @Success  204
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 404 {object} echo.HTTPError "can't find goal with such id"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Router   /goal/{id} [delete]
func (delivery *Delivery) DeleteGoal(c echo.Context) error {

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

	err = delivery.GoalUC.DeleteGoal(id, userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMyGoals godoc
// @Summary      Get my goals
// @Description  Get my goals. Acl: admin,
// @Tags     goal
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespGoal} "success get goals"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /me/goals [get]
func (delivery *Delivery) GetMyGoals(c echo.Context) error {
	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	goals, err := delivery.GoalUC.GetUserGoals(userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respEnties := dto.GetResponseFromModelGoals(goals)

	return c.JSON(http.StatusOK, pkg.Response{Body: respEnties})
}

// GetUserGoals godoc
// @Summary      Get user goals
// @Description  Get user goals. Acl: admin, friends
// @Tags     goal
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespGoal} "success get goals"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /user/{user_id}/goals [get]
func (delivery *Delivery) GetUserGoals(c echo.Context) error {
	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	goals, err := delivery.GoalUC.GetUserGoals(userId)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respEnties := dto.GetResponseFromModelGoals(goals)

	return c.JSON(http.StatusOK, pkg.Response{Body: respEnties})
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

func NewDelivery(e *echo.Echo, eu goalUsecase.UsecaseI, aclM *middleware.AclMiddleware) {
	handler := &Delivery{
		GoalUC: eu,
	}

	e.POST("/goal/create", handler.CreateGoal)
	e.POST("/goal/edit", handler.UpdateGoal)
	e.GET("/goal/:id", handler.GetGoal)
	e.DELETE("/goal/:id", handler.DeleteGoal)
	e.GET("/me/goals", handler.GetMyGoals)
	e.GET("/user/:user_id/goals", handler.GetUserGoals, aclM.FriendsOrAdminOnly)
}
