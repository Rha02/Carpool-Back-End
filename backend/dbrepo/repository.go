package dbrepo

import "github.com/Rha02/carpool_app/models"

// DatabaseRepository is an interface for DBRepo
type DatabaseRepository interface {
	GetAllUsers() []models.User
	GetUserByID(id string) (*models.User, error)
	CreateUser(u models.User) error
	DeleteUserByID(id string) error
}
