package server

import (
	"eniqilo-store/controller"
	"eniqilo-store/repo"
	"eniqilo-store/service"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1")

	registerHealthRoute(mainRoute, s.db)
}

func registerHealthRoute(e *echo.Group, db *sqlx.DB) {
	ctr := controller.NewController(service.NewService(repo.NewRepo(db)))

	e.GET("/health", ctr.HealthCheck)
}
