package tests

import (
	"context"
	domain "task_manager/Domain"
	usecases "task_manager/Usecases"
	"task_manager/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskUsecaseSuite struct {
	suite.Suite
	repo    *mocks.TaskRepository
	usecase domain.TaskUsecase
}

func (suite *TaskUsecaseSuite) SetupTest() {
	suite.repo = new(mocks.TaskRepository)
	timeout := 2 * time.Second
	suite.usecase = usecases.NewTaskUsecase(suite.repo, timeout)
}

func (suite *TaskUsecaseSuite) TestGetAllTasks_RegularUser_Success() {
	tasks := []*domain.Task{
		{ID: 100, Title: "Title 1", Description: "Test title", DueDate: "2025-7-7", Status: domain.Completed, CreatedBy: "John Doe"},
	}

	suite.repo.Mock.On("GetAllUserTasks", mock.Anything, "John Doe").Return(tasks, nil)

	result, err := suite.usecase.GetAllTasks(context.Background(), "John Doe", string(domain.RegularUser))

	suite.Nil(err)
	suite.Equal(tasks, result)
	suite.repo.AssertCalled(suite.T(), "GetAllUserTasks", mock.Anything, "John Doe")
}

func (suite *TaskUsecaseSuite) TestGetAllTasks_AdminUser_Success() {
	taskAdmin := []*domain.Task{
		{ID: 100, Title: "Title 1", Description: "Test title", DueDate: "2025-7-7", Status: domain.Completed, CreatedBy: "John Doe"},
		{ID: 102, Title: "Title 2", Description: "Test title 2", DueDate: "2025-3-12", Status: domain.Pending, CreatedBy: "Jane Doe"},
	}

	suite.repo.Mock.On("GetAllTasks", mock.Anything).Return(taskAdmin, nil)

	result, err := suite.usecase.GetAllTasks(context.Background(), "admin", string(domain.Admin))

	suite.Nil(err)
	suite.Equal(taskAdmin, result)
	suite.repo.AssertCalled(suite.T(), "GetAllTasks", mock.Anything)
}

func (suite *TaskUsecaseSuite) TestGetTask_Admin_Success() {
	task := domain.Task{ID: 1, Title: "Admin Task", CreatedBy: "admin"}
	suite.repo.On("GetTask", mock.Anything, 1).Return(task, nil)

	result, err := suite.usecase.GetTask(context.Background(), "1", "admin", string(domain.Admin))
	
	suite.Nil(err)
	suite.Equal(task, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TaskUsecaseSuite) TestGetTask_Regular_Success() {
	task := domain.Task{ID: 2, Title: "User Task", CreatedBy: "john"}
	suite.repo.On("GetTaskByIDForUser", mock.Anything, 2, "john").Return(task, nil)

	result, err := suite.usecase.GetTask(context.Background(), "2", "john", string(domain.RegularUser))

	suite.Nil(err)
	suite.Equal(task, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TaskUsecaseSuite) TestUpdateTask_Admin_Success() {
	oldTask := domain.Task{ID: 1, CreatedBy: "someone"}
	updated := domain.Task{
		Title:       "Updated",
		Description: "New Desc",
		DueDate:     "2025-12-12",
		Status:      domain.Completed,
	}

	suite.repo.On("GetTask", mock.Anything, 1).Return(oldTask, nil)
	finalTask := updated
	finalTask.ID = oldTask.ID
	finalTask.CreatedBy = oldTask.CreatedBy
	suite.repo.On("UpdateTask", mock.Anything, 1, finalTask).Return(finalTask, nil)

	result, err := suite.usecase.UpdateTask(context.Background(), "1", "admin", string(domain.Admin), updated)

	suite.Nil(err)
	suite.Equal(finalTask, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TaskUsecaseSuite) TestUpdateTask_InvalidStatus() {
	task := domain.Task{
		Title:       "A",
		Description: "B",
		DueDate:     "2025-01-01",
		Status:      "invalid",
	}
	_, err := suite.usecase.UpdateTask(context.Background(), "1", "john", string(domain.RegularUser), task)

	suite.Error(err)
	_, isBadReq := err.(*domain.BadRequestError)
	suite.True(isBadReq)
}

func (suite *TaskUsecaseSuite) TestAddTask_Success() {
	task := domain.Task{
		ID:          100,
		Title:       "Title",
		Description: "Description",
		DueDate:     "2025-08-01",
		Status:      domain.Pending,
		CreatedBy:   "John Doe",
	}

	username := "John Doe"
	newID := 100
	expectedTask := task

	suite.repo.On("GetNewID", mock.Anything).Return(newID, nil)
	suite.repo.On("AddTask", mock.Anything, expectedTask).Return(expectedTask, nil)

	result, err := suite.usecase.AddTask(context.Background(), username, task)

	suite.Nil(err)
	suite.Equal(expectedTask, result)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TaskUsecaseSuite) TestAddTask_ValidationError() {
	task := domain.Task{
		Title:       "",
		Description: "something",
		DueDate:     "2025-08-01",
		Status:      domain.Pending,
	}
	username := "Jane Doe"

	result, err := suite.usecase.AddTask(context.Background(), username, task)

	suite.Error(err)
	suite.Contains(err.Error(), "Fields cannot be empty")
	suite.Equal(domain.Task{}, result)
}

func (suite *TaskUsecaseSuite) TestDeleteTask_Admin_Success() {
	task := domain.Task{ID: 5, CreatedBy: "john"}
	suite.repo.On("GetTask", mock.Anything, 5).Return(task, nil)
	suite.repo.On("DeleteTask", mock.Anything, 5).Return(nil)

	err := suite.usecase.DeleteTask(context.Background(), "5", "admin", string(domain.Admin))

	suite.Nil(err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *TaskUsecaseSuite) TestDeleteTask_RegularUserForbidden() {
	task := domain.Task{ID: 6, CreatedBy: "otherUser"}
	suite.repo.On("GetTask", mock.Anything, 6).Return(task, nil)

	err := suite.usecase.DeleteTask(context.Background(), "6", "john", string(domain.RegularUser))

	suite.Error(err)
	_, isUnauthorized := err.(*domain.UnauthorizedError)
	suite.True(isUnauthorized)
}

func (suite *TaskUsecaseSuite) TestDeleteTask_InvalidID() {
	err := suite.usecase.DeleteTask(context.Background(), "xyz", "john", string(domain.RegularUser))

	suite.Error(err)
	_, isBadReq := err.(*domain.BadRequestError)
	suite.True(isBadReq)
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}