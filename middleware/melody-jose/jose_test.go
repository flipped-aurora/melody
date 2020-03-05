package jose

import (
	"net/http"
	"testing"

	"gopkg.in/square/go-jose.v2/jwt"
)

func nopExtractor(_ string) func(r *http.Request) (*jwt.JSONWebToken, error) {
	return func(_ *http.Request) (*jwt.JSONWebToken, error) { return nil, nil }
}

func Test_NewValidator_unkownAlg(t *testing.T) {
	_, err := NewValidator(&SignatureConfig{
		Alg: "random",
	}, nopExtractor)
	if err == nil || err.Error() != "JOSE: unknown algorithm random" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCanAccess(t *testing.T) {
	for _, v := range []struct {
		name         string
		roleKey      string
		claims       map[string]interface{}
		requirements []string
		expected     bool
	}{
		{
			name:         "simple_success",
			roleKey:      "role",
			claims:       map[string]interface{}{"role": []interface{}{"a", "b"}},
			requirements: []string{"a"},
			expected:     true,
		},
		{
			name:         "simple_sfail",
			roleKey:      "role",
			claims:       map[string]interface{}{"role": []interface{}{"c", "b"}},
			requirements: []string{"a"},
			expected:     false,
		},
		{
			name:         "multiple_success",
			roleKey:      "role",
			claims:       map[string]interface{}{"role": []interface{}{"c"}},
			requirements: []string{"a", "b", "c"},
			expected:     true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			if res := CanAccess(v.roleKey, v.claims, v.requirements); res != v.expected {
				t.Errorf("'%s' have %v, want %v", v.name, res, v.expected)
			}
		})
	}
}

func TestCanAccessNested(t *testing.T) {
	for _, v := range []struct {
		name         string
		roleKey      string
		claims       map[string]interface{}
		requirements []string
		expected     bool
	}{
		{
			name:         "simple_success",
			roleKey:      "role",
			claims:       map[string]interface{}{"role": []interface{}{"a", "b"}},
			requirements: []string{"a"},
			expected:     true,
		},
		{
			name:         "simple_sfail",
			roleKey:      "role",
			claims:       map[string]interface{}{"role": []interface{}{"c", "b"}},
			requirements: []string{"a"},
			expected:     false,
		},
		{
			name:         "multiple_success",
			roleKey:      "role",
			claims:       map[string]interface{}{"role": []interface{}{"c"}},
			requirements: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "struct_success",
			roleKey:      "data.role",
			claims:       map[string]interface{}{"data": map[string]interface{}{"role": []interface{}{"c"}}},
			requirements: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:    "complex_struct_success",
			roleKey: "data.data.data.data.data.data.data.role",
			claims: map[string]interface{}{
				"data": map[string]interface{}{
					"data": map[string]interface{}{
						"data": map[string]interface{}{
							"data": map[string]interface{}{
								"data": map[string]interface{}{
									"data": map[string]interface{}{
										"data": map[string]interface{}{
											"role": []interface{}{"c"},
										},
									},
								},
							},
						},
					},
				},
			},
			requirements: []string{"a", "b", "c"},
			expected:     true,
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			if res := CanAccessNested(v.roleKey, v.claims, v.requirements); res != v.expected {
				t.Errorf("'%s' have %v, want %v", v.name, res, v.expected)
			}
		})
	}
}
