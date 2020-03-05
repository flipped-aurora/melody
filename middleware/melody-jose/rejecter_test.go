package jose

import (
	"melody/config"
	"melody/logging"
	"testing"
)

func TestChainedRejecterFactory(t *testing.T) {
	rf := ChainedRejecterFactory([]RejecterFactory{
		NopRejecterFactory{},
		RejecterFactoryFunc(func(_ logging.Logger, _ *config.EndpointConfig) Rejecter {
			return RejecterFunc(func(in map[string]interface{}) bool {
				v, ok := in["key"].(int)
				return ok && v == 42
			})
		}),
	})

	rejecter := rf.New(nil, nil)

	for _, tc := range []struct {
		name     string
		in       map[string]interface{}
		expected bool
	}{
		{
			name: "empty",
			in:   map[string]interface{}{},
		},
		{
			name:     "reject",
			in:       map[string]interface{}{"key": 42},
			expected: true,
		},
		{
			name: "pass-1",
			in:   map[string]interface{}{"key": "42"},
		},
		{
			name: "pass-2",
			in:   map[string]interface{}{"key": 9876},
		},
	} {
		if v := rejecter.Reject(tc.in); tc.expected != v {
			t.Errorf("unexpected result for %s. have %v, want %v", tc.name, v, tc.expected)
		}
	}
}
