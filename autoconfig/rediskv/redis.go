package rediskv

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"thingworks.net/thingworks/jarvis-boot/utils/strings2"
	"time"
)

type RedisConfig struct {
	Host        string
	Database    int
	Password    string
	MaxActive   int
	MaxIdle     int
	IdleTimeout time.Duration
}

type RedisConnection interface {
	Reconnect()
	Ping() bool
	HSet(key string, subKey string, value interface{}) error
	HDelete(key string, subKey string) error
	HGet(key string, subKey string) (string, error)
	HGetByKey(key string) (map[string]string, error)
	ScanAll() (*[]string, error)
	Scan(func(*[]string) bool)
	CheckAndDel(string) error
	SetIfNotExistWithExpiryTime(key string, value interface{}, seconds int) error
	Get(key string) (string, error)
	Set(key string, value interface{}) (interface{}, error)
}

type DefaultRedisConnection struct {
	pool   *redis.Pool
	config RedisConfig
}

func (repo *DefaultRedisConnection) Ping() bool {
	result, err := redis.String(repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("PING")
	}))

	if err != nil || !strings2.EqualCaseIgnored(result, "PONG") {
		return false
	}

	return true
}

func (repo *DefaultRedisConnection) Set(key string, value interface{}) (interface{}, error) {
	return redis.String(repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("SET", key, value)
	}))
}

func (repo *DefaultRedisConnection) Get(key string) (string, error) {
	return redis.String(repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("GET", key)
	}))
}

func (repo *DefaultRedisConnection) SetIfNotExistWithExpiryTime(key string, value interface{}, seconds int) error {
	result, err := repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("SET", key, value, "EX", seconds, "NX")
	})

	if err != nil || result != "OK" {
		return errors.New(fmt.Sprintf("Key [%s] has already exists", key))
	}

	return nil
}

func (repo *DefaultRedisConnection) CheckAndDel(key string) error {
	_, err := repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		retry := 3

		for {
			result, err := doCheckAndDelete(conn, key)

			if err == nil || retry <= 0 {
				return result, err
			}

			retry--
		}

	})

	return err
}

func doCheckAndDelete(conn redis.Conn, key string) (interface{}, error) {
	err := conn.Send("WATCH", key)
	if err != nil {
		log.Println("Failed to watch key: ", key)
		return nil, err
	}

	conn.Send("MULTI", key)
	conn.Send("DEL", key)

	return conn.Do("EXEC")
}

func (repo *DefaultRedisConnection) Scan(callback func(*[]string) bool) {
	_, _ = repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		var keys []string

		index := 0
		for {
			array, err := redis.Values(conn.Do("SCAN", index))

			if err != nil {
				log.Printf("Failed to batch scan, becuase of [%s]", err)
				return nil, err
			}

			index, _ = redis.Int(array[0], nil)
			scannedKey, _ := redis.Strings(array[1], nil)
			keys = append(keys, scannedKey...)

			if callback(&scannedKey) || index == 0 {
				break
			}
		}

		return keys, nil
	})
}

func (repo *DefaultRedisConnection) ScanAll() (*[]string, error) {
	result, err := repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		var keys []string

		index := 0
		for {
			array, err := redis.Values(conn.Do("SCAN", index))

			if err != nil {
				log.Printf("Failed to batch scan, becuase of [%s]", err)
				return nil, err
			}

			index, _ = redis.Int(array[0], nil)
			scannedKey, _ := redis.Strings(array[1], nil)
			keys = append(keys, scannedKey...)

			if index == 0 {
				break
			}
		}

		return keys, nil
	})

	if err != nil {
		return nil, err
	}

	strings := result.([]string)

	return &strings, nil
}

func (repo *DefaultRedisConnection) openAndClose(callback func(redis.Conn) (interface{}, error)) (interface{}, error) {
	connection := repo.pool.Get()
	defer func() {
		_ = connection.Close()
	}()

	return callback(connection)
}

func (repo *DefaultRedisConnection) HSet(key string, subKey string, object interface{}) error {
	_, err := repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		value, _ := json.Marshal(object)

		_, err := conn.Do("HSET", key, subKey, value)

		return nil, err
	})

	return err
}

func (repo *DefaultRedisConnection) HGet(key string, subKey string) (string, error) {
	return redis.String(repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		return redis.String(conn.Do("HGET", key, subKey))
	}))
}

func (repo *DefaultRedisConnection) HGetByKey(key string) (map[string]string, error) {
	result, err := repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		keys, err := redis.Strings(conn.Do("HKEYS", key))

		if err != nil {
			return nil, err
		}

		var hash = make(map[string]string)

		for _, subKey := range keys {
			if value, err := redis.String(conn.Do("HGET", key, subKey)); err == nil {
				hash[subKey] = value
			}
		}

		return hash, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(map[string]string), nil
}

func (repo *DefaultRedisConnection) HDelete(key string, subKey string) error {
	_, err := repo.openAndClose(func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HDEL", key, subKey)
	})

	return err
}

func (repo *DefaultRedisConnection) Reconnect() {
	panic("implement me")
}

func NewBidder(config RedisConfig) *DistributedRightBidder {
	pool := redis.Pool{
		MaxActive:   4,
		MaxIdle:     2,
		IdleTimeout: 120 * time.Second,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			connection, err := redis.Dial(
				"tcp",
				config.Host,
				redis.DialDatabase(config.Database),
				redis.DialPassword(config.Password))

			if err != nil {
				log.Println("Failed to connect to redis, because of " + err.Error())
			}

			return connection, err
		},
	}

	return &DistributedRightBidder{
		conn: &DefaultRedisConnection{
			pool:   &pool,
			config: config,
		},
	}
}

func NewRedis(config RedisConfig) DefaultRedisConnection {
	pool := redis.Pool{
		MaxActive:   config.MaxActive,
		MaxIdle:     config.MaxIdle,
		IdleTimeout: config.IdleTimeout,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			connection, err := redis.Dial(
				"tcp",
				config.Host,
				redis.DialDatabase(config.Database),
				redis.DialPassword(config.Password))

			if err != nil {
				log.Println("Failed to connect to redis, because of " + err.Error())
			}

			return connection, err
		},
	}

	return DefaultRedisConnection{
		pool:   &pool,
		config: config,
	}
}
