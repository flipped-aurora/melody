package melodyrate

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestMemoryBackend(t *testing.T) {
	for _, tc := range []struct {
		name string
		s    int
		f    func(context.Context, time.Duration) Backend
	}{
		{name: "memory", s: 1, f: func(ctx context.Context, ttl time.Duration) Backend { return NewMemoryBackend(ctx, ttl) }},
		{name: "sharded", s: 256, f: func(ctx context.Context, ttl time.Duration) Backend {
			return NewShardedMemoryBackend(ctx, 256, ttl, PseudoFNV64a)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testBackend(t, tc.s, tc.f)
		})
	}
}

func testBackend(t *testing.T, storesInit int, f func(context.Context, time.Duration) Backend) {
	ttl := time.Second
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mb := f(ctx, ttl)
	total := 1000 * runtime.NumCPU()

	<-time.After(ttl)

	wg := new(sync.WaitGroup)

	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func() {
			for i := 0; i < total; i++ {
				mb.Store(fmt.Sprintf("key-%d", i), i)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	noResult := func() interface{} { return nil }

	for i := 0; i < total; i++ {
		v := mb.Load(fmt.Sprintf("key-%d", i), noResult)
		if v == nil {
			t.Errorf("key %d not present", i)
			return
		}
		if res, ok := v.(int); !ok || res != i {
			t.Errorf("unexpected value. want: %d, have: %v", i, v)
			return
		}
	}

	<-time.After(2 * ttl)

	for i := 0; i < total; i++ {
		if v := mb.Load(fmt.Sprintf("key-%d", i), noResult); v != nil {
			t.Errorf("key %d present after 2 TTL: %v", i, v)
			return
		}
	}
}
