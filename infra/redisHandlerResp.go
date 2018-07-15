package infra

import (
	"bufio"
	"net"
	"fmt"
	"log"
	"io"
	"bytes"
	"strconv"
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
func (s *RedisHandlerResp) Get(key string) string {
	log.Println("in advanced Get method", s.addr)
	//conn, err := net.Dial("tcp", s.addr)
	//defer conn.Close()
	rw, _ := s.Open(s.addr)
	cmdStr := "*2\r\n$3\r\nGET\r\n$" + fmt.Sprint(len(key)) + "\r\n" + key + "\r\n"

	// send to socket
	//fmt.Fprintf(conn, cmdStr)
	// listen for reply
	//message, err := bufio.NewReader(conn).ReadString('\n')
	rw.WriteString(cmdStr)
	rw.Flush()
	log.Println("in Resp Get method after flushing")

	line, _, err := rw.ReadLine()
	if err != nil {
		return ""
	}

	b, err := s.readBulkString(rw, line)

	//resp, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Error in getting from TCP connection:", err.Error())
		return ""
	}

    log.Println("finally read string is : ", b)
	return string(b[:])
}

func (s *RedisHandlerResp) readBulkString(r *bufio.ReadWriter, line []byte) ([]byte, error) {
	log.Println("in readBulkString method after flushing")

	end := bytes.IndexByte(line, '\r')
	log.Println("in readBulkString atoi:", 	string(line[1:end]))

	count, err := strconv.Atoi(string(line[1:end]))

	if err != nil {
		return nil, err
	}

	log.Println("count is : ", count)

	buf := make([]byte, len(line)+count+2)
	copy(buf, line)

	log.Println("uffer is copied to buf from line ")

	_, err = io.ReadFull(r, buf[len(line):])
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *RedisHandlerResp) Open(addr string) (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Error in opening TCP connection:", err.Error())
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}