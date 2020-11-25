package pkg

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Web handles the web view stuff.
type Web struct {
	cfg  Config
	echo *echo.Echo
}

// NewWeb returns a new instance of the Web struct.
func NewWeb(cfg Config) Web {
	return Web{
		cfg:  cfg,
		echo: echo.New(),
	}
}

// Run runs the web server.
func (w Web) Run() {
	w.echo.Logger.Fatal(w.echo.Start(fmt.Sprintf(":%d", w.cfg.Port)))
}
