package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository() *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "avito_pass",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error Redis connection: %v", err)
	}

	log.Println("Successfully connected to the Redis.")
	return &RedisRepository{Client: client}
}

func (r *RedisRepository) CacheToken(username, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := r.Client.Set(ctx, "auth_token:"+username, token, 72*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) GetCachedToken(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	token, err := r.Client.Get(ctx, "auth_token:"+username).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		return "", err
	}
	return token, nil
}
