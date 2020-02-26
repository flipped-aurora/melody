package bf

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHasher(t *testing.T) {
	for _, hash := range defaultHashers {
		array1 := []byte{1, 2, 3}
		if !reflect.DeepEqual(hash(array1), hash(array1)) {
			t.Error("undeterministic")
		}
	}
}

func BenchmarkHasher(b *testing.B) {
	for k, hash := range defaultHashers {
		b.Run(fmt.Sprintf("hasher %d", k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				array1 := []byte{1, 2, 3}
				hash(array1)
			}
		})
	}
}

func TestDefaultHashFactory(t *testing.T) {
	for _, hash := range DefaultHashFactory(23) {
		array1 := []byte{1, 2, 3}
		if !reflect.DeepEqual(hash(array1), hash(array1)) {
			t.Error("undeterministic")
		}
	}
}

func TestOptimalHashFactory(t *testing.T) {
	for _, hash := range OptimalHashFactory(23) {
		array1 := []byte{1, 2, 3}
		if !reflect.DeepEqual(hash(array1), hash(array1)) {
			t.Error("undeterministic")
		}
	}
}

func BenchmarkOptimalHashFactory(b *testing.B) {
	for k, hash := range OptimalHashFactory(23) {
		b.Run(fmt.Sprintf("hasher %d", k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				array1 := []byte{1, 2, 3}
				hash(array1)
			}
		})
	}
}
