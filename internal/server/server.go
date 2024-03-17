package server

import (
	"Todo-app/internal/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
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
	router   *mux.Router
	services service.Service
	sessions sessions.Store
}

type ctxKey int8

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func newServer(services service.Service, sessionStore sessions.Store) *server {
	s := &server{
		router:   mux.NewRouter(),
		services: services,
		sessions: sessionStore,
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
	items.HandleFunc("/", s.createItem()).Methods("POST")
	items.HandleFunc("/{id}", s.updateItem()).Methods("PUT")
	items.HandleFunc("/{id}", s.deleteItem()).Methods("DELETE")
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
