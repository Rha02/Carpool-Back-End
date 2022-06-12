package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/Rha02/carpool_app/driver"
	"github.com/Rha02/carpool_app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBRepo struct {
	DB *driver.DB
}

func NewDatabaseRepo(db *driver.DB) DatabaseRepository {
	return &DBRepo{db}
}

// GetAllUsers returns an array of all users
func (m *DBRepo) GetAllUsers() ([]models.User, error) {
	var res []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := m.DB.Conn.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (m *DBRepo) GetUserByID(id string) (*models.User, error) {
	var res models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	err = m.DB.Conn.Collection("users").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &res, nil
}

func (m *DBRepo) CreateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Conn.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBRepo) DeleteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.DB.Conn.Collection("users").DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	return nil
}

func (m *DBRepo) UpdateUserByID(id string, updatedUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	query := bson.D{{Key: "$set", Value: updatedUser}}

	_, err = m.DB.Conn.Collection("users").UpdateByID(ctx, objectId, query)
	if err != nil {
		return err
	}

	return nil
}
