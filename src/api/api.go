package api

import (
	"github.com/eduardor2m/work-with-sqlc/src/api/router"
	"github.com/labstack/echo/v4"
)

type API interface {
	Serve()
	loadRoutes()
}

type Options struct{}

type api struct {
	options      *Options
	group        *echo.Group
	echoInstance *echo.Echo
}

func NewAPI(options *Options) API {
	echoInstance := echo.New()
	return &api{options: options, group: echoInstance.Group("/api"), echoInstance: echoInstance}
}
func (instance *api) Serve() {
	instance.loadRoutes()
	instance.echoInstance.Logger.Fatal(instance.echoInstance.Start(":8080"))
}

func (instance *api) loadRoutes() {
	router := router.New()
	router.Load(instance.group)
}
