package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {

	connectionString := "host=localhost port=5432 user=mahid password=root dbname=mydb sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Lifetime of a connection
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}
	log.Println("Connection pool initialized successfully")

}
