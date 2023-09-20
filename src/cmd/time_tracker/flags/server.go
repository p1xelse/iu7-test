package flags

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type ServerFlags struct {
	Addr              string        `toml:"addr"`
	ReadTimeout       time.Duration `toml:"read-timeout"`
	ReadHeaderTimeout time.Duration `toml:"read-header-timeout"`
	WriteTimeout      time.Duration `toml:"write-timeout"`
}

func (f ServerFlags) Init(e *echo.Echo) *http.Server {
	return &http.Server{
		Addr:              f.Addr,
		Handler:           e,
		ReadTimeout:       f.ReadTimeout,
		ReadHeaderTimeout: f.ReadHeaderTimeout,
		WriteTimeout:      f.WriteTimeout,
	}
}
