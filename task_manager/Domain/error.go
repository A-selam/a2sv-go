package domain

import "fmt"

type BadRequestError struct {
	Reason string
}

func (err *BadRequestError) Error() string {
	return fmt.Sprintf("Bad request: %s", err.Reason)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found!", err.Resource, err.ID)
}

type UserNotFoundError struct {
	Resource string
	Username string
}

func (err *UserNotFoundError) Error() string {
	return fmt.Sprintf("%s with username %s not found!", err.Resource, err.Username)
}

type UnauthorizedError struct{}

func (err *UnauthorizedError) Error() string {
	return "User not authorized to perform this action."
}

type ForbiddenError struct{}

func (err *ForbiddenError) Error() string {
	return "You don't have permission to access this resource."
}
