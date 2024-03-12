package server

import (
	"Todo-app/internal/repository"
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
)

func Start() error {
	connStr := "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable"
	db, err := newDB(connStr)
	if err != nil {
		return err
	}

	defer db.Close()
	//temp key
	userRepository := repository.UserRepository{}
	err = userRepository.Open()
	if err != nil {
		return err
	}

	sessionStore := sessions.NewCookieStore([]byte("239239"))
	srv := newServer(userRepository, sessionStore)
	constAddr := "localhost:8080"

	return http.ListenAndServe(constAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
