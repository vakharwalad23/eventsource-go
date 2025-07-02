package app

import (
	"context"
	"errors"
	"time"

	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/minio"
	"github.com/vakharwalad23/eventsource-starter-go/internal/infrastructure/redis"
)

type AccountService struct {
	minio *minio.MinioClient
	redis *redis.RedisClient
}

func NewAccountService(minioClient *minio.MinioClient, redisClient *redis.RedisClient) *AccountService {
	return &AccountService{
		minio: minioClient,
		redis: redisClient,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, accountId string) error {
	event := domain.Event{
		Type:      domain.AccountCreated,
		AccountID: accountId,
		Time:      time.Now(),
	}
	return s.minio.AppendEvent(ctx, event)
}

func (s *AccountService) Deposit(ctx context.Context, accountId string, amount float64) error {
	event := domain.Event{
		Type:      domain.MoneyDeposited,
		AccountID: accountId,
		Amount:    amount,
		Time:      time.Now(),
	}
	if err := s.minio.AppendEvent(ctx, event); err != nil {
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
	if err = s.minio.AppendEvent(ctx, event); err != nil {
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

	// If not found in Redis, fetch from MinIO
	events, err := s.minio.GetEvents(ctx, accountId)
	if err != nil {
		return nil, err
	}

	acc = &domain.Account{ID: accountId, Balance: 0}
	for _, e := range events {
		switch e.Type {
		case domain.AccountCreated:
			acc.Balance = 0
		case domain.MoneyDeposited:
			acc.Balance += e.Amount
		case domain.MoneyWithdrawn:
			acc.Balance -= e.Amount
		}
	}

	// hydrate Redis cache
	_ = s.redis.SetAccount(ctx, acc)

	return acc, nil
}
