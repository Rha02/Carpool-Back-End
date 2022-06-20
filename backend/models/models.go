package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `json:"email"`
	Name     string             `json:"name"`
	Password string             `bson:"password,omitempty" json:"-"`
}

type Thread struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title  string             `json:"title"`
	Body   string             `json:"body"`
}
