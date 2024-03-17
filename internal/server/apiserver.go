package server

import (
	"Todo-app/internal/repository"
	"Todo-app/internal/service"
	"database/sql"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"net/http"
)

func Start() error {
	connStr := "postgres://postgres:postgres@localhost:5437/postgres?sslmode=disable"
	db, err := newDB(connStr)
	if err != nil {
		return err
	}

	defer db.Close()

	repos := repository.NewRepository(db)
	//temp key
	sessionStore := sessions.NewCookieStore([]byte("239239"))
	services := service.NewService(repos)
	srv := newServer(*services, sessionStore)
	constAddr := ":8080"

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
