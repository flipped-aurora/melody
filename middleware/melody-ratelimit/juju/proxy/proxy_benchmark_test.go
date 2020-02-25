package proxy

import (
	"context"
	"testing"

	"melody/config"
	"melody/proxy"
)

func BenchmarkNewMiddleware_ok(b *testing.B) {
	p := NewMiddleware(&config.Backend{
		ExtraConfig: map[string]interface{}{Namespace: map[string]interface{}{"maxRate": 10000000000000.0, "capacity": 100000000000.0}},
	})(dummyProxy(&proxy.Response{}, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(context.Background(), &proxy.Request{
			Path: "/tupu",
		})
	}
}

func BenchmarkNewCMiddleware_ko(b *testing.B) {
	p := NewMiddleware(&config.Backend{
		ExtraConfig: map[string]interface{}{Namespace: map[string]interface{}{"maxRate": 1.0, "capacity": 1.0}},
	})(dummyProxy(&proxy.Response{}, nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(context.Background(), &proxy.Request{
			Path: "/tupu",
		})
	}
}
