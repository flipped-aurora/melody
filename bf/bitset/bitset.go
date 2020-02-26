package bitset

import (
	"github.com/tmthrgd/go-bitset"
	"melody/bf"
)

type BitSet struct {
	bs     bitset.Bitset
	hasher bf.Hash
}

// NewBitSet constructor for BitSet with an array of m bits
func NewBitSet(m uint) *BitSet {
	return &BitSet{bitset.New(m), bf.MD5}
}

// Add element to bitset
func (bs *BitSet) Add(elem []byte) {
	bs.bs.Set(bs.hasher(elem)[0] % bs.bs.Len())
}

// Check element in bitset
func (bs *BitSet) Check(elem []byte) bool {
	return bs.bs.IsSet(bs.hasher(elem)[0] % bs.bs.Len())
}

// Union 组合两个 bitSet
func (bs *BitSet) Union(that interface{}) (float64, error) {
	other, ok := that.(*BitSet)
	if !ok {
		return bs.getCount(), bf.ErrImpossibleToTreat
	}

	bs.bs.Union(bs.bs, other.bs)
	return bs.getCount(), nil
}

func (bs *BitSet) getCount() float64 {
	return float64(bs.bs.Count()) / float64(bs.bs.Len())
}
