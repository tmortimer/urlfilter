# urlfilter
REST service to filter malicious URLs.

# Requirements
```
github.com/google/go-cmp/cmp
github.com/gomodule/redigo/redis
```

# Docker Compose
## Requirements
[Docker](https://www.docker.com/get-started)
[Docker Compose](https://docs.docker.com/compose/)

## Run With The Fake Filter
docker-compose -f docker/docker-compose-fake.yaml up -d
...
docker-compose -f docker/docker-compose-fake.yaml stop

## Run With Redis Filter
docker-compose -f docker/docker-compose-redis.yaml up -d
...
docker-compose -f docker/docker-compose-redis.yaml stop

# Testing
Run **go test ./..** to run unit tests.

# Docs
[![GoDoc](https://godoc.org/github.com/tmortimer/urlfilter?status.svg)](https://godoc.org/github.com/tmortimer/urlfilter)
