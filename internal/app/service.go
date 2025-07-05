package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/kafka"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/minio"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/redis"
)

type AccountService struct {
	minio  *minio.MinioClient
	redis  *redis.RedisClient
	kafka  *kafka.Producer
	readDB *sql.DB
}

func NewAccountService(minioClient *minio.MinioClient, redisClient *redis.RedisClient, kafkaProducer *kafka.Producer, readDB *sql.DB) *AccountService {
	return &AccountService{
		minio:  minioClient,
		kafka:  kafkaProducer,
		redis:  redisClient,
		readDB: readDB,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, accountId string) error {
	event := domain.Event{
		Type:      domain.AccountCreated,
		AccountID: accountId,
		Time:      time.Now(),
	}
	return s.kafka.PublishEvents(ctx, event)
}

func (s *AccountService) Deposit(ctx context.Context, accountId string, amount float64) error {
	event := domain.Event{
		Type:      domain.MoneyDeposited,
		AccountID: accountId,
		Amount:    amount,
		Time:      time.Now(),
	}
	if err := s.kafka.PublishEvents(ctx, event); err != nil {
		return err
	}
	return s.redis.DeleteAccount(ctx, accountId)
}

func (s *AccountService) Withdraw(ctx context.Context, accountId string, amount float64) error {
	acc, err := s.GetAccount(ctx, accountId)
	if err != nil {
		return err
	}
	if acc.Balance < amount {
		return errors.New("insufficient funds")
	}
	event := domain.Event{
		Type:      domain.MoneyWithdrawn,
		AccountID: accountId,
		Amount:    amount,
		Time:      time.Now(),
	}
	if err = s.kafka.PublishEvents(ctx, event); err != nil {
		return err
	}
	return s.redis.DeleteAccount(ctx, accountId)
}

func (s *AccountService) GetAccount(ctx context.Context, accountId string) (*domain.Account, error) {
	// Try to get the account from Redis first
	acc, err := s.redis.GetAccount(ctx, accountId)
	if err == nil {
		return acc, nil
	}

	// If not found in Redis, read from readDB
	row := s.readDB.QueryRowContext(ctx, "SELECT id, balance FROM accounts WHERE id = $1", accountId)
	var account domain.Account
	if err := row.Scan(&account.ID, &account.Balance); err != nil {
		return nil, err
	}

	// hydrate Redis cache
	_ = s.redis.SetAccount(ctx, &account)

	return &account, nil
}
