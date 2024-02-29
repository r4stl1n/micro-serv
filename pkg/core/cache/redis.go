package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// Redis creates and manages the redis connection information
type Redis struct {
	config     *RedisConfig
	rdb        *redis.Client
	rdOptions  *redis.Options
	lockPrefix string
}

// Init create a new redis manager with a given config
func (r *Redis) Init(database int) (*Redis, error) {

	config, configError := new(RedisConfig).FromEnv()

	if configError != nil {
		return nil, configError
	}

	*r = Redis{
		config: config,
		rdOptions: &redis.Options{
			Addr:     config.Host,
			Password: config.Pass, // no password set
			DB:       database,    // use default DB
		},
		lockPrefix: "lock-",
	}

	return r, nil
}

// GetConfig get the redis config
func (r *Redis) GetConfig() RedisConfig {
	return *r.config
}

// Connect to the redis server
func (r *Redis) Connect() error {

	// Check if our nats configuration is set properly
	if !r.config.IsSet() {
		return fmt.Errorf("redis not configured")
	}

	// Creating the redis client starts a connection
	r.rdb = redis.NewClient(r.rdOptions)

	// Attempt to ping the redis database to confirm working connection
	timeoutContext, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := r.rdb.Ping(timeoutContext).Result()

	zap.L().Info("connected to redis server", zap.String("host", r.config.Host))

	timeoutCancel()
	return err
}

// Disconnect from the redis server
func (r *Redis) Disconnect() error {
	closeError := r.rdb.Close()

	zap.L().Info("disconnected from redis server", zap.String("host", r.config.Host))

	return closeError
}

// AddCache inserts an object into the cache
func (r *Redis) AddCache(key string, data any, expiration time.Duration) (string, error) {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Marshall the data to a json string to store
	marshalBytes, marshalError := json.Marshal(&data)
	if marshalError != nil {
		return "", marshalError
	}

	// Store the data into the cache
	setError := r.rdb.Set(ctx, key, marshalBytes, expiration).Err()

	if setError != nil {
		return key, setError
	}

	zap.L().Info("added data to redis", zap.Int("size", len(marshalBytes)))

	return key, nil
}

// RetrieveCache retrieve an object from the cache
func (r *Redis) RetrieveCache(key string, result any) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Grab the data from the cache
	val, getError := r.rdb.Get(ctx, key).Result()

	if getError != nil {
		return getError
	}

	// Unmarshall the object for the result
	unmarshallError := json.Unmarshal([]byte(val), &result)

	if unmarshallError != nil {
		return unmarshallError
	}

	zap.L().Info("retrieved data from redis", zap.Int("size", len([]byte(val))))

	return nil
}

// CacheExists retrieve an object from the cache
func (r *Redis) CacheExists(key string) (bool, error) {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Grab the data from the cache
	amount, checkError := r.rdb.Exists(ctx, key).Result()

	if checkError != nil {
		return false, checkError
	}

	zap.L().Info("checked if data exists in redis", zap.String("key", key))

	if amount <= 0 {
		return false, nil
	}

	return true, nil
}

// RemoveCache inserts an object into the cache
func (r *Redis) RemoveCache(key string) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Empty the key and set the expiration to one second for full removal
	setError := r.rdb.Set(ctx, key, "", 1*time.Second).Err()

	if setError != nil {
		return setError
	}

	zap.L().Info("removed data from redis cache", zap.String("key", key))

	return nil
}

// RequestKeyLock irequests to use a key for locking purposes
// Will return false if the key lock already exists
func (r *Redis) RequestKeyLock(key string, expiration time.Duration) (bool, error) {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store the data into the cache
	value, setError := r.rdb.SetArgs(ctx, r.lockPrefix+key, r.lockPrefix+key, redis.SetArgs{Get: true, TTL: expiration}).Result()

	if setError != nil && setError.Error() != "redis: nil" {
		return false, setError
	}

	if value != "" {
		return false, nil
	}

	zap.L().Info("key lock added to redis cache", zap.String("key", key))

	return true, nil
}

// RemoveKeyLock irequests to use a key for locking purposes
// Will return false if the key lock already exists
func (r *Redis) RemoveKeyLock(key string) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store the data into the cache
	_, removeError := r.rdb.Set(ctx, r.lockPrefix+key, "", 1*time.Second).Result()

	if removeError != nil {
		return removeError
	}

	zap.L().Info("key lock removed from redis cache", zap.String("key", key))

	return nil
}
