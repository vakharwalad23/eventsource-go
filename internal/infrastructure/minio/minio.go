package minio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
)

type MinioClient struct {
	client *minio.Client
	bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey string) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Set to true while using HTTPS
	})
	if err != nil {
		return nil, err
	}

	bucket := "events" // Default bucket name
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	if err != nil {
		return nil, err
	}

	return &MinioClient{client: client, bucket: bucket}, nil
}

func (m *MinioClient) AppendEvent(ctx context.Context, event domain.Event) error {
	key := fmt.Sprintf("%s.jsonl", event.AccountID)
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	data = append(data, '\n') // Append newline for JSON Lines format

	// Get the existing object
	obj, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	var existingData []byte
	if err == nil {
		existingData, _ = io.ReadAll(obj)
	}

	// Append new event data to existing data
	all := append(existingData, data...)
	reader := strings.NewReader(string(all))
	_, err = m.client.PutObject(ctx, m.bucket, key, reader, int64(len(all)), minio.PutObjectOptions{})

	return err
}

func (m *MinioClient) GetEvents(ctx context.Context, accountId string) ([]domain.Event, error) {
	key := fmt.Sprintf("%s.jsonl", accountId)
	obj, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	var events []domain.Event

	dec := json.NewDecoder(obj)
	for {
		var e domain.Event
		if err := dec.Decode(&e); err != nil {
			break
		}
		events = append(events, e)
	}
	return events, nil
}

func (m *MinioClient) Close(ctx context.Context) error {
	// MinIO client does not have a specific close method, but resources can be cleaned up if needed.
	// This is a placeholder for any cleanup logic if required.
	return nil
}
