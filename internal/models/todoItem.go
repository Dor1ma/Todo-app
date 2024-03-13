package models

type ToDoItem struct {
	ID          int
	Title       string
	Description string
	Done        bool
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}
