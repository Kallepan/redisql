package mock

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestMockRedis(t *testing.T) {
	// Simple ping test
	ctx := context.Background()
	client := InitMockDB()

	if _, err := client.Ping(ctx).Result(); err != nil {
		t.Errorf("Ping(): expected nil, actual %v", err)
	}

	// Set and Get
	client.Set(ctx, "test", "value", 0)

	value, err := client.Get(ctx, "test").Result()
	if err != nil {
		t.Errorf("Get(): expected nil, actual %v", err)
	}
	if value != "value" {
		t.Errorf("Get(): expected value, actual %v", value)
	}

	secondClient := InitMockDB()
	if err := secondClient.Ping(ctx).Err(); err != nil {
		t.Errorf("Ping(): expected nil, actual %v", err)
	}

	// Get from another DB
	if _, err := secondClient.Get(ctx, "test").Result(); err != redis.Nil {
		t.Errorf("Get(): expected nil, actual %v", err)
	}
}
