package router

import (
	task_controllers "task_manager/Delivery/controllers"
	domain "task_manager/Domain"
	infrastructure "task_manager/Infrastructure"
	repositories "task_manager/Repositories"
	usecases "task_manager/Usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Setup(timeout time.Duration, db mongo.Database, engine *gin.Engine, jwtSecret string) {
	// Create shared services
	jwtService := infrastructure.NewJWTService(jwtSecret)
	passwordService := infrastructure.NewPasswordService()

	// ============ Public Routes ============
	publicRouter := engine.Group("")
	NewUserRouter(timeout, db, jwtService, passwordService, publicRouter)

	// ============ Protected Routes ============
	protectedRouter := engine.Group("")
	protectedRouter.Use(infrastructure.NewAuthMiddleware(jwtService))
	NewTaskRouter(timeout, db, protectedRouter)
}

func NewUserRouter(timeout time.Duration, db mongo.Database, jwtService domain.JWTService, passwordService domain.PasswordService, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, "users")
	uc := &task_controllers.UserController{
		UserUsecase: usecases.NewUserUsecase(ur, jwtService, passwordService, timeout),
	}

	group.POST("/register", uc.RegisterUser)
	group.POST("/login", uc.LoginUser)
}

func NewTaskRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repositories.NewTaskRepository(db, "tasks")
	tc := &task_controllers.TaskController{
		TaskUseCase: usecases.NewTaskUsecase(tr, timeout),
	}

	group.GET("/tasks", tc.GetTasks)
	group.GET("/tasks/:id", tc.GetTask)
	group.PUT("/tasks/:id", tc.UpdateTask)
	group.DELETE("/tasks/:id", tc.DeleteTask)
	group.POST("/tasks", tc.PostTask)
}
