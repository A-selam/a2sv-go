package models

type status string

const (
	Pending   status = "Pending"
	Completed status = "Completed"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      status `json:"status"` // e.g., "pending", "completed"
}
