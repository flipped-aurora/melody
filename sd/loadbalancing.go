package sd

import (
	"errors"
	"runtime"
	"sync/atomic"

	"github.com/valyala/fastrand"
)

// Balancer 用一个负载均衡策略去选择使用哪一个backend
type Balancer interface {
	Host() (string, error)
}

var ErrNoHosts = errors.New("no hosts available")

// NewBalancer 按照可用 processors 的数量来返回最合适的 balancer
// If GOMAXPROCS = 1, 返回一个轮询LB，因为原子计数器没有竞争
// If GOMAXPROCS > 1, 返回基于CPU数量而优化的伪随机LB
func NewBalancer(subscriber Subscriber) Balancer {
	if p := runtime.GOMAXPROCS(-1); p == 1 { // runtime.GOMAXPROCS 返回可同时执行的最大CPU数
		return NewRoundRobinLB(subscriber)
	}
	return NewRandomLB(subscriber)
}

// NewRoundRobinLB returns a new balancer using a round robin strategy
func NewRoundRobinLB(subscriber Subscriber) Balancer {
	if s, ok := subscriber.(FixedSubscriber); ok && len(s) == 1 {
		return nopBalancer(s[0])
	}
	return &roundRobinLB{
		balancer: balancer{subscriber: subscriber},
		counter:  0,
	}
}

type roundRobinLB struct {
	balancer
	counter uint64
}

// Host implements the balancer interface
func (r *roundRobinLB) Host() (string, error) {
	hosts, err := r.hosts()
	if err != nil {
		return "", err
	}
	// atomic.AddUint64 原子性的将val的值添加到*addr并返回新值
	offset := (atomic.AddUint64(&r.counter, 1) - 1) % uint64(len(hosts))
	return hosts[offset], nil
}

// NewRandomLB 使用 fastrand 伪随机数生成器
func NewRandomLB(subscriber Subscriber) Balancer {
	if s, ok := subscriber.(FixedSubscriber); ok && len(s) == 1 {
		return nopBalancer(s[0])
	}
	return &randomLB{
		balancer: balancer{subscriber: subscriber},
		rand:     fastrand.Uint32n,
	}
}

type randomLB struct {
	balancer
	rand func(uint32) uint32
}

// Host implements the balancer interface
func (r *randomLB) Host() (string, error) {
	hosts, err := r.hosts()
	if err != nil {
		return "", err
	}
	return hosts[int(r.rand(uint32(len(hosts))))], nil
}

type balancer struct {
	subscriber Subscriber
}

func (b *balancer) hosts() ([]string, error) {
	hs, err := b.subscriber.Hosts()
	if err != nil {
		return hs, err
	}
	if len(hs) <= 0 {
		return hs, ErrNoHosts
	}
	return hs, nil
}

type nopBalancer string

func (b nopBalancer) Host() (string, error) { return string(b), nil }
