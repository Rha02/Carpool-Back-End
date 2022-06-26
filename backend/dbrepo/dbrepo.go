package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Rha02/carpool_app/driver"
	"github.com/Rha02/carpool_app/models"
	"github.com/Rha02/carpool_app/utils"
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

// TODO: replace err.Error() before deploying

func (m *DBRepo) Authenticate(email string, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u models.User

	filter := bson.M{"email": email}

	if err := m.DB.Conn.Collection("users").FindOne(ctx, filter).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &utils.DBError{Msg: "error: user with this email is not found", Code: http.StatusNotFound}
		}

		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, &utils.DBError{Msg: "error: invalid credentials", Code: http.StatusUnauthorized}
	}

	return &u, nil
}

func (m *DBRepo) RegisterUser(u models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": u.Email}
	if err := m.DB.Conn.Collection("users").FindOne(ctx, filter).Err(); !errors.Is(err, mongo.ErrNoDocuments) {
		if err != nil {
			return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
		}
		return nil, &utils.DBError{Msg: "error: this email is already in use", Code: http.StatusBadRequest}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	u.Password = string(hashedPassword)

	res, err := m.DB.Conn.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
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
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	defer cur.Close(ctx)

	if err = cur.All(ctx, &res); err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return res, nil
}

func (m *DBRepo) GetUserByID(id string) (*models.User, error) {
	var res models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.M{"_id": objectId}
	if err = m.DB.Conn.Collection("users").FindOne(ctx, filter).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &utils.DBError{Msg: fmt.Sprintf("error: user with id %s does not exist", id), Code: http.StatusNotFound}
		}
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return &res, nil
}

func (m *DBRepo) DeleteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.M{"_id": objectId}
	res, err := m.DB.Conn.Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if res.DeletedCount == 0 {
		return &utils.DBError{Msg: "error: user with id %s does not exist", Code: http.StatusNotFound}
	}

	return nil
}

func (m *DBRepo) UpdateUserByID(id string, updatedUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.D{{Key: "$set", Value: updatedUser}}

	res, err := m.DB.Conn.Collection("users").UpdateByID(ctx, objectId, filter)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	if res.MatchedCount == 0 {
		return &utils.DBError{Msg: fmt.Sprintf("error: user with id %s does not exist", id), Code: http.StatusNotFound}
	}

	return nil
}

func (m *DBRepo) GetAllThreads() ([]models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := m.DB.Conn.Collection("threads").Find(ctx, bson.D{})
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	defer cur.Close(ctx)

	var res []models.Thread

	if err = cur.All(ctx, &res); err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return res, nil
}

func (m *DBRepo) GetThreadByID(id string) (*models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res models.Thread

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.M{"_id": objectID}
	if err = m.DB.Conn.Collection("threads").FindOne(ctx, filter).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusNotFound}
		}
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return &res, nil
}

func (m *DBRepo) CreateThread(t models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.DB.Conn.Collection("threads").InsertOne(ctx, t); err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return nil
}

func (m *DBRepo) DeleteThreadByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.M{"_id": objectID}
	res, err := m.DB.Conn.Collection("threads").DeleteOne(ctx, filter)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if res.DeletedCount == 0 {
		return &utils.DBError{Msg: fmt.Sprintf("error: thread with id %s does not exist", id), Code: http.StatusNotFound}
	}

	return nil
}

func (m *DBRepo) UpdateThreadByID(id string, ut models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.D{{"$set", ut}}
	res, err := m.DB.Conn.Collection("threads").UpdateByID(ctx, objectID, filter)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if res.MatchedCount == 0 {
		return &utils.DBError{Msg: fmt.Sprintf("error: thread with id %s does not exist", id), Code: http.StatusNotFound}
	}

	return nil
}

func (m *DBRepo) GetUserThreads(u_id string) ([]models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(u_id)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.D{{"user_id", objectID}}

	cur, err := m.DB.Conn.Collection("threads").Find(ctx, filter)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	defer cur.Close(ctx)

	var res []models.Thread

	if err = cur.All(ctx, &res); err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return res, nil
}
