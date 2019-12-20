package redis

import (
	"encoding/json"
	"github.com/Els-y/coupons/server/pkgs/setting"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

var redisConn *redigo.Pool

// Setup Initialize the Redis instance
func Setup() error {
	redisConn = &redigo.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", setting.RedisSetting.Host)
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
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
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

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := redisConn.Get()
	defer conn.Close()

	exists, err := redigo.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := redisConn.Get()
	defer conn.Close()

	reply, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redigo.Bool(conn.Do("DEL", key))
}

func Incr(key string) (int, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redigo.Int(conn.Do("INCR", key))
}

func Decr(key string) (int, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redigo.Int(conn.Do("DECR", key))
}

func IncrBy(key string, value int) (int, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redigo.Int(conn.Do("INCRBY", key, value))
}

func SAdd(key, value string) (bool, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redigo.Bool(conn.Do("SADD", key, value))
}

func SIsmember(key, value string) (bool, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redigo.Bool(conn.Do("SISMEMBER", key, value))
}
