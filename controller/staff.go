package controller

import (
	"eniqilo-store/model"
	"eniqilo-store/pkg/customErr"
	"eniqilo-store/service"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type StaffController struct {
	svc      service.StaffService
	validate *validator.Validate
}

func NewStaffContoller(svc service.StaffService, validate *validator.Validate) *StaffController {
	validate.RegisterValidation("phone_number", validatePhoneNumber)

	return &StaffController{svc: svc, validate: validate}
}

func (c *StaffController) Register(ctx echo.Context) error {
	var newStaffReq model.RegisterStaffRequest
	if err := ctx.Bind(&newStaffReq); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	if err := c.validate.Struct(&newStaffReq); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	newStaff := model.Staff{
		Name:        newStaffReq.Name,
		PhoneNumber: newStaffReq.PhoneNumber,
		Password:    newStaffReq.Password,
	}

	serviceRes, err := c.svc.Register(newStaff)
	if err != nil {
		switch err.Error() {
		case "User already exist":
			return ctx.JSON(http.StatusConflict, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	registerStaffResponse := model.RegisterStaffResponse{
		Message: "User registered successfully",
		Data: model.StaffWithToken{
			UserId:      serviceRes.UserId,
			Name:        newStaff.Name,
			PhoneNumber: newStaff.PhoneNumber,
			AccessToken: serviceRes.AccessToken,
		},
	}

	return ctx.JSON(http.StatusCreated, registerStaffResponse)
}

func (c *StaffController) Login(ctx echo.Context) error {
	var loginReq model.LoginStaffRequest
	if err := ctx.Bind(&loginReq); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	if err := c.validate.Struct(&loginReq); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	serviceRes, err := c.svc.Login(loginReq)
	if err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	registerStaffResponse := model.RegisterStaffResponse{
		Message: "User registered successfully",
		Data: model.StaffWithToken{
			UserId:      serviceRes.UserId,
			Name:        serviceRes.Name,
			PhoneNumber: serviceRes.PhoneNumber,
			AccessToken: serviceRes.AccessToken,
		},
	}

	return ctx.JSON(http.StatusOK, registerStaffResponse)
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// Regular expression to match the phone number pattern
	regex := `^\+(?:[1-9]\d{0,2}-?[1-9]\d{1,3}-?)\d{1,9}$`
	return regexp.MustCompile(regex).MatchString(phoneNumber)
}
