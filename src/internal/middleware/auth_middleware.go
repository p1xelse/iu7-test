package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	authUsecase "timetracker/internal/Auth/usecase"
)

const session_name = "session_token"

type Middleware struct {
	authUC authUsecase.UsecaseI
}

func NewMiddleware(authUC authUsecase.UsecaseI) *Middleware {
	return &Middleware{authUC: authUC}
}

func (m *Middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/signup" || c.Request().URL.Path == "/signin" ||
			c.Request().URL.Path == "/auth" || c.Request().URL.Path == "/prometheus" ||
			c.Request().URL.Path == "/favicon.ico" {
			return next(c)
		}

		cookie, err := c.Cookie(session_name)
		if err == http.ErrNoCookie {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		} else if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		user, err := m.authUC.Auth(cookie.Value)
		if err != nil {
			causeErr := errors.Cause(err)
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, causeErr.Error())
		}

		c.Set("user_id", user.ID)
		c.Set("user", user)

		return next(c)
	}
}
