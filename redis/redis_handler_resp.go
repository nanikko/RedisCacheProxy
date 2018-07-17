package redis

import (
	"bufio"
	"net"
	"log"
	"strings"
	"fmt"
)

// RedisHandler implements the JobHandler interface
type RedisHandlerResp struct {
	addr string
}

// NewRedisHandlerResp creates new redis handler
func NewRedisHandlerResp(addr string) (*RedisHandlerResp, error) {
	redisHandlerResp := &RedisHandlerResp{
		addr:   addr,
	}

	return redisHandlerResp, nil
}

// Gets the data based on key
func (s *RedisHandlerResp) Get(cmdStr string) string {
	rw, _ := s.open(s.addr)

	rw.WriteString(cmdStr)
	rw.Flush()
	log.Println("in Resp Get method after flushing")

	b, err := s.readBulkString(rw)

	if err != nil {
		log.Println("Error in getting from TCP connection:", err.Error())
		return ""
	}

	return b
}

// Get gets the data based on key
func (s *RedisHandlerResp) Set(cmdStr string) bool {
	rw, _ := s.open(s.addr)
	rw.WriteString(cmdStr)
	rw.Flush()

	log.Println("in Resp Set method after flushing")

	status, err := s.readResp(rw)

	if status != "+OK" || err != nil {
		log.Println("Error in setting :", err.Error())
		return false
	}

	return true
}

func (s *RedisHandlerResp) readResp(r *bufio.ReadWriter) (string, error) {
	// ignoring this value as it is the count of number of bytes in the value
	result, err := r.ReadString('\r')
	if err != nil {
		return "", err
	}

	return result[:len(result)-1], nil
}

func (s *RedisHandlerResp) readBulkString(r *bufio.ReadWriter) (string, error) {
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
	resultStr := fmt.Sprintf("%s\n%s\r\n", resultChars, resultVal)
	return resultStr, nil
}

// returns substring between two strings
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


func (s *RedisHandlerResp) open(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Error in opening TCP connection:", err.Error())
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}