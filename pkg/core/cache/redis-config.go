package cache

import (
	"fmt"
	"os"
)

// RedisConfig configuration structure for the redis connection
type RedisConfig struct {
	Host string
	Pass string
}

// FromEnv creates a nats config using environmental variables
func (r *RedisConfig) FromEnv() (*RedisConfig, error) {

	redisHost, redisHostOk := os.LookupEnv("REDIS_HOST")
	redisPass, redisPassOk := os.LookupEnv("REDIS_PASS")

	if !redisHostOk {
		return nil, fmt.Errorf("REDIS_HOST enviornment variable not set")
	}

	if !redisPassOk {
		return nil, fmt.Errorf("REDIS_PASS enviornment variable not set")
	}

	*r = RedisConfig{
		Host: redisHost,
		Pass: redisPass,
	}

	return r, nil
}

// IsSet checks if the configuration has been set
// Note: All fields required
func (r *RedisConfig) IsSet() bool {
	return r.Host != "" && r.Pass != ""
}
