package proxy

import (
	"fmt"
    "log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"RedisProxy/redis"
	"RedisProxy/cache"
)

type RedisProxy struct {
	Cache *cache.MinimalCache
	RedisHandler *redis.RedisHandler
	RedisHandlerResp *redis.RedisHandlerResp
	Scheduler *redis.Scheduler
}

func (s *RedisProxy) Get(key string) (bool, string) {
	found, v := s.Cache.Get(key)

	if found {
		log.Printf("RedisProxy.Get.HIT: got value:%s for key:%s", v, key)
		return true, v
	} else {
		resp := make(chan string)
		defer close(resp)
		work := redis.Job{
			Request: key,
			JobHandler: s.RedisHandler,
			Resp: resp,
		}
		s.Scheduler.JobQueue <- work

		v = <-resp
		if v != "" {
			s.Cache.Add(key, v)
		}
		log.Printf("RedisProxy.Get.MISS: got value:%s for key:%s", v, key)
		return false, v
	}

}

// Function to handle Get requests
func (s *RedisProxy) GetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")

	if key == "" {
		log.Println("[RedisProxy.GetHandler]: Key is empty")
		return
	}
	log.Println("[RedisProxy.GetHandler]: Searching for key", key)
	isHit, v := s.Get(key)

	var cacheStatus = "MISS";
	if isHit {
		cacheStatus = "HIT"
	}

	fmt.Fprintf(w, "%s with cache %s\n", v, cacheStatus)
}

// Function to handle Set requests for a key/value pair
func (s *RedisProxy) SetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value := params.ByName("value")

	if key == "" || value == "" {
		log.Println("[RedisProxy.SetHandler]: Key/value is empty")
		return
	}
	log.Println("[RedisProxy.GetHandler]: Setting key", key)
	ok := s.RedisHandler.Set(key, value)

	if !ok {
		fmt.Fprintf(w, "Could not set value:%s against key:%s\n", value, key)
	}

	fmt.Fprintf(w, "Successfully set value %s against key %s\n", value, key)
}