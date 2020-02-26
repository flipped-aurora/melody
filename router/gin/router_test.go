// +build !race

package gin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"melody/config"
	"melody/logging"
	"melody/proxy"
	"melody/router"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestDefaultFactory_ok(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	logger, err := logging.NewLogger("DEBUG", buff, "pref")
	if err != nil {
		t.Error("building the logger:", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(5 * time.Millisecond)
	}()

	r := DefaultFactory(noopProxyFactory(map[string]interface{}{"supu": "tupu"}), logger).NewWithContext(ctx)

	serviceCfg := config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint: "/some",
				Method:   "GET",
				Timeout:  10,
				Backends: []*config.Backend{
					{},
				},
			},
			{
				Endpoint: "/some",
				Method:   "post",
				Timeout:  10,
				Backends: []*config.Backend{
					{},
				},
			},
			{
				Endpoint: "/some",
				Method:   "put",
				Timeout:  10,
				Backends: []*config.Backend{
					{},
				},
			},
			{
				Endpoint: "/some",
				Method:   "PATCH",
				Timeout:  10,
				Backends: []*config.Backend{
					{},
				},
			},
			{
				Endpoint: "/some",
				Method:   "DELETE",
				Timeout:  10,
				Backends: []*config.Backend{
					{},
				},
			},
		},
	}

	go func() { r.Run(serviceCfg) }()

	time.Sleep(5 * time.Millisecond)

	for _, endpoint := range serviceCfg.Endpoints {
		req, _ := http.NewRequest(strings.ToTitle(endpoint.Method), fmt.Sprintf("http://127.0.0.1:8080%s", endpoint.Endpoint), nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error("Making the request:", err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.Header.Get("Cache-Control") != "" {
			t.Error("Cache-Control error:", resp.Header.Get("Cache-Control"))
		}
		if resp.Header.Get(router.HeaderCompleteKey) != router.HeaderCompleteResponseValue {
			t.Error(router.HeaderCompleteKey, "error:", resp.Header.Get(router.HeaderCompleteKey))
		}
		if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Error("Content-Type error:", resp.Header.Get("Content-Type"))
		}
		if resp.Header.Get("X-Melody") != "Version 0.0.1" {
			t.Error("X-Melody error:", resp.Header.Get("X-Melody"))
		}
		if resp.StatusCode != http.StatusOK {
			t.Error("Unexpected status code:", resp.StatusCode)
		}
	}
}

func TestDefaultFactory_ko(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	logger, err := logging.NewLogger("ERROR", buff, "pref")
	if err != nil {
		t.Error("building the logger:", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(5 * time.Millisecond)
	}()

	r := DefaultFactory(noopProxyFactory(map[string]interface{}{"supu": "tupu"}), logger).NewWithContext(ctx)

	serviceCfg := config.ServiceConfig{
		Debug: true,
		Port:  8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint: "/ignored",
				Method:   "GETTT",
				Backends: []*config.Backend{
					{},
				},
			},
			{
				Endpoint: "/empty",
				Method:   "GETTT",
				Backends:  []*config.Backend{},
			},
			{
				Endpoint: "/also-ignored",
				Method:   "PUTt",
				Backends: []*config.Backend{
					{},
					{},
				},
			},
		},
	}

	go func() { r.Run(serviceCfg) }()

	time.Sleep(5 * time.Millisecond)

	for _, subject := range [][]string{
		{"GET", "ignored"},
		{"GET", "empty"},
		{"PUT", "also-ignored"},
	} {
		req, _ := http.NewRequest(subject[0], fmt.Sprintf("http://127.0.0.1:8080/%s", subject[1]), nil)
		req.Header.Set("Content-Type", "application/json")
		checkResponseIs404(t, req)
	}
}

func TestDefaultFactory_proxyFactoryCrash(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	logger, err := logging.NewLogger("ERROR", buff, "pref")
	if err != nil {
		t.Error("building the logger:", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(5 * time.Millisecond)
	}()

	r := DefaultFactory(erroredProxyFactory{fmt.Errorf("%s", "crash!!!")}, logger).NewWithContext(ctx)

	serviceCfg := config.ServiceConfig{
		Debug: true,
		Port:  8074,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint: "/ignored",
				Method:   "GET",
				Timeout:  10,
				Backends: []*config.Backend{
					{},
				},
			},
		},
	}

	go func() { r.Run(serviceCfg) }()

	time.Sleep(5 * time.Millisecond)

	for _, subject := range [][]string{{"GET", "ignored"}, {"PUT", "also-ignored"}} {
		req, _ := http.NewRequest(subject[0], fmt.Sprintf("http://127.0.0.1:8074/%s", subject[1]), nil)
		req.Header.Set("Content-Type", "application/json")
		checkResponseIs404(t, req)
	}
}

func TestRunServer_ko(t *testing.T) {
	buff := new(bytes.Buffer)
	logger, err := logging.NewLogger("ERROR", buff, "")
	if err != nil {
		t.Error("building the logger:", err.Error())
		return
	}

	errorMsg := "runServer error"
	runServerFunc := func(_ context.Context, _ config.ServiceConfig, _ http.Handler) error {
		return errors.New(errorMsg)
	}

	pf := noopProxyFactory(map[string]interface{}{"supu": "tupu"})
	r := NewFactory(
		Config{
			Engine:         gin.Default(),
			MiddleWares:    []gin.HandlerFunc{},
			HandlerFactory: EndpointHandler,
			ProxyFactory:   pf,
			Logger:         logger,
			RunServer:      runServerFunc,
		},
	).New()

	serviceCfg := config.ServiceConfig{}
	r.Run(serviceCfg)
	re := regexp.MustCompile(errorMsg)
	if !re.MatchString(string(buff.Bytes())) {
		t.Errorf("the logger doesn't contain the expected msg: %s", buff.Bytes())
	}
}

func checkResponseIs404(t *testing.T, req *http.Request) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Making the request:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.Header.Get("Cache-Control") != "" {
		t.Error(req.URL.String(), "Cache-Control error:", resp.Header.Get("Cache-Control"))
	}
	if resp.Header.Get(router.HeaderCompleteKey) != router.HeaderInCompleteResponseValue {
		t.Error(req.URL.String(), router.HeaderCompleteKey, "error:", resp.Header.Get(router.HeaderCompleteKey))
	}
	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Error(req.URL.String(), "Content-Type error:", resp.Header.Get("Content-Type"))
	}
	if resp.Header.Get("X-Melody") != "" {
		t.Error(req.URL.String(), "X-Melody error:", resp.Header.Get("X-Melody"))
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error(req.URL.String(), "Unexpected status code:", resp.StatusCode)
	}
}

type noopProxyFactory map[string]interface{}

func (n noopProxyFactory) New(_ *config.EndpointConfig) (proxy.Proxy, error) {
	return func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return &proxy.Response{
			IsComplete: true,
			Data:       n,
		}, nil
	}, nil
}

type erroredProxyFactory struct {
	Error error
}

func (e erroredProxyFactory) New(_ *config.EndpointConfig) (proxy.Proxy, error) {
	return proxy.NoopProxy, e.Error
}
