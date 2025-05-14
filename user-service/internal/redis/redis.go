package redis

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
    Ctx = context.Background()

    Client = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"), // немесе config арқылы
        Password: "",
        DB:       0,
    })
)

func GetFromCache[T any](key string) (*T, error) {
    val, err := Client.Get(Ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    var result T
    if err := json.Unmarshal([]byte(val), &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func SetToCache(key string, value any, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return Client.Set(Ctx, key, data, ttl).Err()
}

func DeleteCache(key string) error {
    return Client.Del(Ctx, key).Err()
}
