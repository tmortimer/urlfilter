package filters

import (
	"github.com/gomodule/redigo/redis"
	"github.com/tmortimer/urlfilter/config"
	"log"
)

// Redis based filter. Depending on config can be used as a cache.
// If this is being used as a cache a secondary filter must be set.
type Redis struct {
	next   Filter
	config config.Redis
	conn   redis.Conn
}

func NewRedis(config config.Redis) *Redis {
	return &Redis{
		config: config,
	}
}

func (f *Redis) ConnectToRedis() {
	conn, err := redis.Dial("tcp", f.config.Host+":"+f.config.Port)
	if err != nil {
		log.Fatalf("Unable to connect to Redis instance %s:%s", f.config.Host, f.config.Port)
	}
	//TOM need to do something about properly closing this connection on application exit.
	f.conn = conn

	f.ConfigureRedis()
}

func (f *Redis) ConfigureRedis() {

}

func (f *Redis) AddSecondaryFilter(filter Filter) {
	f.next = filter
}

func (f *Redis) ContainsURL(url string) (bool, error) {
	// Since each next step in the chain has better information
	// we can go further down the chain if we have an error. Otherwise
	// return if we know it's found already, or if this is the end of
	// the line.
	found, err := f.RedisContainsURL(url)
	if err != nil {
		// Not sure what the state of found will be after a failed
		// call to the Redis library, so be sure it's false.
		found = false
	}

	if found || f.next == nil {
		return found, err
	}

	return f.next.ContainsURL(url)
}

func (f *Redis) RedisContainsURL(url string) (bool, error) {
	return true, nil
}
