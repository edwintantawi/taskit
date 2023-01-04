package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

type Config struct {
	Port     string
	Host     string
	DB       string
	User     string
	Password string
	SSLMode  string
}

const Driver = "postgres"

// New create new postgres connection.
func New(config *Config) (*sql.DB, func(autoMigrate bool) error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DB,
		config.SSLMode,
	)

	db, err := sql.Open(Driver, dsn)
	if err != nil {
		log.Fatalf("Failed open postgres connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping postgres database: %v", err)
	}

	migrate := func(autoMigrate bool) error {
		if !autoMigrate {
			return nil
		}
		iofsDriver, err := iofs.New(os.DirFS("migrations"), ".")
		if err != nil {
			return err
		}
		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
		if err != nil {
			return err
		}
		err = migrator.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Migration no change, nothing to migrate")
			return nil
		} else if err != nil {
			return err
		}
		log.Println("Successfully migrate database")
		return nil
	}

	return db, migrate
}
