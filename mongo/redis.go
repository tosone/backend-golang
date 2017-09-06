package mongo

import (
	"github.com/garyburd/redigo/redis"
)

// RedisPool redispool
var RedisPool *redis.Pool

func init() {
	RedisPool = redis.NewPool(func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1") }, 10)
}
