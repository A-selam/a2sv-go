package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	task_controllers "task_manager/Delivery/controllers"
	domain "task_manager/Domain"
	"task_manager/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userControllerSuite struct {
    suite.Suite
    usecase *mocks.UserUsecase
    server  *httptest.Server
}

func (suite *userControllerSuite) SetupTest() {
    gin.SetMode(gin.TestMode)
    suite.usecase = new(mocks.UserUsecase)
    controller := task_controllers.NewUserController(suite.usecase)

    router := gin.Default()
    router.POST("/register", controller.RegisterUser)
    router.POST("/login", controller.LoginUser)

    suite.server = httptest.NewServer(router)
}

func (suite *userControllerSuite) TearDownTest() {
    suite.server.Close()
}

func (suite *userControllerSuite) TestRegisterUser_Success() {
	user := domain.User{ID:10,Username: "test", Password: "pass", Role: domain.RegularUser}
    suite.usecase.On("Register", mock.Anything, user).Return(nil)

    body, _ := json.Marshal(user)
    resp, err := http.Post(suite.server.URL+"/register", "application/json", bytes.NewBuffer(body))
    suite.NoError(err)
    defer resp.Body.Close()

    var got map[string]string
    json.NewDecoder(resp.Body).Decode(&got)
    suite.Equal(http.StatusCreated, resp.StatusCode)
    suite.Equal("User registered successfully", got["message"])
    suite.usecase.AssertCalled(suite.T(), "Register", mock.Anything, user)
}

func (suite *userControllerSuite) TestRegisterUser_InvalidJSON() {
    resp, err := http.Post(suite.server.URL+"/register", "application/json", bytes.NewBuffer([]byte("{invalid")))
    suite.NoError(err)
    defer resp.Body.Close()
    suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *userControllerSuite) TestRegisterUser_Error() {
    user := domain.User{Username: "test", Password: "pass", Role: domain.RegularUser}
    suite.usecase.On("Register", mock.Anything, user).Return(errors.New("fail"))

    body, _ := json.Marshal(user)
    resp, err := http.Post(suite.server.URL+"/register", "application/json", bytes.NewBuffer(body))
    suite.NoError(err)
    defer resp.Body.Close()

    var got map[string]string
    json.NewDecoder(resp.Body).Decode(&got)
    suite.Equal(http.StatusBadRequest, resp.StatusCode)
    suite.Equal("fail", got["error"])
}

func (suite *userControllerSuite) TestLoginUser_Success() {
    login := domain.Login{Username: "test", Password: "pass"}
    suite.usecase.On("Login", mock.Anything, login).Return("token123", nil)

    body, _ := json.Marshal(login)
    resp, err := http.Post(suite.server.URL+"/login", "application/json", bytes.NewBuffer(body))
    suite.NoError(err)
    defer resp.Body.Close()

    var got map[string]string
    json.NewDecoder(resp.Body).Decode(&got)
    suite.Equal(http.StatusOK, resp.StatusCode)
    suite.Equal("token123", got["token"])
}

func (suite *userControllerSuite) TestLoginUser_InvalidJSON() {
    resp, err := http.Post(suite.server.URL+"/login", "application/json", bytes.NewBuffer([]byte("{invalid")))
    suite.NoError(err)
    defer resp.Body.Close()
    suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *userControllerSuite) TestLoginUser_Error() {
    login := domain.Login{Username: "test", Password: "pass"}
    suite.usecase.On("Login", mock.Anything, login).Return("", errors.New("unauthorized"))

    body, _ := json.Marshal(login)
    resp, err := http.Post(suite.server.URL+"/login", "application/json", bytes.NewBuffer(body))
    suite.NoError(err)
    defer resp.Body.Close()

    var got map[string]string
    json.NewDecoder(resp.Body).Decode(&got)
    suite.Equal(http.StatusUnauthorized, resp.StatusCode)
    suite.Equal("unauthorized", got["error"])
}

func TestUserControllerSuite(t *testing.T) {
    suite.Run(t, new(userControllerSuite))
}