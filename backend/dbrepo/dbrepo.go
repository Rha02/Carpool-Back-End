package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rha02/carpool_app/driver"
	"github.com/Rha02/carpool_app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type DBRepo struct {
	DB *driver.DB
}

type TestDBRepo struct{}

func NewDatabaseRepo(db *driver.DB) DatabaseRepository {
	return &DBRepo{db}
}

func NewTestingRepo() DatabaseRepository {
	return &TestDBRepo{}
}

func (m *DBRepo) Authenticate(email string, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u models.User

	filter := bson.M{"email": email}
	err := m.DB.Conn.Collection("users").FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *DBRepo) RegisterUser(u models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": u.Email}
	err := m.DB.Conn.Collection("users").FindOne(ctx, filter).Err()
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("error: this email is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)

	res, err := m.DB.Conn.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}

	u.ID = res.InsertedID.(primitive.ObjectID)

	return &u, nil
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
		return nil, err
	}

	return &res, nil
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

func (m *DBRepo) GetAllThreads() ([]models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := m.DB.Conn.Collection("threads").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var res []models.Thread

	err = cur.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *DBRepo) GetThreadByID(id string) (*models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res models.Thread

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	err = m.DB.Conn.Collection("threads").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (m *DBRepo) CreateThread(t models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Conn.Collection("threads").InsertOne(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBRepo) DeleteThreadByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = m.DB.Conn.Collection("threads").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBRepo) UpdateThreadByID(id string, ut models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	query := bson.D{{"$set", ut}}
	_, err = m.DB.Conn.Collection("threads").UpdateByID(ctx, objectID, query)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBRepo) GetUserThreads(u_id string) ([]models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(u_id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"user_id", objectID}}

	cur, err := m.DB.Conn.Collection("threads").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var res []models.Thread

	err = cur.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
