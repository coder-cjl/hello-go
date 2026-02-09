package redis

import (
	"hello-go/src/logger"

	"github.com/garyburd/redigo/redis"
)

var RedisConn redis.Conn

func init() {
	var err error
	RedisConn, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		logger.Errorf("Failed to connect to Redis:", err)
		return
	}
	logger.Info("Redis 链接成功")
}

func SetValue(key string, value string) error {
	_, err := RedisConn.Do("SET", key, value)
	return err
}

func GetValue(key string) (string, error) {
	return redis.String(RedisConn.Do("GET", key))
}

func SetValueWithExpire(key string, value string, expire int) error {
	_, err := RedisConn.Do("SETEX", key, expire, value)
	return err
}

// 批量
func MSet(pairs map[string]string) error {
	args := redis.Args{}
	for k, v := range pairs {
		args = args.Add(k, v)
	}
	_, err := RedisConn.Do("MSET", args...)
	return err
}

// 队列
func LPush(queue string, value string) error {
	_, err := RedisConn.Do("LPUSH", queue, value)
	return err
}
