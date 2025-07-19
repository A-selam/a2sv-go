package models

type role string

const (
	Admin       role = "Admin"
	RegularUser role = "User"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     role   `json:"role"`
}