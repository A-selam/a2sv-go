package usecases

import (
	"context"
	"fmt"
	domain "task_manager/Domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userUsecase struct{
	userRepository domain.UserRepository
	jwtServices domain.JWTService
	passwordServices domain.PasswordService
	contextTimeout time.Duration
}

func NewUserUsecase(
	userRepository domain.UserRepository,
	jwtServices domain.JWTService,
	passwordServices domain.PasswordService,
	timeout time.Duration,
) domain.UserUsecase {
	return &userUsecase{
		userRepository:   userRepository,
		jwtServices:      jwtServices,
		passwordServices: passwordServices,
		contextTimeout:   timeout,
	}
}


func (u *userUsecase) Register(c context.Context, user domain.User) error{
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if user.Username == "" || user.Password == ""{
		return &domain.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	if user.Role != domain.Admin && user.Role != domain.RegularUser{
		return &domain.BadRequestError{Reason: "Invalid role input."}
	}

	_, err := u.userRepository.GetUser(c, user.Username)
	if err == nil {
		return &domain.BadRequestError{Reason: "Username is already taken."}
	}
	if err != mongo.ErrNoDocuments {
		return err 
	}

	newUserID, err := u.userRepository.GetNewUserID(c)
	if err != nil {
		return err
	}

	user.ID = newUserID

	hashedPassword, err := u.passwordServices.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	
	err = u.userRepository.Register(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) Login(c context.Context, user domain.Login) (string, error){
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if user.Username == "" || user.Password == ""{
		return "", &domain.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	existingUser, err := u.userRepository.GetUser(ctx, user.Username)
	if err == mongo.ErrNoDocuments || u.passwordServices.ComparePassword(existingUser.Password, user.Password) != nil{
		return "", &domain.BadRequestError{Reason: "Incorrect username or password"}
	}
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return u.jwtServices.GenerateToken(existingUser.Username, string(existingUser.Role))
}

