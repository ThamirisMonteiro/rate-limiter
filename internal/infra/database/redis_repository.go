package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const KeyDoesNotExist = time.Duration(-2)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(address string) *RedisRepository {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})
	return &RedisRepository{client: client, ctx: ctx}
}

func (r RedisRepository) SetKey(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r RedisRepository) GetKey(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r RedisRepository) IncrCounter(counter string) error {
	return r.client.Incr(r.ctx, counter).Err()
}

func (r RedisRepository) ExpireKey(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r RedisRepository) TTLKey(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

func (r RedisRepository) Close() error {
	return r.client.Close()
}

func (r RedisRepository) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}
