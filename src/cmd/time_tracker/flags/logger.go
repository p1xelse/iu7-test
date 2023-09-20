package flags

import (
	"os"

	"github.com/labstack/echo/v4"
	elog "github.com/labstack/gommon/log"
)

type LoggerFlags struct {
	LogLevel      uint8  `toml:"level"`
	LogHeader     string `toml:"header"`
	LogHttpFormat string `toml:"log-http-format"`
	LogFilePath   string `toml:"log-file-path"`
}

func (f LoggerFlags) Init(e *echo.Echo) echo.Logger {
	e.Logger.SetLevel(elog.Lvl(f.LogLevel))
	e.Logger.SetHeader(f.LogHeader)

	file, err := os.Create(f.LogFilePath)
	if err != nil {
		e.Logger.Fatal("Failed to create log file:", err)
	}

	e.Logger.SetOutput(file)
	return e.Logger
}
