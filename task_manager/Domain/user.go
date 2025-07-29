package domain

import "context"

type Role string

const (
	Admin       Role = "Admin"
	RegularUser Role = "User"
)

type User struct {
	ID       int    `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Role     Role   `json:"role" bson:"role"`
}

type Login struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type UserRepository interface {
	Register(c context.Context, user User) error
	GetUser(c context.Context, username string) (*User, error)
	GetNewUserID(c context.Context) (int, error)
}

type UserUsecase interface {
	Register(c context.Context, user User) error
	Login(c context.Context, user Login) (string, error)
}