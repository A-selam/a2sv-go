package router

import (
	task_controllers "task_manager/controllers"

	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine{
	router := gin.Default()

	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	authorized.GET("/tasks", task_controllers.GetTasks)
	authorized.GET("/tasks/:id", task_controllers.GetATask)
	authorized.PUT("/tasks/:id", task_controllers.UpdateATask)
	authorized.DELETE("/tasks/:id", task_controllers.DeleteATask)
	authorized.POST("/tasks", task_controllers.PostTask)

	// user routes
	router.POST("/register", task_controllers.RegisterUser)
	router.POST("/login", task_controllers.LoginUser)

	return router 
}