package auth

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func SetToken(ctx context.Context, token, userID string) error {
	return rdb.Set(ctx, token, userID, time.Hour*24*7).Err()
}

func GetToken(ctx context.Context, userID string) (string, error) {
	return rdb.Get(ctx, userID).Result()
}
