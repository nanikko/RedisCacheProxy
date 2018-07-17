FROM golang:1.9

ADD . /go/src/RedisProxy

RUN go get github.com/go-redis/redis
RUN go get github.com/hashicorp/golang-lru
RUN go get github.com/julienschmidt/httprouter

RUN go test -v RedisProxy/cache
RUN go install -v ./...

ENTRYPOINT ["/go/bin/RedisProxy", "-redisIpAndPort=172.17.0.1:6379","-expiry=10","-capacity=100","-port=10000","-concurrentJobs=1000","-workers=100"]

EXPOSE 10000