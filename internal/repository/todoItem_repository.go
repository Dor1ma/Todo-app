package repository

import (
	"Todo-app/internal/models"
	"database/sql"
	"fmt"
)

type TodoItemPostgres struct {
	db *sql.DB
}

func NewTodoItemPostgres(db *sql.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item *models.ToDoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO todo_items (title, description) values ($1, $2) RETURNING id")

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO lists_items (list_id, item_id) values ($1, $2)")
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]*models.ToDoItem, error) {
	var items []*models.ToDoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti INNER JOIN lists_items li on li.item_id = ti.id
									INNER JOIN users_lists ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`)
	rows, err := r.db.Query(query, listId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.ToDoItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Done); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (*models.ToDoItem, error) {
	var item models.ToDoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM todo_items ti INNER JOIN lists_items li on li.item_id = ti.id
									INNER JOIN users_lists ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`)
	err := r.db.QueryRow(query, itemId, userId).Scan(&item.ID, &item.Title, &item.Description, &item.Done)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM todo_items ti USING lists_items li, users_lists ul 
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) Update(userId, itemId int, input *models.UpdateItemInput) error {
	query := `UPDATE todo_items ti SET title = $1, description = $2, done = $3 FROM lists_items li, users_lists ul
	WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $4 AND ti.id = $5`
	_, err := r.db.Exec(query, input.Title, input.Description, input.Done, userId, itemId)
	return err
}
