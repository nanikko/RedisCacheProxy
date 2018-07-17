# Overview
The system supports worker pool and LRU Cache

The LRU cache being used is a wrapper over an external library that supports moving the entries based on update/access times. However, expiration
is not supported. On cache miss, proxy searches Redis and updates the value against the access key in cache

It supports worker pool in that a fixed number of workers are used to perform the submitted jobs

## Tests:
Manually tested using http endpoints. However, there are tests to verify cache evictions and redis proxy GET happy paths


## Algorithmic complexity of the cache operations.
O(1). Lru Cache internally uses map to store.


## Steps to run.
```
make deps
make test
```

### Run:
docker-compose up

If you just want to run the proxy in docker, and mannualy point to a redis backend:
docker run -p 10000:10000 myproxy -redisIpAndPort=127.0.0.1:6379

### Options

| name | descr |
|---|---|
| `capacity` | Proxy cache capacity (defaults to 100) |
| `port` | Server Port (defaults to 10000 ) |
| `concurrentJobs` | Max number of concurrent connections allowed |
| `workers` | Max number of requests can be executed in parallel |
| `protocol` | Whether the request can be Http based or RESP based request |
| `redisIpAndPort` | Redis Ip and Port. Default is localhost:6379 |
| `expiry` | Proxy cache global expiry in second (defaults to 10 sec) |


### Test by using curl:
curl -X G http://<IP>:<Port>/<your_word>

### How long you spent on each part of the project.
Coding 2.5h
Unit Testing 0.5h
Docker,makefile, integration testing 1h

### A list of the requirements that you did not implement and the reasons for omitting them.
Redis client protocol - Prefer http way. Easy to use, test, no integration effort.



