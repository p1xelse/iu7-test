package delivery

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	userUsecase "timetracker/internal/User/usecase"
	"timetracker/internal/middleware"
	"timetracker/models"
	"timetracker/models/dto"
	"timetracker/pkg"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	UserUC userUsecase.UsecaseI
}

// GetUser godoc
// @Summary      GetProfile
// @Description  get user's profile
// @Tags     users
// @Produce  application/json
// @Param id path int true "User ID"
// @Success  200 {object} pkg.Response{body=dto.RespUser} "success get user"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 404 {object} echo.HTTPError "can't find user with such id"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /users/{user_id} [get]
func (del *Delivery) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}
	user, err := del.UserUC.GetUser(id)

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respUser := dto.GetResponseFromModelUser(user)

	return c.JSON(http.StatusOK, pkg.Response{Body: respUser})
}

// GetUsers godoc
// @Summary      GetUsers
// @Description  get all users. Acl: admin
// @Tags     users
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=[]dto.RespUser} "success get users"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /users [get]
func (del *Delivery) GetUsers(c echo.Context) error {
	users, err := del.UserUC.GetUsers()

	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	respUsers := dto.GetResponseFromModelUsers(users)

	return c.JSON(http.StatusOK, pkg.Response{Body: respUsers})
}

// GetMe godoc
// @Summary      get info about me.
// @Description  get info about me.
// @Tags     users
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=dto.RespUser} "success get users"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /me [get]
func (del *Delivery) GetMe(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		c.Logger().Error("can't get user from context")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	respUser := dto.GetResponseFromModelUser(user)
	return c.JSON(http.StatusOK, pkg.Response{Body: respUser})
}

// UpdateUser godoc
// @Summary      UpdateUser
// @Description  update user's profile. Acl: user(owner account)
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param user body dto.ReqUpdateUser true "user data"
// @Success  204 "success update"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 404 {object} echo.HTTPError "can't find user with such id"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Router   /me/edit [put]
func (del *Delivery) UpdateUser(c echo.Context) error {
	var reqUser dto.ReqUpdateUser
	err := c.Bind(&reqUser)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqUser); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		c.Logger().Error("can't get user from context")
		return echo.NewHTTPError(http.StatusInternalServerError, models.ErrInternalServerError.Error())
	}

	modelUser := reqUser.ToModelUser()
	modelUser.ID = user.ID

	err = del.UserUC.UpdateUser(modelUser)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}
	return c.NoContent(http.StatusNoContent)
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

func NewDelivery(e *echo.Echo, uc userUsecase.UsecaseI, aclM *middleware.AclMiddleware) {
	handler := &Delivery{
		UserUC: uc,
	}

	e.GET("/users/:user_id", handler.GetUser, aclM.FriendsOrAdminOnly)
	e.GET("/me", handler.GetMe)
	e.GET("/users", handler.GetUsers, aclM.AdminOnly)
	e.PUT("/me/edit", handler.UpdateUser)
}
