package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func newRedisClient() *redis.Client {
	redisHost := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func cacheSet(cfg *apiConfig, key string, messages []Message) {

	jsonmessages, err := json.Marshal(messages)
	if err != nil {
		log.Printf("  ERROR : failed to marshal json %v", err)
		return
	}
	if err := cfg.rdb.Set(cfg.ctx, key, jsonmessages, 10*time.Hour).Err(); err != nil {
		log.Printf("  ERROR : failed to set to redis : %v", err)
		return
	}

}

func cacheGet(cfg *apiConfig, key string) ([]Message, error) {

	cachedTweet, err := cfg.rdb.Get(cfg.ctx, key).Result()
	if err == redis.Nil {
		log.Printf("cach miss for %v",key)
        return nil , err
	} else if err != nil {
		log.Printf("  ERROR request from redis : %v", err)
        return nil , err
	}

	var messages []Message

	if err := json.NewDecoder(strings.NewReader(cachedTweet)).Decode(&messages); err != nil {
		log.Printf("  ERROR : failed to parse json in redis response %v", err)
		return nil, err
	}

    log.Println("returning messages from cache ")
	return messages, nil
}

