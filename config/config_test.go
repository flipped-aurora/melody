package config

import (
	"fmt"
	"testing"
	"time"
)

func TestConfig_uniqueOutput_result(t *testing.T) {
	params := []string{"id", "name", "name"}

	outputParams, outputParamsSize := uniqueOutput(params)

	if outputParamsSize != 1 {
		t.Errorf("expected 1, actually :%d", outputParamsSize)
	}

	fmt.Println(outputParams, outputParamsSize)
}

func TestConfig_initBackendURLMappings_ok(t *testing.T) {
	samples := []string{
		"supu/{tupu}",
		"/supu/{tupu1}",
		"/supu.local/",
		"supu/{tupu_56}/{supu-5t6}?a={foo}&b={foo}",
		"supu/{tupu_56}{supu-5t6}?a={foo}&b={foo}",
		"supu/tupu{supu-5t6}?a={foo}&b={foo}",
		"{resp0_x}/{tupu1}/{tupu_56}{supu-5t6}?a={tupu}&b={foo}",
	}

	expected := []string{
		"/supu/{{.Tupu}}",
		"/supu/{{.Tupu1}}",
		"/supu.local/",
		"/supu/{{.Tupu_56}}/{{.Supu-5t6}}?a={{.Foo}}&b={{.Foo}}",
		"/supu/{{.Tupu_56}}{{.Supu-5t6}}?a={{.Foo}}&b={{.Foo}}",
		"/supu/tupu{{.Supu-5t6}}?a={{.Foo}}&b={{.Foo}}",
		"/{{.Resp0_x}}/{{.Tupu1}}/{{.Tupu_56}}{{.Supu-5t6}}?a={{.Tupu}}&b={{.Foo}}",
	}

	backend := Backend{}
	endpoint := EndpointConfig{Backends: []*Backend{&backend}}
	subject := ServiceConfig{Endpoints: []*EndpointConfig{&endpoint}, uriParser: NewURIParser()}

	inputSet := map[string]interface{}{
		"tupu":     nil,
		"tupu1":    nil,
		"tupu_56":  nil,
		"supu-5t6": nil,
		"foo":      nil,
	}

	for i := range samples {
		backend.URLPattern = samples[i]
		if err := subject.initBackendsURLMappings(0, 0, inputSet); err != nil {
			t.Error(err)
		}
		if backend.URLPattern != expected[i] {
			t.Errorf("want: %s, have: %s\n", expected[i], backend.URLPattern)
		}
	}
}

func TestConfig_initBackendURLMappings_tooManyOutput(t *testing.T) {
	backend := Backend{URLPattern: "supu/{tupu_56}/{supu-5t6}?a={foo}&b={foo}"}
	endpoint := EndpointConfig{
		Method:   "GET",
		Endpoint: "/some/{tupu}",
		Backends: []*Backend{&backend},
	}
	subject := ServiceConfig{Endpoints: []*EndpointConfig{&endpoint}, uriParser: NewURIParser()}

	inputSet := map[string]interface{}{
		"tupu": nil,
	}

	expectedErrMsg := "input and output params do not match. endpoint: GET /some/{tupu}, backend: 0. input: [tupu], output: [foo supu-5t6 tupu_56]"

	err := subject.initBackendsURLMappings(0, 0, inputSet)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestConfig_initBackendURLMappings_undefinedOutput(t *testing.T) {
	backend := Backend{URLPattern: "supu/{tupu_56}/{supu-5t6}?a={foo}&b={foo}"}
	endpoint := EndpointConfig{Endpoint: "/", Method: "GET", Backends: []*Backend{&backend}}
	subject := ServiceConfig{Endpoints: []*EndpointConfig{&endpoint}, uriParser: NewURIParser()}

	inputSet := map[string]interface{}{
		"tupu": nil,
		"supu": nil,
		"foo":  nil,
	}

	expectedErrMsg := "Undefined output param 'supu-5t6'! endpoint: GET /, backend: 0. input: [foo supu tupu], output: [foo supu-5t6 tupu_56]"
	err := subject.initBackendsURLMappings(0, 0, inputSet)
	if err == nil || err.Error() != expectedErrMsg {
		t.Errorf("error expected. have: %v", err)
	}
}

func TestConfig_init(t *testing.T) {
	supuBackend := Backend{
		URLPattern: "/__debug/supu",
	}
	supuEndpoint := EndpointConfig{
		Endpoint:       "/supu",
		Method:         "post",
		Timeout:        1500 * time.Millisecond,
		CacheTTL:       6 * time.Hour,
		Backends:       []*Backend{&supuBackend},
		OutputEncoding: "some_render",
	}

	githubBackend := Backend{
		URLPattern: "/",
		Host:       []string{"https://api.github.com"},
		Whitelist:  []string{"authorizations_url", "code_search_url"},
	}
	githubEndpoint := EndpointConfig{
		Endpoint: "/github",
		Timeout:  1500 * time.Millisecond,
		CacheTTL: 6 * time.Hour,
		Backends: []*Backend{&githubBackend},
	}

	userBackend := Backend{
		URLPattern: "/users/{user}",
		Host:       []string{"https://jsonplaceholder.typicode.com"},
		Mapping:    map[string]string{"email": "personal_email"},
	}
	rssBackend := Backend{
		URLPattern: "/users/{user}",
		Host:       []string{"https://jsonplaceholder.typicode.com"},
		Encoding:   "rss",
	}
	postBackend := Backend{
		URLPattern: "/posts/{user}",
		Host:       []string{"https://jsonplaceholder.typicode.com"},
		Group:      "posts",
		Encoding:   "xml",
	}
	userEndpoint := EndpointConfig{
		Endpoint: "/users/{user}",
		Backends: []*Backend{&userBackend, &rssBackend, &postBackend},
	}

	subject := ServiceConfig{
		Version:   1,
		Timeout:   5 * time.Second,
		CacheTTL:  30 * time.Minute,
		Host:      []string{"http://127.0.0.1:8080"},
		Endpoints: []*EndpointConfig{&supuEndpoint, &githubEndpoint, &userEndpoint},
	}

	if err := subject.Init(); err != nil {
		t.Error("Error at the configuration init:", err.Error())
	}

	if len(supuBackend.Host) != 1 || supuBackend.Host[0] != subject.Host[0] {
		t.Error("Default hosts not applied to the supu backend", supuBackend.Host)
	}

	for level, method := range map[string]string{
		"userBackend":  userBackend.Method,
		"postBackend":  postBackend.Method,
		"userEndpoint": userEndpoint.Method,
	} {
		if method != "GET" {
			t.Errorf("Default method not applied at %s. Get: %s", level, method)
		}
	}

	if supuBackend.Method != "post" {
		t.Error("unexpected supuBackend")
	}

	if userBackend.Timeout != subject.Timeout {
		t.Error("default timeout not applied to the userBackend")
	}

	if userEndpoint.CacheTTL != subject.CacheTTL {
		t.Error("default CacheTTL not applied to the userEndpoint")
	}
}

func TestConfig_init_NoBackends(t *testing.T) {

	supuEndpoint := EndpointConfig{
		Endpoint:       "/supu",
		Method:         "post",
		Timeout:        1500 * time.Millisecond,
		CacheTTL:       6 * time.Hour,
		Backends:       []*Backend{},
		OutputEncoding: "some_render",
	}

	githubEndpoint := EndpointConfig{
		Endpoint: "/github",
		Timeout:  1500 * time.Millisecond,
		CacheTTL: 6 * time.Hour,
		Backends: []*Backend{},
	}

	userEndpoint := EndpointConfig{
		Endpoint: "/users/{user}",
		Backends: []*Backend{},
	}

	subject := ServiceConfig{
		Version:   1,
		Timeout:   5 * time.Second,
		CacheTTL:  30 * time.Minute,
		Host:      []string{"http://127.0.0.1:8080"},
		Endpoints: []*EndpointConfig{&supuEndpoint, &githubEndpoint, &userEndpoint},
	}

	if err := subject.Init(); err != nil {
		if _, ok := err.(*NoBackendsError); !ok {
			t.Errorf("expected: %v, have :%v", NoBackendsError{}, err)
		}
	}
}
