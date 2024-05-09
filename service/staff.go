package service

import (
	"database/sql"
	"eniqilo-store/config"
	"eniqilo-store/model"
	"eniqilo-store/pkg/crypto"
	"eniqilo-store/pkg/customErr"
	"eniqilo-store/repo"

	"github.com/google/uuid"
)

type StaffService interface {
	Register(newStaff model.Staff) (model.StaffWithToken, error)
	Login(loginReq model.LoginStaffRequest) (model.StaffWithToken, error)
}

type staffSvc struct {
	cfg  *config.Config
	repo repo.StaffRepo
}

func NewStaffService(cfg *config.Config, r repo.StaffRepo) StaffService {
	return &staffSvc{
		repo: r,
	}
}

func (s *staffSvc) Register(newStaff model.Staff) (model.StaffWithToken, error) {
	existingData, err := s.repo.GetStaff(newStaff.PhoneNumber)

	if err != nil && err != sql.ErrNoRows {
		return model.StaffWithToken{}, customErr.NewInternalServerError("Internal server error")
	}

	if existingData != nil {
		return model.StaffWithToken{}, customErr.NewConflictError("User already exist")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newStaff.Password, s.cfg.BcryptSalt)
	if err != nil {
		return model.StaffWithToken{}, err
	}

	id := uuid.New()
	newStaff.UserId = id

	err = s.repo.CreateStaff(newStaff, hashedPassword)
	if err != nil {
		return model.StaffWithToken{}, err
	}

	token, err := crypto.GenerateToken(id, newStaff.PhoneNumber, newStaff.Name, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, err
	}

	serviceResponse := model.
		StaffWithToken{
		UserId:      id,
		AccessToken: token,
	}

	return serviceResponse, err
}

func (s *staffSvc) Login(loginReq model.LoginStaffRequest) (model.StaffWithToken, error) {
	user, err := s.repo.GetStaff(loginReq.PhoneNumber)

	if err != nil && err != sql.ErrNoRows {
		return model.StaffWithToken{}, customErr.NewInternalServerError("Internal server error")
	}

	err = crypto.VerifyPassword(loginReq.Password, user.Password)
	if err != nil {
		return model.StaffWithToken{}, customErr.NewBadRequestError("Invalid phone or password")
	}

	token, err := crypto.GenerateToken(user.UserId, user.PhoneNumber, user.Name, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, customErr.NewBadRequestError(err.Error())
	}

	serviceResponse := model.StaffWithToken{
		UserId:      user.UserId,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		AccessToken: token,
	}

	return serviceResponse, nil
}
