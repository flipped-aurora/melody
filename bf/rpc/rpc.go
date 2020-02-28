package rpc_bf

import (
	"context"
	"fmt"
	"melody/bf/rotate"
)

var (
	ErrNoBloomFilterInitialized = fmt.Errorf("BloomFilter not initialized")
	bf                          *rotate.BloomFilter
)

type Config struct {
	rotate.Config
	Port int
}

type BloomFilterRPC int

type BloomFilter struct {
	BloomFilterRPC
}

func (b *BloomFilter) Get() *rotate.BloomFilter {
	return bf
}

func (b *BloomFilter) Close() {
	if bf != nil {
		bf.Close()
	}
}

func New(ctx context.Context, cfg Config) *BloomFilter {
	if bf != nil {
		bf.Close()
	}

	bf = rotate.New(ctx, cfg.Config)

	return new(BloomFilter)
}

// Add 函数的 input
// TODO why
type AddInput struct {
	Elems [][]byte
}

// Add 函数的 output
// TODO why
type AddOutput struct {
	Count int
}

func (b *BloomFilterRPC) Add(in AddInput, out *AddOutput) error {
	fmt.Println("add:", in.Elems)
	defer func() { fmt.Println("added elements:", out.Count) }()

	if bf == nil {
		out.Count = 0
		return ErrNoBloomFilterInitialized
	}

	k := 0
	for _, elem := range in.Elems {
		bf.Add(elem)
		k++
	}
	out.Count = k

	return nil
}

// CheckInput 函数的 input
type CheckInput struct {
	Elems [][]byte
}

// CheckOutput 函数的 output
type CheckOutput struct {
	Checks []bool
}

func (b *BloomFilterRPC) Check(in CheckInput, out *CheckOutput) error {
	fmt.Println("check:", in.Elems)
	defer func() { fmt.Println("checked elements:", out.Checks) }()

	checkRes := make([]bool, len(in.Elems))

	if bf == nil {
		out.Checks = checkRes
		return ErrNoBloomFilterInitialized
	}

	for i, elem := range in.Elems {
		checkRes[i] = bf.Check(elem)
	}
	out.Checks = checkRes

	return nil
}

// union 函数的 input
type UnionInput struct {
	BF *rotate.BloomFilter
}

// union 函数的 output
type UnionOutput struct {
	Capacity float64
}

func (b *BloomFilterRPC) Union(in UnionInput, out *UnionOutput) error {
	fmt.Println("union:", in.BF)
	defer func() { fmt.Println("union resulting capacity:", out.Capacity) }()

	if bf == nil {
		out.Capacity = 0
		return ErrNoBloomFilterInitialized
	}

	var err error
	out.Capacity, err = bf.Union(in.BF)

	return err
}
