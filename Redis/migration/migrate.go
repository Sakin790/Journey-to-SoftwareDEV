package migration

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) {
	// SQL to create products table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		stock INT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
