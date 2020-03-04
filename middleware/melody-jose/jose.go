package jose

//JOSE是一个框架，旨在提供一种在各方之间安全地转移声明（如授权信息）的方法。JOSE框架提供了一系列规范来实现此目的。它是由一组规范构成:
//
//JSON Web Token (JWT)：JSON Web令牌（RFC7519），定义了一种可以签名或加密的标准格式；
//JSON Web Signature (JWS)：JSON Web签名（RFC7515），定义对JWT进行数字签名的过程；
//JSON Web Encryption (JWE)：JSON Web加密（RFC7516），定义加密JWT的过程；
//JSON Web Algorithm（JWA）：JSON Web算法（RFC7518），定义用于数字签名或加密的算法列表；
//JSON Web Key (JWK)：JSON Web密钥（RFC7517），定义加密密钥和密钥集的表示方式；

import (
	"fmt"
	"github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"melody/proxy"
	"net/http"
	"strings"
)

type ExtractorFactory func(string) func(r *http.Request) (*jwt.JSONWebToken, error)

func NewValidator(signatureConfig *SignatureConfig, ef ExtractorFactory) (*auth0.JWTValidator, error) {
	sa, ok := supportedAlgorithms[signatureConfig.Alg]
	if !ok {
		return nil, fmt.Errorf("JOSE: unknown algorithm %s", signatureConfig.Alg)
	}
	te := auth0.FromMultiple(
		auth0.RequestTokenExtractorFunc(auth0.FromHeader),
		auth0.RequestTokenExtractorFunc(ef(signatureConfig.CookieKey)),
	)

	decodedFs, err := DecodeFingerprints(signatureConfig.Fingerprints)
	if err != nil {
		return nil, err
	}

	cfg := SecretProviderConfig{
		URI:           signatureConfig.URI,
		CacheEnabled:  signatureConfig.CacheEnabled,
		Cs:            signatureConfig.CipherSuites,
		Fingerprints:  decodedFs,
		LocalCA:       signatureConfig.LocalCA,
		AllowInsecure: signatureConfig.DisableJWKSecurity,
	}

	sp, err := SecretProvider(cfg, te)
	if err != nil {
		return nil, err
	}

	return auth0.NewValidator(
		auth0.NewConfiguration(
			sp,
			signatureConfig.Audience,
			signatureConfig.Issuer,
			sa,
		),
		te,
	), nil
}

func CanAccessNested(roleKey string, claims map[string]interface{}, required []string) bool {
	if len(required) == 0 {
		return true
	}

	tmp := claims
	keys := strings.Split(roleKey, ".")

	for _, key := range keys[:len(keys)-1] {
		v, ok := tmp[key]
		if !ok {
			return false
		}
		tmp, ok = v.(map[string]interface{})
		if !ok {
			return false
		}
	}
	return CanAccess(keys[len(keys)-1], tmp, required)
}

func CanAccess(roleKey string, claims map[string]interface{}, required []string) bool {
	if len(required) == 0 {
		return true
	}

	tmp, ok := claims[roleKey]
	if !ok {
		return false
	}

	roles, ok := tmp.([]interface{})
	if !ok {
		return false
	}

	for _, role := range required {
		for _, r := range roles {
			if r.(string) == role {
				return true
			}
		}
	}
	return false
}

func SignFields(keys []string, signer Signer, response *proxy.Response) error {
	for _, key := range keys {
		tmp, ok := response.Data[key]
		if !ok {
			continue
		}
		data, ok := tmp.(map[string]interface{})
		if !ok {
			continue
		}
		token, err := signer(data)
		if err != nil {
			return err
		}
		response.Data[key] = token
	}
	return nil
}

var supportedAlgorithms = map[string]jose.SignatureAlgorithm{
	"EdDSA": jose.EdDSA,
	"HS256": jose.HS256,
	"HS384": jose.HS384,
	"HS512": jose.HS512,
	"RS256": jose.RS256,
	"RS384": jose.RS384,
	"RS512": jose.RS512,
	"ES256": jose.ES256,
	"ES384": jose.ES384,
	"ES512": jose.ES512,
	"PS256": jose.PS256,
	"PS384": jose.PS384,
	"PS512": jose.PS512,
}
