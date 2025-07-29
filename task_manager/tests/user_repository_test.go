package tests

import (
	"context"
	"log"
	"os"
	domain "task_manager/Domain"
	repositories "task_manager/Repositories"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepositorySuite struct {
	suite.Suite

	repository domain.UserRepository
	db *mongo.Database
}

func (s *userRepositorySuite) SetupSuite(){
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: .env file not found or failed to load.", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == ""{
		log.Fatal("MONGODB_URI is not set")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOptions)
	s.Require().NoError(err)

	db := client.Database("task_manager_test")

	s.db = db
	s.repository = repositories.NewUserRepositoryFromDB(db)
}

func (s *userRepositorySuite) TearDownTest(){
	// err := s.db.Collection("Users").Drop(context.Background())
	err := s.db.Drop(context.TODO())
	s.Require().NoError(err)
}

func (s *userRepositorySuite) TestRegister() {
	user := domain.User{
		ID : 1,
		Username : "tester",
		Password : "12345678",
		Role : domain.Admin,
	}

	err := s.repository.Register(context.TODO(), user)
	s.NoError(err)

	registeredUser, err := s.repository.GetUser(context.TODO(), user.Username)
	s.NoError(err)
	s.Equal(registeredUser.Password, user.Password)
	s.Equal(registeredUser.Role, user.Role)
}

func TestUserRepositorySuite(t *testing.T){
	suite.Run(t, new(userRepositorySuite))
}