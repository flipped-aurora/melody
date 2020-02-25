/*
Package juju provides a set of rate-limit proxy and router middlewares using the github.com/juju/ratelimit lib.

The juju package provides an efficient token bucket implementation. See https://github.com/juju/ratelimit
and http://en.wikipedia.org/wiki/Token_bucket for more details.
*/
package juju

import (
	"context"

	"github.com/juju/ratelimit"

	melodyrate "melody/middleware/melody-ratelimit"
)

// NewLimiter creates a new Limiter
func NewLimiter(maxRate float64, capacity int64) Limiter {
	return Limiter{ratelimit.NewBucketWithRate(maxRate, capacity)}
}

// Limiter is a simple wrapper over the ratelimit.Bucket struct
type Limiter struct {
	limiter *ratelimit.Bucket
}

// Allow checks if its possible to extract 1 token from the bucket
func (l Limiter) Allow() bool {
	return l.limiter.TakeAvailable(1) > 0
}

// NewLimiterStore returns a LimiterStore using the received backend for persistence
func NewLimiterStore(maxRate float64, capacity int64, backend melodyrate.Backend) melodyrate.LimiterStore {
	f := func() interface{} { return NewLimiter(maxRate, capacity) }
	return func(t string) melodyrate.Limiter {
		return backend.Load(t, f).(Limiter)
	}
}

// NewMemoryStore returns a LimiterStore using the memory backend
func NewMemoryStore(maxRate float64, capacity int64) melodyrate.LimiterStore {
	return NewLimiterStore(maxRate, capacity, melodyrate.DefaultShardedMemoryBackend(context.Background()))
}
