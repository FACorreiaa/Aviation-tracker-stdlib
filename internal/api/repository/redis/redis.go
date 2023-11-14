package redis

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
	"syscall"
)

type Redis struct {
	db *redis.Client
}

type Config struct {
	redisHost     string
	redisPassword string
	redisDb       int
}

func NewRedisConfig(
	redisHost string,
	redisPassword string,
	redisDb int,
) Config {
	return Config{
		redisHost:     redisHost,
		redisPassword: redisPassword,
		redisDb:       redisDb,
	}
}

func newDB(config Config) (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // Default value if environment variable is not set
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.redisHost,
		Password: config.redisPassword, // no password set
		DB:       config.redisDb,       // use default DB
	})

	return rdb, nil
}
func NewRedis(config Config) *Redis {
	db, err := newDB(config)
	if err != nil {
		zap.L().Fatal("Error on postgres init")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
	return &Redis{db: db}
}

func (p *Redis) GetDB() *redis.Client {
	return p.db
}
