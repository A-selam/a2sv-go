package task_controllers

import (
	"net/http"
	domain "task_manager/Domain"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUseCase domain.TaskUsecase
}

func NewTaskController(tu domain.TaskUsecase) *TaskController{
	return &TaskController{
		TaskUseCase: tu,
	}
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")
	
	tasks, err := tc.TaskUseCase.GetAllTasks(c, username, role)
	if err != nil {
		errorHandler(c, err)
		return
	}
	
	c.IndentedJSON(http.StatusOK, tasks)
}

func (tc *TaskController) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	username := c.GetString("username")
	role := c.GetString("role") 

	task, err := tc.TaskUseCase.GetTask(c, taskID, username, role)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	username := c.GetString("username")
	role := c.GetString("role")

	var updatedTask domain.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		errorHandler(c, &domain.BadRequestError{Reason: "Invalid JSON"})
		return
	}

	task, err := tc.TaskUseCase.UpdateTask(c, id, username, role, updatedTask)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (tc *TaskController) PostTask(c *gin.Context) {
	username := c.GetString("username")

	var newTask domain.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		errorHandler(c, &domain.BadRequestError{Reason: "Invalid JSON"})
		return
	}

	task, err := tc.TaskUseCase.AddTask(c, username, newTask)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	username := c.GetString("username")
	role := c.GetString("role")

	err := tc.TaskUseCase.DeleteTask(c, id, username, role)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func errorHandler(c *gin.Context, err error) {
	switch err.(type) {
	case *domain.NotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case *domain.BadRequestError:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case *domain.UnauthorizedError:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case *domain.ForbiddenError:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
	}
}