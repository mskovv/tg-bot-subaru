package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, pass string) *RedisStorage {
	return &RedisStorage{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
		}),
	}
}

func (r *RedisStorage) GetState(ctx context.Context, userId int64) (string, error) {
	state, err := r.client.Get(ctx, userIdToString(userId)).Result()
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
