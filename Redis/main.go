package main

import (
	"backend/db"
	migrate "backend/migration"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var DB *sqlx.DB
var REDIS *redis.Client
var Ctx = context.Background()

func main() {
	DB = db.ConnectDB() // Initialize DB
	migrate.Migrate(DB)
	utils.ConnectRedis()

	mux := http.NewServeMux()

	mux.HandleFunc("/products", getProductsHandler)          // GET all products
	mux.HandleFunc("/products/create", createProductHandler) // POST create product

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", mux)
}

// GET /products
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	var products []models.Product
	err := DB.Select(&products, "SELECT * FROM products ORDER BY id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST /products/create
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body")
		return
	}

	query := `INSERT INTO products (name, stock) VALUES ($1, $2) RETURNING id`
	err = DB.QueryRow(query, p.Name, p.Stock).Scan(&p.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database insert error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
