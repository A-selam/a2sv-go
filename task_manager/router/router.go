package router

import (
	task_controllers "task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine{
	router := gin.Default()

	router.GET("/tasks", task_controllers.GetTasks)
	router.GET("/tasks/:id", task_controllers.GetATask)
	router.PUT("/tasks/:id", task_controllers.UpdateATask)
	router.DELETE("/tasks/:id", task_controllers.DeleteATask)
	router.POST("/tasks", task_controllers.PostTask)

	// user routes
	router.POST("/register", task_controllers.RegisterUser)
	router.POST("/login", task_controllers.LoginUser)

	return router 
}