package jose

import (
	"gopkg.in/square/go-jose.v2"
	"melody/config"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_getSignatureConfig(t *testing.T) {
	server := httptest.NewServer(jwkEndpoint("private"))
	defer server.Close()

	scfg, err := GetSignatureConfig(newVerifierEndpointCfg("RS256", server.URL, []string{}))
	if err != nil {
		t.Error(err.Error())
		return
	}

	if scfg.Issuer != "http://example.com" {
		t.Errorf("unexpected issuer: %s", scfg.Issuer)
	}

	if scfg.Audience[0] != "http://api.example.com" {
		t.Errorf("unexpected audience: %v", scfg.Audience)
	}
}

func Test_getSignatureConfig_unsecure(t *testing.T) {
	cfg := &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/private",
		Backends: []*config.Backend{
			{
				URLPattern: "/",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			ValidatorNamespace: map[string]interface{}{
				"alg":      "RS256",
				"jwk-url":  "http://jwk.example.com",
				"audience": []string{"http://api.example.com"},
				"issuer":   "http://example.com",
				"roles":    []string{},
				"cache":    false,
			},
		},
	}

	_, err := GetSignatureConfig(cfg)
	if err != ErrInsecureJWKSource {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_getSignatureConfig_wrongStruct(t *testing.T) {
	cfg := &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/private",
		Backends: []*config.Backend{
			{
				URLPattern: "/",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			ValidatorNamespace: true,
		},
	}

	_, err := GetSignatureConfig(cfg)
	if err == nil || err.Error() != "json: cannot unmarshal bool into Go value of type jose.SignatureConfig" {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_newSigner(t *testing.T) {
	server := httptest.NewServer(jwkEndpoint("private"))
	defer server.Close()

	_, signer, err := NewSigner(newSignerEndpointCfg("RS256", "2011-04-29", server.URL), nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	msg, err := signer(map[string]interface{}{
		"aud": "http://api.example.com",
		"iss": "http://example.com",
		"sub": "1234567890qwertyuio",
		"jti": "mnb23vcsrt756yuiomnbvcx98ertyuiop",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}

	expected := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjIwMTEtMDQtMjkifQ.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cDovL2V4YW1wbGUuY29tIiwianRpIjoibW5iMjN2Y3NydDc1Nnl1aW9tbmJ2Y3g5OGVydHl1aW9wIiwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9.TWdsBQPqfDV1IFe1iD0KFu-E_wqeFXgNJXIoESl9smg2W_Snh2GwwktwlvHCSAGvUdkKU6Js6LQ594e6HZ3eAdj3mfCdCxerhuodb6GS-rZ2OrMv44VaC_YnzoOjCWUrU3ivzhYjEFBxgDgWc0G9qFdQVaZPOLPohd_mXpeM5jAS-vFzudOlJz8rtK9KfVDPiAWnGxih5fa3MF1b19vnnsfyN1Y8hTeen3j24thQbuh61vkqu8TLoG2NrETyC9zqCuL3IQnPld3IBolYJhqEcka95cCNZ1dQnqsgrP4q325JmRxXsn0GJM3VtFpKbfJCcQgdpixCohQ-_xHmTUpXng"
	if msg != expected {
		t.Errorf("unexpected signed payload: %s", msg)
	}
}

func Test_newSigner_unsecure(t *testing.T) {
	cfg := &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/token",
		Method:   "POST",
		Backends: []*config.Backend{
			{
				URLPattern: "/token",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			SignerNamespace: map[string]interface{}{
				"alg":          "RS256",
				"kid":          "2011-04-29",
				"jwk-url":      "http://jwk.example.com",
				"keys-to-sign": []string{"access_token", "refresh_token"},
			},
		},
	}
	_, _, err := NewSigner(cfg, nil)
	if err != ErrInsecureJWKSource {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_newSigner_wrongStruct(t *testing.T) {
	cfg := &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/token",
		Method:   "POST",
		Backends: []*config.Backend{
			{
				URLPattern: "/token",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			SignerNamespace: true,
		},
	}
	_, _, err := NewSigner(cfg, nil)
	if err == nil || err.Error() != "json: cannot unmarshal bool into Go value of type jose.SignerConfig" {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_newSigner_unknownKey(t *testing.T) {
	server := httptest.NewServer(jwkEndpoint("private"))
	defer server.Close()

	_, _, err := NewSigner(newSignerEndpointCfg("RS256", "unknown key", server.URL), nil)
	if err == nil || err.Error() != "no Keys has been found" {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_RSAPrivateSigner(t *testing.T) {
	testPrivateSigner(
		t,
		"private",
		"2011-04-29",
		`{"payload":"eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cDovL2V4YW1wbGUuY29tIiwianRpIjoibW5iMjN2Y3NydDc1Nnl1aW9tbmJ2Y3g5OGVydHl1aW9wIiwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9","protected":"eyJhbGciOiJSUzI1NiJ9","signature":"Cz7OEXmH6CsjFYFnGyrGMe7QsjrTk-QLTfP4VL6CZVpKKeVYKSI0NlquzlEGgwY3pujhdpQGVV2md3MvrccY6-a7-C8nRjyv4TnYkAk0lQcdmaG4hd38SwG0jZ6LpzgyL5l51txQATnayZgbRuUVzco-AZTPfTw9xS15CDPtFjoQHAKe9w9kvCGR6RmyOP29-YgqAk20hVqUj5EiFD_q-m2lTGYIAAYiWqYon661Ep9vfRNO1acq9ch_7qe1UBSmWu1BGiN2u7sq0rlkJ0Z4WKQL914eEiVBLRAUEsYpV-W-OBKmM3NMsm_2Ems0CcCaax0OrS0nHViuI4ZeT_molw"}`,
		"eyJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cDovL2V4YW1wbGUuY29tIiwianRpIjoibW5iMjN2Y3NydDc1Nnl1aW9tbmJ2Y3g5OGVydHl1aW9wIiwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9.Cz7OEXmH6CsjFYFnGyrGMe7QsjrTk-QLTfP4VL6CZVpKKeVYKSI0NlquzlEGgwY3pujhdpQGVV2md3MvrccY6-a7-C8nRjyv4TnYkAk0lQcdmaG4hd38SwG0jZ6LpzgyL5l51txQATnayZgbRuUVzco-AZTPfTw9xS15CDPtFjoQHAKe9w9kvCGR6RmyOP29-YgqAk20hVqUj5EiFD_q-m2lTGYIAAYiWqYon661Ep9vfRNO1acq9ch_7qe1UBSmWu1BGiN2u7sq0rlkJ0Z4WKQL914eEiVBLRAUEsYpV-W-OBKmM3NMsm_2Ems0CcCaax0OrS0nHViuI4ZeT_molw",
	)
}

func Test_HSAPrivateSigner(t *testing.T) {
	testPrivateSigner(
		t,
		"symmetric",
		"sim2",
		`{"payload":"eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cDovL2V4YW1wbGUuY29tIiwianRpIjoibW5iMjN2Y3NydDc1Nnl1aW9tbmJ2Y3g5OGVydHl1aW9wIiwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9","protected":"eyJhbGciOiJIUzI1NiJ9","signature":"2eGKzqRiIJE5TJ4WcgnmopwhUczIdTFuQkp9ZVuFyUk"}`,
		"eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJodHRwOi8vYXBpLmV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cDovL2V4YW1wbGUuY29tIiwianRpIjoibW5iMjN2Y3NydDc1Nnl1aW9tbmJ2Y3g5OGVydHl1aW9wIiwic3ViIjoiMTIzNDU2Nzg5MHF3ZXJ0eXVpbyJ9.2eGKzqRiIJE5TJ4WcgnmopwhUczIdTFuQkp9ZVuFyUk",
	)
}

func testPrivateSigner(t *testing.T, keyType, keyName, full, compact string) {
	server := httptest.NewServer(jwkEndpoint(keyType))
	defer server.Close()

	sp, err := SecretProvider(SecretProviderConfig{URI: server.URL}, nil)
	if err != nil {
		t.Error(err)
		return
	}
	key, err := sp.GetKey(keyName)
	if err != nil {
		t.Errorf("getting the key: %s", err.Error())
		return
	}

	signingKey := jose.SigningKey{
		Key:       key.Key,
		Algorithm: jose.SignatureAlgorithm(key.Algorithm),
	}
	s, err := jose.NewSigner(signingKey, nil)
	if err != nil {
		t.Errorf("building the signer: %s", err.Error())
		return
	}

	payload := map[string]interface{}{
		"aud": "http://api.example.com",
		"iss": "http://example.com",
		"sub": "1234567890qwertyuio",
		"jti": "mnb23vcsrt756yuiomnbvcx98ertyuiop",
	}
	for _, tc := range []struct {
		Name     string
		Signer   Signer
		Expected string
	}{
		{
			Name:     keyType + "-full",
			Signer:   fullSerializeSigner{signer{s}}.Sign,
			Expected: full,
		},
		{
			Name:     keyType + "-compact",
			Signer:   compactSerializeSigner{signer{s}}.Sign,
			Expected: compact,
		},
	} {
		data, err := tc.Signer(payload)
		if err != nil {
			t.Errorf("[%s] signing the payload: %s", tc.Name, err.Error())
			return
		}
		if data != tc.Expected {
			t.Errorf("[%s] unexpected signed payload: %s", tc.Name, data)
		}
	}
}

func newSignerEndpointCfg(alg, ID, URL string) *config.EndpointConfig {
	return &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/token",
		Method:   "POST",
		Backends: []*config.Backend{
			{
				URLPattern: "/token",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			SignerNamespace: map[string]interface{}{
				"alg":                  alg,
				"kid":                  ID,
				"jwk-url":              URL,
				"keys-to-sign":         []string{"access_token", "refresh_token"},
				"disable_jwk_security": true,
				"cache":                true,
			},
		},
	}
}

func newVerifierEndpointCfg(alg, URL string, roles []string) *config.EndpointConfig {
	return &config.EndpointConfig{
		Timeout:  time.Second,
		Endpoint: "/private",
		Backends: []*config.Backend{
			{
				URLPattern: "/",
				Host:       []string{"http://example.com/"},
				Timeout:    time.Second,
			},
		},
		ExtraConfig: config.ExtraConfig{
			ValidatorNamespace: map[string]interface{}{
				"alg":                  alg,
				"jwk-url":              URL,
				"audience":             []string{"http://api.example.com"},
				"issuer":               "http://example.com",
				"roles":                roles,
				"disable_jwk_security": true,
				"cache":                true,
			},
		},
	}
}
