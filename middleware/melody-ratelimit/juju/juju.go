package juju

import (
	"context"

	"github.com/juju/ratelimit"

	melodyrate "melody/middleware/melody-ratelimit"
)

// NewLimiter 创建了一个新的 Limiter
func NewLimiter(maxRate float64, capacity int64) Limiter {
	return Limiter{ratelimit.NewBucketWithRate(maxRate, capacity)}
}

// Limiter 是对 ratelimit.Bucket struct 的简单包装
type Limiter struct {
	limiter *ratelimit.Bucket
}

// Allow 检查是否可以从 bucket 中提取一个token
func (l Limiter) Allow() bool {
	return l.limiter.TakeAvailable(1) > 0
}

// NewLimiterStore 使用接收的后端返回一个用于持久性的LimiterStore
func NewLimiterStore(maxRate float64, capacity int64, backend melodyrate.Backend) melodyrate.LimiterStore {
	f := func() interface{} { return NewLimiter(maxRate, capacity) }
	return func(t string) melodyrate.Limiter {
		return backend.Load(t, f).(Limiter)
	}
}

// NewMemoryStore 使用内存后端返回一个 LimiterStore
func NewMemoryStore(maxRate float64, capacity int64) melodyrate.LimiterStore {
	return NewLimiterStore(maxRate, capacity, melodyrate.DefaultShardedMemoryBackend(context.Background()))
}
