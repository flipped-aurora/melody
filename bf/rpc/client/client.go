// Package client implements an rpc client for the BloomFilter, along with Add and Check methods.
package client

import (
	"errors"
	"fmt"
	"melody/bf/rotate"
	rpc_bf "melody/bf/rpc"
	"net/rpc"
)

// BloomFilter rpc client type
type BloomFilter struct {
	client *rpc.Client
}

// New creates a new bloomfilter rpc client with address
func New(address string) (*BloomFilter, error) {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &BloomFilter{client}, nil
}

// Add element through bloomfilter rpc client
func (b *BloomFilter) Add(elem []byte) {
	var addOutput rpc_bf.AddOutput
	if err := b.client.Call("BloomFilterRPC.Add", rpc_bf.AddInput{[][]byte{elem}}, &addOutput); err != nil {
		fmt.Println("error on adding bloomfilter:", err.Error())
	}
}

// Check present element through bloomfilter rpc client
func (b *BloomFilter) Check(elem []byte) bool {
	var checkOutput rpc_bf.CheckOutput
	if err := b.client.Call("BloomFilterRPC.Check", rpc_bf.CheckInput{[][]byte{elem}}, &checkOutput); err != nil {
		fmt.Println("error on check bloomfilter:", err.Error())
		return false
	}
	for _, v := range checkOutput.Checks {
		if !v {
			return false
		}
	}
	return true
}

// Union element through bloomfilter rpc client with sliding bloomfilters
func (b *BloomFilter) Union(that interface{}) (float64, error) {
	v, ok := that.(*rotate.BloomFilter)
	if !ok {
		return -1.0, errors.New("invalide argument to Union, expected rotate.BloomFilter")
	}
	var unionOutput rpc_bf.UnionOutput
	if err := b.client.Call("BloomFilterRPC.Union", rpc_bf.UnionInput{v}, &unionOutput); err != nil {
		fmt.Println("error on union bloomfilter:", err.Error())
		return -1.0, err
	}

	return unionOutput.Capacity, nil
}

func (b *BloomFilter) Close() {
	b.client.Close()
}
