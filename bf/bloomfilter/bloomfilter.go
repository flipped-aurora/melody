package bloomfilter

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/tmthrgd/go-bitset"
	"melody/bf"
	"reflect"
)

type BloomFilter struct {
	bs  bitset.Bitset
	m   uint
	k   uint
	h   []bf.Hash
	cfg bf.Config
}

func New(cfg bf.Config) *BloomFilter {
	m := bf.M(cfg.N, cfg.P)
	k := bf.K(m, cfg.N)
	return &BloomFilter{
		bs:  bitset.New(m),
		m:   m,
		k:   k,
		h:   bf.HashFactoryNames[cfg.HashName](k),
		cfg: cfg,
	}
}

func (b *BloomFilter) Add(elem []byte) {
	for _, h := range b.h {
		for _, x := range h(elem) {
			// 标记到bitSet中
			b.bs.Set(x % b.m)
		}
	}
}

// 不在bitSet中 返回 false
func (b *BloomFilter) Check(elem []byte) bool {
	for _, h := range b.h {
		for _, x := range h(elem) {
			if !b.bs.IsSet(x % b.m) {
				return false
			}
		}
	}
	return true
}

// Union 合并
func (b *BloomFilter) Union(that interface{}) (float64, error) {
	other, ok := that.(*BloomFilter)
	if !ok {
		return b.Capacity(), bf.ErrImpossibleToTreat
	}

	// 确保 m 一样
	if b.m != other.m {
		return b.Capacity(), fmt.Errorf("m1(%d) != m2(%d)", b.m, other.m)
	}
	// 确保 k 一样
	if b.k != other.k {
		return b.Capacity(), fmt.Errorf("k1(%d) != k2(%d)", b.k, other.k)
	}

	hf0 := b.hashFactoryNameK(b.cfg.HashName)
	hf1 := other.hashFactoryNameK(other.cfg.HashName)
	// 确保 hash函数一样
	subject := make([]byte, 1000)
	rand.Read(subject) // 填充随机数
	for i, f := range hf0 {
		if !reflect.DeepEqual(f(subject), hf1[i](subject)) {
			return b.Capacity(), errors.New("error: different hashers")
		}
	}

	b.bs.Union(b.bs, other.bs)

	return b.Capacity(), nil
}

// Capacity 返回bitSet占用率
func (b *BloomFilter) Capacity() float64 {
	return float64(b.bs.Count()) / float64(b.m)
}

func (b *BloomFilter) hashFactoryNameK(hashName string) []bf.Hash {
	return bf.HashFactoryNames[hashName](b.k)
}

type SerializibleBloomfilter struct {
	BS       bitset.Bitset
	M        uint
	K        uint
	HashName string
	Cfg      bf.Config
}

// MarshalBinary 序列化
func (b *BloomFilter) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(&SerializibleBloomfilter{
		BS:       b.bs,
		M:        b.m,
		K:        b.k,
		HashName: b.cfg.HashName,
		Cfg:      b.cfg,
	})
	//zip buf.Bytes

	return buf.Bytes(), err
}

// UnmarshalBinary 反序列化
func (b *BloomFilter) UnmarshalBinary(data []byte) error {
	//unzip data
	buf := bytes.NewBuffer(data)
	target := SerializibleBloomfilter{}

	if err := gob.NewDecoder(buf).Decode(&target); err != nil {
		return err
	}
	*b = BloomFilter{
		bs:  target.BS,
		m:   target.M,
		k:   target.K,
		h:   bf.HashFactoryNames[target.HashName](target.K),
		cfg: target.Cfg,
	}

	return nil
}
