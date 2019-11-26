package redis

import (
	"encoding/json"
	"github.com/Els-y/coupons/server/pkgs/setting"
	"github.com/gomodule/redigo/redis"
	"time"
)

var redisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	redisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func Set(key string, data interface{}) error {
	conn := redisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := redisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := redisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func Incr(key string) (int, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

func Decr(key string) (int, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redis.Int(conn.Do("DECR", key))
}
