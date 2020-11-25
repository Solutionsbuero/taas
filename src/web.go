package ttrn

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
)

// TemplateRenderer implements the Echo renderer interface.
type TemplateRenderer struct {
	tpls *template.Template
}

// Implement the Echo renderer interface.
func (t TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.tpls.ExecuteTemplate(w, name, data)
}

// Web handles the web view stuff.
type Web struct {
	cfg  Config
	echo *echo.Echo
}

// NewWeb returns a new instance of the Web struct. When debug parameter is true, debugging
// is enabled.
func NewWeb(cfg Config, doDebug bool) Web {
	e := echo.New()

	if doDebug {
		log.Logger().SetLevel(echoLog.DEBUG)
	} else {
		log.Logger().SetLevel(echoLog.INFO)
	}
	e.Logger = log.Logger()
	e.Use(middleware.Logger())

	e.Renderer = &TemplateRenderer{
		tpls: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Static("/", "public/assets")
	addRoutes(e, cfg)

	return Web{
		cfg:  cfg,
		echo: e,
	}
}

// Run runs the web server.
func (w *Web) Run() {
	w.echo.Logger.Fatal(w.echo.Start(fmt.Sprintf(":%d", w.cfg.Port)))
}

// addRoutes adds the routes to the echo element.
func addRoutes(e *echo.Echo, cfg Config) {
	e.GET("/", getIndex)
}

// getIndex handles the GET request on /.
func getIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}
