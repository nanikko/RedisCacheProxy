package proxy

import (
	"fmt"
	"testing"
	"RedisProxy/cache"
	"RedisProxy/redis"
)

func TestProxyRespGet(t *testing.T) {
	cache, _ := cache.NewMinimalCache(10, 1000)
	redisHandlerResp, _ := redis.NewRedisHandlerResp("127.0.0.1:6379")
	scheduler := redis.NewScheduler(10, 100)

	rpTcp := &RedisProxyTcp{
		Cache : cache,
		RedisHandlerResp : redisHandlerResp,
		Scheduler : scheduler,
		Port : "10000",
	}

	// Set some value against key
	redisHandlerResp.Set("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")

	// Get the same key for which we set a value earlier
	cmdStr := "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"

	//resp := redisHandlerResp.Get(cmdStr)
	_, resp :=  rpTcp.Get(cmdStr)

	fmt.Println(resp)

	if resp != "$5\r\nvalue\r\n" {
		t.Fatalf("Wrong return value from proxy")
	}
}