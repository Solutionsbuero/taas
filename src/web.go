package ttrn

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
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
	upgrader              websocket.Upgrader
	turnoutPositionEvents chan TurnoutPositionEvent
	trainSpeedEvents      chan TrainSpeedEvent
	trainPositionEvents   chan TrainPositionEvent
	updateFrontend        chan FrontendState
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
		upgrader:              websocket.Upgrader{},
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
	e.POST("/api/turnout/:id/change", w.postTournoutChange)
	e.POST("/api/train/:id/speed", w.postTrainSpeed)
	e.GET("/ws", w.websocket)
}

// getIndex handles the GET request on /.
func getIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

// getImpressum handles the GET request on /impressum.
func getImpressum(c echo.Context) error {
	return c.Render(http.StatusOK, "impressum.html", map[string]interface{}{})
}

// websocket provides the websocket duh.
func (w *Web) websocket(c echo.Context) error {
	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer ws.Close()

	for {
		state := <-w.updateFrontend
		data, err := json.Marshal(state)
		if err != nil {
			err := fmt.Errorf("couldn't marshal state to JSON, %s", err)
			logrus.Error(err)
			return err
		}

		if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
			err := fmt.Errorf("error occured while write status message to websocket, %s", err)
			logrus.Error(err)
			return err
		}
	}
}

// postTournoutPosition handles the POST request on /api/tunrout/.id/position.
func (w *Web) postTournoutChange(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := fmt.Errorf("couldn't parse id parameter %s as int", c.Param("id"))
		logrus.Error(err)
		return err
	}
	nPos, err := w.state.SwitchTurnout(id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	w.turnoutPositionEvents <- TurnoutPositionEvent{
		Id:          id,
		NewPosition: nPos,
	}
	return nil
}

type trainSpeedRequest struct {
	SpeedDelta int `json:"speed_delta"`
}

// postTrainSpeed handles the POST request on /api/train/speed.
func (w *Web) postTrainSpeed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := fmt.Errorf("couldn't parse id parameter %s as int", c.Param("id"))
		logrus.Error(err)
		return err
	}
	d := new(trainSpeedRequest)
	if err := c.Bind(d); err != nil {
		err := fmt.Errorf("fail to bind train speed data, %s", err)
		logrus.Error(err)
		return err
	}
	nSpeed, err := w.state.ChangeTrainSpeed(id, d.SpeedDelta)
	if err != nil {
		logrus.Error(err)
		return err
	}
	w.trainSpeedEvents <- TrainSpeedEvent{
		Id:       id,
		NewSpeed: nSpeed,
	}
	return nil
}
