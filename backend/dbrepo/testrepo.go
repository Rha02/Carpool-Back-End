package dbrepo

import (
	"errors"

	"github.com/Rha02/carpool_app/models"
)

func (m *TestDBRepo) GetAllUsers() ([]models.User, error) {
	return []models.User{}, nil
}

func (m *TestDBRepo) GetUserByID(id string) (*models.User, error) {
	if id == "error" {
		return nil, errors.New("error")
	}

	return &models.User{}, nil
}

func (m *TestDBRepo) CreateUser(u models.User) error {
	if u.Name == "error" {
		return errors.New("error")
	}
	return nil
}

func (m *TestDBRepo) UpdateUserByID(id string, u models.User) error {
	if id == "error" {
		return errors.New("error")
	}
	return nil
}

func (m *TestDBRepo) DeleteUserByID(id string) error {
	if id == "error" {
		return errors.New("error")
	}
	return nil
}
