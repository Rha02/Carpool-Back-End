package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Rha02/carpool_app/driver"
	"github.com/Rha02/carpool_app/models"
	"go.mongodb.org/mongo-driver/bson"
)

type DBRepo struct {
	DB *driver.DB
}

func NewDatabaseRepo(db *driver.DB) DatabaseRepository {
	return &DBRepo{db}
}

// GetAllUsers returns an array of all users
func (m *DBRepo) GetAllUsers() []models.User {
	var res []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := m.DB.Conn.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result models.User

		err = cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		res = append(res, result)
	}

	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}

	return res
}

func (m *DBRepo) GetUserByID(id string) (*models.User, error) {
	fmt.Println("Get User By ID")
	return nil, nil
}

func (m *DBRepo) CreateUser(u models.User) error {
	fmt.Println("Create User")
	return nil
}

func (m *DBRepo) DeleteUserByID(id string) error {
	fmt.Println("Delete User By ID")
	return nil
}
