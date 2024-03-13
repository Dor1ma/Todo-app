package repository

import (
	"Todo-app/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

type TodoListPostgres struct {
	db *sql.DB
}

func NewTodoListPostgres(db *sql.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list *models.ToDoList) (int, error) {
	if err := list.Validate(); err != nil {
		return 0, err
	}

	err := r.db.QueryRow("INSERT INTO todo_lists (title, description) VALUES ($1, $2) RETURNING id",
		list.Title,
		list.Description,
	).Scan(&list.ID)
	if err != nil {
		return 0, err
	}

	_, err = r.db.Exec("INSERT INTO users_lists (user_id, list_id) VALUES ($1, $2)", userId, list.ID)
	if err != nil {
		return 0, err
	}

	return list.ID, nil
}

func (r *TodoListPostgres) GetAll(userId int) ([]*models.ToDoList, error) {
	var lists []*models.ToDoList

	query := "SELECT tl.id, tl.title, tl.description FROM todo_lists tl INNER JOIN users_lists ul on tl.id = ul.list_id WHERE ul.user_id = $1"
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var list models.ToDoList
		if err := rows.Scan(&list.ID, &list.Title, &list.Description); err != nil {
			return nil, err
		}
		lists = append(lists, &list)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoListPostgres) GetById(userId, listId int) (*models.ToDoList, error) {
	list := &models.ToDoList{}

	query := "SELECT tl.id, tl.title, tl.description FROM todo_lists tl INNER JOIN users_lists ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2"
	err := r.db.QueryRow(query, userId, listId).Scan(&list.ID, &list.Title, &list.Description)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := "DELETE FROM todo_lists tl USING users_lists ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2"
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input *models.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE todo_lists tl FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d", setQuery, argId, argId+1)
	args = append(args, userId, listId)

	_, err := r.db.Exec(query, args...)
	return err
}
