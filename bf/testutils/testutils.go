// Package testutils contains utils for the tests.
package testutils

import (
	"melody/bf"
	"testing"
)

var (
	TestCfg = bf.Config{
		N:        100,
		P:        0.001,
		HashName: bf.HASHER_OPTIMAL,
	}

	TestCfg2 = bf.Config{
		N:        100,
		P:        0.00001,
		HashName: bf.HASHER_OPTIMAL,
	}

	TestCfg3 = bf.Config{
		N:        100,
		P:        0.001,
		HashName: bf.HASHER_DEFAULT,
	}
)

func CallSet(t *testing.T, set bf.Bloomfilter) {
	set.Add([]byte{1, 2, 3})
	if !set.Check([]byte{1, 2, 3}) {
		t.Error("failed check")
	}

	if set.Check([]byte{1, 2, 4}) {
		t.Error("unexpected check")
	}
}

func CallSetUnion(t *testing.T, set1, set2 bf.Bloomfilter) {
	elem := []byte{1, 2, 3}
	set1.Add(elem)
	if !set1.Check(elem) {
		t.Error("failed add set1 before union")
		return
	}

	if set2.Check(elem) {
		t.Error("unexpected check to union of set2")
		return
	}

	if _, err := set2.Union(set1); err != nil {
		t.Error("failed union set1 to set2", err.Error())
		return
	}

	if !set2.Check(elem) {
		t.Error("failed union check of set1 to set2")
		return
	}
}
