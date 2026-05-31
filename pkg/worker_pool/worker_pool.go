package worker_pool

import (
	"errors"
	"sync"
)

var ErrWorkerIsFull = errors.New("task channel is full")

type WorkerPool struct {
	closed  bool
	taskCh  chan func()
	closeCh chan bool
	mu      sync.Mutex
}

func New(nums int, max int) (*WorkerPool, error) {
	if nums <= 0 {
		return nil, errors.New("invalid number of workers or max tasks")
	}

	wp := WorkerPool{
		closed:  false,
		taskCh:  make(chan func(), max),
		closeCh: make(chan bool),
	}

	go wp.process(nums)

	return &wp, nil
}

func (wp *WorkerPool) process(nums int) {
	wg := sync.WaitGroup{}
	wg.Add(nums)

	for range nums {
		go func() {
			defer wg.Done()
			for task := range wp.taskCh {
				task()
			}
		}()
	}

	wg.Wait()
	close(wp.closeCh)
}

func (wp *WorkerPool) AddTask(task func()) error {
	if task == nil {
		return errors.New("task is nil")
	}

	if wp.closed {
		return errors.New("worker pool is closed")
	}

	select {
	case wp.taskCh <- task:
		return nil
	default:
		return ErrWorkerIsFull
	}
}

func (wp *WorkerPool) Close() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.closed {
		return
	}

	wp.closed = true
	close(wp.taskCh)

	<-wp.closeCh
}
