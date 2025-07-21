package data

import (
	"context"
	"log"
	"os"
	"strconv"
	"task_manager/customError"
	"task_manager/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var taskCollection *mongo.Collection
var userCollection *mongo.Collection
var Client *mongo.Client
var JwtSecret = []byte("jwtSecret")

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

	taskCollection = Client.Database(dbName).Collection("tasks")
	userCollection = Client.Database(dbName).Collection(("users"))

	count, err := taskCollection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}

	adminCount, err := userCollection.CountDocuments(context.TODO(), bson.D{{Key: "role", Value: "Admin"}})
	if err != nil{
		return err
	}

	hashedSeedPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	if adminCount == 0{
		userSeed := models.User{
			ID: 1,
			Username: "Mr. Admin",
			Password: string(hashedSeedPassword),
			Role: models.Admin,
		}

		_, err = userCollection.InsertOne(context.TODO(), userSeed)
		if err != nil{
			return err
		}
	}

	if count == 0{
		seed := models.Task{
			ID: 1,
			Title: "Learn Go", 
			Description: "Practice structs and interfaces", 
			DueDate: "2025-08-01", 
			Status: models.Pending,
			CreatedBy: "Admin",
		}
		
		_, err = taskCollection.InsertOne(context.TODO(), seed)
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

func GetAllTasks(username string, role string) ([]*models.Task, error) {
	var allTasks []*models.Task

	var filter interface{} 
	if role == string(models.Admin){
		filter = bson.D{{}}
	} else {
		filter = bson.D{{Key:"createdby",Value: username}}
	}

	cursor, err := taskCollection.Find(context.TODO(),filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			log.Printf("Error decoding task: %v", err)
			continue 
		}
		allTasks = append(allTasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	println(allTasks)
	return allTasks, nil
}

func GetTask(taskId, username, role string) (models.Task, error) {
	taskID, err := strconv.Atoi(taskId)
	if err != nil {
		return models.Task{}, &customError.BadRequestError{Reason: "Invalid format of ID!"}
	}

	var taskFilter bson.D
	if role == string(models.Admin) {
		taskFilter = bson.D{{Key: "id", Value: taskID}} 
	} else {
		taskFilter = bson.D{{Key: "id", Value: taskID}, {Key: "createdby", Value: username}}
	}

	var task models.Task
	err = taskCollection.FindOne(context.TODO(), taskFilter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, &customError.NotFoundError{ID: taskID}
		}
		return models.Task{}, err
	}

	return task, nil
}

func UpdateTask(taskId, username, role string, updatedTask models.Task) (models.Task, error) {
	taskID, err := strconv.Atoi(taskId)
	if err != nil {
		return models.Task{}, &customError.BadRequestError{Reason: "Invalid format of ID!"}
	}

	if updatedTask.Title == "" || updatedTask.Description == "" || updatedTask.DueDate == ""{
		return models.Task{}, &customError.BadRequestError{Reason: "Fields can not be empty!"}
	}

	if updatedTask.Status != models.Pending && updatedTask.Status != models.Completed {
		return models.Task{}, &customError.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}

	var taskFilter bson.D
	if role == string(models.Admin){
		taskFilter = bson.D{{Key: "id", Value: taskID}}
	} else {
		taskFilter = bson.D{{Key: "id", Value: taskID}, {Key: "createdby", Value: username}}
	}

	// check if task exists
	var oldTask models.Task
	err = taskCollection.FindOne(context.TODO(), taskFilter).Decode(&oldTask)
	if err != nil{
		if err == mongo.ErrNoDocuments{
			return models.Task{}, &customError.NotFoundError{ID: taskID}
		}
		return models.Task{}, err
	}

	updatedTask.ID = oldTask.ID
	updatedTask.CreatedBy = oldTask.CreatedBy

	_, err = taskCollection.ReplaceOne(context.TODO(), taskFilter, updatedTask)
	if err != nil{
		return models.Task{}, err
	}

	err = taskCollection.FindOne(context.TODO(), taskFilter).Decode(&updatedTask)
	if err != nil {
	    return models.Task{}, err
	}

	return updatedTask, nil
}

func AddATask(username string, task models.Task) (models.Task, error) {
	if task.Title == "" || task.Description == "" || task.DueDate == "" {
		return models.Task{}, &customError.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	if task.Status != models.Pending && task.Status != models.Completed {
		return models.Task{}, &customError.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}

	// Find the highest current ID
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var lastTask models.Task
	err := taskCollection.FindOne(context.TODO(), bson.M{}, opts).Decode(&lastTask)
	if err != nil && err != mongo.ErrNoDocuments {
		return models.Task{}, err
	}

	// Set the new ID
	if err == mongo.ErrNoDocuments {
		task.ID = 1
	} else {
		task.ID = lastTask.ID + 1
	}

	task.CreatedBy = username

	_, err = taskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func DeleteTask(id, username, role string) (error){
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return &customError.BadRequestError{Reason: "Invalid format of ID!"} 
	}
	
	var taskFilter bson.D
	if role == string(models.Admin){
		taskFilter = bson.D{{Key:"id", Value: taskID}}
	} else {
		taskFilter = bson.D{{Key:"id", Value: taskID}, {Key:"createdby", Value: username}}
	}
	
	deleteResult, err := taskCollection.DeleteOne(context.TODO(), taskFilter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0{
		return &customError.NotFoundError{ID: taskID}
	}

	return nil
}