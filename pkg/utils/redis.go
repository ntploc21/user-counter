package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	//"usercounter/config"

	"github.com/redis/go-redis/v9"
)

func getUserCountKey(userID uint) string {
	return fmt.Sprintf("user:count:%d", userID)
}

// GetUserCount retrieves the counter value from Redis
func GetUserCount(ctx context.Context, redisDB *redis.Client, userID uint) (int64, bool, error) {
	key := getUserCountKey(userID)

	val, err := redisDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, false, nil // Cache miss
	}
	if err != nil {
		return 0, false, fmt.Errorf("redis get error: %w", err)
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("failed to parse counter value: %w", err)
	}

	return count, true, nil
}

// SetUserCount sets the counter value in Redis with TTL
func SetUserCount(ctx context.Context, redisDB *redis.Client, userID uint, count int64) error {
	key := getUserCountKey(userID)

	err := redisDB.Set(ctx, key, count, time.Duration(10*60*60)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// IncrementUserCount atomically increments the counter in Redis
func IncrementUserCount(
	ctx context.Context,
	redisDB *redis.Client,
	userID uint,
	amount int64,
) (int64, error) {
	key := getUserCountKey(userID)

	newValue, err := redisDB.IncrBy(ctx, key, amount).Result()
	if err != nil {
		return 0, fmt.Errorf("redis increment error: %w", err)
	}

	// Reset TTL after increment
	redisDB.Expire(ctx, key, time.Duration(10*60*60)*time.Second)

	return newValue, nil
}

// DeleteUserCount removes the counter from Redis
func DeleteUserCount(ctx context.Context, redisDB *redis.Client, userID uint) error {
	key := getUserCountKey(userID)
	return redisDB.Del(ctx, key).Err()
}

// Ping checks if Redis is alive
func Ping(ctx context.Context, redisDB *redis.Client) error {
	return redisDB.Ping(ctx).Err()
}
