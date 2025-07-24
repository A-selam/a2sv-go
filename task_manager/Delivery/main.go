package main

import (
	"context"
	"log"
	"os"
	"time"

	"task_manager/Delivery/router"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)

	// Setup Gin engine
	engine := gin.Default()

	// Pass jwtSecret to router so it can pass to JWT service (you'll need to update your router.Setup accordingly)
	router.Setup(30*time.Second, *db, engine, jwtSecret)

	if err := engine.Run("localhost:3000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

