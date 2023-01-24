package database

import (
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migration runs the database migrations.
// The dsn is the database connection string.
// The migrationDir is the path to the directory containing the migrations.
func Migration(dsn, migrationDir string) error {
	iofsDriver, err := iofs.New(os.DirFS(migrationDir), ".")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
	if err != nil {
		return err
	}
	err = migrator.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
