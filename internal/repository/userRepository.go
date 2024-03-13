package repository

import (
	"Todo-app/internal/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (s *UserRepository) Close() {
	err := s.db.Close()
	if err != nil {
		return
	}
}

func (r *UserRepository) Create(u *models.User) (*models.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		u.Name,
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.QueryRow(
		"SELECT id, name, email, password_hash FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Find(id int) (*models.User, error) {
	user := &models.User{}
	if err := r.db.QueryRow(
		"SELECT id, name, email, password_hash FROM users WHERE id = $1",
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		return nil, err
	}

	return user, nil
}
