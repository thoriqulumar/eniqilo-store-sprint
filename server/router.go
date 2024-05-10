package server

import (
	"eniqilo-store/config"
	"eniqilo-store/controller"
	"eniqilo-store/middleware"
	"eniqilo-store/repo"
	"eniqilo-store/service"

	"eniqilo-store/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoute(cfg *config.Config) {
	mainRoute := s.app.Group("/v1")

	registerHealthRoute(mainRoute, s.db)
	registerStaffRoute(mainRoute, s.db, cfg, s.validator)
	registerCustomerRoute(mainRoute, s.db, cfg)
	registerProductRoute(mainRoute, s.db, cfg)
}

func registerHealthRoute(e *echo.Group, db *sqlx.DB) {
	ctr := controller.NewController(service.NewService(repo.NewRepo(db)))

	e.GET("/health", ctr.HealthCheck)

}

func registerCustomerRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config) {
	ctr := controller.NewCheckoutController(service.NewCheckoutService(repo.NewCheckoutRepo(db)))
	e.POST("/customer/register", ctr.PostCustomer, middleware.Authentication(cfg.JWTSecret))
	e.POST("/product/checkout", ctr.PostCheckout, middleware.Authentication(cfg.JWTSecret))
	e.GET("/customer", ctr.GetCustomer, middleware.Authentication(cfg.JWTSecret))
}

func registerStaffRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewStaffContoller(service.NewStaffService(cfg, repo.NewStaffRepo(db)), validate)

	e.POST("/staff/login", ctr.Login)
	e.POST("/staff/register", ctr.Register)
}

func registerProductRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config) {
	ctr := controller.NewProductController(service.NewProductService(repo.NewProductRepo(db)))
	e.POST("/product", ctr.PostProduct, middleware.Authentication(cfg.JWTSecret))
	e.PUT("/product/:id", ctr.UpdateProduct, middleware.Authentication(cfg.JWTSecret))
	e.DELETE("/product/:id", ctr.DeleteProduct, middleware.Authentication(cfg.JWTSecret))
	e.GET("/product", ctr.GetProduct, middleware.Authentication(cfg.JWTSecret))
}
