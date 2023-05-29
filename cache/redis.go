package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var RedisInstance *Redis

type Redis struct {
	Conn redis.Conn
}
type Config struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	Expire int    `json:"expire"`
}

func (r *Redis) SetValue(key string, data string) {

	_, err := r.Conn.Do(
		"HMSET",
		key,
		"data",
		data,
	)

	if err != nil {
		fmt.Println("not set in redis")
	}
}

func (r *Redis) GetValue(key string) string {
	data, err := redis.String(r.Conn.Do("HGET", key, "data"))
	if err != nil {
		// fmt.Println("Redis Connection Error", err) // handle error
		return ""
	}
	return data
}

func (r *Redis) GetKeys(key string) []string {
	data, _ := redis.Strings(r.Conn.Do("KEYS", key))
	return data
}

// New ...
func New(cnf Config) (*Redis, error) {

	redisConnectionString := fmt.Sprintf("%s:%s", cnf.Host, cnf.Port)
	conn, err := redis.Dial("tcp", redisConnectionString)
	//defer conn.Close()
	RedisInstance = &Redis{Conn: conn}
	return RedisInstance, err

}

//// GetInstance ...
//func GetInstance() redis.Conn {
//	return _conn
//}

// GetInstance ...
func GetInstance() *Redis {
	return RedisInstance
}
