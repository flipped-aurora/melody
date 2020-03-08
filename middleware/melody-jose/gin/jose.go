package gin

import (
	"context"
	"fmt"
	"melody/config"
	"melody/logging"
	melodyjose "melody/middleware/melody-jose"
	"melody/proxy"
	melodygin "melody/router/gin"
	"net/http"
	"strings"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/jwt"
)

func HandlerFactory(hf melodygin.HandlerFactory, logger logging.Logger, rejecterF melodyjose.RejecterFactory) melodygin.HandlerFactory {
	return TokenSigner(TokenSignatureValidator(hf, logger, rejecterF), logger)
}

func TokenSigner(hf melodygin.HandlerFactory, logger logging.Logger) melodygin.HandlerFactory {
	return func(cfg *config.EndpointConfig, prxy proxy.Proxy) gin.HandlerFunc {
		signerCfg, signer, err := melodyjose.NewSigner(cfg, nil)
		// 如果是签名接口 则返回签名后的 token
		// 如果不是签名接口 则判断是否需要签名 和 签名的正确性
		if err == melodyjose.ErrNoSignerCfg {
			logger.Info("JOSE: singer disabled for the endpoint", cfg.Endpoint)
			return hf(cfg, prxy)
		}
		if err != nil {
			logger.Error(err.Error(), cfg.Endpoint)
			return hf(cfg, prxy)
		}

		logger.Info("JOSE: singer enabled for the endpoint", cfg.Endpoint)

		return func(c *gin.Context) {
			proxyReq := melodygin.NewRequest(cfg.HeadersToPass)(c, cfg.QueryString)
			ctx, cancel := context.WithTimeout(c, cfg.Timeout)
			defer cancel()

			response, err := prxy(ctx, proxyReq)
			if err != nil {
				logger.Error("proxy response error:", err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if response == nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if err := melodyjose.SignFields(signerCfg.KeysToSign, signer, response); err != nil {
				logger.Error(err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			for k, v := range response.Metadata.Headers {
				c.Header(k, v[0])
			}
			c.JSON(response.Metadata.StatusCode, response.Data)
		}
	}
}

func TokenSignatureValidator(hf melodygin.HandlerFactory, logger logging.Logger, rejecterF melodyjose.RejecterFactory) melodygin.HandlerFactory {
	return func(cfg *config.EndpointConfig, prxy proxy.Proxy) gin.HandlerFunc {
		if rejecterF == nil {
			rejecterF = new(melodyjose.NopRejecterFactory)
		}
		rejecter := rejecterF.New(logger, cfg)

		handler := hf(cfg, prxy) // 完成juju的准备工作
		signatureCfg, err := melodyjose.GetSignatureConfig(cfg)
		if err == melodyjose.ErrNoValidatorCfg {
			logger.Info("JOSE: validator disabled for the endpoint", cfg.Endpoint)
			return handler
		}
		if err != nil {
			logger.Warning(fmt.Sprintf("JOSE: validator for %s: %s", cfg.Endpoint, err.Error()))
			return handler
		}

		validator, err := melodyjose.NewValidator(signatureCfg, FromCookie)
		if err != nil {
			logger.Fatal("%s: %s", cfg.Endpoint, err.Error())
		}

		var aclCheck func(string, map[string]interface{}, []string) bool

		if strings.Contains(signatureCfg.RolesKey, ".") {
			aclCheck = melodyjose.CanAccessNested
		} else {
			aclCheck = melodyjose.CanAccess
		}

		logger.Info("JOSE: validator enabled for the endpoint", cfg.Endpoint)

		return func(c *gin.Context) {
			token, err := validator.ValidateRequest(c.Request)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

			claims := map[string]interface{}{}
			err = validator.Claims(c.Request, token, &claims)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

			if rejecter.Reject(claims) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if !aclCheck(signatureCfg.RolesKey, claims, signatureCfg.Roles) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			handler(c)
		}
	}
}

func FromCookie(key string) func(r *http.Request) (*jwt.JSONWebToken, error) {
	if key == "" {
		key = "access_token"
	}
	return func(r *http.Request) (*jwt.JSONWebToken, error) {
		cookie, err := r.Cookie(key)
		if err != nil {
			return nil, auth0.ErrTokenNotFound
		}
		return jwt.ParseSigned(cookie.Value)
	}
}
