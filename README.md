# urlfilter
REST service to filter malicious URLs.

[![GoDoc](https://godoc.org/github.com/tmortimer/urlfilter?status.svg)](https://godoc.org/github.com/tmortimer/urlfilter)

## The Basics
The application accepts requests against **/urlinfo/1/somefull.url.com/here?query=string** and returns **200-OK** if they are safe to visit and **403-Forbidden** if the URL has been flagged and should not be visited. It is a standard Golang application that can be executed from it's project folder by running **go run urlfilter.go --config=configs/bloom-redis-mysql.json**, hower it is simplest to run this with [Docker Compose](#Docker Compose)

### Notes About The API
The requirements stated *"The caller wants to know if it is safe to access that URL or not. As the implementer you get to choose the response format and structure. These lookups are blocking users from accessing the URL until the caller receives a response from your service."*.

Based on this I kept the API as simple as possible. It may be that additional information is stored and retrievable about any flagged URL, but for the use outlined we're looking for a simple and fast YES or NO. As a result I limited the information to straightforward status codes.

A separate endpoint, or even service, could provide additional information where necessary.

### Notes About The URL Format
After an initial question about query strings I decided to handle the URL passed in verbatim. It's very possible/likely that in the real world this would be insufficient. Based on how I set up the Filter Chain (more on this later) one could add a filter to the front of the chain which sanitizes the URL, ie: consistently removes *www.*.

Additionally the query strings are likely to change order, if not content and format. This could be handled by breaking them up and handling them separately from the URL.

My intuition is that you would actually want to completely break-up the URL so that you could block entire domains, only specific paths within domains, or only specific paths with specific query strings (not dependant on order). However without actual data and requirements it didn't make sense to pre-maturely tackle this given the timeframe/nature of the project.

## The Design
I've implemented a configurable Filter Chain so that different storage and retrieval techniques could be evaluated. Full disclosure I didn't evaluate them, I picked what seemed the most complete and flexible. In the real world you would need to collect or estimate load patterns and run load tests against different configurations.

# Requirements
## Golang
[Installing Golang](https://golang.org/doc/install)

### Golang Dependencies
```
github.com/google/go-cmp/cmp
github.com/gomodule/redigo/redis
github.com/tjarratt/babble
```

## Docker Compose
Running with Docker Compose removes the need to install, configure, and manage MySQL and Redis instances for testing. Using the provided docker-compose file individual containers will be spun up for the application, MySQL, the Redis cache, and the Redis Bloom Filter.

[Docker](https://www.docker.com/get-started)\
[Docker Compose](https://docs.docker.com/compose/) - Should come with Docker installation.

### Run With Docker Compose
Passing -d to the up command causes the containers to run in the background. If you want to see logs don't do this. Run these from the urlfitler directory. For some reason the compose dependencies aren't quite working right, so the first up may take to long to initialize MySQL and the main app container will fail. Stopping them and bringing them up again will solve that.

Stop will stop the containers, down will remove them.

The data will persist until you run the docker volume command.

Run the following from the root urlfilter directory where the docker-compose.yaml file is found.
```
docker-compose up [-d]
...
curl 'http://localhost:8080/urlinfo/1/wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism' -v
...
docker-compose logs bloom-redis-mysql

docker-compose stop

docker-compose down

docker volume rm $(docker volume ls -q)
```

## MySQL (Not Needed With Docker Compose Workflow)
[Getting Started With MySQL](https://dev.mysql.com/doc/mysql-getting-started/en/)

## Redis (Not Needed With Docker Compose Workflow)
[Redis Quickstart](https://redis.io/topics/quickstart)

# Populate Some Data
Run the following with the docker-compose containers up and running. This will take the domains and add 0-5 path components, and 0-3 query components. It will populate more than 26k urls. It spams the console so you know what to test with. It uses a domain list from [http://mirror1.malwaredomains.com/files/domains.txt](http://mirror1.malwaredomains.com/files/domains.txt) and then adds 0-5 path components, and 0-3 query components generated with [github.com/tjarratt/babble](github.com/tjarratt/babble).
```
go run mysqlloader/mysqlloader.go --config=configs/mysql-loader.json --list mysqlloader/domains-only.txt -mpdepth 5 -mqdepth 3
```

# Some Examples Of Key Application Functionality In Action
## Filter Chaining
```
bloom-redis-mysql_1          | 2019/03/21 04:56:10 The Bloom Filter loaded an additional 26839 urls for a total of 26839.
bloom-redis-mysql_1          | 2019/03/21 04:56:14 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1          | 2019/03/21 04:56:15 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in MySQL.
bloom-redis-mysql_1          | 2019/03/21 04:56:15 Adding URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism to Redis cache.
bloom-redis-mysql_1          | 2019/03/21 04:56:32 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1          | 2019/03/21 04:56:32 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis cache.
```

## Waiting Until The Filter Is Ready On A New Worker
```
bloom-redis-mysql_1             | 2019/03/21 18:08:58 Redis Bloom Filter is not yet loaded, checking the next filter.
bloom-redis-mysql_1             | 2019/03/21 18:08:58 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in MySQL.
bloom-redis-mysql_1             | 2019/03/21 18:08:58 Adding URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism to Redis cache.
bloom-redis-mysql_1             | 2019/03/21 18:09:07 The Bloom Filter loaded 53678 urls for a total of 53678.
bloom-redis-mysql_1             | 2019/03/21 18:09:11 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1             | 2019/03/21 18:09:11 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis cache.
```

## Bloom Filter Updates
```
bloom-redis-mysql_1             | 2019/03/22 05:11:11 Redis Bloom Filter is not yet loaded, checking the next filter.
bloom-redis-mysql_1             | 2019/03/22 05:11:11 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in MySQL.
bloom-redis-mysql_1             | 2019/03/22 05:11:11 Adding URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism to Redis cache.
bloom-redis-mysql_1             | 2019/03/22 05:11:18 The Bloom Filter loaded 53678 urls for a total of 53678.
bloom-redis-mysql_1             | 2019/03/22 05:11:21 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1             | 2019/03/22 05:11:21 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis cache.
bloom-redis-mysql_1             | 2019/03/22 05:12:11 The Bloom Filter loaded 6000 urls for a total of 59678.
bloom-redis-mysql_1             | 2019/03/22 05:13:14 The Bloom Filter loaded 18000 urls for a total of 77678.
bloom-redis-mysql_1             | 2019/03/22 05:14:08 The Bloom Filter loaded 2839 urls for a total of 80517.
bloom-redis-mysql_1             | 2019/03/22 05:15:08 The Bloom Filter loaded 0 urls for a total of 80517.
```

## Performance Of Querying MySQL With The CRC Index
```
mysql> explain select url from crcurls where url_crc=1076669273 AND url="9oxigfyv1n.bradul.creatory.org";
+----+-------------+---------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
| id | select_type | table   | partitions | type | possible_keys | key     | key_len | ref   | rows | filtered | Extra       |
+----+-------------+---------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
|  1 | SIMPLE      | crcurls | NULL       | ref  | url_crc       | url_crc | 4       | const |    1 |    10.00 | Using where |
+----+-------------+---------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
```

After adding two rounds of URLs, based off the same domain names, there were 126 collisions, and nothing more than 2 deep. With more URLs in the picture you could switch to CRC64.

# Unit Tests
Run **go test ./..** to run unit tests.

## Coverage
```
go test ./... -cover
?   	github.com/tmortimer/urlfilter	[no test files]
ok  	github.com/tmortimer/urlfilter/config	0.011s	coverage: 94.4% of statements
?   	github.com/tmortimer/urlfilter/connectors	[no test files]
ok  	github.com/tmortimer/urlfilter/filters	68.064s	coverage: 90.0% of statements
ok  	github.com/tmortimer/urlfilter/handlers	0.015s	coverage: 100.0% of statements
?   	github.com/tmortimer/urlfilter/mysqlloader	[no test files]
ok  	github.com/tmortimer/urlfilter/server	0.016s	coverage: 80.0% of statements
```
