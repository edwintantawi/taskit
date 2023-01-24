package test

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

var Tables = []string{"users"}

func TruncateTables(db *sql.DB) func() {
	return func() {
		q := fmt.Sprintf("TRUNCATE TABLE %s", strings.Join(Tables, ", "))
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("Could not truncate tables: %s", err)
		}
	}
}
