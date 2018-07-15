package service

import (
	"fmt"
    "log"
	"net/http"
	"redisproxy/util"
	"redisproxy/infra"
	"github.com/julienschmidt/httprouter"
)

type ProxyService struct {
	Cache *ProxyCache
	RedisHandler *infra.RedisHandler
	RedisHandlerResp *infra.RedisHandlerResp
	Scheduler *util.Scheduler
}

func (s *ProxyService) Get(key string) (bool, string) {
	found, v := s.Cache.Get(key)

	if found {
		log.Printf("ProxyService.Get.HIT: got value:%s for key:%s", v, key)
		return true, v
	} else {
		resp := make(chan string)
		defer close(resp)
		work := util.Job{
			Request: key,
			JobHandler: s.RedisHandlerResp,
			Resp: resp,
		}
		s.Scheduler.JobQueue <- work
		
		v = <-resp
		if v != "" {
			s.Cache.Add(key, v)
		}
		log.Printf("ProxyService.Get.MISS: got value:%s for key:%s", v, key)
		return false, v
	}

}

func (s * ProxyService) GetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("In gethandler method")
	key := params.ByName("key")

	if key == "" {
		log.Println("[ProxyService.GetHandler]: Key is empty")
		return
	}
	log.Println("[ProxyService.GetHandler]: Searching for key", key)
	isHit, v := s.Get(key)

	var cacheStatus = "MISS";
	if isHit {
		cacheStatus = "HIT"
	}
	fmt.Fprintf(w, "%s with cache %s\n", v, cacheStatus)
}

func (s * ProxyService) SetHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value := params.ByName("value")

	if key == "" || value == "" {
		log.Println("[ProxyService.SetHandler]: Key/value is empty")
		return
	}
	log.Println("[ProxyService.GetHandler]: Setting key", key)
	ok := s.RedisHandler.Add(key, value)

	if !ok {
		fmt.Fprintf(w, "Could not set value:%s against key:%s\n", value, key)
	}

}