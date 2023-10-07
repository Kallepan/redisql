package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	DB   *redis.Client
	once sync.Once
)

func initDB() {
	/**
	* In the initDB function, we create a new Redis client and assign it to the DB variable.
	* The sync.Once type is used to perform initialization exactly once.
	**/
	ctx := context.Background()
	once.Do(func() {
		// Get the connection details from the environment variables.
		addr, pass, db := getConnectionDetails()

		// Create a new Redis client.
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       db,
		})
		DB = client

		// Close the connection when the application exits.
		go func() {
			<-ctx.Done()
			if err := DB.Close(); err != nil {
				panic(err)
			}
		}()

		// Ping the Redis server and check if any errors occurred.
		if _, err := DB.Ping(ctx).Result(); err != nil {
			panic(err)
		}
	})
}

func getConnectionDetails() (string, string, int) {
	/**
	* The getConnectionDetails function returns the connection details for the Redis client.
	* The function returns the host, password, and database number.
	**/
	host := os.Getenv("REDIS_DSN")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASSWORD")
	db := os.Getenv("REDIS_DB")

	addr := fmt.Sprintf("%s:%s", host, port)
	intDB := convertToInt(db)

	return addr, pass, intDB
}

func convertToInt(db string) int {
	/**
	* The convertToInt function converts the database number to an integer.
	**/
	intDB, err := strconv.Atoi(db)
	if err != nil {
		panic(err)
	}

	return intDB
}
