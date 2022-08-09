package app

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	logger "git.raoulh.pw/raoulh/my_greenhouse_backend/log"
	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"
	"github.com/sirupsen/logrus"
)

const (
	maxFileSize = 1 * 1024 * 1024 * 1024
)

type AppServer struct {
	quitHeartbeat chan interface{}
	wgDone        sync.WaitGroup

	appFiber *fiber.App
}

var logging *logrus.Entry

func init() {
	logging = logger.NewLogger("app")
}

// Init the app
func NewApp() (a *AppServer, err error) {
	logging.Infoln("Init server")

	a = &AppServer{
		quitHeartbeat: make(chan interface{}),
		appFiber: fiber.New(fiber.Config{
			ServerHeader:          "Greenhouse (Linux)",
			ReadTimeout:           time.Second * 20,
			AppName:               "Greenhouse",
			DisableStartupMessage: true,
			EnablePrintRoutes:     false,
			BodyLimit:             maxFileSize,
			Views:                 html.New(config.Config.String("general.data"), ".tpl"),
		}),
	}

	a.appFiber.
		Use(fiberLog.New(fiberLog.Config{})).
		Use(NewTokenMiddleware())

	a.appFiber.Use(cors.New(cors.Config{
		AllowOrigins: "https://greenhouse.raoulh.pw",
	}))

	a.appFiber.Hooks().OnShutdown(func() error {
		a.wgDone.Done()
		return nil
	})

	//API
	api := a.appFiber.Group("/api")

	api.Post("/auth/login", func(c *fiber.Ctx) error {
		return a.apiLogin(c)
	})
	api.Get("/auth/check", func(c *fiber.Ctx) error {
		return a.apiCheckToken(c)
	})
	api.Get("/data/full", func(c *fiber.Ctx) error {
		return a.apiGetDataFull(c)
	})
	api.Get("/data/refresh", func(c *fiber.Ctx) error {
		return a.apiDataRefresh(c)
	})

	api.Post("/notif/id", func(c *fiber.Ctx) error {
		return a.apiNotifId(c)
	})

	api.Get("/notif/ph", func(c *fiber.Ctx) error {
		return a.apiNotifGet(c, models.NotifTypePh)
	})
	api.Get("/notif/watertemp", func(c *fiber.Ctx) error {
		return a.apiNotifGet(c, models.NotifTypeWaterTemp)
	})
	api.Get("/notif/airtemp", func(c *fiber.Ctx) error {
		return a.apiNotifGet(c, models.NotifTypeAirTemp)
	})
	api.Get("/notif/humidity", func(c *fiber.Ctx) error {
		return a.apiNotifGet(c, models.NotifTypeHumidity)
	})

	api.Get("/notif/ph", func(c *fiber.Ctx) error {
		return a.apiNotifSet(c, models.NotifTypePh)
	})
	api.Get("/notif/watertemp", func(c *fiber.Ctx) error {
		return a.apiNotifSet(c, models.NotifTypeWaterTemp)
	})
	api.Get("/notif/airtemp", func(c *fiber.Ctx) error {
		return a.apiNotifSet(c, models.NotifTypeAirTemp)
	})
	api.Get("/notif/humidity", func(c *fiber.Ctx) error {
		return a.apiNotifSet(c, models.NotifTypeHumidity)
	})

	a.appFiber.Get("/apple-app-site-association", func(c *fiber.Ctx) error {
		return c.Render("apple-app-site-association", fiber.Map{})
	})

	return
}

// Run the app
func (a *AppServer) Start() {
	addr := config.Config.String("general.address") + ":" + strconv.Itoa(config.Config.Int("general.port"))

	logging.Infoln("\u21D2 Server listening on", addr)

	go func() {
		if err := a.appFiber.Listen(addr); err != nil {
			logging.Fatalf("Failed to listen http server: %v", err)
		}
	}()
	a.wgDone.Add(1)
}

// Stop the app
func (a *AppServer) Shutdown() {
	close(a.quitHeartbeat)
	a.appFiber.Shutdown()
	a.wgDone.Wait()
}

const (
	heartbeatTime = 5 * time.Second
)
