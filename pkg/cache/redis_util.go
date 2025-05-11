package cache

import (
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetJSON(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Rdb.Set(Ctx, key, data, ttl).Err()
}

func GetJSON(key string, dest any) (bool, error) {
	val, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(val), dest)
	return true, err
}
