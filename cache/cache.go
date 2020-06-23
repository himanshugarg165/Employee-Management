package cache

import (
	"errors"
	"time"
)

// CacheConfig is the config used to provide cache configurations
type CacheConfig struct {
	URL      string
	Password string
}

// Cache interface
type Cache interface {
	Get(key string, value interface{}) error

	Set(key string, value interface{}, expiry time.Duration) error

	Delete(key string) error

	Flush() error
}

// Instance is the global cache instance which can be used anywhere in the application.
var Instance Cache

// Cache errors
var (
	ErrCacheMiss error = errors.New("cache: miss")
)
