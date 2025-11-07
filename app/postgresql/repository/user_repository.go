package repository

import (
	"alumni-app/app/postgresql/model"
	"alumni-app/database/postgresql"
	"database/sql"
)

type UserRepository interface {
	GetByID(id int) (model.User, error)
	GetByUsername(username string) (model.User, error)
	Create(user *model.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Get user by ID
func (r *userRepository) GetByID(id int) (model.User, error) {
	var u model.User
	row := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`, id)

	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return model.User{}, sql.ErrNoRows
	}
	return u, err
}

func (r *userRepository) GetByUsername(username string) (model.User, error) {
	var u model.User
	row := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE username = $1
	`, username)

	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return model.User{}, sql.ErrNoRows
	}
	return u, err
}

func (r *userRepository) Create(user *model.User) error {
    query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`
    return database.DB.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)
}

