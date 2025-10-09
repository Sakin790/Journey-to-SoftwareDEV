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
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var (
	DB    *sqlx.DB
	REDIS *redis.Client
	RMQ   *amqp091.Connection
	Ctx   = context.Background()
)

func main() {
	// Initialize dependencies
	DB = db.ConnectDB()
	migrate.Migrate(DB)
	REDIS = utils.ConnectRedis()
	RMQ = utils.ConnectRabbitMQ()

	// Start background worker
	go startWorker(DB, RMQ)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/products", getProductsHandler)
	mux.HandleFunc("/products/create", createProductHandler)

	log.Println("🚀 Server running on port 8080")
	http.ListenAndServe(":8080", mux)
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cacheKey := "products:all"
	cachedData, err := REDIS.Get(Ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		log.Println("✅ Returned from Redis cache")
		return
	}

	var products []models.Product
	err = DB.Select(&products, "SELECT * FROM products ORDER BY id")
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(products)
	REDIS.Set(Ctx, cacheKey, jsonData, 30*time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	log.Println("🗄️ Fetched from DB and cached in Redis")
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ch, err := RMQ.Channel()
	if err != nil {
		http.Error(w, "RabbitMQ channel error", http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"product_insert_queue",
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		http.Error(w, "Queue declare error", http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(p)
	err = ch.PublishWithContext(
		Ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "✅ Product queued successfully")
	log.Println("📦 Queued product for insertion:", p.Name)
}

// 🧠 Background worker function
func startWorker(DB *sqlx.DB, RMQ *amqp091.Connection) {
	ch, err := RMQ.Channel()
	if err != nil {
		log.Fatalf("❌ Worker channel error: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"product_insert_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Queue declare error: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Consume error: %v", err)
	}

	log.Println("👷 Worker started: waiting for messages...")

	for msg := range msgs {
		var p models.Product
		if err := json.Unmarshal(msg.Body, &p); err != nil {
			log.Printf("❌ Invalid message format: %v", err)
			msg.Nack(false, false)
			continue
		}

		query := `INSERT INTO products (name, stock) VALUES ($1, $2)`
		_, err := DB.Exec(query, p.Name, p.Stock)
		if err != nil {
			log.Printf("❌ DB insert error: %v", err)
			msg.Nack(false, true) // retry message
			continue
		}

		log.Printf("✅ Product inserted into DB: %s", p.Name)
		msg.Ack(false)
	}
}
