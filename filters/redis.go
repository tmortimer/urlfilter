package filters

import (
	"log"
	"github.com/gomodule/redigo/redis"
	"github.com/tmortimer/urlfilter/config"
)

// Redis based filter. Depending on config can be used as a cache.
// If this is being used as a cache a secondary filter must be set.
type Redis struct {
	next Filter
	config config.Redis
	conn redis.Conn
}

func NewRedis(config config.Redis) *Redis {
	return &Redis{
		config: config,
	}
}

func (f *Redis) ConnectToRedis() {
	conn, err := redis.Dial("tcp", f.config.Host + ":" + f.config.Port)
	if err != nil {
	    log.Fatalf("Unable to connect to Redis instance %s:%s", f.config.Host, f.config.Port)
	}
	f.conn = conn
	//TOM need to do something about properly closing this connection on application exit.
}

func (f *Redis) AddSecondaryFilter(filter Filter) {
	f.next = filter
}

func (f *Redis) ContainsURL(url string) bool {
	found := f.RedisContainsURL(url)

	if found || f.next == nil {
		return found
	}

	return f.next.ContainsURL(url)
}

func (f *Redis) RedisContainsURL(url string) bool {
	return true
}
