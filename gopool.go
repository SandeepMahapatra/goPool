package gopool

import (
	"context"
	"math"
	"sync"
)

// GoPool limits the number of concurrently running goroutines using a buffered channel and sync.WaitGroup.
type GoPool struct {
	pool chan struct{}  // buffered channel to control concurrency
	wg   sync.WaitGroup // internal wait group to wait for all tasks

	PoolSize int // the maximum number of concurrent goroutines
}

// New creates a new GoPool with the given concurrency limit.
// If size is zero or negative, it defaults to math.MaxInt32 (no effective limit).
func New(size int32) *GoPool {
	poolSize := math.MaxInt32
	if size > 0 {
		poolSize = int(size)
	}
	return &GoPool{
		pool:     make(chan struct{}, poolSize),
		wg:       sync.WaitGroup{},
		PoolSize: poolSize,
	}
}

// Add blocks if the pool is full. It should be followed by a Done() when the goroutine completes.
func (p *GoPool) Add() error {
	return p.AddWithContext(context.Background())
}

// Done signals that a goroutine has finished and frees a slot in the pool.
func (p *GoPool) Done() {
	<-p.pool
	p.wg.Done()
}

// Wait blocks until all goroutines added with Add() or AddWithContext() have called Done().
func (p *GoPool) Wait() {
	p.wg.Wait()
}

// AddWithContext is like Add, but allows cancellation via context.
// Returns an error if the context is canceled before a slot is acquired.
func (p *GoPool) AddWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.pool <- struct{}{}:
		p.wg.Add(1)
		return nil
	}
}
