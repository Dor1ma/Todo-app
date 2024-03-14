package server

import (
	"Todo-app/internal/models"
	"Todo-app/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

const (
	sessionName        = "tempSessionName"
	ctxKeyUser  ctxKey = iota
)

type server struct {
	router     *mux.Router
	repository repository.Repository
	sessions   sessions.Store
}

type ctxKey int8

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func newServer(repository repository.Repository, sessionStore sessions.Store) *server {
	s := &server{
		router:     mux.NewRouter(),
		repository: repository,
		sessions:   sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")

	todos := private.PathPrefix("/todos").Subrouter()
	todos.HandleFunc("/", s.handleTodosCreate()).Methods("POST")
	todos.HandleFunc("/{id}", s.handleTodosUpdate()).Methods("PUT")
	todos.HandleFunc("/{id}", s.handleTodosDelete()).Methods("DELETE")
	todos.HandleFunc("/{id}", s.getListById()).Methods("GET")
	todos.HandleFunc("/", s.getAllLists()).Methods("GET")

	items := todos.PathPrefix("/{id}/items").Subrouter()
	items.HandleFunc("/", s.getAllItems()).Methods("GET")
	items.HandleFunc("/{id}", s.getItemById()).Methods("GET")
	//items.HandleFunc("/", s.createItem()).Methods("POST")
	items.HandleFunc("/{id}", s.updateItem()).Methods("PUT")
	items.HandleFunc("/{id}", s.deleteItem()).Methods("DELETE")
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(writer, r, http.StatusBadRequest, err)
			return
		}

		u := &models.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}

		if _, err := s.repository.Authorization.Create(u); err != nil {
			s.error(writer, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(writer, r, http.StatusCreated, u)
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessions.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.repository.Authorization.Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*models.User))
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.repository.Authorization.FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessions.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessions.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

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

/*func (s *server) createItem() http.HandlerFunc {
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

		u := r.Context().Value(ctxKeyUser).(*models.User)
		userId := u.ID

		vars := mux.Vars(r)
		listId, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		item := &models.ToDoItem{
			Title: req.Title,
		}

		id, err := s.repository.TodoItem.Create(listId, item)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}*/

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

		items, err := s.repository.TodoItem.GetAll(userId, listId)
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

		item, err := s.repository.TodoItem.GetById(userId, itemId)
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

		if err := s.repository.TodoItem.Update(userId, itemId, input); err != nil {
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

		err = s.repository.TodoItem.Delete(userId, itemId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
