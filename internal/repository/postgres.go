package repository

import "database/sql"

func NewPostgresDB() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
