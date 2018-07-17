package redis

import (
	"github.com/go-redis/redis"
	"log"
)

// RedisHandler implements the JobHandler interface
type RedisHandler struct {
	client *redis.Client
}

// NewRedisHandler creates new redis handler
func NewRedisHandler(addr string) (*RedisHandler, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisHandler := new(RedisHandler)
	redisHandler.client = client

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisHandler, nil
}

// Gets the data based on key
func (s *RedisHandler) Get(key string) string {
	log.Println("In regular Get method")
	val, err := s.client.Get(key).Result()
	if err != nil {
		log.Println("[RedisHandler.Get]: Failed to found key", key, "from redis")
		return ""
	}
	return val
}

// Get gets the data based on key
func (s *RedisHandler) Set(key string, value string) bool {
	status := s.client.Set(key, value, 0)

	log.Println("set status in basic redis handler: ", status.String())
	return true
}