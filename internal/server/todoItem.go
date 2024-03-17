package server

import (
	"Todo-app/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *server) createItem() http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		vars := mux.Vars(r)
		listId, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		item := &models.ToDoItem{
			Title:       req.Title,
			Description: req.Description,
			Done:        false,
		}

		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		id, err := s.services.TodoItem.Create(userId, listId, item)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}

func (s *server) getAllItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		vars := mux.Vars(r)
		listId, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		items, err := s.services.TodoItem.GetAll(userId, listId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, items)
	}
}

func (s *server) getItemById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		vars := mux.Vars(r)
		itemId, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		item, err := s.services.TodoItem.GetById(userId, itemId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, item)
	}
}

func (s *server) updateItem() http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Done        bool   `json:"done"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		vars := mux.Vars(r)
		itemId, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		input := &models.UpdateItemInput{
			Title:       &req.Title,
			Description: &req.Description,
			Done:        &req.Done,
		}

		if err := s.services.TodoItem.Update(userId, itemId, input); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) deleteItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		vars := mux.Vars(r)
		itemId, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.services.TodoItem.Delete(userId, itemId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}
