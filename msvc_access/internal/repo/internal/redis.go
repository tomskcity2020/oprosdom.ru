package repo_internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, addr string) (*Redis, error) {

	// таймаут 10 сек если ctx не придет быстрее
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := client.Ping(ctxTimeout).Result(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Redis{client: client}, nil
}

func (p *Redis) Close() {
	p.client.Close()
}

func (p *Redis) Set(ctx context.Context, k string, v any, ttl time.Duration) error {
	if err := p.client.Set(ctx, k, v, ttl).Err(); err != nil {
		log.Printf("redis SET failed: %v", err)
		return err
	}
	return nil
}

// func (p *Redis) Del(ctx context.Context, k string) (int64, error) {
// 	count, err := p.client.Del(ctx, k).Result()
// 	if err != nil {
// 		log.Printf("redis DEL failed: %v", err)
// 		return 0, fmt.Errorf("redis DEL failed: %w", err)
// 	}

// 	return count, nil
// }
