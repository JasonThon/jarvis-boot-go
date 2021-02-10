package config

import "fmt"

func (redis *RedisConfig) GetHost() string {
	return fmt.Sprintf("%s:%d", redis.Host, redis.Port)
}

type RedisConfig struct {
	Database Database
	Host     string
	Port     int
	Password string
}

type Database struct {
	Monitor  int
	Registry int
	Bidder   int
}

