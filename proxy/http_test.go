package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"melody/config"
	"melody/encoding"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewHTTPProxy_ok(t *testing.T) {
	expectedMethod := "GET"
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 11 {
			t.Errorf("unexpected request size. Want: 11. Have: %d", r.ContentLength)
		}
		if h := r.Header.Get("Content-Length"); h != "11" {
			t.Errorf("unexpected content-length header. Want: 11. Have: %s", h)
		}
		if r.Method != expectedMethod {
			t.Errorf("Wrong request method. Want: %s. Have: %s", expectedMethod, r.Method)
		}
		if h := r.Header.Get("X-First"); h != "first" {
			t.Errorf("unexpected first header: %s", h)
		}
		if h := r.Header.Get("X-Second"); h != "second" {
			t.Errorf("unexpected second header: %s", h)
		}
		r.Header.Del("X-Second")
		fmt.Fprintf(w, "{\"supu\":42, \"tupu\":true, \"foo\": \"bar\"}")
	}))
	defer backendServer.Close()

	rpURL, _ := url.Parse(backendServer.URL)
	backend := config.Backend{
		Decoder: encoding.JSONDecoder(),
	}
	request := Request{
		Method: expectedMethod,
		Path:   "/",
		URL:    rpURL,
		Body:   newDummyReadCloser(`{"abc": 42}`),
		Headers: map[string][]string{
			"X-First":        {"first"},
			"X-Second":       {"second"},
			"Content-Length": {"11"},
		},
	}
	mustEnd := time.After(time.Duration(150) * time.Millisecond)

	result, err := HTTPProxyFactory(http.DefaultClient)(&backend)(context.Background(), &request)
	if err != nil {
		t.Errorf("The proxy returned an unexpected error: %s\n", err.Error())
		return
	}
	if result == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("Error: expected response")
		return
	default:
	}

	tmp, ok := result.Data["supu"]
	if !ok {
		t.Errorf("The proxy returned an unexpected result: %v\n", result)
	}
	supuValue, err := tmp.(json.Number).Int64()
	if err != nil || supuValue != 42 {
		t.Errorf("The proxy returned an unexpected result: %v\n", supuValue)
	}
	if v, ok := result.Data["tupu"]; !ok || !v.(bool) {
		t.Errorf("The proxy returned an unexpected result: %v\n", result)
	}
	if v, ok := result.Data["foo"]; !ok || v.(string) != "bar" {
		t.Errorf("The proxy returned an unexpected result: %v\n", result)
	}
	if v, ok := request.Headers["X-Second"]; !ok || len(v) != 1 {
		t.Errorf("the proxy request headers were changed: %v", request.Headers)
	}
}
