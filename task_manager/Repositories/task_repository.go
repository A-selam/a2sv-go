package repositories

import (
	"context"
	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type taskRepository struct {
	database mongo.Database
	collection string
}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository{
	return &taskRepository{
		database: db,
		collection: collection,
	}
}

// // /////////////////////////////////////////////////////////////////////////////////////////////////
// func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
// 	collection := tr.database.Collection(tr.collection)

// 	_, err := collection.InsertOne(c, task)

// 	return err
// }

// func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
// 	collection := tr.database.Collection(tr.collection)

// 	var tasks []domain.Task

// 	idHex, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return tasks, err
// 	}

// 	cursor, err := collection.Find(c, bson.M{"userID": idHex})
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = cursor.All(c, &tasks)
// 	if tasks == nil {
// 		return []domain.Task{}, err
// 	}

// 	return tasks, err
// }
// /////////////////////////////////////////////////////////////////////////////////////////////////
func (r *taskRepository)GetAllTasks(c context.Context) ([]*domain.Task, error){
	collection := r.database.Collection(r.collection)

	var allTasks []*domain.Task

	cursor, err := collection.Find(c, bson.D{{}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(c)
	
	for cursor.Next(c){
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			continue // skip the hard to decode task
		}
		allTasks = append(allTasks, &task)
	}

	if err := cursor.Err(); err != nil{
		return nil, err
	}

	return allTasks, err
}

func (r *taskRepository)GetAllUserTasks(c context.Context, username string)([]*domain.Task, error){
	collection := r.database.Collection(r.collection)
	filter := bson.D{{Key: "created_by", Value: username}}

	var allTasks []*domain.Task

	cursor, err := collection.Find(c, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(c)
	
	for cursor.Next(c){
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			continue // skip the hard to decode task
		}
		allTasks = append(allTasks, &task)
	}

	if err := cursor.Err(); err != nil{
		return nil, err
	}

	return allTasks, err
}

func (r *taskRepository)GetTask(c context.Context, taskID int) (domain.Task, error){
	collection := r.database.Collection(r.collection)

	filter := bson.D{{Key: "id", Value: taskID}}
	
	var task domain.Task
	err := collection.FindOne(c, filter).Decode(&task)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *taskRepository)GetTaskByIDForUser(c context.Context, taskID int, username string) (domain.Task, error){
	collection := r.database.Collection(r.collection)

	filter := bson.D{{Key: "id", Value: taskID}, {Key: "created_by", Value: username}}

	var task domain.Task
	err := collection.FindOne(c, filter).Decode(&task)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *taskRepository)UpdateTask(c context.Context, taskID int, updatedTask domain.Task) (domain.Task, error){
	collection := r.database.Collection(r.collection)

	filter := bson.D{{Key: "id", Value: taskID}}
	_, err := collection.ReplaceOne(c, filter, updatedTask)
	if err != nil{
		return domain.Task{}, err
	}

	err = collection.FindOne(c, filter).Decode(&updatedTask)
	if err != nil {
		return domain.Task{}, err
	}

	return updatedTask, nil
}

func (r *taskRepository)AddTask(c context.Context, newTask domain.Task) (domain.Task, error){
	collection := r.database.Collection(r.collection)

	_, err := collection.InsertOne(c, newTask)
	if err != nil {
		return domain.Task{}, err
	}

	return newTask, nil
}

func (r *taskRepository)GetNewID(c context.Context)(int, error){
	collection := r.database.Collection(r.collection)

	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})

	var lastTask domain.Task
	err := collection.FindOne(c, bson.M{}, opts).Decode(&lastTask)
	if err != nil && err != mongo.ErrNoDocuments{
		return -1, err
	}

	if err == mongo.ErrNoDocuments{
		return 1, nil
	} 

	return lastTask.ID + 1, nil
}

func (r *taskRepository)DeleteTask(c context.Context, taskID int) error{
	collection := r.database.Collection(r.collection)
		
	filter := bson.D{{Key:"id", Value: taskID}}
	deleteResult, err := collection.DeleteOne(c, filter)
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0{
		return &domain.NotFoundError{}
	}

	return nil
}