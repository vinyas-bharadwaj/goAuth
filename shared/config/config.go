package config

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	Rdb *redis.Client
)

func LoadEnv() error {
	// Try loading from current directory first
	err := godotenv.Load(".env")
	if err != nil {
		// Try relative path for tests
		err = godotenv.Load("../.env")
		if err != nil {
			return err
		}
	}
	return nil
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Connected to redis")
}