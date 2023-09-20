package delivery

import (
	"net/http"
	"time"
	"timetracker/models"
	"timetracker/models/dto"
	"timetracker/pkg"

	"github.com/pkg/errors"

	authUsecase "timetracker/internal/Auth/usecase"

	"github.com/labstack/echo/v4"
)

const sessionName = "session_token"

type Delivery struct {
	AuthUC authUsecase.UsecaseI
}

// SignUp godoc
// @Summary      SignUp
// @Description  user sign up
// @Tags     auth
// @Accept	 application/json
// @Produce  application/json
// @Param    user body dto.ReqUserSignUp true "user data"
// @Success 201 {object} pkg.Response{body=dto.RespUser} "user created"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 409 {object} echo.HTTPError "nickname already exists"
// @Failure 409 {object} echo.HTTPError "email already exists"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /signup [post]
func (del *Delivery) SignUp(c echo.Context) error {
	var reqUser dto.ReqUserSignUp
	err := c.Bind(&reqUser)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqUser); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	if reqUser.Role == models.Admin.String() {
		if reqUser.AdminToken != "secret_token" {
			c.Logger().Error("invalid secret_token")
			return echo.NewHTTPError(http.StatusBadRequest, "invalid secret_token")
		}
	}

	user := reqUser.ToModelUser()
	createdCookie, err := del.AuthUC.SignUp(user)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	c.SetCookie(&http.Cookie{
		Name:     sessionName,
		Value:    createdCookie.SessionToken,
		MaxAge:   int(createdCookie.MaxAge),
		HttpOnly: true,
	})

	respUser := dto.GetResponseFromModelUser(user)

	return c.JSON(http.StatusCreated, pkg.Response{Body: respUser})
}

// SignIn godoc
// @Summary      SignIn
// @Description  user sign in
// @Tags     auth
// @Accept	 application/json
// @Produce  application/json
// @Param    user body dto.ReqUserSignIn true "user info"
// @Success  200 {object} pkg.Response{body=dto.RespUser} "success sign in"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 404 {object} echo.HTTPError "user doesn't exist"
// @Failure 401 {object} echo.HTTPError "invalid password"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /signin [post]
func (del *Delivery) SignIn(c echo.Context) error {
	var reqUser dto.ReqUserSignIn
	err := c.Bind(&reqUser)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	if ok, err := pkg.IsRequestValid(&reqUser); !ok {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrBadRequest.Error())
	}

	user := reqUser.ToModelUser()
	gotUser, createdCookie, err := del.AuthUC.SignIn(user)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)

	}

	c.SetCookie(&http.Cookie{
		Name:     sessionName,
		Value:    createdCookie.SessionToken,
		MaxAge:   int(createdCookie.MaxAge),
		HttpOnly: true,
	})

	respUser := dto.GetResponseFromModelUser(gotUser)

	return c.JSON(http.StatusOK, pkg.Response{Body: respUser})
}

// Logout godoc
// @Summary      Logout
// @Description  user logout
// @Tags     auth
// @Produce  application/json
// @Success  204 "success logout, body is empty"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Failure 403 {object} echo.HTTPError "invalid csrf"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Router   /logout [post]
func (del *Delivery) Logout(c echo.Context) error {
	cookie, err := c.Cookie(sessionName)
	if err == http.ErrNoCookie {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = del.AuthUC.DeleteCookie(cookie.Value)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	c.SetCookie(&http.Cookie{
		Name:    sessionName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
	})
	return c.NoContent(http.StatusNoContent)
}

// Auth godoc
// @Summary      Auth
// @Description  check user auth
// @Tags     auth
// @Produce  application/json
// @Success  200 {object} pkg.Response{body=models.User} "success auth"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no cookie"
// @Router   /auth [get]
func (del *Delivery) Auth(c echo.Context) error {
	cookie, err := c.Cookie(sessionName)
	if err == http.ErrNoCookie {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	gotUser, err := del.AuthUC.Auth(cookie.Value)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.JSON(http.StatusOK, pkg.Response{Body: gotUser})
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

func NewDelivery(e *echo.Echo, uc authUsecase.UsecaseI) {
	handler := &Delivery{
		AuthUC: uc,
	}

	e.POST("/signin", handler.SignIn)
	e.POST("/signup", handler.SignUp)
	e.POST("/logout", handler.Logout)
	e.GET("/auth", handler.Auth)
}
