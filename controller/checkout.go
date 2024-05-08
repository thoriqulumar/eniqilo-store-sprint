package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CheckoutController struct {
	service service.CheckoutService
}

func NewCheckoutController(service service.CheckoutService) *CheckoutController {
	return &CheckoutController{
		service: service,
	}
}

func (c *CheckoutController) PostCustomer(ctx echo.Context) error {
	var customerRequest model.CustomerRequest
	if err := ctx.Bind(&customerRequest); err != nil {
		return err
	}

	// Call the service to create a new customer
	newCustomer, err := c.service.CreateNewCustomer(ctx.Request().Context(), customerRequest)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, model.CustomerResponse{
		Message: "Customer registered successfully",
		Data:    newCustomer,
	})
}
