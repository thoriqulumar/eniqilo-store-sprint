package service

import (
	"database/sql"
	"eniqilo-store/model"
	"eniqilo-store/pkg/crypto"
	"eniqilo-store/pkg/customErr"
	"eniqilo-store/repo"
)

type StaffService interface {
	Register(newStaff model.RegisterStaffRequest) (model.StaffServiceResponse, error)
}

type staffSvc struct {
	repo repo.StaffRepo
}

func NewStaffService(r repo.StaffRepo) StaffService {
	return &staffSvc{
		repo: r,
	}
}

func (s *staffSvc) Register(newStaff model.RegisterStaffRequest) (model.StaffServiceResponse, error) {
	existingData, err := s.repo.GetStaff(newStaff.PhoneNumber)

	if err != nil && err != sql.ErrNoRows {
		return model.StaffServiceResponse{}, err
	}

	if existingData != nil {
		return model.StaffServiceResponse{}, customErr.NewConflictError("User already exist")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newStaff.Password)
	if err != nil {
		return model.StaffServiceResponse{}, err
	}

	id, err := s.repo.CreateStaff(&newStaff, hashedPassword)
	if err != nil {
		return model.StaffServiceResponse{}, err
	}

	token, err := crypto.GenerateToken(id, newStaff.PhoneNumber, newStaff.Name)
	if err != nil {
		return model.StaffServiceResponse{}, err
	}

	serviceResponse := model.StaffServiceResponse{
		ID:          id,
		AccessToken: token,
	}

	return serviceResponse, err
}
