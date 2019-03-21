# urlfilter
REST service to filter malicious URLs.

# Requirements
```
github.com/google/go-cmp/cmp
github.com/gomodule/redigo/redis
github.com/tjarratt/babble
```

# Filter Chaining
Blurb about filter chain...
```
bloom-redis-mysql_1          | 2019/03/21 04:56:10 The Bloom Filter loaded an additional 26839 urls for a total of 26839.
bloom-redis-mysql_1          | 2019/03/21 04:56:14 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1          | 2019/03/21 04:56:15 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in MySQL.
bloom-redis-mysql_1          | 2019/03/21 04:56:15 Adding URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism to Redis cache.
bloom-redis-mysql_1          | 2019/03/21 04:56:32 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1          | 2019/03/21 04:56:32 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis cache.
```

# Waiting Until The Filter Is Ready On A New Worker
```
bloom-redis-mysql_1             | 2019/03/21 18:08:58 Redis Bloom Filter is not yet loaded, checking the next filter.
bloom-redis-mysql_1             | 2019/03/21 18:08:58 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in MySQL.
bloom-redis-mysql_1             | 2019/03/21 18:08:58 Adding URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism to Redis cache.
bloom-redis-mysql_1             | 2019/03/21 18:09:07 The Bloom Filter loaded 53678 urls for a total of 53678.
bloom-redis-mysql_1             | 2019/03/21 18:09:11 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis Bloom Filter, checking the next filter.
bloom-redis-mysql_1             | 2019/03/21 18:09:11 URL wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism found in Redis cache.
```

# Docker Compose
## Requirements
[Docker](https://www.docker.com/get-started)\
[Docker Compose](https://docs.docker.com/compose/)

## Run With Docker Compose
Passing -d to the up command causes the containers to run in the background. If you want to see logs don't do this. Run these from the urlfitler directory. For some reason the compose dependencies aren't quite working right, so the first up may take to long to initialize MySQL and the main app container will fail. Stopping them and bringing them up again will solve that.

Stop will stop the containers, down will remove them.

The data will persist until you run the docker volume command.
```
docker-compose up [-d]
...
curl 'http://localhost:8080/urlinfo/1/www.facebook.com:9090/peww/what/who/merp.html?face=ac&w' -v
...
docker-compose logs bloom-redis-mysql

docker-compose stop

docker-compose down

docker volume rm $(docker volume ls -q)
```

# Populate Some Data
Run the following with the docker-compose containers up and running. This will take the domains and add 0-5 path components, and 0-3 query components. It will populate more than 26k urls. It spams the console so you know what to test with.
```
go run mysqlloader/mysqlloader.go --config=configs/mysql-loader.json --list mysqlloader/domains-only.txt -mpdepth 5 -mqdepth 3
```

# Performance Of Querying MySQL With The CRC Index
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

# Docs
[![GoDoc](https://godoc.org/github.com/tmortimer/urlfilter?status.svg)](https://godoc.org/github.com/tmortimer/urlfilter)
