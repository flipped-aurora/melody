//go:generate go run $GOROOT/src/crypto/tls/generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,localhost --ca --start-date "Jan 1 00:00:00 2000" --duration=1000000h
package jose

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWK(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		t.Error(err)
		return
	}

	for _, tc := range []struct {
		Name string
		Alg  string
		ID   []string
	}{
		{
			Name: "public",
			ID:   []string{"2011-04-29"},
			Alg:  "RS256",
		},
		{
			Name: "public",
			ID:   []string{"1"},
		},
		{
			Name: "private",
			ID:   []string{"2011-04-29"},
			Alg:  "RS256",
		},
		{
			Name: "private",
			ID:   []string{"1"},
		},
		{
			Name: "symmetric",
			ID:   []string{"sim2"},
			Alg:  "HS256",
		},
	} {
		server := httptest.NewUnstartedServer(jwkEndpoint(tc.Name))
		defer server.Close()
		server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
		server.StartTLS()

		secretProvidr, err := SecretProvider(SecretProviderConfig{URI: server.URL, LocalCA: "cert.pem"}, nil)
		if err != nil {
			t.Error(err)
		}
		for _, k := range tc.ID {
			key, err := secretProvidr.GetKey(k) // 发get请求
			if err != nil {
				t.Errorf("[%s] extracting the key %s: %s", tc.Name, k, err.Error())
			}
			if key.Algorithm != tc.Alg {
				t.Errorf("wrong alg. have: %s, want: %s", key.Algorithm, tc.Alg)
			}
		}
	}
}

func TestDialer_DialTLS_ko(t *testing.T) {
	d := NewDialer(SecretProviderConfig{})
	c, err := d.DialTLS("\t", "addr")
	if err == nil {
		t.Error(err)
	}
	if c != nil {
		t.Errorf("unexpected connection: %v", c)
	}
}

func Test_decodeFingerprints(t *testing.T) {
	_, err := DecodeFingerprints([]string{"not_encoded_message"})
	if err == nil {
		t.Error(err)
	}
}

func jwkEndpoint(name string) http.HandlerFunc {
	data, err := ioutil.ReadFile("./fixtures/" + name + ".json")
	return func(rw http.ResponseWriter, _ *http.Request) {
		if err != nil {
			rw.WriteHeader(500)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(data)
	}
}
