package dbrepo

import (
	"net/http"

	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/Rha02/Carpool-Back-End/utils"
)

func (m *TestDBRepo) Authenticate(email, password string) (*models.User, error) {
	return nil, nil
}

func (m *TestDBRepo) RegisterUser(u models.User) (*models.User, error) {
	return nil, nil
}

func (m *TestDBRepo) GetAllUsers() ([]models.User, error) {
	return []models.User{}, nil
}

func (m *TestDBRepo) GetUserByID(id string) (*models.User, error) {
	if id == "error" {
		return nil, &utils.DBError{Msg: "error", Code: http.StatusNotFound}
	}

	return &models.User{}, nil
}

func (m *TestDBRepo) UpdateUserByID(id string, u models.User) error {
	if id == "error" {
		return &utils.DBError{Msg: "error", Code: http.StatusNotFound}
	}
	return nil
}

func (m *TestDBRepo) DeleteUserByID(id string) error {
	if id == "error" {
		return &utils.DBError{Msg: "error", Code: http.StatusNotFound}
	}
	return nil
}

func (m *TestDBRepo) GetAllThreads() ([]models.Thread, error) {
	return nil, nil
}

func (m *TestDBRepo) GetUserThreads(id string) ([]models.Thread, error) {
	return nil, nil
}

func (m *TestDBRepo) GetThreadByID(id string) (*models.Thread, error) {
	return nil, nil
}

func (m *TestDBRepo) CreateThread(t models.Thread) error {
	return nil
}

func (m *TestDBRepo) UpdateThreadByID(id string, t models.Thread) error {
	return nil
}

func (m *TestDBRepo) DeleteThreadByID(id string) error {
	return nil
}
