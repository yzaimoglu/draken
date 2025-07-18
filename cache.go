package draken

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Cache interface {
	Init(bool)
	Stop()
	Get(key string) (*string, error)
	Set(key string, value any, ttl time.Duration) error
	Exists(key string) bool
	Expire(key string, ttl time.Duration) error
	Push(key string, value any) error
	Pop(key string) (string, error)
	Len(key string) (int64, error)
}

type Redis struct {
	Client  *redis.Client
	Context context.Context
	Cancel  context.CancelFunc
}

func (d *Draken) initCache() {
	if !d.Config.Cache.Enabled {
		log.Debug().Msgf("Cache is disabled in the config, skipping...")
		return
	}
	log.Debug().Msgf("Initializing cache...")

	switch d.Config.Cache.Type {
	case CacheTypeRedis:
		d.Cache = NewRedis(d.Config.Cache.DSN)
	}
	log.Info().Msgf("Cache initialized.")
}

// NewRedis creates a new Redis object
func NewRedis(dsn string) *Redis {
	ctx, cancel := context.WithCancel(context.Background())
	log.Debug().Msgf("Connecting to redis...")
ConnectionStart:
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		log.Error().Msgf("Could not parse dsn: %v", err)
		time.Sleep(10 * time.Second)
		log.Warn().Msgf("Waiting for 10 seconds before trying to establish a new connection to redis...")
		goto ConnectionStart
	}
	client := redis.NewClient(opt)

	res := client.Ping(ctx)
	if res.Err() != nil {
		log.Error().Msgf("Could not connect to redis: %v", res.Err())
		time.Sleep(10 * time.Second)
		log.Warn().Msgf("Waiting for 10 seconds before trying to establish a new connection to redis...")
		goto ConnectionStart
	}

	rd := &Redis{
		Client:  client,
		Context: ctx,
		Cancel:  cancel,
	}
	log.Info().Msgf("Connected to redis.")
	return rd
}

// Check if the Redis struct implements all Cache methods
var _ = (*Redis)(nil)

func (r *Redis) Init(e bool) {
	if !e {
		log.Debug().Msgf("Cache is disabled in the config, skipping initialization...")
		return
	}

	log.Debug().Msgf("Redis cache initialized.")
}

func (r *Redis) Stop() {
	if r.Cancel != nil {
		r.Cancel()
	}
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			log.Error().Msgf("Error closing Redis client: %v", err)
		} else {
			log.Info().Msgf("Redis client closed successfully.")
		}
	}
}

func (r *Redis) Get(key string) (*string, error) {
	var result string

	cmd := r.Client.Get(r.Context, key)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Scan(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Redis) Set(key string, value any, ttl time.Duration) error {
	cmd := r.Client.Set(r.Context, key, value, ttl)
	return cmd.Err()
}

func (r *Redis) Exists(key string) bool {
	cmd := r.Client.Exists(r.Context, key)
	return cmd.Val() == 1
}

func (r *Redis) Expire(key string, ttl time.Duration) error {
	cmd := r.Client.Expire(r.Context, key, ttl)
	return cmd.Err()
}

// Push pushes a single value onto the tail of the list at key.
func (r *Redis) Push(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Debug().Err(err).Msg("json  marshal failed")
		return err
	}

	if r == nil || r.Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// RPush returns *redis.IntCmd; Err() reflects any underlying error.
	return r.Client.LPush(r.Context, key, string(data)).Err()
}

// Pop removes and returns the head element of the list at key.
// If the list is empty, it returns ("", nil), but you could also choose
// to return ("", redis.Nil) and let the caller distinguish that yourself.
func (r *Redis) Pop(key string) (string, error) {
	if r == nil || r.Client == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	if !r.Exists(key) {
		return "", nil
	}

	str, err := r.Client.RPop(r.Context, key).Result()
	if err == redis.Nil {
		// List is empty
		return "", nil
	}
	if err != nil {
		// Real Redis error (network, wrong data type, etc.)
		return "", err
	}
	return str, nil
}

// Len returns the length of the list stored at key.
func (r *Redis) Len(key string) (int64, error) {
	if r == nil || r.Client == nil {
		return -1, fmt.Errorf("redis client not initialized")
	}

	if !r.Exists(key) {
		return -1, nil
	}

	size, err := r.Client.LLen(r.Context, key).Result()
	if err == redis.Nil {
		// List is empty
		return -1, nil
	}
	if err != nil {
		// Real Redis error (network, wrong data type, etc.)
		return -1, err
	}
	return size, nil
}
