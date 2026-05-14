package redis

import (
	"context"
	"demo/internal/common/config"
	metricsMocks "demo/internal/common/metrics/mocks"
	"demo/pkg/redis"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var testKey = "redisTestKey"

func TestNewRedis(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}

type RedisSuite struct {
	suite.Suite
	wrappedRedis redis.IRedis
}

func (s *RedisSuite) SetupTest() {
	cfg, err := config.New()
	s.Require().NoError(err)

	metricsMock := metricsMocks.NewIMetrics(s.T())
	metricsMock.On("RecordExternalCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	r := redis.NewRedis(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		Prefix:   cfg.Redis.Prefix,
	})

	s.wrappedRedis = NewRedis(r, metricsMock)
}

func (s *RedisSuite) TestGetSet() {
	err := s.wrappedRedis.Set(context.Background(), testKey, 1, time.Second)
	s.Require().NoError(err)

	// int
	resInt, exists, err := s.wrappedRedis.GetInt(context.Background(), testKey)
	s.Require().NoError(err)
	s.True(exists)
	s.Equal(1, resInt)

	_, _, err = s.wrappedRedis.GetString(context.Background(), testKey)
	s.Require().Error(err)

	err = s.wrappedRedis.Set(context.Background(), testKey, 1.99, time.Second)
	s.Require().NoError(err)

	// float
	resFloat, exists, err := s.wrappedRedis.GetFloat(context.Background(), testKey)
	s.Require().NoError(err)
	s.True(exists)
	s.InEpsilon(1.99, resFloat, 0)

	err = s.wrappedRedis.Set(context.Background(), testKey, "str", time.Second)
	s.Require().NoError(err)

	// str
	resStr, exists, err := s.wrappedRedis.GetString(context.Background(), testKey)
	s.Require().NoError(err)
	s.True(exists)
	s.Equal("str", resStr)
	resStruct := struct{ Val string }{}
	_, err = s.wrappedRedis.GetStruct(context.Background(), testKey, &resStruct)
	s.Require().Error(err)

	// struct
	inputStruct := struct{ Val string }{Val: "test"}
	err = s.wrappedRedis.Set(context.Background(), testKey, inputStruct, time.Second)
	s.Require().NoError(err)

	resStruct = struct{ Val string }{}
	exists, err = s.wrappedRedis.GetStruct(context.Background(), testKey, &resStruct)
	s.Require().NoError(err)
	s.True(exists)
	s.Equal(resStruct.Val, inputStruct.Val)
}

func (s *RedisSuite) TestGetAndSet() {
	var res int
	err := s.wrappedRedis.SetAndGet(context.Background(), &res, func() (interface{}, error) {
		return nil, errors.New("err")
	}, "redis_setandget", time.Second)
	s.Require().Error(err)

	err = s.wrappedRedis.SetAndGet(context.Background(), &res, func() (interface{}, error) {
		return 1, nil
	}, "redis_setandget", time.Second)
	s.Require().NoError(err)
	s.Equal(1, res)

	var res2 string
	err = s.wrappedRedis.SetAndGet(context.Background(), &res2, func() (interface{}, error) {
		return 1, nil
	}, "redis_setandget", time.Second)
	s.Require().Error(err)

	err = s.wrappedRedis.SetAndGet(context.Background(), &res, func() (interface{}, error) {
		return math.Inf(1), nil
	}, "redis_setandget2", -time.Second)
	s.Require().Error(err)
}

func (s *RedisSuite) TestLockUnlock() {
	ok, err := s.wrappedRedis.Lock(context.Background(), testKey, time.Second)
	s.Require().NoError(err)
	s.True(ok)

	ok, err = s.wrappedRedis.Lock(context.Background(), testKey, time.Second)
	s.Require().NoError(err)
	s.False(ok)

	err = s.wrappedRedis.Unlock(context.Background(), testKey)
	s.Require().NoError(err)
}
