package proxy

import (
	"net"
	"strings"
	"RedisProxy/redis"
	"RedisProxy/cache"
	"log"
	"bufio"
	"fmt"
)

type RedisProxyTcp struct {
	Cache   *cache.MinimalCache
	RedisHandler *redis.RedisHandler
	RedisHandlerResp *redis.RedisHandlerResp
	Scheduler *redis.Scheduler
	Workers int
	lock chan struct{}
	Port string
}

// Function that allows the proxy to continuously listen to new requests and asynchronously put them into the job queue
func (rptcp *RedisProxyTcp) Listen() {
	var pool = redis.NewWorkerPool(rptcp.Workers)
	ln, _ := net.Listen("tcp", rptcp.Port)

	// deferring closing the connection
	defer func() {
		ln.Close()
		log.Println("[RedisProxyTcp] listener closed]")
	}()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}

		w := redis.Work {
			Input: conn,
			Workable: func(c redis.InputData) redis.Result {
				rptcp.handleConnection(c.(net.Conn))
				return 0
			},
		}

		// Put the work into the pool
		pool.Input <- &w
	}
}

func (rptcp *RedisProxyTcp) handleConnection(conn net.Conn) {
	message, err := getStringFromReader(bufio.NewReader(conn))

	if err != nil {
		fmt.Println(err)
		return
	}

	_, resp := rptcp.Get(message)
	conn.Write([]byte(resp + "\n"))
}

func getStringFromReader(r *bufio.Reader) (string, error) {
	// ignoring this value as it is the count of number of bytes in the value
	resultChars, err := r.ReadString('\r')
	if err != nil {
		return "", err
	}


	resultVal, err := r.ReadString('\r')
	if err != nil {
		return "", err
	}

	resultVal = between(resultVal, "\n", "\r")

	resultStr := fmt.Sprintf("%s\r\n%s\r\n", resultChars, resultVal)
	return resultStr, nil
}

// Returns substring between two strings
func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func (rp *RedisProxyTcp) Get(key string) (bool, string) {
	log.Println("========= GET is ", key)
	_, cacheKey := rp.parseKey(key)
	found, v := rp.Cache.Get(cacheKey)

	if found {
		log.Printf("RedisProxyTcp.Get.HIT: got value:%s for key:%s", v, key)
		return true, v
	} else {

        resultString := rp.RedisHandlerResp.Get(key)
		return false, resultString
	}

}

func (rp *RedisProxyTcp) parseKey(str string) (bool, string) {
	words := strings.Fields(str)

	if strings.ToUpper(words[2]) == "GET" {
		return true, strings.Fields(str)[4]
	} else {
		return false, ""
	}
}

func (rp *RedisProxyTcp) ParseValue(str string) (string) {
	words := strings.Fields(str)

	return words[1]
}
