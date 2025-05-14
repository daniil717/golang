package redis

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    "inventory-service/config" // Config пакетін қосу
)

var (
    Ctx = context.Background()
    Client *redis.Client
)

// Redis клиентін инициализациялау
func InitRedis() {
    cfg := config.Load() // Конфигурацияны жүктеу

    Client = redis.NewClient(&redis.Options{
        Addr:     cfg.RedisAddr,
        Password: cfg.RedisPassword, // Пароль қажет болса
        DB:       0,                  // Redis қолданатын база
    })

    // Redis-ке қосылуды тексеру
    _, err := Client.Ping(Ctx).Result()
    if err != nil {
        log.Fatalf("❌ Redis connection failed: %v", err)
    }

    log.Println("✅ Redis connection established")
}

func GetFromCache[T any](key string) (*T, error) {
    val, err := Client.Get(Ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Кэште жоқ болса
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
