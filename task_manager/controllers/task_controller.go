package task_controllers

import (
	"net/http"
	"task_manager/customError"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context){
	tasks := data.GetAllTasks()
	c.IndentedJSON(http.StatusOK, tasks)
}

func GetATask(c *gin.Context){
	id := c.Param("id")
	task, err := data.GetTask(id)
	
	if err != nil{
		errorHandler(c, err)
		return
	}
	
	c.IndentedJSON(http.StatusOK, task)
}

func UpdateATask(c *gin.Context){
	id := c.Param("id")

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask) ; err != nil{
		errorHandler(c, &customError.BadRequestError{Reason: "Invalid JSON"})
		return 
	}

	_, err := data.UpdateTask(id, updatedTask)
	if err != nil {
		errorHandler(c, err)
		return 
	}

	c.IndentedJSON(http.StatusOK, updatedTask)
}

func DeleteATask(c *gin.Context){
	id := c.Param("id")

	err := data.DeleteTask(id)
	if err != nil{
		errorHandler(c, err)
		return 
	} 

	c.Status(http.StatusNoContent)
}

func PostTask(c *gin.Context){
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil{
		errorHandler(c, &customError.BadRequestError{Reason: "Invalid JSON"})
		return
	}
	
	task, err := data.AddATask(newTask)
	if err != nil {
		errorHandler(c, err)
		return
	}
	
	c.IndentedJSON(http.StatusCreated, task)
}

func errorHandler(c *gin.Context, err error){
	switch err.(type){
		case *customError.NotFoundError:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error":err.Error()})
		case *customError.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		}
}