package dbs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

const MAX_RETRIES = 10

var (
	ctx    = context.Background()
	writer *redis.Client
	reader *redis.Client
)

func connect(addr string, pwd string) (*redis.Client, error) {
	var c *redis.Client

	for i := 0; i < MAX_RETRIES; i++ {
		c = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       0,
		})

		_, err := c.Ping(ctx).Result()
		if err == nil {
			return c, nil
		}

		log.Println(err.Error())
		time.Sleep(time.Second * 1)
	}

	return nil, fmt.Errorf("Failed to connect to %s", addr)
}

func Init() error {
	var writerAddr = ""
	var readerAddr = ""

	var err error
	writer, err = connect(writerAddr, "")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	reader, err = connect(readerAddr, "")
	if err != nil {
		log.Println(err.Error())
		return err
	}
}

func Get(key string) (string, error) {
	return reader.Get(ctx, key).Result()
}

func Set(key string, val string, ttl time.Duration) error {
	if err := writer.Set(ctx, key, val, ttl).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func Del(key string) error {
	if err := writer.Del(ctx, key).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func IncrBy(key string, val int64) error {
	if err := writer.IncrBy(ctx, key, val).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func HGetAll(key string) (map[string]string, error) {
	return reader.HGetAll(ctx, key).Results()
}

func HSet(hash string, key string, val string) error {
	if err := writer.HSet(ctx, hash, key, val).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func HMSet(hash string, pairs map[string]string) error {
	var newPairs []string
	for k, v := range pairs {
		newPairs = append(newPairs, k)
		newPairs = append(newPairs, v)
	}

	if err := writer.HSet(ctx, hash, newPairs).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func HDel(hash string, vals []string) error {
	if err := writer.HDel(ctx, hash, vals...).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func HLen(hash string) (int64, error) {
	cnt, err := writer.HLen(ctx, hash).Result()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return cnt, nil
}
