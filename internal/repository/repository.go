package repository

import (
	"Todo-app/internal/models"
	"database/sql"
)

type Authorization interface {
	Create(u *models.User) (*models.User, error)
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
	Create(listId int, item *models.ToDoItem) (int, error)
	GetAll(userId, listId int) ([]*models.ToDoItem, error)
	GetById(userId, itemId int) (*models.ToDoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input *models.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewUserRepository(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
