package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/pkg/customErr"
	"eniqilo-store/service"
	cerr "eniqilo-store/utils/error"
	"net/http"
	"strconv"

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
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	// Call the service to create a new customer
	newCustomer, err := c.service.CreateNewCustomer(ctx.Request().Context(), customerRequest)
	if err != nil {
		resErr := customErr.NewNotFoundError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	return ctx.JSON(http.StatusCreated, model.CustomerResponse{
		Message: "Customer registered successfully",
		Data:    newCustomer,
	})
}

func (c *CheckoutController) PostCheckout(ctx echo.Context) error {
	var orderRequest model.OrderRequest
	if err := ctx.Bind(&orderRequest); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	//validate user
	_, err := c.service.ValidateUser(ctx.Request().Context(), orderRequest.CustomerId)
	if err != nil {
		resErr := customErr.NewNotFoundError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	//validate product
	totalPrice, err := c.service.ValidateProduct(ctx.Request().Context(), orderRequest.ProductDetails)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.ErrorMessageOrder{
			Message:    err.Error(),
			StatusCode: cerr.GetCode(err),
		})
	}

	if totalPrice > float32(orderRequest.Paid) {
		resErr := customErr.NewBadRequestError("Paid amount is not enough based on all bought products")
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	change := float32(orderRequest.Paid) - float32(totalPrice)
	if float32(orderRequest.Change) != change {
		resErr := customErr.NewBadRequestError("Change is not correct based on all bought products and what is paid")
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	err = c.service.CheckoutProduct(ctx.Request().Context(), orderRequest.ProductDetails)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.ErrorMessageOrder{
			Message:    err.Error(),
			StatusCode: cerr.GetCode(err),
		})
	}

	return ctx.JSON(http.StatusOK, "Successfully Checkout")
}

func (c *CheckoutController) GetCustomer(ctx echo.Context) error {
	// Retrieve query parameters
	phoneNumber := ctx.QueryParam("phoneNumber")
	name := ctx.QueryParam("name")
	limitStr := ctx.QueryParam("limit")
	offsetStr := ctx.QueryParam("offset")

	// Convert limit and offset parameters to integers
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// get all customer
	listCustomer, err := c.service.GetAllCustomer(ctx.Request().Context(), name, phoneNumber, limit, offset)
	if err != nil {
		resErr := customErr.NewInternalServerError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	return ctx.JSON(http.StatusOK, model.ResponseCustomerList{
		Message: "Customer registered successfully",
		Data:    listCustomer,
	})
}
