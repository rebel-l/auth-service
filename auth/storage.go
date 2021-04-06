package auth

import (
	"time"

	"github.com/go-redis/redis/v7"
)

// Storage defines an interface to store tokens.
type Storage interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}
