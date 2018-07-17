test: docker
	go test .

docker:
	docker-compose up -d

deps:
	go get "github.com/go-redis/redis"
	go get github.com/julienschmidt/httprouter
	go get github.com/hashicorp/golang-lru
