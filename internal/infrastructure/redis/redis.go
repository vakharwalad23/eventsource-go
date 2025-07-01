package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vakharwalad23/eventsource-starter-go/internal/domain"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) GetAccount(ctx context.Context, accountId string) (*domain.Account, error) {
	val, err := r.client.Get(ctx, accountId).Result()
	if err != nil {
		return nil, err
	}

	var acc domain.Account
	if err = json.Unmarshal([]byte(val), &acc); err != nil {
		return nil, err
	}

	return &acc, nil
}

func (r *RedisClient) SetAccount(ctx context.Context, account *domain.Account) error {
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, account.ID, data, 10*time.Minute).Err()
}

func (r *RedisClient) DeleteAccount(ctx context.Context, accountId string) error {
	return r.client.Del(ctx, accountId).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
