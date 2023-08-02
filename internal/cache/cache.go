package cache

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedis(addr, password string, db int, expiration time.Duration) *RedisCache {
	return &RedisCache{
		redis.NewClient(&redis.Options{
			Addr:     addr, //"localhost:6379"
			Password: password,
			DB:       db,
		}),
		expiration,
	}
}

func (rc *RedisCache) GetFloat64(key string) (float64, bool) {
	strValue, err := rc.client.Get(key).Result()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return 0, false
	}

	if err == redis.Nil {
		return 0, false
	}

	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		log.Println(err)
		return 0, false
	}

	return floatValue, true
}

func (rc *RedisCache) SetFloat64(key string, value float64) bool {
	_, err := rc.client.Set(key, fmt.Sprint(value), rc.expiration).Result()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
