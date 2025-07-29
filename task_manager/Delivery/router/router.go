package router

import (
	task_controllers "task_manager/Delivery/controllers"
	domain "task_manager/Domain"
	infrastructure "task_manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func Setup(uc *task_controllers.UserController, tc *task_controllers.TaskController, engine *gin.Engine, jwtService domain.JWTService) {
	

	// ============ Public Routes ============
	publicRouter := engine.Group("")
	NewUserRouter(uc, publicRouter)

	// ============ Protected Routes ============
	protectedRouter := engine.Group("")
	protectedRouter.Use(infrastructure.NewAuthMiddleware(jwtService))
	NewTaskRouter(tc, protectedRouter)
}

func NewUserRouter(uc *task_controllers.UserController, group *gin.RouterGroup) {
	group.POST("/register", uc.RegisterUser)
	group.POST("/login", uc.LoginUser)
}

func NewTaskRouter(tc *task_controllers.TaskController, group *gin.RouterGroup) {
	group.GET("/tasks", tc.GetTask)
	group.GET("/tasks/:id", tc.GetTask)
	group.PUT("/tasks/:id", tc.UpdateTask)
	group.DELETE("/tasks/:id", tc.DeleteTask)
	group.POST("/tasks", tc.PostTask)
}
