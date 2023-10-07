package mocks

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var dbTracker = []int{}

func InitMockDB() *redis.Client {
	/**
	* In the initDB function, we create a new Redis client and assign it to the DB variable.
	* The sync.Once type is used to perform initialization exactly once.
	**/
	ctx := context.Background()
	// Get the connection details from the environment variables.
	addr, pass := getConnectionDetails()

	// Create a new Redis client.
	dbID := len(dbTracker)
	dbTracker = append(dbTracker, dbID)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbID,
	})

	// Close the connection when the application exits.
	go func() {
		<-ctx.Done()
		if err := client.Close(); err != nil {
			panic(err)
		}
	}()

	// Ping the Redis server and check if any errors occurred.
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return client
}

func getConnectionDetails() (string, string) {
	/**
	* The getConnectionDetails function returns the connection details for the Redis client.
	* The function returns the host, password, and database number.
	**/
	host := os.Getenv("REDIS_DSN")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASSWORD")

	addr := fmt.Sprintf("%s:%s", host, port)

	return addr, pass
}
