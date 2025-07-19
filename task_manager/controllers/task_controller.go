package task_controllers

import (
	"net/http"
	"task_manager/customError"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")

	tasks, err := data.GetAllTasks(username, role)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func GetATask(c *gin.Context) {
	taskID := c.Param("id")
	username := c.GetString("username")
	role := c.GetString("role") 

	task, err := data.GetTask(taskID, username, role)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func UpdateATask(c *gin.Context) {
	id := c.Param("id")
	username := c.GetString("username")
	role := c.GetString("role")

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		errorHandler(c, &customError.BadRequestError{Reason: "Invalid JSON"})
		return
	}

	task, err := data.UpdateTask(id, username, role, updatedTask)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func PostTask(c *gin.Context) {
	username := c.GetString("username")

	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		errorHandler(c, &customError.BadRequestError{Reason: "Invalid JSON"})
		return
	}

	task, err := data.AddATask(username, newTask)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, task)
}

func DeleteATask(c *gin.Context) {
	id := c.Param("id")
	username := c.GetString("username")
	role := c.GetString("role")

	err := data.DeleteTask(id, username, role)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func errorHandler(c *gin.Context, err error) {
	switch err.(type) {
	case *customError.NotFoundError:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case *customError.BadRequestError:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
	}
}