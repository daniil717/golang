package redis

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
)

var rdb *redis.Client

// Redis клиентін инициализациялау
func Init(redisURL string) {
    rdb = redis.NewClient(&redis.Options{
        Addr: redisURL,
    })
}

// Кэшке деректерді жазу
func SetToCache[T any](key string, value T, expiration time.Duration) error {
    // Сериализациялау
    data, err := json.Marshal(value)
    if err != nil {
        log.Printf("Failed to marshal value for key %s: %v", key, err)
        return err
    }

    err = rdb.Set(context.Background(), key, data, expiration).Err()
    if err != nil {
        log.Printf("Failed to set key %s to cache: %v", key, err)
        return err
    }
    return nil
}

// Кэштен деректерді алу
func GetFromCache[T any](key string) (*T, error) {
    val, err := rdb.Get(context.Background(), key).Result()
    if err == redis.Nil {
        return nil, nil // Егер кэште деректер болмаса
    } else if err != nil {
        log.Printf("Failed to get key %s from cache: %v", key, err)
        return nil, err
    }

    var result T
    // Десериализациялау
    err = json.Unmarshal([]byte(val), &result)
    if err != nil {
        log.Printf("Failed to unmarshal value for key %s: %v", key, err)
        return nil, err
    }

    return &result, nil
}

// Кэшті тазалау
func DeleteCache(key string) error {
    err := rdb.Del(context.Background(), key).Err()
    if err != nil {
        log.Printf("Failed to delete key %s from cache: %v", key, err)
        return err
    }
    return nil
}
