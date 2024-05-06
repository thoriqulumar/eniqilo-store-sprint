package controller

import (
	"eniqilo-store/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	svc service.Service
}

func NewController(svc service.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (c *Controller) HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Healthy")
}
