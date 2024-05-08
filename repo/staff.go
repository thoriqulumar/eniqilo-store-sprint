package repo

import (
	"eniqilo-store/model"

	"github.com/jmoiron/sqlx"
)

type StaffRepo interface {
	GetStaff(phoneNumber string) (*model.Staff, error)
	CreateStaff(newStaff *model.RegisterStaffRequest, hashPassword string) (string, error)
}

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepo(db *sqlx.DB) StaffRepo {
	return &staffRepo{db}
}

func (r *staffRepo) GetStaff(phoneNumber string) (*model.Staff, error) {
	var staff model.Staff

	query := `SELECT * FROM staff WHERE "phoneNumber" = $1`

	err := r.db.Get(&staff, query, phoneNumber)
	if err != nil {
		return nil, err
	}

	return &staff, nil
}

func (r *staffRepo) CreateStaff(newStaff *model.RegisterStaffRequest, hashPassword string) (string, error) {
	var userId string

	query := `INSERT INTO staff (name, "phoneNumber", password) VALUES ($1, $2, $3) RETURNING userId`

	row := r.db.QueryRowx(query, newStaff.Name, newStaff.PhoneNumber, hashPassword)

	if err := row.Scan(&userId); err != nil {
		return "", nil
	}

	return userId, nil
}
