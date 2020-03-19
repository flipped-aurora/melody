package martian

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"melody/config"
	"melody/logging"
	"melody/proxy"

	"github.com/google/martian/parse"
)

func TestHTTPRequestExecutor_ok(t *testing.T) {
	re := func(_ context.Context, req *http.Request) (resp *http.Response, err error) {
		if req.Header.Get("Content-Type") != "application/json" {
			err = fmt.Errorf("unexpected content type header: %s\n", req.Header.Get("Content-Type"))
			t.Error(err)
			return
		}
		if req.Header.Get("X-Neptunian") != "no!" {
			err = fmt.Errorf("unexpected X-Neptunian header: %s\n", req.Header.Get("X-Neptunian"))
			t.Error(err)
			return
		}
		if req.Header.Get("X-Martian-New") != "some value" {
			err = fmt.Errorf("unexpected X-Martian-New header: %s\n", req.Header.Get("X-Martian-New"))
			t.Error(err)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}

		if string(body) != `{"msg":"you rock!"}` {
			err = fmt.Errorf("unexpected request body: %s\n", string(body))
			t.Error(err)
			return
		}

		resp = &http.Response{
			Request: req,
			Header: http.Header{
				"X-Custom":     {"one"},
				"Content-Type": {"application/json"},
			},
			Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
			Status:     http.StatusText(200),
			StatusCode: 200,
		}
		return
	}

	r, err := parse.FromJSON([]byte(definition))
	if err != nil {
		t.Error(err)
		return
	}

	re = HTTPRequestExecutor(r, re)

	req, _ := http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	req.Header.Set("X-Neptunian", "no!")
	resp, err := re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Errorf("unexpected response: %v", *resp)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected response status: %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Custom") != "one" {
		t.Errorf("unexpected custom header: %s", resp.Header.Get("X-Custom"))
	}

	req, _ = http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	req.Header.Set("X-Neptunian", "no!")
	req.Header.Set("X-Martian-New", "some value")
	resp, err = re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Errorf("unexpected response: %v", *resp)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected response status: %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Custom") != "one" {
		t.Errorf("unexpected custom header: %s", resp.Header.Get("X-Custom"))
	}
}

func TestHTTPRequestExecutor_koEmptyResponse(t *testing.T) {
	r, err := parse.FromJSON([]byte(definition))
	if err != nil {
		t.Error(err)
		return
	}

	re := HTTPRequestExecutor(r, func(_ context.Context, _ *http.Request) (resp *http.Response, err error) { return })

	req, _ := http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err := re(context.Background(), req)
	if err != ErrEmptyResponse {
		t.Error(err)
	}
	if resp != nil {
		t.Errorf("unexpected response: %v", *resp)
	}
}

func TestHTTPRequestExecutor_koEmptyResponseBody(t *testing.T) {
	r, err := parse.FromJSON([]byte(definition))
	if err != nil {
		t.Error(err)
		return
	}

	re := HTTPRequestExecutor(r, func(_ context.Context, req *http.Request) (resp *http.Response, err error) {
		resp = &http.Response{
			Request:    req,
			StatusCode: 200,
		}
		return
	})

	req, _ := http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err := re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected response: %v", *resp)
	}
}

func TestHTTPRequestExecutor_koErroredResponse(t *testing.T) {
	r, err := parse.FromJSON([]byte(definition))
	if err != nil {
		t.Error(err)
		return
	}

	expectedErr := fmt.Errorf("some error")
	re := HTTPRequestExecutor(r, func(_ context.Context, _ *http.Request) (resp *http.Response, err error) { return nil, expectedErr })

	req, _ := http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err := re(context.Background(), req)
	if err != expectedErr {
		t.Error(err)
	}
	if resp != nil {
		t.Errorf("unexpected response: %v", *resp)
	}
}

func TestHTTPRequestExecutor_koErroredRequest(t *testing.T) {
	r, err := parse.FromJSON([]byte(definition))
	if err != nil {
		t.Error(err)
		return
	}

	expectedErr := fmt.Errorf("some error")
	re := HTTPRequestExecutor(r, func(_ context.Context, _ *http.Request) (resp *http.Response, err error) { return nil, expectedErr })

	req := &http.Request{
		Method:     "GET",
		RequestURI: "/url",
		Host:       "someHost",
	}
	resp, err := re(context.Background(), req)
	if err != expectedErr {
		t.Error(err)
	}
	if resp != nil {
		t.Errorf("unexpected response: %v", *resp)
	}
}

func TestNewBackendFactory_noExtra(t *testing.T) {
	expectedErr := fmt.Errorf("some error")
	re := func(_ context.Context, _ *http.Request) (resp *http.Response, err error) { return nil, expectedErr }
	buf := bytes.NewBuffer(make([]byte, 1024))
	l, _ := logging.NewLogger("DEBUG", buf, "")
	bf := NewBackendFactory(l, re)
	p := bf(&config.Backend{})
	resp, err := p(context.Background(), &proxy.Request{
		URL: &url.URL{
			Host: "example.com",
			Path: "/",
		},
		Body: ioutil.NopCloser(bytes.NewBufferString("")),
	})
	if resp != nil {
		t.Error("unexpected response:", resp)
	}
	if err != expectedErr {
		t.Error("unexpected error:", err)
	}
}

func TestNewBackendFactory_wrongExtra(t *testing.T) {
	expectedErr := fmt.Errorf("some error")
	re := func(_ context.Context, _ *http.Request) (resp *http.Response, err error) { return nil, expectedErr }
	buf := bytes.NewBuffer(make([]byte, 1024))
	l, _ := logging.NewLogger("DEBUG", buf, "")
	bf := NewBackendFactory(l, re)
	p := bf(&config.Backend{
		ExtraConfig: config.ExtraConfig{
			Namespace: 42,
		},
	})
	resp, err := p(context.Background(), &proxy.Request{
		URL: &url.URL{
			Host: "example.com",
			Path: "/",
		},
		Body: ioutil.NopCloser(bytes.NewBufferString("")),
	})
	if resp != nil {
		t.Error("unexpected response:", resp)
	}
	if err != expectedErr {
		t.Error("unexpected error:", err)
	}
}

func TestNewBackendFactory_ok(t *testing.T) {
	expectedErr := fmt.Errorf("some error")
	re := func(_ context.Context, _ *http.Request) (resp *http.Response, err error) { return nil, expectedErr }
	buf := bytes.NewBuffer(make([]byte, 1024))
	l, _ := logging.NewLogger("DEBUG", buf, "")
	bf := NewBackendFactory(l, re)
	p := bf(&config.Backend{
		ExtraConfig: config.ExtraConfig{
			Namespace: map[string]interface{}{
				"fifo.Group": map[string]interface{}{
					"scope":           []interface{}{"request", "response"},
					"aggregateErrors": true,
					"modifiers": []map[string]interface{}{
						{
							"header.Modifier": map[string]interface{}{
								"scope": []interface{}{"request", "response"},
								"name":  "X-Martian",
								"value": "ouh yeah!",
							},
						},
					},
				},
			},
		},
	})
	resp, err := p(context.Background(), &proxy.Request{
		URL: &url.URL{
			Host: "example.com",
			Path: "/",
		},
		Body: ioutil.NopCloser(bytes.NewBufferString("")),
	})
	if resp != nil {
		t.Error("unexpected response:", resp)
	}
	if err != expectedErr {
		t.Error("unexpected error:", err)
	}
}

func TestHTTPRequestExecutor_NoPanicWhenScopeLimitedToResponse(t *testing.T) {
	r, err := parse.FromJSON([]byte(responseOnlyDefinition))
	if err != nil {
		t.Error(err)
		return
	}

	re := HTTPRequestExecutor(r, func(_ context.Context, req *http.Request) (resp *http.Response, err error) {
		resp = &http.Response{
			Request:    req,
			StatusCode: 200,
		}
		return
	})

	req, _ := http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err := re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected response: %v", *resp)
	}
}

func TestHTTPRequestExecutor_NoPanicWhenScopeLimitedToRequest(t *testing.T) {
	r, err := parse.FromJSON([]byte(requestOnlyDefinition))
	if err != nil {
		t.Error(err)
		return
	}

	re := HTTPRequestExecutor(r, func(_ context.Context, req *http.Request) (resp *http.Response, err error) {
		resp = &http.Response{
			Request:    req,
			StatusCode: 200,
		}
		return
	})

	req, _ := http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err := re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected response: %v", *resp)
	}
}

func TestHTTPRequestExecutor_static(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_explicit_path_mapping_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	if err := os.MkdirAll(path.Join(tmpdir, "explicit/path"), 0777); err != nil {
		t.Fatalf("os.Mkdir(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "sfmtest.txt"), []byte("dont return"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	re := func(_ context.Context, req *http.Request) (resp *http.Response, err error) {
		t.Error("the request executor should not be called")
		return
	}

	msg := []byte(fmt.Sprintf(`{
		"static.Modifier": {
			"scope": ["request", "response"],
			"explicitPaths": {"/foo/bar.baz": "/subdir/sfmtest.txt"},
			"rootPath": %q
		}
	}`, tmpdir))

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Error(err)
		return
	}

	re = HTTPRequestExecutor(r, re)

	req, _ := http.NewRequest("GET", "/sfmtest.txt", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err := re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Errorf("unexpected response: %v", *resp)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected response status: %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("unexpected custom header: %s", resp.Header.Get("Content-Type"))
	}

	req, _ = http.NewRequest("GET", "url", ioutil.NopCloser(bytes.NewBufferString("")))
	resp, err = re(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Errorf("unexpected response: %v", *resp)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected response status: %d", resp.StatusCode)
	}
}

const definition = `{
    "fifo.Group": {
        "scope": ["request", "response"],
        "aggregateErrors": true,
        "modifiers": [
        {
            "header.Modifier": {
                "scope": ["request", "response"],
                "name" : "X-Martian",
                "value" : "ouh yeah!"
            }
        },
        {
            "body.Modifier": {
                "scope": ["request"],
                "contentType" : "application/json",
                "body" : "eyJtc2ciOiJ5b3Ugcm9jayEifQ=="
            }
        },
        {
            "header.RegexFilter": {
                "scope": ["request"],
                "header" : "X-Neptunian",
                "regex" : "no!",
                "modifier": {
                    "header.Modifier": {
                        "scope": ["request"],
                        "name" : "X-Martian-New",
                        "value" : "some value"
                    }
                }
            }
        }
        ]
    }
}`

const requestOnlyDefinition = `{
    "fifo.Group": {
        "scope": ["request"],
        "aggregateErrors": true,
        "modifiers": [
        {
            "header.Modifier": {
                "scope": ["request"],
                "name" : "X-Martian",
                "value" : "ouh yeah!"
            }
        },
        {
            "body.Modifier": {
                "scope": ["request"],
                "contentType" : "application/json",
                "body" : "eyJtc2ciOiJ5b3Ugcm9jayEifQ=="
            }
        },
        {
            "header.RegexFilter": {
                "scope": ["request"],
                "header" : "X-Neptunian",
                "regex" : "no!",
                "modifier": {
                    "header.Modifier": {
                        "scope": ["request"],
                        "name" : "X-Martian-New",
                        "value" : "some value"
                    }
                }
            }
        }
        ]
    }
}`

const responseOnlyDefinition = `{
    "fifo.Group": {
        "scope": ["response"],
        "aggregateErrors": true,
        "modifiers": [
        {
            "header.Modifier": {
                "scope": ["response"],
                "name" : "X-Martian",
                "value" : "ouh yeah!"
            }
        }
        ]
    }
}`
