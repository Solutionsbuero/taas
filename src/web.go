package ttrn

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
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
	cfg                   Config
	state                 State
	echo                  *echo.Echo
	turnoutPositionEvents chan TurnoutPositionEvent
	trainSpeedEvents      chan TrainSpeedEvent
	trainPositionEvents   chan TrainPositionEvent
}

// NewWeb returns a new instance of the Web struct. When debug parameter is true, debugging
// is enabled.
func NewWeb(cfg Config, doDebug bool, turnoutPositionEvents chan TurnoutPositionEvent, trainSpeedEvents chan TrainSpeedEvent, trainPositionEvents chan TrainPositionEvent) Web {
	e := echo.New()

	if doDebug {
		log.Logger().SetLevel(echoLog.DEBUG)
	} else {
		log.Logger().SetLevel(echoLog.WARN)
	}
	e.Logger = log.Logger()
	e.Use(middleware.Logger())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("lala"))))

	e.Renderer = &TemplateRenderer{
		tpls: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Static("/", "public/assets")

	rsl := Web{
		cfg:                   cfg,
		echo:                  e,
		state:                 DefaultState(),
		turnoutPositionEvents: turnoutPositionEvents,
		trainSpeedEvents:      trainSpeedEvents,
		trainPositionEvents:   trainPositionEvents,
	}

	rsl.addRoutes(e, cfg)
	return rsl
}

// Run runs the web server.
func (w *Web) Run() {
	w.echo.Logger.Fatal(w.echo.Start(fmt.Sprintf(":%d", w.cfg.Port)))
}

// addRoutes adds the routes to the echo element.
func (w Web) addRoutes(e *echo.Echo, cfg Config) {
	e.GET("/", getIndex)
	e.GET("/impressum", getImpressum)
	e.POST("/api/turnout/:id/position", w.postTournoutPosition)
	e.POST("/api/train/:id/speed", w.postTrainSpeed)
}

// getIndex handles the GET request on /.
func getIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

// getImpressum handles the GET request on /impressum.
func getImpressum(c echo.Context) error {
	return c.Render(http.StatusOK, "impressum.html", map[string]interface{}{})
}

type turnoutPositionRequest struct {
	Position int `json:"position"`
}

// postTournoutPosition handles the POST request on /api/tunrout/.id/position.
func (w Web) postTournoutPosition(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := fmt.Errorf("couldn't parse id parameter %s as int", c.Param("id"))
		logrus.Error(err)
		return err
	}
	if id < 0 || id > 4 {
		err := fmt.Errorf("got invalid turnout id %d", id)
		logrus.Error(err)
		return err
	}
	d := new(turnoutPositionRequest)
	if err := c.Bind(d); err != nil {
		err := fmt.Errorf("fail to bind turnout position data, %s", err)
		logrus.Error(err)
		return err
	}
	if d.Position != -1 && d.Position != 1 {
		err := fmt.Errorf("%d is not a valid position for a turnout", d.Position)
		logrus.Error(err)
		return err
	}
	w.turnoutPositionEvents <- TurnoutPositionEvent{
		Id:          id,
		NewPosition: d.Position,
	}
	return nil
}

type trainSpeedRequest struct {
	Speed int `json:"speed"`
}

// postTrainSpeed handles the POST request on /api/train/speed.
func (w Web) postTrainSpeed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := fmt.Errorf("couldn't parse id parameter %s as int", c.Param("id"))
		logrus.Error(err)
		return err
	}
	if id < 1 || id > 2 {
		err := fmt.Errorf("got invalid train id %d", id)
		logrus.Error(err)
		return err
	}
	d := new(trainSpeedRequest)
	if err := c.Bind(d); err != nil {
		err := fmt.Errorf("fail to bind train speed data, %s", err)
		logrus.Error(err)
		return err
	}
	if d.Speed < -4 || d.Speed > 4 {
		err := fmt.Errorf("%d is not a valid speed for a train", d.Speed)
		logrus.Error(err)
		return err
	}
	w.trainSpeedEvents <- TrainSpeedEvent{
		Id:       id,
		NewSpeed: d.Speed,
	}
	return nil
}
