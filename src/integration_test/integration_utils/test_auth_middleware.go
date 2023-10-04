package integrationutils

import (
	"timetracker/models"

	"github.com/labstack/echo/v4"
)

var (
	DefaultUser = models.User{
		ID:       1,
		Name:     "name",
		Email:    "email",
		About:    "about",
		Role:     "admin",
		Password: "password",
	}
)

func AuthMiddlewareAction(ctx echo.Context)  {
	ctx.Set("user", &DefaultUser)
	ctx.Set("user_id", DefaultUser.ID)
}
