package tests

import (
	"context"
	"errors"
	domain "task_manager/Domain"
	usecases "task_manager/Usecases"
	"task_manager/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserUsecaseSuite struct {
	suite.Suite
	repo    *mocks.UserRepository
	jwtSvc  *mocks.JWTService
	passSvc *mocks.PasswordService
	usecase domain.UserUsecase
}

func (suite *UserUsecaseSuite) SetupTest() {
	suite.repo = new(mocks.UserRepository)
	suite.jwtSvc = new(mocks.JWTService)
	suite.passSvc = new(mocks.PasswordService)
	suite.usecase = usecases.NewUserUsecase(suite.repo, suite.jwtSvc, suite.passSvc, 2*time.Second)
}

func (suite *UserUsecaseSuite) TestRegister_Success() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
		Role:     domain.RegularUser,
	}

	suite.repo.On("GetUser", mock.Anything, user.Username).Return(nil, mongo.ErrNoDocuments)
	suite.repo.On("GetNewUserID", mock.Anything).Return(42, nil)
	suite.passSvc.On("HashPassword", user.Password).Return("hashedpass", nil)
	suite.repo.On("Register", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Username == user.Username && u.Password == "hashedpass" && u.Role == user.Role && u.ID == 42
	})).Return(nil)

	err := suite.usecase.Register(context.Background(), user)
	suite.NoError(err)
	suite.repo.AssertExpectations(suite.T())
	suite.passSvc.AssertExpectations(suite.T())
}

func (suite *UserUsecaseSuite) TestRegister_EmptyFields() {
	user := domain.User{
		Username: "",
		Password: "",
		Role:     domain.RegularUser,
	}

	err := suite.usecase.Register(context.Background(), user)
	suite.Error(err)
	suite.IsType(&domain.BadRequestError{}, err)
}

func (suite *UserUsecaseSuite) TestRegister_InvalidRole() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
		Role:     "invalid_role",
	}

	err := suite.usecase.Register(context.Background(), user)
	suite.Error(err)
	suite.IsType(&domain.BadRequestError{}, err)
}

func (suite *UserUsecaseSuite) TestRegister_UsernameTaken() {
	user := domain.User{
		Username: "takenuser",
		Password: "password123",
		Role:     domain.RegularUser,
	}

	suite.repo.On("GetUser", mock.Anything, user.Username).Return(&domain.User{}, nil)

	err := suite.usecase.Register(context.Background(), user)
	suite.Error(err)
	suite.IsType(&domain.BadRequestError{}, err)
}

func (suite *UserUsecaseSuite) TestRegister_HashPasswordError() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
		Role:     domain.RegularUser,
	}

	suite.repo.On("GetUser", mock.Anything, user.Username).Return(nil, mongo.ErrNoDocuments)
	suite.repo.On("GetNewUserID", mock.Anything).Return(0, nil)
	suite.passSvc.On("HashPassword", user.Password).Return("", errors.New("hash error"))

	err := suite.usecase.Register(context.Background(), user)
	suite.Error(err)
}

func (suite *UserUsecaseSuite) TestRegister_RepoRegisterError() {
	user := domain.User{
		Username: "testuser",
		Password: "password123",
		Role:     domain.RegularUser,
	}

	suite.repo.On("GetUser", mock.Anything, user.Username).Return(nil, mongo.ErrNoDocuments)
	suite.repo.On("GetNewUserID", mock.Anything).Return(1, nil)
	suite.passSvc.On("HashPassword", user.Password).Return("hashedpass", nil)
	suite.repo.On("Register", mock.Anything, mock.Anything).Return(errors.New("db error"))

	err := suite.usecase.Register(context.Background(), user)
	suite.Error(err)
}

func (suite *UserUsecaseSuite) TestLogin_Success() {
	login := domain.Login{
		Username: "user1",
		Password: "pass1",
	}

	existingUser := &domain.User{
		Username: "user1",
		Password: "hashedpass",
		Role:     domain.RegularUser,
	}

	suite.repo.On("GetUser", mock.Anything, login.Username).Return(existingUser, nil)
	suite.passSvc.On("ComparePassword", existingUser.Password, login.Password).Return(nil)
	suite.jwtSvc.On("GenerateToken", existingUser.Username, string(existingUser.Role)).Return("token123", nil)

	token, err := suite.usecase.Login(context.Background(), login)
	suite.NoError(err)
	suite.Equal("token123", token)
}

func (suite *UserUsecaseSuite) TestLogin_EmptyFields() {
	login := domain.Login{
		Username: "",
		Password: "",
	}

	token, err := suite.usecase.Login(context.Background(), login)
	suite.Error(err)
	suite.Empty(token)
	suite.IsType(&domain.BadRequestError{}, err)
}

func (suite *UserUsecaseSuite) TestLogin_IncorrectUserOrPassword() {
	login := domain.Login{
		Username: "user1",
		Password: "wrongpass",
	}

	existingUser := &domain.User{
		Username: "user1",
		Password: "hashedpass",
		Role:     domain.RegularUser,
	}

	// Simulate user not found
	suite.repo.On("GetUser", mock.Anything, login.Username).Return(nil, mongo.ErrNoDocuments)
	token, err := suite.usecase.Login(context.Background(), login)
	suite.Error(err)
	suite.Empty(token)

	// Simulate password mismatch
	suite.repo.ExpectedCalls = nil
	suite.passSvc.ExpectedCalls = nil
	suite.repo.On("GetUser", mock.Anything, login.Username).Return(existingUser, nil)
	suite.passSvc.On("ComparePassword", existingUser.Password, login.Password).Return(errors.New("password mismatch"))

	token, err = suite.usecase.Login(context.Background(), login)
	suite.Error(err)
	suite.Empty(token)
}

func (suite *UserUsecaseSuite) TestLogin_GenerateTokenError() {
	login := domain.Login{
		Username: "user1",
		Password: "pass1",
	}

	existingUser := &domain.User{
		Username: "user1",
		Password: "hashedpass",
		Role:     domain.RegularUser,
	}

	suite.repo.On("GetUser", mock.Anything, login.Username).Return(existingUser, nil)
	suite.passSvc.On("ComparePassword", existingUser.Password, login.Password).Return(nil)
	suite.jwtSvc.On("GenerateToken", existingUser.Username, string(existingUser.Role)).Return("", errors.New("token error"))

	token, err := suite.usecase.Login(context.Background(), login)
	suite.Error(err)
	suite.Empty(token)
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}
