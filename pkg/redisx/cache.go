package redisx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Rdb    *redis.Client
	Prefix string
}

func (c Cache) key(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return c.Prefix + hex.EncodeToString(sum[:])
}

func (c Cache) SetEx(ctx context.Context, rawKey string, val any, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.Rdb.SetEx(ctx, c.key(rawKey), b, ttl).Err()
}

func (c Cache) GetJSON(ctx context.Context, rawKey string, out any) (bool, error) {
	b, err := c.Rdb.Get(ctx, c.key(rawKey)).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(b, out)
}

func (c Cache) Del(ctx context.Context, rawKey string) error {
	return c.Rdb.Del(ctx, c.key(rawKey)).Err()
}
