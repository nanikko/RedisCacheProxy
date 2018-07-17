test: docker
	go test ./cache
	go test ./proxy

docker:
	docker-compose build
	docker-compose up -d

deps:
	go get "github.com/go-redis/redis"
	go get github.com/julienschmidt/httprouter
	go get github.com/hashicorp/golang-lru
