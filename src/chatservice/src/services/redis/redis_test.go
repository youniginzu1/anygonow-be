package redis

import (
	"context"
	"testing"
	"time"

	"github.com/aqaurius6666/chatservice/src/internal/var/c"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func mockRedis() (Redis, error) {
	client, err := NewRedisClient(context.Background(), REDIS_URI("localhost:6380"), REDIS_USER(""), REDIS_PASS(""))
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

func TestRedisImpl_AddMemeberSortedSet(t *testing.T) {
	r, err := mockRedis()
	// assert.Nil(t, err)
	// err = r.Clear(context.Background())
	assert.Nil(t, err)
	err = r.AddMemeberSortedSet(context.Background(), "inactive-set", "25483915-a70c-4161-bb15-617e0c3ddf2c", float64(time.Now().Add(c.ACTIVE_TIMEOUT).UnixMilli()))
	assert.Nil(t, err)
}
