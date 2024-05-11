package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/pkg/customErr"
	"eniqilo-store/service"
	cerr "eniqilo-store/utils/error"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CheckoutController struct {
	service  service.CheckoutService
	validate *validator.Validate
}

func NewCheckoutController(service service.CheckoutService, validate *validator.Validate) *CheckoutController {
	_ = validate.RegisterValidation("phone_number", validatePhoneNumber)

	return &CheckoutController{
		service:  service,
		validate: validate,
	}
}

func (c *CheckoutController) PostCustomer(ctx echo.Context) error {
	var customerRequest model.CustomerRequest
	if err := ctx.Bind(&customerRequest); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	if err := c.validate.Struct(&customerRequest); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	// Call the service to create a new customer
	newCustomer, err := c.service.CreateNewCustomer(ctx.Request().Context(), customerRequest)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.GenericResponse{Message: err.Error()})
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

	err := c.validate.Struct(&orderRequest)
	if err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	//validate user
	_, err = c.service.ValidateUser(ctx.Request().Context(), *orderRequest.CustomerId)
	if err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
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

	if totalPrice > float32(*orderRequest.Paid) {
		resErr := customErr.NewBadRequestError("Paid amount is not enough based on all bought products")
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	change := float32(*orderRequest.Paid) - float32(totalPrice)
	if float32(*orderRequest.Change) != change {
		resErr := customErr.NewBadRequestError("Change is not correct based on all bought products and what is paid")
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	transaction := model.Transaction{
		TransactionId:  uuid.New(),
		CustomerId:     uuid.MustParse(*orderRequest.CustomerId),
		ProductDetails: orderRequest.ProductDetails,
		Paid:           *orderRequest.Paid,
		Change:         *orderRequest.Change,
	}

	err = c.service.CheckoutProduct(ctx.Request().Context(), transaction)
	if err != nil {
		return ctx.JSON(cerr.GetCode(err), model.ErrorMessageOrder{
			Message:    err.Error(),
			StatusCode: cerr.GetCode(err),
		})
	}

	return ctx.JSON(http.StatusOK, model.GenericResponse{
		Message: "Successfully Checkout",
	})
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

func (c *CheckoutController) GetHistoryTransaction(ctx echo.Context) error {
	value, err := ctx.FormParams()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	// query to service
	data, err := c.service.GetAllTransaction(ctx.Request().Context(), parseGetHistoryParams(value))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, model.GenericResponse{
		Message: "success",
		Data:    data,
	})
}

func parseGetHistoryParams(params url.Values) model.GetHistoryParam {
	var result model.GetHistoryParam

	for key, values := range params {
		switch key {
		case "customerId":
			customerId, err := uuid.Parse(values[0])
			if err == nil {
				result.CustomerId = &customerId
			}
		case "limit":
			limit, err := strconv.Atoi(values[0])
			if err == nil {
				result.Limit = limit
			}
		case "offset":
			offset, err := strconv.Atoi(values[0])
			if err == nil {
				result.Offset = offset
			}
		case "createdAt":
			result.CreatedAt = &values[0]
		}
	}

	return result
}
