package redis

import (
	"errors"
	"time"

	"qms.mgmt.api/base/log"

	"github.com/garyburd/redigo/redis"
)

const (
	EXPIRES_DEFAULT = time.Duration(0)
	EXPIRES_FOREVER = time.Duration(-1)
)

// RedisStore Wraps the Redis client to meet the Cache interface.
type RedisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

func (c *RedisStore) Set(key string, value interface{}, expires time.Duration) error {
	return c.invoke(c.pool.Get().Do, key, value, expires)
}

func (c *RedisStore) Get(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("GET", key)
	if raw == nil {
		return errors.New("Get: key not found.")
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return deserialize(item, ptrValue)
}

// Exists 检查kay是否存在
func Exists(conn redis.Conn, key string) bool {
	retval, _ := redis.Bool(conn.Do("EXISTS", key))
	return retval
}

// Delete 删除key
func (c *RedisStore) Delete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	if !Exists(conn, key) {
		log.Logger.Info("Delete: key not found.")
		return nil
	}
	_, err := conn.Do("DEL", key)
	return err
}

// Increment key值增加
func (c *RedisStore) Increment(key string, delta uint64) (uint64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	// Check for existance *before* increment as per the cache contract.
	// redis will auto create the key, and we don't want that. Since we need to do increment
	// ourselves instead of natively via INCRBY (redis doesn't support wrapping), we get the value
	// and do the exists check this way to minimize calls to Redis
	val, err := conn.Do("GET", key)
	if val == nil {
		return 0, errors.New("Get: key not found.")
	}
	if err == nil {
		currentVal, err := redis.Int64(val, nil)
		if err != nil {
			return 0, err
		}
		var sum int64 = currentVal + int64(delta)
		_, err = conn.Do("SET", key, sum)
		if err != nil {
			return 0, err
		}
		return uint64(sum), nil
	} else {
		return 0, err
	}
}

// Decrement key值减小
func (c *RedisStore) Decrement(key string, delta uint64) (newValue uint64, err error) {
	conn := c.pool.Get()
	defer conn.Close()
	// Check for existance *before* increment as per the cache contract.
	// redis will auto create the key, and we don't want that, hence the exists call
	if !Exists(conn, key) {
		return 0, errors.New("Get: key not found.")
	}
	// Decrement contract says you can only go to 0
	// so we go fetch the value and if the delta is greater than the amount,
	// 0 out the value
	currentVal, err := redis.Int64(conn.Do("GET", key))
	if err == nil && delta > uint64(currentVal) {
		tempint, err := redis.Int64(conn.Do("DECRBY", key, currentVal))
		return uint64(tempint), err
	}
	tempint, err := redis.Int64(conn.Do("DECRBY", key, delta))
	return uint64(tempint), err
}

// Flush 清空整个 Redis 服务器的数据(删除所有数据库的所有 key )
func (c *RedisStore) Flush() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

func (c *RedisStore) invoke(f func(string, ...interface{}) (interface{}, error),
	key string, value interface{}, expires time.Duration) error {

	switch expires {
	case EXPIRES_DEFAULT:
		expires = c.defaultExpiration
	case EXPIRES_FOREVER:
		expires = time.Duration(0)
	}

	b, err := serialize(value)
	if err != nil {
		return err
	}
	conn := c.pool.Get()
	defer conn.Close()
	if expires > 0 {
		_, err := f("SETEX", key, int32(expires/time.Second), b)
		return err
	} else {
		_, err := f("SET", key, b)
		return err
	}
}

// RedisCache redis连接池对象
var RedisCache *RedisStore

// NewRedisCache 初始化redis连接池对象，连接池存储在RedisCache
func NewRedisCache(host string, password string, defaultExpiration time.Duration) {
	defer log.Logger.Sync()

	var pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// the redis protocol should probably be made sett-able
			c, err := redis.Dial("tcp", host)
			if err != nil {
				log.Logger.Error("NewRedisCache error-->" + err.Error())
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					log.Logger.Error("NewRedisCache error-->" + err.Error())
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					log.Logger.Error("NewRedisCache error-->" + err.Error())
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				log.Logger.Error("NewRedisCache error-->" + err.Error())
				return err
			}
			return nil
		},
	}

	RedisCache = &RedisStore{pool, defaultExpiration}
}
