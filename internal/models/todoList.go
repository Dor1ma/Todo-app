package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type ToDoList struct {
	ID          int
	Title       string
	Description string
}

func (t *ToDoList) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Title, validation.Required, validation.Length(2, 100)),
	)
}
