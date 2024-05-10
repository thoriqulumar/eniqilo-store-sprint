package server

import (
	"eniqilo-store/config"
	"eniqilo-store/controller"
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
	registerCustomerRoute(mainRoute, s.db)
	registerProductRoute(mainRoute, s.db, cfg.JWTSecret)
}

func registerHealthRoute(e *echo.Group, db *sqlx.DB) {
	ctr := controller.NewController(service.NewService(repo.NewRepo(db)))

	e.GET("/health", ctr.HealthCheck)

}

func registerCustomerRoute(e *echo.Group, db *sqlx.DB) {
	ctr := controller.NewCheckoutController(service.NewCheckoutService(repo.NewCheckoutRepo(db)))
	e.POST("/customer/register", ctr.PostCustomer)
	e.POST("/product/checkout", ctr.PostCheckout)
	e.GET("/customer", ctr.GetCustomer)
}

func registerStaffRoute(e *echo.Group, db *sqlx.DB, cfg *config.Config, validate *validator.Validate) {
	ctr := controller.NewStaffContoller(service.NewStaffService(cfg, repo.NewStaffRepo(db)), validate)

	e.POST("/staff/login", ctr.Login)
	e.POST("/staff/register", ctr.Register)
}

func registerProductRoute(e *echo.Group, db *sqlx.DB, secretKey string) {
	ctr := controller.NewProductController(service.NewProductService(repo.NewProductRepo(db)))
	e.POST("/api/product", ctr.PostProduct, middleware.Authentication(secretKey))
	e.DELETE("/product/:id", ctr.DeleteProduct, middleware.Authentication(secretKey))
}
