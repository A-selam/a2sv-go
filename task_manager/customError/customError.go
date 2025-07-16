package customError

import "fmt"

type BadRequestError struct {
	Reason string
}

func (err *BadRequestError) Error() string {
	return fmt.Sprintf("Bad request: %s", err.Reason)
}

type NotFoundError struct {
	ID int
}

func (err *NotFoundError) Error() string{
	return fmt.Sprintf("Task with ID %d not found!", err.ID)
}