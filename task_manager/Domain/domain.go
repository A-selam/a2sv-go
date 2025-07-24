package domain

import (
	"context"
)

type Status string

const (
	Pending   Status = "Pending"
	Completed Status = "Completed"
)

type Task struct {
	ID          int    `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	DueDate     string `json:"due_date" bson:"due_date"`
	Status      Status `json:"status" bson:"status"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
}

type TaskRepository interface{
	GetAllTasks(c context.Context) ([]*Task, error)
	GetAllUserTasks(c context.Context, username string) ([]*Task, error)
	GetTask(c context.Context, taskID int) (Task, error)
	GetTaskByIDForUser(c context.Context, taskID int, username string) (Task, error)
	UpdateTask(c context.Context, taskID int, updatedTask Task) (Task, error)
	AddTask(c context.Context, newTask Task) (Task, error)
	GetNewID(c context.Context)(int, error)
	DeleteTask(c context.Context, taskID int) error
}

type TaskUsecase interface {
	GetAllTasks(c context.Context, username, role string)([]*Task, error)
	GetTask(c context.Context, taskID, username, role string)(Task, error)
	UpdateTask(c context.Context, taskId, username, role string, updatedTask Task) (Task, error)
	AddTask(c context.Context, username string, task Task) (Task, error)
	DeleteTask(c context.Context, id, username, role string) error
}

type JWTService interface {
	ParseToken(tokenString string) (string, string, error)
	GenerateToken(username, role string) (string, error)
}

type PasswordService interface{
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, plainPassword string) error
}