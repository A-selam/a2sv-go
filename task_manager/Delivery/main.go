package main

import (
	"context"
	"log"
	"os"
	task_controllers "task_manager/Delivery/controllers"
	"task_manager/Delivery/router"
	infrastructure "task_manager/Infrastructure"
	repositories "task_manager/Repositories"
	usecases "task_manager/Usecases"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
    // Load .env file
    err := godotenv.Load("../.env") 
    if err != nil {
        log.Println("Warning: .env file not found or failed to load", err)
    }

    // Retrieve environment variables
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("MONGODB_URI is not set")
    }

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        log.Fatal("DB_NAME is not set")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET is not set")
    }

    // Initialize MongoDB client
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().
        ApplyURI(mongoURI).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second)

    client, err := mongo.Connect(clientOptions)
    if err != nil {
        log.Fatal("MongoDB connection error:", err)
    }
    defer client.Disconnect(context.Background()) 

    // Verify connection
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    // Initialize database
    db := client.Database(dbName)

    // Initialize services and repositories
    timeOut := 30 * time.Second
    jwtService := infrastructure.NewJWTService(jwtSecret)
    passwordService := infrastructure.NewPasswordService()

    ur := repositories.NewUserRepositoryFromDB(db) 
    uu := usecases.NewUserUsecase(ur, jwtService, passwordService, timeOut)
    uc := task_controllers.NewUserController(uu)

    tr := repositories.NewTaskRepositoryFromDB(db) 
    tu := usecases.NewTaskUsecase(tr, timeOut)
    tc := task_controllers.NewTaskController(tu)

    // Set up Gin router
    engine := gin.Default()
    router.Setup(uc, tc, engine, jwtService)

    // Start server
    if err := engine.Run("localhost:3000"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}