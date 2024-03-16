package service

import (
	"Todo-app/internal/models"
	"Todo-app/internal/repository"
)

type Authorization interface {
	CreateUser(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Find(id int) (*models.User, error)
}

type TodoList interface {
	Create(userId int, list *models.ToDoList) (int, error)
	GetAll(userId int) ([]*models.ToDoList, error)
	GetById(userId, listId int) (*models.ToDoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input *models.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item *models.ToDoItem) (int, error)
	GetAll(userId, listId int) ([]*models.ToDoItem, error)
	GetById(userId, itemId int) (*models.ToDoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input *models.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
