package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const Driver = "postgres"

// New create new postgres connection.
func New(dsn string) *sql.DB {
	db, err := sql.Open(Driver, dsn)
	if err != nil {
		log.Fatalf("Failed open postgres connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping postgres database: %v", err)
	}
	return db
}
