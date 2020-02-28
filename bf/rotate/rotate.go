// Package rotate 实现了由三个 bf 组成的滑动组： previous current next
// 当添加一个元素时，它被保存在 current next
// 当滑动时， current 变成 previous；next 变成 current
package rotate

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"melody/bf"
	"melody/bf/bloomfilter"
	"sync"
	"time"
)

type Config struct {
	bf.Config
	TTL uint // 旋转频率TTL（以秒为单位）
}

type BloomFilter struct {
	Previous, Current, Next *bloomfilter.BloomFilter
	Config                  Config
	mutex                   *sync.RWMutex
	ctx                     context.Context
	cancel                  context.CancelFunc
}

func New(ctx context.Context, cfg Config) *BloomFilter {
	localCtx, cancel := context.WithCancel(ctx)
	preCfg := bf.EmptyConfig
	preCfg.HashName = cfg.HashName
	r := &BloomFilter{
		Previous: bloomfilter.New(preCfg),
		Current:  bloomfilter.New(cfg.Config),
		Next:     bloomfilter.New(cfg.Config),
		Config:   cfg,
		mutex:    &sync.RWMutex{},
		ctx:      localCtx,
		cancel:   cancel,
	}
	go r.keepRotating(localCtx, time.NewTicker(time.Duration(cfg.TTL)*time.Second).C)
	return r
}

func (b *BloomFilter) Close() {
	if b != nil && b.cancel != nil {
		b.cancel()
	}
}

func (b *BloomFilter) Add(elem []byte) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	b.Current.Add(elem)
	b.Next.Add(elem)
}

// elem 在过滤器中 返回 true
func (b *BloomFilter) Check(elem []byte) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.Previous.Check(elem) || b.Current.Check(elem)
}

func (b *BloomFilter) Union(that interface{}) (float64, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	other, ok := that.(*BloomFilter)
	if !ok {
		return b.capacity(), bf.ErrImpossibleToTreat
	}
	if other.Config.N != b.Config.N {
		return b.capacity(), fmt.Errorf("error: diferrent n values %d vs. %d", other.Config.N, b.Config.N)
	}

	if other.Config.P != b.Config.P {
		return b.capacity(), fmt.Errorf("error: diferrent p values %.2f vs. %.2f", other.Config.P, b.Config.P)
	}

	if _, err := b.Previous.Union(other.Previous); err != nil {
		return b.capacity(), err
	}

	if _, err := b.Current.Union(other.Current); err != nil {
		return b.capacity(), err
	}

	if _, err := b.Next.Union(other.Next); err != nil {
		return b.capacity(), err
	}

	return b.capacity(), nil

}

func (b *BloomFilter) keepRotating(ctx context.Context, c <-chan time.Time) {
	for {
		select {
		case <-c:
		case <-ctx.Done():
			return
		}

		b.mutex.Lock()

		b.Previous = b.Current
		b.Current = b.Next
		b.Next = bloomfilter.New(bf.Config{
			N:        b.Config.N,
			P:        b.Config.P,
			HashName: b.Config.HashName,
		})

		b.mutex.Unlock()
	}
}

func (b *BloomFilter) capacity() float64 {
	return (b.Previous.Capacity() + b.Current.Capacity() + b.Next.Capacity()) / 3.0
}

type SerializableBloomFilter struct {
	Previous, Current, Next *bloomfilter.BloomFilter
	Config                  Config
}

func (b *BloomFilter) MarshalBinary() ([]byte, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	buf := new(bytes.Buffer)
	w := compressor.NewWriter(buf)

	err := gob.NewEncoder(w).Encode(SerializableBloomFilter{
		Previous: b.Previous,
		Current:  b.Current,
		Next:     b.Next,
		Config:   b.Config,
	})

	w.Close()
	return buf.Bytes(), err
}

func (b *BloomFilter) UnmarshalBinary(data []byte) error {
	if b != nil && b.cancel != nil {
		b.cancel()

		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	buf := bytes.NewBuffer(data)
	r, err := compressor.NewReader(buf)
	if err != nil {
		return err
	}

	target := &SerializableBloomFilter{}
	if err := gob.NewDecoder(r).Decode(target); err != nil && err != io.EOF {
		return err
	}

	ctx := context.Background()
	if b != nil && b.ctx != nil {
		ctx = b.ctx
	}

	localCtx, cancel := context.WithCancel(ctx)

	*b = BloomFilter{
		Previous: target.Previous,
		Next:     target.Next,
		Current:  target.Current,
		Config:   target.Config,
		ctx:      ctx,
		cancel:   cancel,
		mutex:    new(sync.RWMutex),
	}

	go b.keepRotating(localCtx, time.NewTicker(time.Duration(b.Config.TTL)*time.Second).C)

	return nil
}
