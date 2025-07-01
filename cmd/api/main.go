package main

import (
	"context"
	"log"
	"os"

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

	// Redis Client
	redisClient := redis.NewRedisClient(os.Getenv("REDIS_ADDR"))

	defer minioClient.Close(ctx)
	defer redisClient.Close()

}
