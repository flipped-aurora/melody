package rotate

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math/rand"
	"melody/bf"
	"melody/bf/bloomfilter"
	"melody/bf/testutils"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRotate_Union_ok(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	set2 := New(ctx, Config{testutils.TestCfg, 5})

	testutils.CallSetUnion(t, set1, set2)
}

func TestRotate_Union_koIncorrectType(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	set2 := 24

	if _, err := set1.Union(set2); err != bf.ErrImpossibleToTreat {
		t.Errorf("Unexpected error, %v", err)
	}
}

func TestRotate_Union_koIncompatibleN(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	cfg := testutils.TestCfg
	cfg.N = 1
	set2 := New(ctx, Config{cfg, 5})
	if _, err := set1.Union(set2); err == nil || !strings.Contains(err.Error(), "error: diferrent n values") {
		t.Errorf("Unexpected error, %v", err)
	}
}

func TestRotate_Union_koIncompatibleP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	cfg := testutils.TestCfg
	cfg.P = 0.5
	set2 := New(ctx, Config{cfg, 5})
	if _, err := set1.Union(set2); err == nil || !strings.Contains(err.Error(), "error: diferrent p values") {
		t.Errorf("Unexpected error, %v", err)
	}
}

func TestRotate_Union_koIncompatibleCurrentBFs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	set2 := New(ctx, Config{testutils.TestCfg2, 5})
	if _, err := set1.Union(set2); err == nil || !strings.Contains(err.Error(), "error: diferrent p values") {
		t.Errorf("Unexpected error, %v", err)
	}
}

func TestRotate_Union_koDifferentHashFuncsBFs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	set2 := New(ctx, Config{testutils.TestCfg3, 5})
	if _, err := set1.Union(set2); err == nil || !strings.Contains(err.Error(), "error: different hashers") {
		t.Errorf("Unexpected error, %v", err)
	}
}

func TestRotate_Unmarshal_ok(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	elem := []byte("wwwww")
	set1.Add(elem)
	set2 := New(ctx, Config{testutils.TestCfg, 5})
	if set2.Check(elem) {
		t.Errorf("Unexpected elem %s in set2", elem)
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(set1); err != nil {
		t.Errorf("Unexpected error, %v", err)
	}
	if err := gob.NewDecoder(buf).Decode(set2); err != nil {
		t.Errorf("Unexpected error, %v", err)
	}
	if !set2.Check(elem) {
		t.Errorf("Expecting elem %s in set2", elem)
	}
}

func TestRotate_UnmarshalBinary_ko(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set1 := New(ctx, Config{testutils.TestCfg, 5})
	if err := set1.UnmarshalBinary([]byte{}); err == nil {
		t.Error("should have given error")
	}
}

func TestRotate_KeepRotating(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dt := 5 * time.Millisecond

	rotate := &BloomFilter{
		Previous: bloomfilter.New(testutils.TestCfg),
		Current:  bloomfilter.New(testutils.TestCfg),
		Next:     bloomfilter.New(testutils.TestCfg),
		Config:   Config{testutils.TestCfg, 5},
		cancel:   cancel,
		mutex:    &sync.RWMutex{},
		ctx:      ctx,
	}

	ch := make(chan time.Time)
	go rotate.keepRotating(ctx, ch)

	rotate.Add([]byte("test"))
	if !rotate.Check([]byte("test")) {
		t.Error("error: \"test\" not present")
	}

	serialized := new(bytes.Buffer)
	if err := gob.NewEncoder(serialized).Encode(rotate); err != nil {
		t.Errorf("error encoding Rotate, %s", err.Error())
	}

	ch <- time.Now()
	<-time.After(dt)
	if !rotate.Check([]byte("test")) {
		t.Error("error: \"test\" not present after 1 TTL")
	}

	ch <- time.Now()
	<-time.After(dt)
	if !rotate.Check([]byte("test")) {
		t.Error("error: \"test\" not present after 2 TTL")

	}
	ch <- time.Now()
	<-time.After(dt)
	if rotate.Check([]byte("test")) {
		t.Error("error: \"test\" present after 3 TTL")

	}

	rotate2 := new(BloomFilter)
	if err := gob.NewDecoder(serialized).Decode(rotate2); err != nil {
		t.Errorf("error encoding Rotate, %s", err.Error())
	}

	if !rotate2.Check([]byte("test")) {
		t.Error("error: \"test\" not present")
	}
}

func BenchmarkRotate_UnmarshalBinary_GZIP(b *testing.B) {
	compressor = new(Gzip)
	cfg := Config{
		bf.Config{
			N:        1000000,
			P:        1e-7,
			HashName: bf.HASHER_OPTIMAL,
		},
		10000,
	}
	ctx, cancel := context.WithCancel(context.Background())
	benchmarkRotate_UnmarshalBinary(b, New(ctx, cfg))
	cancel()
}

func benchmarkRotate_UnmarshalBinary(b *testing.B, bf bf.BloomFilter) {
	buf := make([]byte, 150*1000*1024)
	rand.Read(buf)

	for _, size := range []int{10, 1024, 1000 * 1024, 10 * 1000 * 1024, 100 * 1000 * 1024} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				offset := i % (len(buf) - size)
				bf.Add(buf[offset : offset+size])
			}
		})
	}
}
