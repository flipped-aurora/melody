// JSON Web Signature (JWS)：JSON Web签名（RFC7515），定义对JWT进行数字签名的过程
// 私钥签名 公钥验证

package jose

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
	"melody/config"
	"strings"
)

const (
	ValidatorNamespace = "melody_jose_validator"
	SignerNamespace    = "melody_jose_signer"
	defaultRolesKey    = "roles"
)

type SignatureConfig struct {
	Alg                string   `json:"alg"`
	URI                string   `json:"jwk-url"`
	CacheEnabled       bool     `json:"cache,omitempty"`
	Issuer             string   `json:"issuer,omitempty"`
	Audience           []string `json:"audience,omitempty"`
	Roles              []string `json:"roles,omitempty"`
	RolesKey           string   `json:"roles_key,omitempty"`
	CookieKey          string   `json:"cookie_key,omitempty"`
	CipherSuites       []uint16 `json:"cipher_suites,omitempty"`
	DisableJWKSecurity bool     `json:"disable_jwk_security"`
	Fingerprints       []string `json:"jwk_fingerprints,omitempty"`
	LocalCA            string   `json:"jwk_local_ca,omitempty"`
}

type SignerConfig struct {
	Alg                string   `json:"alg"`
	KeyID              string   `json:"kid"`
	URI                string   `json:"jwk-url"`
	FullSerialization  bool     `json:"full,omitempty"`
	KeysToSign         []string `json:"keys-to-sign,omitempty"`
	CipherSuites       []uint16 `json:"cipher_suites,omitempty"`
	DisableJWKSecurity bool     `json:"disable_jwk_security"`
	Fingerprints       []string `json:"jwk_fingerprints,omitempty"`
	LocalCA            string   `json:"jwk_local_ca,omitempty"`
}

var (
	ErrNoValidatorCfg = errors.New("JOSE: no validator config")
	ErrNoSignerCfg    = errors.New("JOSE: no signer config")
)

func GetSignatureConfig(cfg *config.EndpointConfig) (*SignatureConfig, error) {
	tmp, ok := cfg.ExtraConfig[ValidatorNamespace]
	if !ok {
		return nil, ErrNoValidatorCfg
	}
	data, _ := json.Marshal(tmp)
	res := new(SignatureConfig)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}

	if res.RolesKey == "" {
		res.RolesKey = defaultRolesKey
	}
	//if !strings.HasPrefix(res.URI, "https://") && !res.DisableJWKSecurity {
	//	return res, ErrInsecureJWKSource
	//}
	return res, nil
}

func NewSigner(cfg *config.EndpointConfig, te auth0.RequestTokenExtractor) (*SignerConfig, Signer, error) {
	signerCfg, err := getSignerConfig(cfg)
	if err != nil {
		return signerCfg, nopSigner, err
	}

	var decodedFs [][]byte
	if !signerCfg.DisableJWKSecurity {
		// 如果不禁止安全性
		// 查看指纹是否格式正确
		decodedFs, err = DecodeFingerprints(signerCfg.Fingerprints)
	}
	if err != nil {
		return signerCfg, nopSigner, err
	}
	spCfg := SecretProviderConfig{
		URI:           signerCfg.URI,
		Fingerprints:  decodedFs,
		Cs:            signerCfg.CipherSuites,
		LocalCA:       signerCfg.LocalCA,
		AllowInsecure: signerCfg.DisableJWKSecurity,
	}

	sp, err := SecretProvider(spCfg, te)
	if err != nil {
		return signerCfg, nopSigner, err
	}
	key, err := sp.GetKey(signerCfg.KeyID) // get 请求 signerCfg.URI
	if err != nil {
		return signerCfg, nopSigner, err
	}
	if key.IsPublic() {
		// TODO: 不用公钥签名
	}
	signingKey := jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(signerCfg.Alg),
		Key:       key.Key,
	}
	opts := &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			jose.HeaderKey("kid"): key.KeyID,
		},
	}
	s, err := jose.NewSigner(signingKey, opts)
	if err != nil {
		return signerCfg, nopSigner, err
	}

	if signerCfg.FullSerialization {
		return signerCfg, fullSerializeSigner{signer{s}}.Sign, nil
	}
	return signerCfg, compactSerializeSigner{signer{s}}.Sign, nil
}

// uri = https://xxx   disableJWKSecurity = true
func getSignerConfig(cfg *config.EndpointConfig) (*SignerConfig, error) {
	tmp, ok := cfg.ExtraConfig[SignerNamespace]
	if !ok {
		return nil, ErrNoSignerCfg
	}
	data, _ := json.Marshal(tmp)
	res := new(SignerConfig)
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(res.URI, "https://") && !res.DisableJWKSecurity {
		return res, ErrInsecureJWKSource
	}
	return res, nil
}

type Signer func(interface{}) (string, error)

func nopSigner(_ interface{}) (string, error) { return "", nil }

type signer struct {
	signer jose.Signer
}

func (s signer) sign(v interface{}) (*jose.JSONWebSignature, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize payload: %s", err.Error())
	}
	return s.signer.Sign(data)
}

type fullSerializeSigner struct {
	signer
}

func (f fullSerializeSigner) Sign(v interface{}) (string, error) {
	obj, err := f.sign(v)
	if err != nil {
		return "", fmt.Errorf("unable to sign payload: %s", err.Error())
	}
	return obj.FullSerialize(), nil
}

type compactSerializeSigner struct {
	signer
}

func (c compactSerializeSigner) Sign(v interface{}) (string, error) {
	obj, err := c.sign(v)
	if err != nil {
		return "", fmt.Errorf("unable to sign payload: %s", err.Error())
	}
	return obj.CompactSerialize()
}
