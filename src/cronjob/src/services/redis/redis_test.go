package redis

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func mockRedis() (Redis, error) {
	client, err := NewRedisClient(context.Background(), REDIS_URI("anygonow-redis-001.anygonow-redis.dmalyf.use1.cache.amazonaws.com:6379"), REDIS_USER(""), REDIS_PASS("QW55R29Ob3cxMjMhNDU2QA=="))
	if err != nil {
		return nil, err
	}
	return &RedisImpl{
		Logger: logrus.New(),
		Client: client,
	}, nil
}

func TestGetKey(t *testing.T) {
	r, err := mockRedis()
	assert.Nil(t, err)
	err = r.Clear(context.Background())
	assert.Nil(t, err)
	_, err = r.GetKey(context.Background(), "test")
	assert.EqualError(t, err, ErrKeyNotFound.Error())

	err = r.SetKey(context.Background(), "test", "test", time.Second)
	assert.Nil(t, err)
	v, err := r.GetKey(context.Background(), "test")
	assert.Nil(t, err)
	assert.Equal(t, "test", v)
}

// TestSetKey tests the SetKey method
func TestSetKey(t *testing.T) {
	r, err := mockRedis()
	assert.Nil(t, err)
	err = r.Clear(context.Background())
	assert.Nil(t, err)
	err = r.SetKey(context.Background(), "test", "test", time.Second)
	assert.Nil(t, err)
	v, err := r.GetKey(context.Background(), "test")
	assert.Nil(t, err)
	assert.Equal(t, "test", v)
}
