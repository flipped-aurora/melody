package bf

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/crc64"
	"hash/fnv"
)

type Hash func([]byte) []uint
type HashFactory func(uint) []Hash

const (
	HASHER_DEFAULT = "default"
	HASHER_OPTIMAL = "optimal"
)

var (
	defaultHashers = []Hash{
		MD5,
		CRC64,
		SHA1,
		FNV64,
		FNV128,
	}

	HashFactoryNames = map[string]HashFactory{
		HASHER_DEFAULT: DefaultHashFactory,
		HASHER_OPTIMAL: OptimalHashFactory,
	}

	ErrImpossibleToTreat = fmt.Errorf("unable to union")

	MD5    = HashWrapper(md5.New())
	SHA1   = HashWrapper(sha1.New())
	CRC64  = HashWrapper(crc64.New(crc64.MakeTable(crc64.ECMA)))
	FNV64  = HashWrapper(fnv.New64())
	FNV128 = HashWrapper(fnv.New128())
)

// 取 K 个hash函数
func DefaultHashFactory(k uint) []Hash {
	if k > uint(len(defaultHashers)) {
		k = uint(len(defaultHashers))
	}
	return defaultHashers[:k]
}

// 最优hash工厂
// FNV能快速hash大量数据并保持较小的冲突率，它的高度分散使它适用于hash一些非常相近的字符串，
// 比如URL，hostname，文件名，text，IP地址等。
func OptimalHashFactory(k uint) []Hash {
	return []Hash{
		func(b []byte) []uint {
			hs := FNV128(b)
			out := make([]uint, k)

			for i := range out {
				// 128位 => len(out) = 2
				out[i] = hs[0] + uint(i)*hs[1]
			}
			return out
		},
	}
}

func HashWrapper(h hash.Hash) Hash {
	return func(elem []byte) []uint {
		h.Reset()
		h.Write(elem)
		result := h.Sum(nil)
		// 例如md5 128位 = 16字节   len(result) = 16
		out := make([]uint, len(result)/8) // len(out) = 2
		for i := 0; i < len(result)/8; i++ {
			// binary.LittleEndian.Uint64(result[i*8 : (i+1)*8])
			// 取 result 中连续8个字节转换成 uint64   why???
			// Little-Endian就是低位字节排放在内存的低地址端，高位字节排放在内存的高地址端。
			out[i] = uint(binary.LittleEndian.Uint64(result[i*8 : (i+1)*8]))
		}
		return out
	}
}
