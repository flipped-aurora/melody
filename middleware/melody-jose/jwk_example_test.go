// +build integration

package jose

import "fmt"

func Example_Auth0Integration() {
	fs, _ := DecodeFingerprints([]string{"--MBgDH5WGvL9Bcn5Be30cRcL0f5O-NyoXuWtQdX1aI="})
	cfg := SecretProviderConfig{
		URI:          "https://albert-test.auth0.com/.well-known/jwks.json",
		Fingerprints: fs,
	}
	client, _ := SecretProvider(cfg, nil)

	k, err := client.GetKey("MDNGMjU2M0U3RERFQUEwOUUzQUMwQ0NBN0Y1RUY0OEIxNTRDM0IxMw")
	fmt.Println("err:", err)
	fmt.Println("is public:", k.IsPublic())
	fmt.Println("alg:", k.Algorithm)
	fmt.Println("id:", k.KeyID)
	// Output:
	// err: <nil>
	// is public: true
	// alg: RS256
	// id: MDNGMjU2M0U3RERFQUEwOUUzQUMwQ0NBN0Y1RUY0OEIxNTRDM0IxMw
}

func Example_Auth0Integration_badFingerprint() {
	cfg := SecretProviderConfig{
		URI:          "https://albert-test.auth0.com/.well-known/jwks.json",
		Fingerprints: [][]byte{make([]byte, 32)},
	}
	client, _ := SecretProvider(cfg, nil)

	_, err := client.GetKey("MDNGMjU2M0U3RERFQUEwOUUzQUMwQ0NBN0Y1RUY0OEIxNTRDM0IxMw")
	fmt.Println("err:", err)
	// Output:
	// err: Get https://albert-test.auth0.com/.well-known/jwks.json: JWK client did not find a pinned key
}
