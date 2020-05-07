package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"melody/config"
	"melody/core"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	versions = map[string]uint16{
		"SSL3.0": tls.VersionSSL30,
		"TLS10":  tls.VersionTLS10,
		"TLS11":  tls.VersionTLS11,
		"TLS12":  tls.VersionTLS12,
	}
	defaultCurves = []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
	defaultCipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}
	errorPublicKey      = errors.New("public key not defined")
	errorPrivateKey     = errors.New("private key not defined")
	ErrorInternalError  = errors.New("internal server error")
	onceTransportConfig sync.Once
	// HeadersToSend 默认放行的请求头
	HeadersToSend = []string{"Content-Type"}
	// UserAgentHeaderValue 添加该请求头表示通过melody去代理了此次请求
	UserAgentHeaderValue = []string{core.MelodyUserAgent}
)

const (
	// HeaderCompleteResponseValue 响应完整时CompleteResponseHeader的值
	HeaderCompleteResponseValue = "true"
	// HeaderIncompleteResponseValue 响应不完整时CompleteResponseHeader的值
	HeaderIncompleteResponseValue = "false"
	HeaderCompleteKey             = "X-Melody-Complete"
)

// ToHTTPError 将错误转换为HTTP状态码
type ToHTTPError func(error) int

// DefaultToHTTPError 是一个错误码转换处理方法
// 它总是返回一个内部服务器错误500
func DefaultToHTTPError(_ error) int {
	return http.StatusInternalServerError
}

// InitHTTPDefaultTransport 确保在每次执行时只配置一次缺省HTTP传输
func InitHTTPDefaultTransport(cfg config.ServiceConfig) {
	onceTransportConfig.Do(func() {
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:       cfg.DialerTimeout,
				KeepAlive:     cfg.DialerKeepAlive,
				FallbackDelay: cfg.DialerFallbackDelay,
				DualStack:     true,
			}).DialContext,
			DisableCompression:    cfg.DisableCompression,
			DisableKeepAlives:     cfg.DisableKeepAlives,
			MaxIdleConns:          cfg.MaxIdleConns,
			MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:       cfg.IdleConnTimeout,
			ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
			ExpectContinueTimeout: cfg.ExpectContinueTimeout,
			TLSHandshakeTimeout:   10 * time.Second,
		}
	})
}

// RunServer 作为默认运行http.Server的函数实现
// 如果需要的话，将配置TLS层
func RunServer(ctx context.Context, cfg config.ServiceConfig, handler http.Handler) error {
	done := make(chan error)
	s := NewServer(cfg, handler)
	if cfg.TLS == nil {
		go func() {
			done <- s.ListenAndServe()
		}()
	} else {
		if cfg.TLS.PublicKey == "" {
			return errorPublicKey
		}
		if cfg.TLS.PrivateKey == "" {
			return errorPrivateKey
		}
		go func() {
			done <- s.ListenAndServeTLS(cfg.TLS.PublicKey, cfg.TLS.PrivateKey)
		}()
	}

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

// NewServer 返回一个默认的http.Server实例
func NewServer(cfg config.ServiceConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: cfg.ReaderHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		TLSConfig:         ParseTLSConfig(cfg.TLS),
	}
}

// ParseTLSConfig 解析TLS配置
func ParseTLSConfig(cfg *config.TLS) *tls.Config {
	if cfg == nil {
		return nil
	}

	if cfg.IsDisabled {
		return nil
	}

	return &tls.Config{
		MinVersion:               parseTLSVersion(cfg.MinVersion),
		MaxVersion:               parseTLSVersion(cfg.MaxVersion),
		CurvePreferences:         parseCurveIDs(cfg),
		CipherSuites:             parseChipherSuites(cfg),
		PreferServerCipherSuites: cfg.PreferServerCipherSuites,
	}
}

func parseTLSVersion(key string) uint16 {
	if v, ok := versions[key]; !ok {
		return v
	}
	// Default use tls version 12
	return tls.VersionTLS12
}

func parseCurveIDs(cfg *config.TLS) []tls.CurveID {
	l := len(cfg.CurvePreferences)
	if l == 0 {
		return defaultCurves
	}

	curves := make([]tls.CurveID, len(cfg.CurvePreferences))
	for i := range curves {
		curves[i] = tls.CurveID(cfg.CurvePreferences[i])
	}
	return curves
}

func parseChipherSuites(cfg *config.TLS) []uint16 {
	l := len(cfg.CipherSuites)
	if l == 0 {
		return defaultCipherSuites
	}

	cs := make([]uint16, l)
	for i := range cs {
		cs[i] = uint16(cfg.CipherSuites[i])
	}
	return cs
}
