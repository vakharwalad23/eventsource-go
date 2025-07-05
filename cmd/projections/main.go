package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
)

func main() {
	ctx := context.Background()
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// ensure table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts(
			id TEXT PRIMARY KEY,
			balance DOUBLE PRECISION NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: "projection-group",
	})
	defer r.Close()

	log.Println("Projection service started")
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka error: %v", err)
			continue
		}
		var event domain.Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Invalid event: %v", err)
			continue
		}

		log.Printf("Processing event: %s for account %s", event.Type, event.AccountID)

		switch event.Type {
		case domain.AccountCreated:
			_, _ = db.Exec(`INSERT INTO accounts (id, balance) VALUES ($1, 0) ON CONFLICT (id) DO NOTHING`, event.AccountID)
		case domain.MoneyDeposited:
			_, _ = db.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, event.Amount, event.AccountID)
		case domain.MoneyWithdrawn:
			_, _ = db.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, event.Amount, event.AccountID)
		}
	}
}
