package server

import (
	"eniqilo-store/controller"
	"eniqilo-store/repo"
	"eniqilo-store/service"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1")

	registerHealthRoute(mainRoute, s.db)
	registerStaffRoute(mainRoute, s.db, s.validator)
}

func registerHealthRoute(e *echo.Group, db *sqlx.DB) {
	ctr := controller.NewController(service.NewService(repo.NewRepo(db)))

	e.GET("/health", ctr.HealthCheck)
}

func registerStaffRoute(e *echo.Group, db *sqlx.DB, validate *validator.Validate) {
	controller.NewStaffContoller(service.NewStaffService(repo.NewStaffRepo(db)), validate)
}
