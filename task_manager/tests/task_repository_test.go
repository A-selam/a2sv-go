package tests

import (
	"context"
	"log"
	"os"
	domain "task_manager/Domain"
	repositories "task_manager/Repositories"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type taskRepositorySuite struct {
	suite.Suite

	repository domain.TaskRepository // the actual repo you’ll test
	db         *mongo.Database       // raw mongo DB to clean collections
}

func (suite *taskRepositorySuite) SetupSuite() {
	// Load env values manually for test (or read from .env if you prefer)
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: .env file not found or failed to load", err)
	}
	
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(nil, options.Client().ApplyURI(mongoURI))
	suite.Require().NoError(err)

	db := client.Database("task_manager_test")

	// Assign fields to suite
	suite.db = db
	suite.repository = repositories.NewTaskRepositoryFromDB(db)
}

func (suite *taskRepositorySuite) TearDownTest() {
	// Drop collection after each test (MongoDB version of truncating)
	// err := suite.db.Collection("tasks").Drop(context.Background())
	err := suite.db.Drop(context.TODO())
	suite.Require().NoError(err)
}

// ✅ Example test: AddTask()
func (suite *taskRepositorySuite) TestAddTask() {
	task := domain.Task{
		ID:          1,
		Title:       "Test",
		Description: "This is a test",
		Status:      domain.Pending,
		CreatedBy:   "tester",
	}

	added, err := suite.repository.AddTask(context.Background(), task)

	suite.NoError(err)
	suite.Equal(task.ID, added.ID)
	suite.Equal(task.Title, added.Title)
	suite.Equal(task.CreatedBy, added.CreatedBy)
}

func (suite *taskRepositorySuite) TestGetTask(){
	task := domain.Task{
		ID: 2,
		Title: "test_task",
		Description: "This is test task",
		Status: domain.Completed,
		CreatedBy: "tester",
	}

	_, err := suite.repository.AddTask(context.Background(), task)
	suite.NoError(err)

	found, err := suite.repository.GetTask(context.Background(), 2)
	suite.NoError(err)
	suite.Equal(found.Title, task.Title)
	suite.Equal(found.Description, task.Description)
	suite.Equal(found.Status, task.Status)
	suite.Equal(found.CreatedBy, task.CreatedBy)
}

func (suite *taskRepositorySuite) TestGetTask_NotFound() {
	_, err := suite.repository.GetTask(context.Background(), 404)
	suite.Error(err)
}

func (suite *taskRepositorySuite) TestGetAllUserTasks_ReturnsOnlyUserTasks() {
	// Add two tasks, one by "user1", one by "user2"
	suite.repository.AddTask(context.Background(), domain.Task{
		ID: 1, Title: "T1", CreatedBy: "user1", Status: domain.Pending,
	})
	suite.repository.AddTask(context.Background(), domain.Task{
		ID: 2, Title: "T2", CreatedBy: "user2", Status: domain.Completed,
	})

	// Now fetch tasks only for "user1"
	user1Tasks, err := suite.repository.GetAllUserTasks(context.Background(), "user1")
	suite.NoError(err)
	suite.Len(user1Tasks, 1)
	suite.Equal("user1", user1Tasks[0].CreatedBy)
}

func (suite *taskRepositorySuite) TestUpdateTask_UpdatesCorrectly() {
	task := domain.Task{
		ID: 10, Title: "Original", Description: "Initial", Status: domain.Pending, CreatedBy: "user",
	}
	suite.repository.AddTask(context.Background(), task)

	// Update title and description
	task.Title = "Updated"
	task.Description = "Now updated"
	updated, err := suite.repository.UpdateTask(context.Background(), task.ID, task)

	suite.NoError(err)
	suite.Equal("Updated", updated.Title)
	suite.Equal("Now updated", updated.Description)
}

func (suite *taskRepositorySuite) TestDeleteTask_DeletesCorrectly() {
	task := domain.Task{ID: 99, Title: "To Delete", CreatedBy: "test", Status: domain.Pending}
	suite.repository.AddTask(context.Background(), task)

	// Delete
	err := suite.repository.DeleteTask(context.Background(), task.ID)
	suite.NoError(err)

	// Try getting it (should fail)
	_, err = suite.repository.GetTask(context.Background(), task.ID)
	suite.Error(err) // Expect an error because it's deleted
}


// ✅ Run the test suite
func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(taskRepositorySuite))
}
