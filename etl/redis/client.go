package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/juliotorresmoreno/iot/etl/config"
	"github.com/juliotorresmoreno/iot/etl/log"
	"github.com/juliotorresmoreno/iot/etl/parser"
	"go.mongodb.org/mongo-driver/bson"
)

var logger = log.Log

type RedisClient struct {
	client *redis.Client
	parser *parser.Parser
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

	result.parser = parser.MakeParser()

	redisConf, err := redis.ParseURL(conf.Redis.DSN)

	if err != nil {
		logger.Trace(err)
		return result, err
	}
	result.client = redis.NewClient(redisConf)

	return result, nil
}

func (el *RedisClient) Set(key string, value bson.M) error {
	b64, err := el.parser.Encode(value)
	if err != nil {
		return err
	}

	c := el.client.Set(key, b64, 24*time.Hour)

	return c.Err()
}

func (el *RedisClient) Get(key string) (bson.M, error) {
	c := el.client.Get(key)
	if c.Err() != nil {
		return map[string]interface{}{}, c.Err()
	}

	result := map[string]interface{}{}
	err := el.parser.Decode(c.Val(), &result)

	return result, err
}

func (el *RedisClient) HSet(key, field string, value map[string]interface{}) (bool, error) {
	b64, err := el.parser.Encode(value)
	if err != nil {
		return false, err
	}

	c := el.client.HSet(key, field, b64)

	return c.Result()
}

func (el *RedisClient) HGet(key, field string) (map[string]interface{}, error) {
	c := el.client.HGet(key, field)

	result := map[string]interface{}{}
	err := el.parser.Decode(c.Val(), &result)

	return result, err
}

func (el *RedisClient) RPush(key string, value map[string]interface{}) (int64, error) {
	b64, err := el.parser.Encode(value)
	if err != nil {
		return 0, err
	}

	c := el.client.RPush(key, b64)

	return c.Result()
}

func (el *RedisClient) LRange(key string, start, stop int64) ([]map[string]interface{}, error) {
	c := el.client.LRange(key, start, stop)
	result := make([]map[string]interface{}, 0)

	items := c.Val()
	for i := 0; i < len(items); i++ {
		item := items[i]

		value := map[string]interface{}{}
		err := el.parser.Decode(item, &value)
		if err != nil {
			return result, err
		}

		result = append(result, value)
	}

	return result, nil
}
