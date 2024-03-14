package models

import (
	"errors"
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

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
