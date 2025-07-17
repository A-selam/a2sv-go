package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"task_manager/customError"
	"task_manager/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var collection *mongo.Collection
var Client *mongo.Client

func InitMongo() error{
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Failed to load env: ", err)
		return err
	}

	connectionString := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)

	Client, err = mongo.Connect(clientOptions)
	if err != nil{
		log.Fatal(err)
		return err
	}
	
	err = Client.Ping(context.TODO(), nil)
	if err != nil{
		log.Fatal(err)
		return err
	}

	collection = Client.Database(dbName).Collection("tasks")

	count, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}

	if count == 0{
		seed := models.Task{
			ID: 1,
			Title: "Learn Go", 
			Description: "Practice structs and interfaces", 
			DueDate: "2025-08-01", 
			Status: models.Pending,
		}
		
		_, err = collection.InsertOne(context.TODO(), seed)
		if err != nil{
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func CloseMongo(){
	if Client != nil{
		_ = Client.Disconnect(context.TODO())
	}
}

func GetAllTasks() ([]*models.Task, error) {
	var allTasks []*models.Task

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	defer func() {
		if err := cursor.Close(context.TODO()); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	for cursor.Next(context.TODO()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			log.Printf("Error decoding task: %v", err)
			continue // Skip problematic documents but continue processing others
		}
		allTasks = append(allTasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return allTasks, nil
}

func GetTask(id string) (models.Task, error){
	taskID, err := strconv.Atoi(id)
	if err != nil{
		return models.Task{}, &customError.BadRequestError{Reason:"Invalid format of ID!"}
	}

	taskFilter := bson.D{{Key: "id", Value: taskID}}
	var task models.Task

	err = collection.FindOne(context.TODO(), taskFilter).Decode(&task)
	if err != nil{
		if err == mongo.ErrNoDocuments{
			return models.Task{}, &customError.NotFoundError{ID: taskID}
		}
		fmt.Println("Failed to fetch single task")
		log.Fatal(err)
		return models.Task{}, err
	}
	return task, nil
}

func UpdateTask(id string, updatedTask models.Task) (models.Task, error) {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return models.Task{}, &customError.BadRequestError{Reason: "Invalid format of ID!"}
	}

	if updatedTask.Title == "" || updatedTask.Description == "" || updatedTask.DueDate == ""{
		return models.Task{}, &customError.BadRequestError{Reason: "Fields can not be empty!"}
	}

	if updatedTask.Status != models.Pending && updatedTask.Status != models.Completed {
		return models.Task{}, &customError.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}

	taskFilter := bson.D{{Key: "id", Value: taskID}}

	// check if task exists
	var oldTask models.Task
	err = collection.FindOne(context.TODO(), taskFilter).Decode(&oldTask)
	if err != nil{
		if err == mongo.ErrNoDocuments{
			return models.Task{}, &customError.NotFoundError{ID: taskID}
		}
		return models.Task{}, err
	}

	updatedTask.ID = oldTask.ID
	_, err = collection.ReplaceOne(context.TODO(), taskFilter, updatedTask)
	if err != nil{
		return models.Task{}, err
	}

	return updatedTask, nil
}

func DeleteTask(id string) (error){
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return &customError.BadRequestError{Reason: "Invalid format of ID!"} 
	}
	
	taskFilter := bson.D{{Key:"id", Value: taskID}}
	
	deleteResult, err := collection.DeleteOne(context.TODO(), taskFilter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0{
		return &customError.NotFoundError{ID: taskID}
	}

	return nil
}

func AddATask(task models.Task) (models.Task, error) {
	if task.Title == "" || task.Description == "" || task.DueDate == "" {
		return models.Task{}, &customError.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	if task.Status != models.Pending && task.Status != models.Completed {
		return models.Task{}, &customError.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}

	// Find the highest current ID
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var lastTask models.Task
	err := collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&lastTask)
	if err != nil && err != mongo.ErrNoDocuments {
		return models.Task{}, err
	}

	// Set the new ID
	if err == mongo.ErrNoDocuments {
		task.ID = 1
	} else {
		task.ID = lastTask.ID + 1
	}

	_, err = collection.InsertOne(context.TODO(), task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}