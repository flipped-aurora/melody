package melodyrate

const (
	offset64 uint64 = 14695981039346656037
	prime64         = 1099511628211
)

// PseudoFNV64a 接受一个string
func PseudoFNV64a(s string) uint64 {
	h := offset64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}
