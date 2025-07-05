package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/minio"
)

func main() {
	ctx := context.Background()

	// Debug environment variables
	log.Printf("MINIO_ENDPOINT: %s", os.Getenv("MINIO_ENDPOINT"))
	log.Printf("KAFKA_BROKER: %s", os.Getenv("KAFKA_BROKER"))
	log.Printf("KAFKA_TOPIC: %s", os.Getenv("KAFKA_TOPIC"))

	// MiniO Client
	minioClient, err := minio.NewMinioClient(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
	)

	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	log.Println("Connected to MinIO successfully")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    os.Getenv("KAFKA_TOPIC"),
		GroupID:  "event-consumer-group",
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
	})

	defer r.Close()

	log.Println("Kafka consumer started, waiting for messages...")
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		log.Printf("Received message: %s", string(m.Value))

		var event domain.Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		log.Printf("Parsed event: %+v", event)

		if err := minioClient.AppendEvent(ctx, event); err != nil {
			log.Printf("Failed to append event to MinIO: %v", err)
		} else {
			log.Printf("Successfully appended event to MinIO")
		}
	}
}
