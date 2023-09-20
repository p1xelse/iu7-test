package delivery

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	friendUsecase "timetracker/internal/Friends/usecase"
	"timetracker/internal/middleware"
	"timetracker/models"
	"timetracker/models/dto"
	"timetracker/pkg"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	FriendsUC friendUsecase.UsecaseI
}

// Subscribe godoc
// @Summary      subscribe
// @Description  subscribe
// @Tags     friends
// @Produce  application/json
// @Param user_id path int true "Friend ID"
// @Success  201 "success subscribe"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "subscribe or user doesn't exist"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 409 {object} echo.HTTPError "subscribe already exists"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Router   /friends/subscribe/{user_id} [post]
func (delivery *Delivery) Subscribe(c echo.Context) error {
	friendId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	friendRel := &models.FriendRelation{SubscriberID: &userId, UserID: &friendId}

	err = delivery.FriendsUC.CreateFriendRelation(friendRel)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.NoContent(http.StatusCreated)
}

// Unsubscribe godoc
// @Summary      Unsubscribe
// @Description  Unsubscribe
// @Tags     friends
// @Produce  application/json
// @Param user_id path int true "Friend ID"
// @Success  204 "success unsubscribe, body is empty"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 404 {object} echo.HTTPError "friend/user/friendship doesn't exist"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Router   /friends/unsubscribe/{user_id} [delete]
func (delivery *Delivery) Unsubscribe(c echo.Context) error {
	friendId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	userId, ok := c.Get("user_id").(uint64)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	friendRel := &models.FriendRelation{SubscriberID: &userId, UserID: &friendId}

	err = delivery.FriendsUC.DeleteFriendRelation(friendRel)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMyFriends godoc
// @Summary      get my friends
// @Description  get my friends
// @Tags     friends
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespUser} "success get profile"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "user doesn't exist"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /me/friends [get]
func (delivery *Delivery) GetMyFriends(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	friends, err := delivery.FriendsUC.GetUserFriends(user.ID)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respFriends := dto.GetResponseFromModelUsers(friends)

	return c.JSON(http.StatusOK, pkg.Response{Body: respFriends})
}

// GetUserFriends godoc
// @Summary      get user friends
// @Description  get user friends
// @Tags     friends
// @Produce  application/json
// @Param user_id path int true "User ID"
// @Success  200 {object} pkg.Response{body=[]dto.RespUser} "success get profile"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "user doesn't exist"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /user/{user_id}/friends [get]
func (delivery *Delivery) GetUserFriends(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	friends, err := delivery.FriendsUC.GetUserFriends(userID)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respFriends := dto.GetResponseFromModelUsers(friends)

	return c.JSON(http.StatusOK, pkg.Response{Body: respFriends})
}

// GetMySubs godoc
// @Summary      get my subs
// @Description  get my subs
// @Tags     friends
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespUser} "success get profile"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "user doesn't exist"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /me/subs [get]
func (delivery *Delivery) GetMySubs(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		c.Logger().Error("can't parse context user_id")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	subs, err := delivery.FriendsUC.GetUserSubs(user.ID)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respFriends := dto.GetResponseFromModelUsers(subs)

	return c.JSON(http.StatusOK, pkg.Response{Body: respFriends})
}

// GetUserSubs godoc
// @Summary      get user subs
// @Description  get user subs. Acl: admin only
// @Tags     friends
// @Produce  application/json
// @Param user_id path int true "User ID"
// @Success  200 {object} pkg.Response{body=[]dto.RespUser} "success get profile"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "user doesn't exist"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /user/{user_id}/subs [get]
func (delivery *Delivery) GetUserSubs(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	subs, err := delivery.FriendsUC.GetUserSubs(userID)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respFriends := dto.GetResponseFromModelUsers(subs)

	return c.JSON(http.StatusOK, pkg.Response{Body: respFriends})
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

func NewDelivery(e *echo.Echo, uc friendUsecase.UsecaseI, aclM *middleware.AclMiddleware) {
	handler := &Delivery{
		FriendsUC: uc,
	}

	e.POST("/friends/subscribe/:user_id", handler.Subscribe)
	e.DELETE("/friends/unsubscribe/:user_id", handler.Unsubscribe)
	e.GET("/user/:user_id/subs", handler.GetUserSubs, aclM.AdminOnly)
	e.GET("/user/:user_id/friends", handler.GetUserFriends, aclM.AdminOnly)
	e.GET("/me/subs", handler.GetMySubs)
	e.GET("/me/friends", handler.GetMyFriends)
}
