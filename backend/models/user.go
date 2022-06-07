package models

import "fmt"

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

var users = []User{
	{ID: 1, Email: "adminovsky@admin.loc", Name: "Admin Adminovsky"},
	{ID: 2, Email: "bruce@wayne.loc", Name: "Bruce Wayne"},
	{ID: 3, Email: "walter@white.loc", Name: "Heisenberg"},
}

// GetExampleUsers returns an example slice of users
func GetExampleUsers() []User {
	return users
}

// GetUserByID returns a pointer to a user by id
func GetUserByID(id int) *User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}

	return nil
}

// InsertUser takes a user and inserts it to the list
func InsertUser(user User) {
	user.ID = users[len(users)-1].ID + 1
	users = append(users, user)
}

// DeleteUser removes a user by id
func DeleteUser(id int) error {
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("could not find user by id %d", id)
}
