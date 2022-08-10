package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Rha02/Carpool-Back-End/driver"
	"github.com/Rha02/Carpool-Back-End/models"
	"github.com/Rha02/Carpool-Back-End/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// DBRepo is a repository of database related methods
type DBRepo struct {
	DB *driver.DB
}

// TestDBRepo imitates DBRepo in unit tests
type TestDBRepo struct{}

// NewDatabaseRepo creates and returns a new DBRepo
func NewDatabaseRepo(db *driver.DB) DatabaseRepository {
	return &DBRepo{db}
}

// NewTestingRepo creates and returns a TestDBRepo
func NewTestingRepo() DatabaseRepository {
	return &TestDBRepo{}
}

// TODO: replace err.Error() before deploying

// Authenticate authenticates a user using given credentials
func (m *DBRepo) Authenticate(phone string, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u models.User

	filter := bson.M{"phone": phone}

	if err := m.DB.Conn.Collection("users").FindOne(ctx, filter).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &utils.DBError{Msg: "error: user with this phone number is not found", Code: http.StatusNotFound}
		}

		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, &utils.DBError{Msg: "error: invalid credentials", Code: http.StatusUnauthorized}
	}

	return &u, nil
}

// RegisterUser creates a new user and authenticates them
func (m *DBRepo) RegisterUser(u models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"phone": u.Phone}
	if err := m.DB.Conn.Collection("users").FindOne(ctx, filter).Err(); !errors.Is(err, mongo.ErrNoDocuments) {
		if err != nil {
			return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
		}
		return nil, &utils.DBError{Msg: "error: this phone number is already in use", Code: http.StatusBadRequest}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	u.Password = string(hashedPassword)

	now := time.Now()
	u.CreatedAt, u.UpdatedAt = now, now

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

// GetUserByID returns a user in records by their id
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

// DeleteUserByID is an unimplemented method meant for User deletion
func (m *DBRepo) DeleteUserByID(id string) error {
	/* TODO:
	- Develop reliable validation in handlers first
	- Before deleting user, delete user-owned resources (threads, comments, etc.)
	*/

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// objectId, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	// }

	// filter := bson.M{"_id": objectId}
	// res, err := m.DB.Conn.Collection("users").DeleteOne(ctx, filter)
	// if err != nil {
	// 	return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	// }

	// if res.DeletedCount == 0 {
	// 	return &utils.DBError{Msg: "error: user with id %s does not exist", Code: http.StatusNotFound}
	// }

	return &utils.DBError{Msg: "error: not implemented yet", Code: http.StatusNotImplemented}
}

// UpdateUserByID updates a user by their id
func (m *DBRepo) UpdateUserByID(id string, updatedUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	updatedUser.UpdatedAt = time.Now()

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

// GetAllThreads returns an array of all existing threads
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

// GetThreadByID returns a thread by an id
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

// CreateThread creates a new thread
func (m *DBRepo) CreateThread(t models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	t.CreatedAt, t.UpdatedAt = now, now

	if _, err := m.DB.Conn.Collection("threads").InsertOne(ctx, t); err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return nil
}

// DeleteThreadByID deletes a thread by a given id
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

// UpdateThreadByID updates a thread by a given id
func (m *DBRepo) UpdateThreadByID(id string, ut models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	ut.UpdatedAt = time.Now()

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

// GetUserThreads returns an array of threads created by a user with a specified id
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

func (m *DBRepo) GetComment(t_id string, c_id string) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	threadID, err := primitive.ObjectIDFromHex(t_id)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	var t models.Thread

	filter := bson.M{"_id": threadID}
	if err = m.DB.Conn.Collection("threads").FindOne(ctx, filter).Decode(&t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusNotFound}
		}
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	commentID, err := primitive.ObjectIDFromHex(c_id)
	if err != nil {
		return nil, &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	// Iterate over comments to find the one with the right ID
	var res *models.Comment
	for i := 0; i < len(t.Comments); i++ {
		if t.Comments[i].ID == commentID {
			res = &t.Comments[i]
		}
	}

	return res, nil
}

func (m *DBRepo) CreateComment(t_id string, comment models.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	threadID, err := primitive.ObjectIDFromHex(t_id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	comment.ID = primitive.NewObjectID()
	now := time.Now()
	comment.CreatedAt, comment.UpdatedAt = now, now

	query := bson.D{{
		"$push", bson.M{
			"comments": comment,
		},
	}}

	res, err := m.DB.Conn.Collection("threads").UpdateByID(ctx, threadID, query)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	if res.MatchedCount == 0 {
		return &utils.DBError{Msg: fmt.Sprintf("error: thread with id %s does not exist", threadID), Code: http.StatusNotFound}
	}

	return nil
}

func (m *DBRepo) UpdateComment(t_id string, c_id string, updatedComment models.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	threadID, err := primitive.ObjectIDFromHex(t_id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	filter := bson.M{"_id": threadID}

	var thread models.Thread

	if err = m.DB.Conn.Collection("threads").FindOne(ctx, filter).Decode(&thread); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &utils.DBError{Msg: err.Error(), Code: http.StatusNotFound}
		}
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	commentID, err := primitive.ObjectIDFromHex(c_id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	for i := 0; i < len(thread.Comments); i++ {
		if thread.Comments[i].ID == commentID {
			thread.Comments[i].Body = updatedComment.Body
			thread.Comments[i].UpdatedAt = time.Now()

			query := bson.M{"$set": thread}
			if _, err := m.DB.Conn.Collection("threads").UpdateByID(ctx, threadID, query); err != nil {
				return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
			}

			return nil
		}
	}

	return &utils.DBError{Msg: "error: comment with specified id does not exist", Code: http.StatusNotFound}
}

func (m *DBRepo) DeleteComment(t_id, c_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	threadID, err := primitive.ObjectIDFromHex(t_id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	commentID, err := primitive.ObjectIDFromHex(c_id)
	if err != nil {
		return &utils.DBError{Msg: err.Error(), Code: http.StatusBadRequest}
	}

	var thread models.Thread
	filter := bson.M{"_id": threadID}

	if err := m.DB.Conn.Collection("threads").FindOne(ctx, filter).Decode(&thread); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &utils.DBError{Msg: err.Error(), Code: http.StatusNotFound}
		}
		return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	for i := 0; i < len(thread.Comments); i++ {
		if thread.Comments[i].ID == commentID {
			thread.Comments = append(thread.Comments[:i], thread.Comments[i+1:]...)

			query := bson.M{"$set": thread}
			if _, err := m.DB.Conn.Collection("threads").UpdateByID(ctx, threadID, query); err != nil {
				return &utils.DBError{Msg: err.Error(), Code: http.StatusInternalServerError}
			}

			return nil
		}
	}

	return &utils.DBError{Msg: "error: comment with specified id does not exist", Code: http.StatusNotFound}
}
