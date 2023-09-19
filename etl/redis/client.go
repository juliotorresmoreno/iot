package redis

import (
	"encoding/base64"
	"time"

	"github.com/go-redis/redis"
	"github.com/juliotorresmoreno/iot/etl/config"
	"github.com/juliotorresmoreno/iot/etl/log"
	"go.mongodb.org/mongo-driver/bson"
)

var logger = log.Log

type RedisClient struct {
	client *redis.Client
}

var DefaultRedisClient *RedisClient

func init() {
	DefaultRedisClient, _ = NewRedisClient()
}

func NewRedisClient() (*RedisClient, error) {
	result := &RedisClient{}

	conf, err := config.GetConfig()
	if err != nil {
		logger.Trace(err)
		return result, err
	}

	redisConf, err := redis.ParseURL(conf.Redis.DSN)

	if err != nil {
		logger.Trace(err)
		return result, err
	}
	result.client = redis.NewClient(redisConf)

	return result, nil
}

func (el *RedisClient) encode(value map[string]interface{}) (string, error) {
	buff, err := bson.Marshal(value)
	if err != nil {
		return "", err
	}
	b64 := base64.RawStdEncoding.EncodeToString(buff)

	return b64, nil
}

func (el *RedisClient) decode(value string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	decoded, err := base64.RawStdEncoding.DecodeString(value)

	if err != nil {
		return result, err
	}
	err = bson.Unmarshal(decoded, result)

	return result, err
}

func (el *RedisClient) Set(key string, value bson.M) error {
	b64, err := el.encode(value)
	if err != nil {
		return err
	}

	c := el.client.Set(key, b64, 24*time.Hour)

	return c.Err()
}

func (el *RedisClient) Get(key string) (bson.M, error) {
	c := el.client.Get(key)
	if c.Err() != nil {
		return bson.M{}, c.Err()
	}

	value, err := el.decode(c.Val())

	return value, err
}

func (el *RedisClient) HSet(key, field string, value map[string]interface{}) (bool, error) {
	b64, err := el.encode(value)
	if err != nil {
		return false, err
	}

	c := el.client.HSet(key, field, b64)

	return c.Result()
}

func (el *RedisClient) HGet(key, field string) (map[string]interface{}, error) {
	c := el.client.HGet(key, field)

	value, err := el.decode(c.Val())

	return value, err
}

func (el *RedisClient) RPUSH(key string, value map[string]interface{}) (int64, error) {
	b64, err := el.encode(value)
	if err != nil {
		return 0, err
	}

	c := el.client.RPush(key, b64)

	return c.Result()
}

func (el *RedisClient) LRANGE(key string, start, stop int64) ([]map[string]interface{}, error) {
	c := el.client.LRange(key, start, stop)
	result := make([]map[string]interface{}, 0)

	items := c.Val()
	for i := 0; i < len(items); i++ {
		item := items[i]
		value, err := el.decode(item)

		if err != nil {
			return result, err
		}
		result = append(result, value)
	}

	return result, nil
}
