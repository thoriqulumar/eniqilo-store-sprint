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
	var newStaff model.RegisterStaffRequest
	if err := ctx.Bind(&newStaff); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	if err := c.validate.Struct(&newStaff); err != nil {
		resErr := customErr.NewBadRequestError(err.Error())
		return ctx.JSON(resErr.StatusCode, resErr)
	}

	serviceRes, err := c.svc.Register(newStaff)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, model.RegisterStaffResponse{
		UserId:      serviceRes.ID,
		Name:        newStaff.Name,
		PhoneNumber: newStaff.PhoneNumber,
		AccessToken: serviceRes.AccessToken,
	})
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// Regular expression to match the phone number pattern
	regex := `^\+(?:[0-9]\d{0,2}|[1-9]\d{1,3}-?)\d{1,14}$`
	return regexp.MustCompile(regex).MatchString(phoneNumber)
}
