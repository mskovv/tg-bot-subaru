package storage

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage() (*RedisStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	log.Printf("Connected to Redis")
	return &RedisStorage{
		client: rdb,
	}, nil
}

func (r *RedisStorage) GetState(ctx context.Context, userId int64) (string, error) {
	state, err := r.client.Get(ctx, userIdToString(userId)).Result()

	if errors.Is(err, redis.Nil) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return state, nil
}

func (r *RedisStorage) SetState(ctx context.Context, userId int64, state string) error {
	return r.client.Set(ctx, userIdToString(userId), state, 0).Err()
}

func userIdToString(userId int64) string {
	return "user_state:" + strconv.FormatInt(userId, 10)
}
