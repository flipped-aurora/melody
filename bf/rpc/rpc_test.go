package rpc_bf

import (
	"context"
	"melody/bf/rotate"
	"melody/bf/testutils"
	"testing"
)

func TestBFAdd_ok(t *testing.T) {
	b := New(context.Background(), Config{rotate.Config{testutils.TestCfg, 5}, 1234})

	var (
		addOutput AddOutput
		elems1    = [][]byte{[]byte("elem1"), []byte("elem2")}
	)

	err := b.Add(AddInput{elems1}, &addOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	b.Close()
}

func TestBFCheck_ok(t *testing.T) {
	b := New(context.Background(), Config{rotate.Config{testutils.TestCfg, 5}, 1234})

	var (
		addOutput   AddOutput
		checkOutput CheckOutput
		elems1      = [][]byte{[]byte("elem1"), []byte("elem2")}
	)

	err := b.Add(AddInput{elems1}, &addOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	err = b.Check(CheckInput{elems1}, &checkOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}
	if len(checkOutput.Checks) != 2 || !checkOutput.Checks[0] || !checkOutput.Checks[1] {
		t.Errorf("checks error, expected true elements")
		return
	}

	b.Close()
}

func TestBFUnion_ok(t *testing.T) {
	b := New(context.Background(), Config{rotate.Config{testutils.TestCfg, 5}, 1234})

	var (
		addOutput   AddOutput
		checkOutput CheckOutput
		unionOutput UnionOutput
		elems1      = [][]byte{[]byte("elem1"), []byte("elem2")}
		elems2      = [][]byte{[]byte("house")}
		elems3      = [][]byte{[]byte("house"), []byte("mouse")}
	)

	err := b.Add(AddInput{elems1}, &addOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	var bf2 = rotate.New(context.Background(), rotate.Config{testutils.TestCfg, 5})
	bf2.Add([]byte("house"))

	err = b.Union(UnionInput{bf2}, &unionOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	err = b.Check(CheckInput{elems2}, &checkOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}
	if len(checkOutput.Checks) != 1 || !checkOutput.Checks[0] {
		t.Error("checks error, expected true element")
		return
	}

	var bf3 = rotate.New(context.Background(), rotate.Config{testutils.TestCfg, 5})
	bf3.Add([]byte("mouse"))

	b.Union(UnionInput{bf3}, &unionOutput)

	err = b.Check(CheckInput{elems3}, &checkOutput)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}
	if len(checkOutput.Checks) != 2 || !checkOutput.Checks[0] || !checkOutput.Checks[1] {
		t.Error("checks error, expected true element")
		return
	}

	bf.Close()
}

func TestBFAdd_ko(t *testing.T) {
	b := new(BloomFilter)
	bf = nil
	var (
		addOutput AddOutput
		elems1    = [][]byte{[]byte("elem1"), []byte("elem2")}
	)

	err := b.Add(AddInput{elems1}, &addOutput)
	if err != ErrNoBloomFilterInitialized {
		t.Error("error, should have been no bloomfilter initialized")
	}
}

func TestBFCheck_ko(t *testing.T) {
	b := new(BloomFilter)
	bf = nil
	var (
		checkOutput CheckOutput
		elems1      = [][]byte{[]byte("elem1"), []byte("elem2")}
	)

	err := b.Check(CheckInput{elems1}, &checkOutput)
	if err != ErrNoBloomFilterInitialized {
		t.Error("error, should have been no bloomfilter initialized")
	}
}

func TestBFUnion_ko(t *testing.T) {
	b := new(BloomFilter)
	bf = nil
	var (
		unionOutput UnionOutput
	)

	var bf2 = rotate.New(context.Background(), rotate.Config{testutils.TestCfg, 5})
	bf2.Add([]byte("house"))

	err := b.Union(UnionInput{bf2}, &unionOutput)
	if err != ErrNoBloomFilterInitialized {
		t.Error("error, should have been no bloomfilter initialized")
	}
}
