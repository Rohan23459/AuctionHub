package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global PostgreSQL database connection.
var DB *gorm.DB

// Rdb is the global Redis client.
var Rdb *redis.Client

// Ctx is the context for Redis operations.
var Ctx = context.Background()

// InitDB initializes the PostgreSQL connection.
func InitDB() {
	dsn := "host=localhost user=postgres password=rohan dbname=auction port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}

// InitRedis initializes the Redis connection.
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
}
