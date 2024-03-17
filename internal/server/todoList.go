package server

import (
	"Todo-app/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *server) handleTodosCreate() http.HandlerFunc {
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

		userID := r.Context().Value(ctxKeyUser).(*models.User).ID

		t := &models.ToDoList{
			Title:       req.Title,
			Description: req.Description,
		}

		if _, err := s.repository.TodoList.Create(userID, t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, t)
	}
}

func (s *server) handleTodosUpdate() http.HandlerFunc {
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
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value(ctxKeyUser).(*models.User).ID
		t := &models.UpdateListInput{
			Title:       &req.Title,
			Description: &req.Description,
		}

		if err := s.repository.TodoList.Update(userID, id, t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, t)
	}
}

func (s *server) handleTodosDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.repository.TodoList.Delete(id, r.Context().Value(ctxKeyUser).(*models.User).ID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) getListById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		list, err := s.repository.TodoList.GetById(userId, id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) getAllLists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		lists, err := s.repository.TodoList.GetAll(userId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, lists)
	}
}
