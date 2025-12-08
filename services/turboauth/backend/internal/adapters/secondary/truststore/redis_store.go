package truststore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"turboauth/internal/domain/auth"
	"turboauth/pkg/metrics"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements Redis-based caching (L2 cache)
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new Redis store
func NewRedisStore(url, password string, db int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       db,
		PoolSize: 100, // High-performance connection pool
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStore{client: client}, nil
}

// Get retrieves cached data from Redis
func (r *RedisStore) Get(ctx context.Context, walletAddress string) (*auth.WalletAuth, error) {
	key := fmt.Sprintf("auth:%s", walletAddress)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		metrics.CacheMisses.WithLabelValues("L2").Inc()
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, err
	}

	var walletAuth auth.WalletAuth
	if err := json.Unmarshal(data, &walletAuth); err != nil {
		return nil, err
	}

	metrics.CacheHits.WithLabelValues("L2").Inc()
	return &walletAuth, nil
}

// Set stores data in Redis
func (r *RedisStore) Set(ctx context.Context, walletAddress string, data *auth.WalletAuth, ttl time.Duration) error {
	key := fmt.Sprintf("auth:%s", walletAddress)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonData, ttl).Err()
}

// Delete removes data from Redis
func (r *RedisStore) Delete(ctx context.Context, walletAddress string) error {
	key := fmt.Sprintf("auth:%s", walletAddress)
	return r.client.Del(ctx, key).Err()
}

// BatchGet retrieves multiple entries using pipeline
func (r *RedisStore) BatchGet(ctx context.Context, walletAddresses []string) (map[string]*auth.WalletAuth, error) {
	if len(walletAddresses) == 0 {
		return make(map[string]*auth.WalletAuth), nil
	}

	// Use pipeline for performance
	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, addr := range walletAddresses {
		key := fmt.Sprintf("auth:%s", addr)
		cmds[addr] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		// Partial failures are OK
	}

	result := make(map[string]*auth.WalletAuth)
	for addr, cmd := range cmds {
		data, err := cmd.Bytes()
		if err == nil {
			var walletAuth auth.WalletAuth
			if json.Unmarshal(data, &walletAuth) == nil {
				result[addr] = &walletAuth
			}
		}
	}

	return result, nil
}

// BatchSet stores multiple entries using pipeline
func (r *RedisStore) BatchSet(ctx context.Context, data map[string]*auth.WalletAuth, ttl time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()

	for addr, walletAuth := range data {
		key := fmt.Sprintf("auth:%s", addr)
		jsonData, err := json.Marshal(walletAuth)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, jsonData, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// HealthCheck verifies Redis connectivity
func (r *RedisStore) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisStore) Close() error {
	return r.client.Close()
}
