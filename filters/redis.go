package filters

import (
	"github.com/gomodule/redigo/redis"
	"github.com/tmortimer/urlfilter/config"
	"log"
	"time"
)

type Runner interface {
	Do(cmd string, keysAndArgs ...interface{}) (interface{}, error)
}

type RedisRunner struct {
	pool *redis.Pool
}

func NewRedisRunner(config config.Redis) *RedisRunner {
	runner := &RedisRunner{}
	runner.pool = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", config.Host+":"+config.Port) },
	}

	runner.ConfigureRedis(config)

	return runner
}

// Configure Redis.
func (r *RedisRunner) ConfigureRedis(config config.Redis) {
	conn := r.pool.Get()
	defer conn.Close()

	for _, command := range config.Config {
		conn.Do(command)
	}
}

func (r *RedisRunner) Do(cmd string, keysAndArgs ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return conn.Do(cmd, keysAndArgs...)
}

// Redis based filter. Depending on config can be used as a cache.
// If this is being used as a cache a secondary filter must be set.
type Redis struct {
	next   Filter
	config config.Redis
	runner Runner
}

// Return a new Redis filter.
func NewRedis(config config.Redis) *Redis {
	return &Redis{
		config: config,
	}
}

// Connect to the configured Redis instance.
func (f *Redis) ConnectToRedis() {
	f.runner = NewRedisRunner(f.config)
}

// Add a secondary filter. Necessary if using Redis as a cache.
func (f *Redis) AddSecondaryFilter(filter Filter) {
	f.next = filter
}

// Return true if the URL is found in Redis. If it's not then return false
// if there are no further filters in the chain, otherwise call the next filter.
// If Redis generates an error and this is only a cache we can continue down the
// filter chain, since each subsequent level should have better information.
func (f *Redis) ContainsURL(url string) (bool, error) {
	found, err := redis.Bool(f.runner.Do("EXISTS", url))
	if err != nil {
		// Not sure what the state of found will be after a failed
		// call to the Redis library, so be sure it's false.
		found = false
		log.Printf("Redis generated an error on GET: %s", err.Error())
	}

	if found || f.next == nil {
		return found, err
	}

	// Not found in the Redis cache based filter, try the next one.
	found, err = f.next.ContainsURL(url)

	if found {
		// Add it to the cache. Use "" as the value since we only care about the key.
		_, err = f.runner.Do("SET ", url, " \"\"")
		log.Printf("Redis generated an error on SET: %s", err.Error())
	}

	return found, err
}
