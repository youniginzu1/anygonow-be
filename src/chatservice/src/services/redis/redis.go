package redis

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

type REDIS_URI string
type REDIS_USER string
type REDIS_PASS string

var Set = wire.NewSet(wire.Struct(new(RedisImpl), "*"), wire.Bind(new(Redis), new(*RedisImpl)), NewRedisClient)

type Redis interface {
	GetKey(ctx context.Context, key string) (string, error)
	UpsertKey(ctx context.Context, key string, value string, ttl time.Duration) (string, error)
	SetKey(ctx context.Context, key string, value string, ttl time.Duration) error
	Clear(ctx context.Context, keys ...string) error
	AddMemeberSortedSet(ctx context.Context, key string, value string, score float64) error
	GetSortedSet(ctx context.Context, key string) ([]redis.Z, error)
	RemoveMemberSortedSet(ctx context.Context, key string, value string) error
}

type RedisImpl struct {
	Logger *logrus.Logger
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, uri REDIS_URI, user REDIS_USER, pass REDIS_PASS) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     string(uri),
		Username: string(user),
		Password: string(pass),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func (r *RedisImpl) GetKey(ctx context.Context, key string) (string, error) {
	v, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrKeyNotFound
		}
		return "", ErrInternal
	}
	return v, nil
}

// Implement RemoveMemberSortedSet
func (r *RedisImpl) RemoveMemberSortedSet(ctx context.Context, key string, value string) error {
	return r.Client.ZRem(ctx, key, value).Err()
}

// Implement UpsertKey
func (r *RedisImpl) UpsertKey(ctx context.Context, key string, value string, ttl time.Duration) (string, error) {
	v, err := r.GetKey(ctx, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return value, r.SetKey(ctx, key, value, ttl)
		}
		return "", err
	}
	_, err = r.Client.Expire(ctx, key, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		return "", ErrInternal
	}
	return v, nil
}

// Implement SetKey
func (r *RedisImpl) SetKey(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

// Implement Clear
func (r *RedisImpl) Clear(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return r.Client.FlushAll(ctx).Err()
	}
	return r.Client.Del(ctx, keys...).Err()
}

// Implement AddMemeberSortedSet
func (r *RedisImpl) AddMemeberSortedSet(ctx context.Context, key string, value string, score float64) error {
	_, err := r.Client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		z, err := p.ZScore(ctx, key, value).Result()
		if err != nil {
			if err == redis.Nil {
				return r.Client.ZAdd(ctx, key, &redis.Z{Score: score, Member: value}).Err()
			}
		}
		if z > score {
			if err := r.Client.ZRem(ctx, key, value).Err(); err != nil {
				return err
			}
			return r.Client.ZAdd(ctx, key, &redis.Z{Score: score, Member: value}).Err()
		}
		return r.Client.ZAdd(ctx, key, &redis.Z{Score: score - z, Member: value}).Err()
	})
	return err
}

// Implement getSortedSet
func (r *RedisImpl) GetSortedSet(ctx context.Context, key string) ([]redis.Z, error) {
	return r.Client.ZRevRangeWithScores(ctx, key, 0, -1).Result()
}
