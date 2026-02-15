package databases

import (
	"fmt"
	"sport-assistance/pkg/configs"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *configs.Config) *redis.Client {
	redisDB, err := strconv.Atoi(cfg.RedisConfig.DBName)
	if err != nil {
		panic(err)
	}

	// TO DO Добавить Write / Read Timeout
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		Password: cfg.RedisConfig.Password,
		DB:       redisDB,
	})

	return redisClient
}
