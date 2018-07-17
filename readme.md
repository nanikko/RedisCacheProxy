# Overview
The system supports worker pool and LRU Cache

The LRU cache being used is a wrapper over an external library that supports moving the entries based on update/access times. However, expiration
is not supported. On cache miss, proxy searches Redis and updates the value against the access key in cache

It supports worker pool in that a fixed number of workers are used to perform the submitted jobs

## Tests:
Manually tested using http endpoints. However, there are tests to verify cache evictions and redis proxy GET happy paths


## Algorithmic complexity of the cache operations.
O(1). Lru Cache internally uses map to store.


## Steps to run unit tests.
```
make deps
make test
```

### Run:
docker-compose build
docker-compose up

### Options
| `port` | Server Port (defaults to 10000 ) |
| `protocol` | Whether the request can be Http based or RESP based request |
| `concurrentJobs` | Max number of concurrent connections allowed |
| `workers` | Max number of requests can be executed in parallel |
| `capacity` | Proxy cache capacity (defaults to 100) |
| `expiry` | Proxy cache global expiry in second (defaults to 10 sec) |
| `redisIpAndPort` | Redis Ip and Port. Default is localhost:6379 |

## Code walkthrough

```
Dockerfile  // Dockerfile for go app
Makefile
README.md
docker-compose.yaml  // Runs redis container and connects main go app
cache // Directory that wraps minimal LRU cache logic
proxy  // Have proxy implementations for both http and Tcp
redis // Has abstractions over Redis Http and Tcp calls
redis/worker // Contains logic for scheduling the work with worker pool
main.go
```

### Testing steps:
```
After running

make deps
make test

puts a value against key - "key" and that can be verified using the following command:

echo "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n" netcat localhost 10000
```

### How long you spent on each part of the project.
Coding 5h
Unit Testing 0.5h
Docker,makefile, integration testing 1.5h

Had minimal exposure to Golang and that perhaps has led to more time spent and also to research on the ways to enqueue the jobs for the Bonus
requirement



