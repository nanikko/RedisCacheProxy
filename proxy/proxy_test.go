package proxy

import (
	"fmt"
	"testing"
	"RedisProxy/cache"
	"RedisProxy/redis"
)

func TestProxyHttpGet(t *testing.T) {
	cache, _ := cache.NewMinimalCache(10, 1000)
	redisHandlerResp, _ := redis.NewRedisHandler("127.0.0.1:6379")
	scheduler := redis.NewScheduler(10, 100)
	scheduler.Run()

	rpHttp := &RedisProxy{
		Cache : cache,
		RedisHandler : redisHandlerResp,
		Scheduler : scheduler,
	}

	// Set some value against key
	redisHandlerResp.Set("key", "value")

	// Get the same key for which we set a value earlier
	_, resp := rpHttp.Get("key")

	fmt.Println(resp)

	if resp != "value" {
		t.Fatalf("Wrong return value from proxy")
	}
}