package repo_internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, addr string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Redis{client: client}, nil
}

func (p *Redis) Close() {
	p.client.Close()
}

func (p *Redis) Set(ctx context.Context, k string, v any, ttl time.Duration) {
	if err := p.client.Set(ctx, k, v, ttl).Err(); err != nil {
		log.Printf("redis SET failed: %v", err)
	}
}

func (p *Redis) Incr(ctx context.Context, k string) (int64, error) {
	count, err := p.client.Incr(ctx, k).Result()
	if err != nil {
		log.Printf("redis INCR failed: %v", err)
		return 0, fmt.Errorf("redis INCR failed: %w", err)
	}

	return count, nil
}

func (p *Redis) Get(ctx context.Context, k string) (string, error) {
	v, err := p.client.Get(ctx, k).Result()
	if err != nil {
		if err == redis.Nil {
			// возвращаем redis.Nil как есть — чтобы вызывающий код мог проверить через errors.Is
			return "", redis.Nil
		}
		log.Printf("redis GET failed: %v", err)
		return "", fmt.Errorf("redis GET failed: %w", err)
	}

	return v, nil
}
