package tests

import (
	"bytes"
	"encoding/json"
	"io"
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

type taskControllerSuite struct {
	suite.Suite
	usecase       *mocks.TaskUsecase
	controller    *task_controllers.TaskController
	server        *httptest.Server
}

func injectUser(role, username string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("role", role)
		c.Set("username", username)
		c.Next()
	}
}

func (suite *taskControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.usecase = new(mocks.TaskUsecase)
	suite.controller = task_controllers.NewTaskController(suite.usecase)

	router := gin.Default()
	// Inject username/role for all requests, can override per test
	router.Use(injectUser(string(domain.Admin), "admin"))

	router.GET("/tasks", suite.controller.GetTasks)
	router.GET("/tasks/:id", suite.controller.GetTask)
	router.POST("/tasks", suite.controller.PostTask)
	router.PUT("/tasks/:id", suite.controller.UpdateTask)
	router.DELETE("/tasks/:id", suite.controller.DeleteTask)

	suite.server = httptest.NewServer(router)
}

func (suite *taskControllerSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *taskControllerSuite) TestGetTasks_Success() {
    expectedTasks := []*domain.Task{
        {ID: 100, Title: "Title 1", Description: "Test title", DueDate: "2025-7-7", Status: domain.Completed, CreatedBy: "admin"},
    }
    
    // Be explicit about expected arguments instead of using mock.Anything
    suite.usecase.On("GetAllTasks", mock.Anything, "admin", string(domain.Admin)).Return(expectedTasks, nil).Once()

    req, _ := http.NewRequest("GET", suite.server.URL+"/tasks", nil)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    suite.NoError(err)
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    var got []*domain.Task
    json.Unmarshal(body, &got)

    suite.Equal(http.StatusOK, resp.StatusCode)
    suite.Equal(expectedTasks, got)
    suite.usecase.AssertExpectations(suite.T())
}

func (suite *taskControllerSuite) TestGetTasks_Error() {
	suite.usecase.On("GetAllTasks", mock.Anything, "admin", string(domain.Admin)).Return(nil, &domain.BadRequestError{Reason: "fail"}).Once()

	req, _ := http.NewRequest("GET", suite.server.URL+"/tasks", nil)
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *taskControllerSuite) TestGetTask_Success() {
	task := domain.Task{ID: 2, Title: "Task2", CreatedBy: "admin"}
	suite.usecase.On("GetTask", mock.Anything, "2", "admin", string(domain.Admin)).Return(task, nil)

	req, _ := http.NewRequest("GET", suite.server.URL+"/tasks/2", nil)
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var got domain.Task
	json.NewDecoder(resp.Body).Decode(&got)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal(task, got)
	suite.usecase.AssertCalled(suite.T(), "GetTask", mock.Anything, "2", "admin", string(domain.Admin))
}

func (suite *taskControllerSuite) TestGetTask_NotFound() {
	suite.usecase.On("GetTask", mock.Anything, "404", "admin", string(domain.Admin)).Return(domain.Task{}, &domain.NotFoundError{Resource: "Task", ID: "404"})

	req, _ := http.NewRequest("GET", suite.server.URL+"/tasks/404", nil)
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *taskControllerSuite) TestPostTask_Success() {
	task := domain.Task{Title: "New", Description: "desc", DueDate: "2025-01-01", Status: domain.Pending}
	created := task
	created.ID = 10
	created.CreatedBy = "admin"

	suite.usecase.On("AddTask", mock.Anything, "admin", task).Return(created, nil)

	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", suite.server.URL+"/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var got domain.Task
	json.NewDecoder(resp.Body).Decode(&got)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	suite.Equal(created, got)
	suite.usecase.AssertCalled(suite.T(), "AddTask", mock.Anything, "admin", task)
}

func (suite *taskControllerSuite) TestPostTask_InvalidJSON() {
	req, _ := http.NewRequest("POST", suite.server.URL+"/tasks", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *taskControllerSuite) TestUpdateTask_Success() {
	updated := domain.Task{Title: "Updated", Description: "desc", DueDate: "2025-01-01", Status: domain.Completed}
	final := updated
	final.ID = 3
	final.CreatedBy = "admin"

	suite.usecase.On("UpdateTask", mock.Anything, "3", "admin", string(domain.Admin), updated).Return(final, nil)

	body, _ := json.Marshal(updated)
	req, _ := http.NewRequest("PUT", suite.server.URL+"/tasks/3", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var got domain.Task
	json.NewDecoder(resp.Body).Decode(&got)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal(final, got)
	suite.usecase.AssertCalled(suite.T(), "UpdateTask", mock.Anything, "3", "admin", string(domain.Admin), updated)
}

func (suite *taskControllerSuite) TestUpdateTask_InvalidJSON() {
	req, _ := http.NewRequest("PUT", suite.server.URL+"/tasks/3", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *taskControllerSuite) TestDeleteTask_Success() {
	suite.usecase.On("DeleteTask", mock.Anything, "5", "admin", string(domain.Admin)).Return(nil)

	req, _ := http.NewRequest("DELETE", suite.server.URL+"/tasks/5", nil)
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusNoContent, resp.StatusCode)
	suite.usecase.AssertCalled(suite.T(), "DeleteTask", mock.Anything, "5", "admin", string(domain.Admin))
}

func (suite *taskControllerSuite) TestDeleteTask_NotFound() {
	suite.usecase.On("DeleteTask", mock.Anything, "404", "admin", string(domain.Admin)).Return(&domain.NotFoundError{Resource: "Task", ID: "404"})

	req, _ := http.NewRequest("DELETE", suite.server.URL+"/tasks/404", nil)
	resp, err := http.DefaultClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestTaskControllerSuite(t *testing.T) {
	suite.Run(t, new(taskControllerSuite))
}
