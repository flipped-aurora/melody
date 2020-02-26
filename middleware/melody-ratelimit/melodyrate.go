package melodyrate

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrLimited 是在超过速率限制时返回的错误
	ErrLimited = errors.New("ERROR: rate limit exceded")

	// DataTTL 是默认的清除时间
	DataTTL        = 10 * time.Minute
	now            = time.Now
	shards  uint64 = 2048
)

// Limiter 定义了一个简单的接口
type Limiter interface {
	Allow() bool
}

// LimiterStore 定义了一个Limiter查找函数的接口
type LimiterStore func(string) Limiter

// Hasher 获取接收到的字符串的hash
type Hasher func(string) uint64

// Backend 是 persistence layer 的接口
type Backend interface {
	Load(string, func() interface{}) interface{}
	Store(string, interface{}) error
}

// ShardedMemoryBackend 是为了避免互斥锁争用而对数据进行分片的memory backend
type ShardedMemoryBackend struct {
	shards []*MemoryBackend
	total  uint64
	hasher Hasher
}

//DefaultShardedMemoryBackend 默认
func DefaultShardedMemoryBackend(ctx context.Context) *ShardedMemoryBackend {
	return NewShardedMemoryBackend(ctx, shards, DataTTL, PseudoFNV64a)
}

// NewShardedMemoryBackend 后端返回带有“碎片”碎片的 ShardedMemoryBackend
func NewShardedMemoryBackend(ctx context.Context, shards uint64, ttl time.Duration, h Hasher) *ShardedMemoryBackend {
	b := &ShardedMemoryBackend{
		shards: make([]*MemoryBackend, shards),
		total:  shards,
		hasher: h,
	}
	var i uint64
	for i = 0; i < shards; i++ {
		b.shards[i] = NewMemoryBackend(ctx, ttl)
	}
	return b
}

func (b *ShardedMemoryBackend) shard(key string) uint64 {
	return b.hasher(key) % b.total
}

// Load 实现 Backend 接口
func (b *ShardedMemoryBackend) Load(key string, f func() interface{}) interface{} {
	return b.shards[b.shard(key)].Load(key, f)
}

// Store 实现 Backend 接口
func (b *ShardedMemoryBackend) Store(key string, v interface{}) error {
	return b.shards[b.shard(key)].Store(key, v)
}

func (b *ShardedMemoryBackend) del(key ...string) {
	buckets := map[uint64][]string{}
	for _, k := range key {
		h := b.shard(k)
		ks, ok := buckets[h]
		if !ok {
			ks = []string{k}
		} else {
			ks = append(ks, k)
		}
		buckets[h] = ks
	}

	for s, ks := range buckets {
		b.shards[s].del(ks...)
	}
}

// NewMemoryBackend 返回一个 MemoryBackend
func NewMemoryBackend(ctx context.Context, ttl time.Duration) *MemoryBackend {
	m := &MemoryBackend{
		data:       map[string]interface{}{},
		lastAccess: map[string]time.Time{},
		mu:         new(sync.RWMutex),
	}

	go m.manageEvictions(ctx, ttl)

	return m
}

// MemoryBackend 通过包装一个sync.Map实现Backend interface
type MemoryBackend struct {
	data       map[string]interface{}
	lastAccess map[string]time.Time
	mu         *sync.RWMutex
}

func (m *MemoryBackend) manageEvictions(ctx context.Context, ttl time.Duration) {
	t := time.NewTicker(ttl)
	for {
		keysToDel := []string{}

		select {
		case <-ctx.Done():
			t.Stop()
			return
		case now := <-t.C:
			m.mu.RLock()
			for k, v := range m.lastAccess {
				if v.Add(ttl).Before(now) {
					keysToDel = append(keysToDel, k)
				}
			}
			m.mu.RUnlock()
		}

		m.del(keysToDel...)
	}
}

// Load 实现 Backend 接口
func (m *MemoryBackend) Load(key string, f func() interface{}) interface{} {
	m.mu.RLock()
	v, ok := m.data[key]
	m.mu.RUnlock()

	n := now()

	if ok {
		go func(t time.Time) {
			m.mu.Lock()
			if t0, ok := m.lastAccess[key]; !ok || t.After(t0) {
				m.lastAccess[key] = t
			}
			m.mu.Unlock()
		}(n)

		return v
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok = m.data[key]
	if ok {
		return v
	}

	v = f()
	m.lastAccess[key] = n
	m.data[key] = v

	return v
}

// Store 实现 Backend 接口
func (m *MemoryBackend) Store(key string, v interface{}) error {
	m.mu.Lock()
	m.lastAccess[key] = now()
	m.data[key] = v
	m.mu.Unlock()
	return nil
}

func (m *MemoryBackend) del(key ...string) {
	m.mu.Lock()
	for _, k := range key {
		delete(m.data, k)
		delete(m.lastAccess, k)
	}
	m.mu.Unlock()
}
