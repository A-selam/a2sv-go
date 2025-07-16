package data

import (
	"fmt"
	"strconv"
	"task_manager/customError"
	"task_manager/models"
)

var idCounter = 0
var tasks = map[int]models.Task{}

func init(){
	idCounter++
	tasks[idCounter] = models.Task{
		ID: idCounter,
		Title: "Learn Go", 
		Description: "Practice structs and interfaces", 
		DueDate: "2025-08-01", 
		Status: models.Pending,
	}
}

func GetAllTasks() []models.Task{
	fmt.Println(tasks)
	var allTask []models.Task
	for _, task := range tasks{
		allTask = append(allTask, task)
	}
	return allTask
}

func GetTask(id string) (models.Task, error){
	taskID, err := strconv.Atoi(id)

	if err != nil{
		return models.Task{}, &customError.BadRequestError{Reason:"Invalid format of ID!"}
	}

	task, ok := tasks[taskID]
	if !ok{
		return models.Task{}, &customError.NotFoundError{ID: taskID}
	}

	return task, nil
}

func UpdateTask(id string, updatedTask models.Task) (models.Task, error) {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return models.Task{}, &customError.BadRequestError{Reason: "Invalid format of ID!"}
	}

	_, ok := tasks[taskID]
	if !ok {
		return models.Task{}, &customError.NotFoundError{ID: taskID}
	}

	if updatedTask.Title == "" || updatedTask.Description == "" || updatedTask.DueDate == ""{
		return models.Task{}, &customError.BadRequestError{Reason: "Fields can not be empty!"}
	}

	if updatedTask.Status != models.Pending && updatedTask.Status != models.Completed {
		return models.Task{}, &customError.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}
	
	updatedTask.ID = taskID
	tasks[taskID] = updatedTask
	
	return updatedTask, nil
}

func DeleteTask(id string) (error){
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return &customError.BadRequestError{Reason: "Invalid format of ID!"} 
	}
	
	_, ok := tasks[taskID]
	if !ok {
		return &customError.NotFoundError{ID: taskID}
	}

	delete(tasks, taskID)
	return nil
}

func AddATask(task models.Task) (models.Task, error){
	idCounter++
	task.ID = idCounter

	if task.Title == "" || task.Description == "" || task.DueDate == ""{
		return models.Task{}, &customError.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	if task.Status != models.Pending && task.Status != models.Completed {
		return models.Task{}, &customError.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}

	tasks[task.ID] = task
	return task, nil
}