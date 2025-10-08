package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB() *sqlx.DB {
	connectionString := "host=localhost port=5432 user=mahid password=root dbname=mydb sslmode=disable"

	// Option 1: use sqlx.Open directly
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Ping DB
	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Connection pool initialized successfully")
	return db
}
