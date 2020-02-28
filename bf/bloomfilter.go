// Package bloomfilter contains common data and interfaces needed to implement bloomfilters.
//
// 理论基础: http://llimllib.github.io/bloomfilter-tutorial/zh_CN/
// 我们实现了三种 bloomfilter ：衍生自bitSet的、 sliding bloomfilters、 rpc bloomfilter
package bf

import "math"

type BloomFilter interface {
	Add([]byte)
	Check([]byte) bool
	Union(interface{}) (float64, error)
}

// P - 失误几率
// N - 存到过滤器的element的个数
// HashName - "default" or "optimal"
type Config struct {
	N        uint
	P        float64
	HashName string
}

// 给 sliding bf 的 previous 用的
var EmptyConfig = Config{
	N: 2,
	P: .5,
}

// M bit array 的长度
func M(n uint, p float64) uint {
	return uint(math.Ceil(-(float64(n) * math.Log(p)) / math.Log(math.Pow(2.0, math.Log(2.0)))))
}

// K hash函数的个数
func K(m, n uint) uint {
	return uint(math.Ceil(math.Log(2.0) * float64(m) / float64(n)))
}

type EmptySet int

// Check implementation for EmptySet
func (e EmptySet) Check(_ []byte) bool { return false }

// Add implementation for EmptySet
func (e EmptySet) Add(_ []byte) {}

// Union implementation for EmptySet
func (e EmptySet) Union(interface{}) (float64, error) { return -1, nil }
