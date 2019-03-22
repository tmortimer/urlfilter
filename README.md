# urlfilter
REST service to filter malicious URLs.

[![GoDoc](https://godoc.org/github.com/tmortimer/urlfilter?status.svg)](https://godoc.org/github.com/tmortimer/urlfilter)

# The Basics
The application accepts requests against **/urlinfo/1/somefull.url.com/here?query=string** and returns **200-OK** if they are safe to visit and **403-Forbidden** if the URL has been flagged and should not be visited. It is a standard Golang application that can be executed from it's project folder by running **go run urlfilter.go --config=configs/bloom-redis-mysql.json**, hower it is simplest to run this with Docker Compose.

### Sample Request-Response
```
curl 'http://localhost:8080/urlinfo/1/wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism' -v
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /urlinfo/1/wsxzsal8.club/crackle/rebute/perfusion/outspill?rodomontade=reg&scolecophagous=militarism HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 403 Forbidden
< Date: Fri, 22 Mar 2019 05:11:21 GMT
< Content-Length: 0
<
```

### Notes About The API
The requirements stated *"The caller wants to know if it is safe to access that URL or not. As the implementer you get to choose the response format and structure. These lookups are blocking users from accessing the URL until the caller receives a response from your service."*.

Based on this I kept the API as simple as possible. It may be that additional information is stored and retrievable about any flagged URL, but for the use outlined we're looking for a simple and fast YES or NO. As a result I limited the information to straightforward status codes.

A separate endpoint, or even service, could provide additional information where necessary.

### Notes About The URL Format
After an initial question about query strings I decided to handle the URL passed in verbatim. It's very possible/likely that in the real world this would be insufficient. Based on how I set up the Filter Chain (more on this later) one could add a filter to the front of the chain which sanitizes the URL, ie: consistently removes *www.*, switches to all lower case, adds or removes trailing slashes, etc.

Additionally the query strings are likely to change order, if not content and format. This could be handled by breaking them up and handling them separately from the URL.

My intuition is that you would actually want to completely break-up the URL so that you could block entire domains, only specific paths within domains, or only specific paths with specific query strings (not dependant on order). However without actual data and requirements it didn't make sense to pre-maturely tackle this given the timeframe/nature of the project.

# The Design
I've implemented a configurable Filter Chain so that different storage and retrieval techniques could be evaluated. Full disclosure I didn't evaluate them, I picked what seemed the most complete and flexible. In the real world you would need to collect or estimate load patterns and run load tests against different configurations.

The filters are chained together, and are accessed one after the other, as necessary to figure out if a URL is flagged. Please see each filter type for more details about how they can fit into the chain. Filters are called in the order they appear in the list.

## Filters
### Fake
***Use:*** Add "fake" to the ["filters" config list](configs/sample-config-defaults.json#L4).

This returns that a URL is found if it has "facebook" anywhere in it. This seems like a good thing to block ;). The next filter in the chain is ignored.

### Redis
***Use:*** Add "redis" to the ["redis" config list](configs/sample-config-defaults.json#L4). Configure the ["redis"](configs/sample-config-defaults.json#L7) section of the config for your Redis instance.

A Redis based filter. This could be local or remote. It would be possible to run even a distributed collection of *urlfilter* workers against a single Redis cluster. This might be a totally sufficient setup, but you'd need to load test it, evaluate latency characteristics, etc.

If there are more filters after the Redis based filter, this will work like a cache. If it's not found it will check the next filter down the line. If implemented as a cache **maxmemory** and **maxmemory-policy** should probably be set on Redis config to control the cache behavior.

### MySQL
***Use:*** Add "mysql" to the ["filters" config list](configs/sample-config-defaults.json#L4). Configure the ["mysql"](configs/sample-config-defaults.json#L16) section of the config for your MySQL instance.

MySQL based filter. This can also be configured as a cache, however it makes more sense as the final stop in the filter chain.

The URL itself is not used as an index, rather a CRC of the URL is computed and stored as the index. This way when searching for a given URL in the database the row is found using an integer based key. Even if there are collisions they should be relatively infrequent. I've implemented this with CRC32, but it would be worth loading real data and measuring the frequency and depth of collisions. It may be worth using CRC64 or another hash all together.

### Bloom Filter
***Use:*** Add "redismysqlbloom" to the ["filters" config list](configs/sample-config-defaults.json#L4). Configure the ["redismysqlbloom"](configs/sample-config-defaults.json#L16) section of the config, including the nested ["redis"](configs/sample-config-defaults.json#L23) and ["mysql"](configs/sample-config-defaults.json#L32) sections.

A Redis based [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter). This should be used as the first filter in the chain, or the benefit is lost. Additionally it can not be the last filter in the chain.

The idea behind a bloom filter is that you can test for existance in an arbitrarily large set of data, without consuming an arbitrarily large amount of memory. Bloom Filters will generate false positives, that means if it's found in the Bloom Filter the next filter in the chain should be checked. They do not however generate false negatives. If a URL is not found in the Bloom Filter, that result can be returned directly.

The Bloom Filter is bypassed until initial data loading is complete. Currently the only supported method of data loading is from MySQL.

The default behavior is to check for new data every one minute. This can be configured. Until new data is loaded into the Bloom Filter, it's as if the data is not present in the DB either, it will still return not found.

The Bloom Filter is configured for 1000000 items out of the box, this can be changed through the config file.

## Default Configuration
The "default" configuration I have settled on, and packaged with Docker Compose, is **Bloom Filter->Redis Cache->MySQL**. We can quickly find out if a URL has not been flagged. However if it is found in the Bloom Filter we then check the Redis Cache, if it's there we can return. If it's not there then we need to check MySQL. This si the final stop and will provide the answer returned to the client. On the way back the URL will be inserted into the Redis Cache.

This addresses several of the stated questions/requirements:

1. **The size of the URL list could grow infinitely, how might you scale this beyond the memory capacity of the system.**

   The Bloom Filter has a relatively small memory footprint even for a large set of data. The Redis cache can then be tuned appropriately. The backing MySQL instance can be whatever it needs to be to support the data set.

   If it was found that a backing Redis cluster was the way to go, this could similarly be built out however large was deemed necessary.

2. **Assume that the number of requests will exceed the capacity of a single system, describe how might you solve this, and how might this change if you have to distribute this workload over geographic regions.**

   With a shared back end, be it Redis, MySQL, or otherwise workers can be scaled as necessary. From there they would be placed behind a load balancer, software or hardware based. Even inside a large Kubernetes deployment. Both Redis and MySQL support multi-site deployments so distributed geographic regions would be handled that way.

   As the Bloom Filter is bypassed until it's loaded the existing data there's no concern about inconsistant results if that's the head of your chain. There would be some additional hurdles with updating Bloom Filters on an army of workers to keep them in sync, but that will be addressed below.

3. **What are some strategies you might use to update the service with new URLs? Updates may be as much as 5 thousand URLs a day with updates arriving every 10 minutes.**

   Given the time frame I stopped short of implementing this, other than the MySQL loader necessary for testing. I would keep the URL loading separate from the main application. Data would be loaded into whatever the final data store is, Redis or MySQL, singular or distributed.

   The catch here is an army of workers with their own Bloom Filters. A Redis based Bloom Filter could be shared, but let's assume the Bloom Filter needs to be local to the worker. Having them all update at their own cadence would introduce inconsistent results.

   This could be addressed by routing traffic based on it's source, rather than in a sequential round robin fashion. IE the same client always hits the same worker. This may be needlessly complex and inflexible if workers go down etc.

   Alternatively I would actually have the data loading process generate a new Bloom Filter locally and then the push it out to the workers triggering the update to the new data. This has the added bonus of easily changing the parameters of the Bloom Filter as the data set grows.

## Available Config Options
Only config options you need to change need to be specified in the file.

[Sample Config File With Defauls](configs/sample-config-defaults.json)

[Bloom-Redis-MySQL Config Used For Docker Compose Execution](configs/bloom-redis-mysql.json)

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

# What's Missing?
In no particular order.

1. Integration tests
2. The mechanism discussed above to do the data loading
3. Proper logging
4. Metrics/statistics collection and reporting
5. Proper error handling and reporting

