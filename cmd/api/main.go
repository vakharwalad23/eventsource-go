package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/vakharwalad23/eventsource-starter-go/internal/api"
	"github.com/vakharwalad23/eventsource-starter-go/internal/app"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/kafka"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/minio"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/redis"
)

func main() {

	ctx := context.Background()

	// MiniO Client
	minioClient, err := minio.NewMinioClient(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
	)

	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// Kafka Producer
	kafkaProducer := kafka.NewProducer([]string{os.Getenv("KAFKA_BROKER")}, os.Getenv("KAFKA_TOPIC"))

	// Redis Client
	redisClient := redis.NewRedisClient(os.Getenv("REDIS_ADDR"))

	// ReadDB
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// App Service
	svc := app.NewAccountService(minioClient, redisClient, kafkaProducer, db)

	// HTTP Handlers
	r := mux.NewRouter()
	api.RgisterHandlers(r, svc)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	// Cleanup resources
	defer minioClient.Close(ctx)
	defer redisClient.Close()
	defer kafkaProducer.Close()

}
