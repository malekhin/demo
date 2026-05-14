package scheduler

import (
	"context"
	"demo/pkg/redis/mocks"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddJobPanicsOnInvalidSpec(t *testing.T) {
	c := cron.New()
	redisMock := mocks.NewIRedis(t)
	sut := NewCron(zap.NewNop(), redisMock, c)

	require.Panics(t, func() {
		sut.AddJob(context.Background(), "invalid spec", func(ctx context.Context) error {
			return nil
		})
	})
}

func TestAddJobSucceedsWithValidSpec(t *testing.T) {
	c := cron.New()
	redisMock := mocks.NewIRedis(t)

	sut := NewCron(zap.NewNop(), redisMock, c)

	require.NotPanics(t, func() {
		sut.AddJob(context.Background(), "@every 1s", func(ctx context.Context) error {
			return nil
		})
	})
}

func TestStartAndStop(t *testing.T) {
	c := cron.New()
	redisMock := mocks.NewIRedis(t)
	sut := NewCron(zap.NewNop(), redisMock, c)

	require.NotPanics(t, func() {
		sut.Start()
		sut.Stop()
	})
}

func TestRunReturnsEntryID(t *testing.T) {
	c := cron.New()
	redisMock := mocks.NewIRedis(t)

	sut := NewCron(zap.NewNop(), redisMock, c)

	entryID, err := sut.run(context.Background(), "@every 1s", func(ctx context.Context) error {
		return nil
	}, &sync.Mutex{})

	require.NoError(t, err)
	require.NotZero(t, entryID)
}

func TestRunReturnsErrorForInvalidSpec(t *testing.T) {
	c := cron.New()
	redisMock := mocks.NewIRedis(t)

	sut := NewCron(zap.NewNop(), redisMock, c)

	entryID, err := sut.run(context.Background(), "invalid", func(ctx context.Context) error {
		return nil
	}, &sync.Mutex{})

	require.Error(t, err)
	require.Zero(t, entryID)
}

func TestNewCronReturnsInstance(t *testing.T) {
	c := cron.New()
	logger := zap.NewNop()
	redisMock := mocks.NewIRedis(t)

	sut := NewCron(logger, redisMock, c)

	require.NotNil(t, sut)
	require.NotNil(t, sut.cron)
	require.NotNil(t, sut.logger)
}

func TestJobExecutesSuccessfully(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	logger := zap.NewNop()
	redisMock := mocks.NewIRedis(t)
	redisMock.On("Lock", context.Background(), mock.Anything, mock.Anything).Return(true, nil)

	sut := NewCron(logger, redisMock, c)

	executed := make(chan bool, 1)
	sut.AddJob(context.Background(), "@every 1s", func(ctx context.Context) error {
		executed <- true
		return nil
	})
	sut.Start()
	defer sut.Stop()

	select {
	case <-executed:
	case <-time.After(2 * time.Second):
		t.Fatal("Job did not execute within timeout")
	}
}

func TestJobExecutesWithError(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	logger := zap.NewNop()
	redisMock := mocks.NewIRedis(t)

	sut := NewCron(logger, redisMock, c)
	redisMock.On("Lock", context.Background(), mock.Anything, mock.Anything).Return(true, nil)

	executed := make(chan bool, 1)
	sut.AddJob(context.Background(), "@every 1s", func(ctx context.Context) error {
		executed <- true
		return fmt.Errorf("test error")
	})
	sut.Start()
	defer sut.Stop()

	select {
	case <-executed:
	case <-time.After(2 * time.Second):
		t.Fatal("Job did not execute within timeout")
	}
}

func TestTwoInstance(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	logger := zap.NewNop()
	redisMock1 := mocks.NewIRedis(t)
	redisMock2 := mocks.NewIRedis(t)
	redisMock1.On("Lock", context.Background(), mock.Anything, mock.Anything).Return(true, nil)
	redisMock2.On("Lock", context.Background(), mock.Anything, mock.Anything).Return(false, nil)
	sut1 := NewCron(logger, redisMock1, c)
	sut2 := NewCron(logger, redisMock2, c)

	executed := make(chan bool, 1)
	notExecuted := make(chan bool, 1)
	sut1.AddJob(context.Background(), "@every 1s", func(ctx context.Context) error {
		executed <- true
		return nil
	})
	sut2.AddJob(context.Background(), "@every 1s", func(ctx context.Context) error {
		return nil
	})
	sut1.Start()
	sut2.Start()
	defer sut1.Stop()
	defer sut2.Stop()

	time.Sleep(time.Second)

	select {
	case <-executed:
		// expected
	case <-notExecuted:
		t.Fatal("not expected channel")
	case <-time.After(time.Second):
		t.Fatal("Job did not execute within timeout")
	}
}
