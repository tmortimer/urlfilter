package connectors

import (
	"github.com/gomodule/redigo/redis"
	"github.com/tmortimer/urlfilter/config"
	"time"
)

const BF_NAME string = "URLFilter"

type ContainsFunc func(url string) (bool, error)
type AddFunc func(url string) error

// Holds the actual Redis connection pool and executes commands against Redis.
type Redis struct {
	// Redis connection pool.
	pool *redis.Pool

	// Redis specific config.
	config config.Redis

	// Function to check for existance of a URL.
	contains ContainsFunc

	// Function to add a URL.
	add AddFunc
}

// Create a new Redis connector and setup the Redis connection pool.
func NewRedisBase(config config.Redis) *Redis {
	connector := &Redis{
		config: config,
	}
	connector.pool = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", config.Host+":"+config.Port) },
	}

	connector.ConfigureRedis()

	return connector
}

// Create a new Redis connector.
func NewRedis(config config.Redis) *Redis {
	connector := NewRedisBase(config)
	connector.SetAccessors(false)
	return connector
}

// Create a new Redis Bloom Filter connector.
func NewRedisBloom(config config.Redis) *Redis {
	connector := NewRedisBase(config)
	connector.SetAccessors(true)
	return connector
}

// Setup the functions used to check if URLs exist, and add them.
// These change if this Redis connector is being used as for a
// Bloom Filter.
func (r *Redis) SetAccessors(bloom bool) {
	if bloom {
		r.contains = func(url string) (bool, error) {
			return redis.Bool(r.Do("BF.EXISTS", BF_NAME, url))
		}
		r.add = func(url string) error {
			_, err := r.Do("BF.ADD", BF_NAME, url)
			return err
		}
	} else {
		r.contains = func(url string) (bool, error) {
			return redis.Bool(r.Do("EXISTS", url))
		}
		r.add = func(url string) error {
			_, err := r.Do("SET", url, "\"\"")
			return err
		}
	}
}

// Configure Redis based on the supplied config file.
func (r *Redis) ConfigureRedis() {
	for _, command := range r.config.Config {
		r.Do(command)
	}
}

// Execute the command with arguments against the Redis connection pool.
func (r *Redis) Do(cmd string, keysAndArgs ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return conn.Do(cmd, keysAndArgs...)
}

// Check if the URL is in Redis.
func (r *Redis) ContainsURL(url string) (bool, error) {
	found, err := redis.Bool(r.Do("EXISTS", url))
	if err != nil {
		// Not sure what the state of found will be after a failed
		// call to the Redis library, so be sure it's false.
		found = false
	}

	return found, err
}

// Add the URL to the Redis. Only used if this DB is being used as a cache.
func (r *Redis) AddURL(url string) error {
	// Use "" as the value since we only care about the key.
	_, err := r.Do("SET", url, "\"\"")

	return err
}

// Return the name Redis for logging.
func (r *Redis) Name() string {
	return "Redis"
}
