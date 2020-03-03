package jose

//JOSE是一个框架，旨在提供一种在各方之间安全地转移声明（如授权信息）的方法。JOSE框架提供了一系列规范来实现此目的。它是由一组规范构成:
//
//JSON Web Token (JWT)：JSON Web令牌（RFC7519），定义了一种可以签名或加密的标准格式；
//JSON Web Signature (JWS)：JSON Web签名（RFC7515），定义对JWT进行数字签名的过程；
//JSON Web Encryption (JWE)：JSON Web加密（RFC7516），定义加密JWT的过程；
//JSON Web Algorithm（JWA）：JSON Web算法（RFC7518），定义用于数字签名或加密的算法列表；
//JSON Web Key (JWK)：JSON Web密钥（RFC7517），定义加密密钥和密钥集的表示方式；

//import (
//	"github.com/auth0-community/go-auth0"
//	jose "gopkg.in/square/go-jose.v2"
//	"gopkg.in/square/go-jose.v2/jwt"
//	"net/http"
//)
//
//type ExtractorFactory func(string) func(r *http.Request) (*jwt.JSONWebToken, error)
//
//func NewValidator(signatureConfig *SigntureConfig)  {
//
//}
