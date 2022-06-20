package dbrepo

import "github.com/Rha02/carpool_app/models"

// DatabaseRepository is an interface for DBRepo
type DatabaseRepository interface {
	Authenticate(email, password string) (*models.User, error)
	RegisterUser(u models.User) (*models.User, error)

	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUserByID(id string, u models.User) error
	DeleteUserByID(id string) error

	GetAllThreads() ([]models.Thread, error)
	GetUserThreads(id string) ([]models.Thread, error)
	GetThreadByID(id string) (*models.Thread, error)
	CreateThread(t models.Thread) error
	UpdateThreadByID(id string, t models.Thread) error
	DeleteThreadByID(id string) error
}
