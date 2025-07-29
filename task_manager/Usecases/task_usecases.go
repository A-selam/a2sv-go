package usecases

import (
	"context"
	"strconv"
	domain "task_manager/Domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase{
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (u *taskUsecase)GetAllTasks(c context.Context, username, role string)([]*domain.Task, error){
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	var tasks []*domain.Task
	var err error
	if role == string(domain.Admin){
		tasks, err = u.taskRepository.GetAllTasks(ctx)
	} else {
		tasks, err = u.taskRepository.GetAllUserTasks(ctx, username)
	}

	return tasks, err
}

func (u *taskUsecase)GetTask(c context.Context, taskID, username, role string)(domain.Task, error){
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	taskId, err := strconv.Atoi(taskID)
	if err != nil {
		return domain.Task{}, &domain.BadRequestError{Reason: "Invalid format of ID!"}
	}

	var task domain.Task
	if role == string(domain.Admin) {
		task, err = u.taskRepository.GetTask(ctx, taskId)
	} else {
		task, err = u.taskRepository.GetTaskByIDForUser(ctx, taskId, username)
	}
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, &domain.NotFoundError{ID: taskID}
		}
		return domain.Task{}, err
	}

	return task, nil
}

func (u *taskUsecase)UpdateTask(c context.Context, taskId, username, role string, updatedTask domain.Task) (domain.Task, error){
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	taskID, err := strconv.Atoi(taskId)
	if err != nil {
		return domain.Task{}, &domain.BadRequestError{Reason: "Invalid ID"}
	}

	if updatedTask.Title == "" || updatedTask.Description == "" || updatedTask.DueDate == ""{
		return domain.Task{}, &domain.BadRequestError{Reason: "Fields can not be empty!"}
	}

	if updatedTask.Status != domain.Pending && updatedTask.Status != domain.Completed {
		return domain.Task{}, &domain.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}

	task, err := u.taskRepository.GetTask(ctx, taskID)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			return domain.Task{}, &domain.NotFoundError{Resource: "Task", ID: taskId}
		}
		return domain.Task{}, err
	}

	if role != string(domain.Admin) && task.CreatedBy != username{
		return domain.Task{}, &domain.ForbiddenError{}
	}

	updatedTask.ID = task.ID
	updatedTask.CreatedBy = task.CreatedBy

	updatedTask, err = u.taskRepository.UpdateTask(ctx, taskID, updatedTask)
	return updatedTask, err
}

func (u *taskUsecase)AddTask(c context.Context, username string, task domain.Task) (domain.Task, error){
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if task.Title == "" || task.Description == "" || task.DueDate == "" {
		return domain.Task{}, &domain.BadRequestError{Reason: "Fields cannot be empty!"}
	}

	if task.Status != domain.Pending && task.Status != domain.Completed {
		return domain.Task{}, &domain.BadRequestError{Reason: "Status must be either 'Pending' or 'Completed'"}
	}	

	ID, err := u.taskRepository.GetNewID(ctx)
	if err != nil {
		return domain.Task{}, err
	}

	task.CreatedBy = username
	task.ID = ID

	return u.taskRepository.AddTask(ctx, task)
}

func (u *taskUsecase)DeleteTask(c context.Context, id, username, role string) error{
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	taskID, err := strconv.Atoi(id)
	if err != nil {
		return &domain.BadRequestError{Reason: "Invalid format of ID!"} 
	}

	task, err := u.taskRepository.GetTask(ctx, taskID)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			return &domain.NotFoundError{Resource: "Task", ID: id}
		}
		return err
	}

	if role != string(domain.Admin) && task.CreatedBy != username {
		return &domain.UnauthorizedError{}
	}

	return u.taskRepository.DeleteTask(ctx, taskID)
}