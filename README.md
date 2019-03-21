# urlfilter
REST service to filter malicious URLs.

# Requirements
```
github.com/google/go-cmp/cmp
github.com/gomodule/redigo/redis
github.com/tjarratt/babble
```

# Docker Compose
## Requirements
[Docker](https://www.docker.com/get-started)\
[Docker Compose](https://docs.docker.com/compose/)

## Run With Docker Compose
Passing -d to the up command causes the containers to run in the background. If you want to see logs don't do this. Run these from the urlfitler directory. Stop will stop the containers, down will remove the.

The data will persist until you run the docker volume command.
```
docker-compose up [-d]
...
curl 'http://localhost:8080/urlinfo/1/www.facebook.com:9090/peww/what/who/merp.html?face=ac&w' -v
...
docker-compose stop

docker-compose down

docker volume rm $(docker volume ls -q)
```

# Populate Some Data
Run the following with the docker-compose containers up and running. This will take the domains and add 0-5 path components, and 0-3 query components. It will populate more than 26k urls. It spams the console so you know what to test with.
```
go run mysqlloader/mysqlloader.go --config=configs/mysql-loader.json --list mysqlloader/domains-only.txt -mpdepth 5 -mqdepth 3
```

# Performance of Querying MySQL With the CRC Index
```
mysql> explain select url from crcurls where url_crc=1076669273 AND url="9oxigfyv1n.bradul.creatory.org";
+----+-------------+---------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
| id | select_type | table   | partitions | type | possible_keys | key     | key_len | ref   | rows | filtered | Extra       |
+----+-------------+---------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
|  1 | SIMPLE      | crcurls | NULL       | ref  | url_crc       | url_crc | 4       | const |    1 |    10.00 | Using where |
+----+-------------+---------+------------+------+---------------+---------+---------+-------+------+----------+-------------+
```

# Testing
Run **go test ./..** to run unit tests.

# Docs
[![GoDoc](https://godoc.org/github.com/tmortimer/urlfilter?status.svg)](https://godoc.org/github.com/tmortimer/urlfilter)
