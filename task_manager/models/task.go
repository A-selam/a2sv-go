package models

type status string

const (
	Pending   status = "Pending"
	Completed status = "Completed"
)

// type Task struct {
// 	ID          int    `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	DueDate     string `json:"due_date"`
// 	Status      status `json:"status"` // e.g., "pending", "completed"
// }

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      status `json:"status"`
	CreatedBy   string `json:"created_by"` // Add this field for user ID or email
}
