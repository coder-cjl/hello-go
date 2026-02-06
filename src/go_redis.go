package main

import (
	"github.com/garyburd/redigo/redis"
)

var RedisConn redis.Conn

// func init() {
// 	var err error
// 	RedisConn, err = redis.Dial("tcp", "localhost:6379")
// 	if err != nil {
// 		Log.Error("Failed to connect to Redis:", err)
// 		return
// 	}
// 	Log.Info("Connected to Redis successfully")
// }

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

type GoRedis struct{}

func go_redis_test() {
	err := SetValue("goredis_key", "Hello, Go Redis!")
	if err != nil {
		Log.Error("SetValue failed:", err)
		return
	}
	val, err := GetValue("goredis_key")
	if err != nil {
		Log.Error("GetValue failed:", err)
		return
	}
	Log.Info("Got value from Redis:", val)
}

func go_redis_test1() {
	pairs := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	err := MSet(pairs)
	if err != nil {
		Log.Error("MSet failed:", err)
		return
	}
	Log.Info("MSet succeeded")
}

func go_redis_test2() {
	err := SetValueWithExpire("temp_key", "temporary value", 10)
	if err != nil {
		Log.Error("SetValueWithExpire failed:", err)
		return
	}
	Log.Info("SetValueWithExpire succeeded")
}

func go_redis_test3() {
	err := LPush("my_queue1", "task22")
	if err != nil {
		Log.Error("LPush failed:", err)
		return
	}
	Log.Info("LPush succeeded")
}

func (g GoRedis) Test() {
	go_redis_test3()
}
