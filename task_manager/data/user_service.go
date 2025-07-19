package data

import (
	"context"
	"task_manager/customError"
	"task_manager/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user models.User) error{
	// validate if fields are not empty 
	if user.Username == "" || user.Password == ""{
		return &customError.BadRequestError{Reason: "Fields cannot be empty!"}
	}
	// validate role field 
	if user.Role != models.Admin && user.Role != models.RegularUser{
		return &customError.BadRequestError{Reason: "Role must be either 'Admin' or 'User'"}
	}
	// validate unique username
	countUsername, err := userCollection.CountDocuments(context.TODO(), bson.D{{Key: "username", Value: user.Username}}) 
	if err != nil{
		return err
	}
	if countUsername > 0{
		return &customError.BadRequestError{Reason: "Must be unique username!"}
	}

	// get the next ID value 
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var lastUser models.User
	err = userCollection.FindOne(context.TODO(), bson.M{}, opts).Decode(&lastUser)
	if err != nil && err != mongo.ErrNoDocuments{
		return err
	}

	if err == mongo.ErrNoDocuments{
		user.ID = 1
	} else {
		user.ID = lastUser.ID + 1
	}

	// hash password 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(user models.Login) (string, error){
	// validate if fields are not empty 
	if user.Username == "" || user.Password == ""{
		return "", &customError.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	// login logic
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.Username}}).Decode(&existingUser)
	if err != nil{
		return "", err
	}

	if err == mongo.ErrNoDocuments || bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)) != nil{
		return "", &customError.BadRequestError{Reason: "Invalid username or password"}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": existingUser.Username,
		"role": existingUser.Role,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	jwtToken, err := token.SignedString(JwtSecret)
	if err != nil{
		return "", err
	}

	return jwtToken, nil
}