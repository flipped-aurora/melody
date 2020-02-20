package gin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInvalidCfg(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	corsMw := New(sampleCfg)
	if corsMw != nil {
		t.Error("The corsMw should be nil.\n")
	}
}

func TestNew(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "melody_cors": {
			"allow_origins": [ "http://foobar.com" ],
			"allow_methods": [ "GET" ],
			"max_age": "2h"
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	e := gin.Default()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/foo", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "http://foobar.com",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "Origin",
		"Access-Control-Max-Age":       "7200",
	})

}

func TestAllowOriginWildcard(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "melody_cors": {
			"allow_origins": [ "*" ]
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	e := gin.Default()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/wildcard", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/wildcard", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "Origin",
	})
}

func TestAllowOriginEmpty(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "melody_cors": {
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	e := gin.Default()
	corsMw := New(sampleCfg)
	if corsMw == nil {
		t.Error("The cors middleware should not be nil.\n")
	}
	e.Use(corsMw)
	e.GET("/foo", func(c *gin.Context) { c.String(200, "Yeah") })
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://foobar.com")
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "origin")
	e.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Invalid status code: %d should be 200", res.Code)
	}

	assertHeaders(t, res.Header(), map[string]string{
		"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
		"Access-Control-Allow-Headers": "Origin",
	})
}

var allHeaders = []string{
	"Vary",
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Credentials",
	"Access-Control-Max-Age",
	"Access-Control-Expose-Headers",
}

func assertHeaders(t *testing.T, resHeaders http.Header, expHeaders map[string]string) {
	for _, name := range allHeaders {
		got := strings.Join(resHeaders[name], ", ")
		want := expHeaders[name]
		if got != want {
			t.Errorf("Response header %q = %q, want %q", name, got, want)
		}
	}
}
