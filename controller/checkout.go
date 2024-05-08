package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/service"
	"net/http"
	cerr "eniqilo-store/utils/error"
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
		return ctx.JSON(http.StatusBadRequest, "request doesn't pass validation")
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

func (c *CheckoutController) PostOrder(ctx echo.Context) error{
	var orderRequest model.OrderRequest
	if err := ctx.Bind(&orderRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, "request doesn't pass validation")
	}

	//validate user
	_, err := c.service.ValidateUser(ctx.Request().Context(), orderRequest.CustomerId)
	if err!=nil{
		return ctx.JSON(cerr.GetCode(err), err.Error())
	}

	//validate product
	totalPrice, err := c.service.ValidateProduct(ctx.Request().Context(), orderRequest.ProductDetails)
	if err!=nil{
		return ctx.JSON(cerr.GetCode(err), err.Error())
	}

	if totalPrice > float32(orderRequest.Paid){
		return ctx.JSON(http.StatusBadRequest, "Paid amount is not enough based on all bought products")
	}

	change := float32(orderRequest.Paid) - float32(totalPrice)
	if float32(orderRequest.Change) != change {
		return ctx.JSON(http.StatusBadRequest, "Change is not correct based on all bought products and what is paid")
	}

	err = c.service.CheckoutProduct(ctx.Request().Context(), orderRequest.ProductDetails)
	if err!=nil{
		return ctx.JSON(cerr.GetCode(err), err.Error())
	}

	return ctx.JSON(http.StatusOK, "Successfully Checkout")
}
