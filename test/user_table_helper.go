package test

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserTableHelper struct {
	db *sql.DB
}

func NewUserTableHelper(db *sql.DB) UserTableHelper {
	return UserTableHelper{db}
}

func (uth UserTableHelper) GetByID(id string) User {
	var user User
	q := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`
	if err := uth.db.QueryRow(q, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		log.Fatalf("Could not get user: %s", err)
	}
	return user
}

func (uth UserTableHelper) Add(user User) {
	q := `INSERT INTO users (id, name, email, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := uth.db.Exec(q, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt); err != nil {
		log.Fatalf("Could not add user: %s", err)
	}
}
