package main

import (
	"flag"
	"log"
	"os"
	"github.com/julienschmidt/httprouter"
	"RedisProxy/proxy"
	"RedisProxy/cache"
	"RedisProxy/redis"
	"net/http"
)

// redis client protocol
var protocol = flag.String("protocol", "http", "Redis communication protocol. Default is RESP")
var port = flag.String("port", "10000", "Server Port")
var redisAddr = flag.String("redisIpAndPort", "172.17.0.1:6379", "Redis Ip and Port. Default is localhost:6379")
// cache vars
var expiryTime = flag.Int64("expiry", 10, "cache global expiry in second")
var capacity = flag.Int("capacity", 100, "cache capacity")
var concurrentJobs = flag.Int("concurrentJobs", 100, "Max number of concurrent jobs allowed")
// scheduler vars
var numWorkers = flag.Int("workers", 10, "Max number of requests can be executed in parallel")


func main() {
	flag.Parse()

	cache, err := cache.NewMinimalCache(*capacity, *expiryTime)

	if err != nil {
		log.Println("Error when creating cache.")
	}

	redisHandler, err := redis.NewRedisHandler(*redisAddr)
	redisHandlerResp, err := redis.NewRedisHandlerResp(*redisAddr)

	if err != nil {
		log.Println("Error when creating Redis handlers", err.Error())
		os.Exit(1)
	}

	// Run Scheduler in the background - Ideally this can as well go into the proxy module
	scheduler := redis.NewScheduler(*numWorkers, *concurrentJobs)
	go scheduler.Run()

	if *protocol == "http" {
		log.Println("------- Starting Http RedisProxy ------")
		rp := proxy.RedisProxy{cache, redisHandler, redisHandlerResp, scheduler}

		// settings paths for key/value getters and setts
		router := httprouter.New()
		router.GET("/get/:key", rp.GetHandler)
		router.GET("/set/:key/:value", rp.SetHandler)

		log.Println("RedisProxy successfully started on port ", *port)

		// Configure httpRouter against the port
		log.Fatal(http.ListenAndServe(":"+*port, router))
	} else {
		log.Println("------ Starting Tcp RedisProxy following RESP protocol for communications -----")
		rpTcp := proxy.RedisProxyTcp{
			Cache : cache,
			RedisHandlerResp : redisHandlerResp,
			Scheduler : scheduler,
			Port : ":"+ *port,
		}

		// Listen to requests
		rpTcp.Listen()
	}

}
