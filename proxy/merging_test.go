package proxy

import (
	"context"
	"melody/config"
	"testing"
	"time"
)

func TestNewMergeDataMiddleware_ok(t *testing.T) {
	timeout := 500
	backend := config.Backend{}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend, &backend},
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	mw := NewMergeDataMiddleware(&endpoint)
	p := mw(
		dummyProxy(&Response{Data: map[string]interface{}{"supu": 42}, IsComplete: true}),
		dummyProxy(&Response{Data: map[string]interface{}{"tupu": true}, IsComplete: true}))
	mustEnd := time.After(time.Duration(2*timeout) * time.Millisecond)
	out, err := p(context.Background(), &Request{})
	if err != nil {
		t.Errorf("The middleware propagated an unexpected error: %s\n", err.Error())
	}
	if out == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("We were expecting a response but we got none\n")
	default:
		if len(out.Data) != 2 {
			t.Errorf("We weren't expecting a partial response but we got %v!\n", out)
		}
		if !out.IsComplete {
			t.Errorf("We were expecting a completed response but we got an incompleted one!\n")
		}
	}
}

func TestNewMergeDataMiddleware_partialTimeout(t *testing.T) {
	timeout := 100
	backend := config.Backend{Timeout: time.Duration(timeout) * time.Millisecond}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend, &backend},
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	mw := NewMergeDataMiddleware(&endpoint)
	p := mw(
		delayedProxy(t, time.Duration(timeout/2)*time.Millisecond, &Response{Data: map[string]interface{}{"supu": 42}, IsComplete: true}),
		delayedProxy(t, time.Duration(5*timeout)*time.Millisecond, nil))
	mustEnd := time.After(time.Duration(2*timeout) * time.Millisecond)
	out, err := p(context.Background(), &Request{})
	if err == nil || err.Error() != "context deadline exceeded" {
		t.Errorf("The middleware propagated an unexpected error: %s\n", err.Error())
	}
	if out == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("We were expecting a response but we got none\n")
	default:
		if len(out.Data) != 1 {
			t.Errorf("We were expecting a partial response but we got %v!\n", out)
		}
		if out.IsComplete {
			t.Errorf("We were expecting an incompleted response but we got a completed one!\n")
		}
	}
}

func TestNewMergeDataMiddleware_partial(t *testing.T) {
	timeout := 100
	backend := config.Backend{Timeout: time.Duration(timeout) * time.Millisecond}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend, &backend},
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	mw := NewMergeDataMiddleware(&endpoint)
	p := mw(
		dummyProxy(&Response{Data: map[string]interface{}{"supu": 42}, IsComplete: true}),
		dummyProxy(&Response{}))
	mustEnd := time.After(time.Duration(2*timeout) * time.Millisecond)
	out, err := p(context.Background(), &Request{})
	if err != nil {
		t.Errorf("The middleware propagated an unexpected error: %s\n", err.Error())
	}
	if out == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("We were expecting a response but we got none\n")
	default:
		if len(out.Data) != 1 {
			t.Errorf("We were expecting a partial response but we got %v!\n", out)
		}
		if out.IsComplete {
			t.Errorf("We were expecting an incompleted response but we got a completed one!\n")
		}
	}
}

func TestNewMergeDataMiddleware_nullResponse(t *testing.T) {
	timeout := 100
	backend := config.Backend{Timeout: time.Duration(timeout) * time.Millisecond}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend, &backend},
	}
	mw := NewMergeDataMiddleware(&endpoint)

	mustEnd := time.After(time.Duration(2*timeout) * time.Millisecond)
	out, err := mw(NoopProxy, NoopProxy)(context.Background(), &Request{})
	if err == nil {
		t.Errorf("The middleware did not propagate the expected error")
	}
	switch mergeErr := err.(type) {
	case mergeError:
		if len(mergeErr.errs) != 2 {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0] != mergeErr.errs[1] {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0] != errNullResult {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
	default:
		t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
	}
	if out == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("We were expecting a response but we got none\n")
	default:
		if len(out.Data) != 0 {
			t.Errorf("We were expecting a partial response but we got %v!\n", out.Data)
		}
		if out.IsComplete {
			t.Errorf("We were expecting an incompleted response but we got a completed one!\n")
		}
	}
}


func TestNewMergeDataMiddleware_timeout(t *testing.T) {
	timeout := 100
	backend := config.Backend{Timeout: time.Duration(timeout) * time.Millisecond}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend, &backend},
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	mw := NewMergeDataMiddleware(&endpoint)
	p := mw(
		delayedProxy(t, time.Duration(5*timeout)*time.Millisecond, nil),
		delayedProxy(t, time.Duration(5*timeout)*time.Millisecond, nil))
	mustEnd := time.After(time.Duration(2*timeout) * time.Millisecond)
	out, err := p(context.Background(), &Request{})
	if err == nil {
		t.Errorf("The middleware did not propagate the expected error")
	}
	switch mergeErr := err.(type) {
	case mergeError:
		if len(mergeErr.errs) != 2 {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0].Error() != mergeErr.errs[1].Error() {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0].Error() != "context deadline exceeded" {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
	default:
		t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
	}
	if out == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("We were expecting a response but we got none\n")
	default:
		if len(out.Data) > 0 {
			t.Errorf("We weren't expecting a response but we got one!\n")
		}
		if out.IsComplete {
			t.Errorf("We were expecting an incompleted response but we got a completed one!\n")
		}
	}
}

func TestNewMergeDataMiddleware_notEnoughBackends(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic\n")
		}
	}()
	backend := config.Backend{}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend},
	}
	mw := NewMergeDataMiddleware(&endpoint)
	mw(explosiveProxy(t), explosiveProxy(t))
}

func TestNewMergeDataMiddleware_notEnoughProxies(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic\n")
		}
	}()
	backend := config.Backend{}
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{&backend, &backend},
	}
	mw := NewMergeDataMiddleware(&endpoint)
	mw(NoopProxy)
}

func TestNewMergeDataMiddleware_noBackends(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic\n")
		}
	}()
	endpoint := config.EndpointConfig{}
	NewMergeDataMiddleware(&endpoint)
}

func Test_incrementalMergeAccumulator_invalidResponse(t *testing.T) {
	acc := newIncrementalMergeAccumulator(3, combineData)
	acc.Merge(nil, nil)
	acc.Merge(nil, nil)
	acc.Merge(nil, nil)
	res, err := acc.Result()
	if res == nil {
		t.Error("response should not be nil")
		return
	}
	if err == nil {
		t.Error("expecting error")
		return
	}
	switch mergeErr := err.(type) {
	case mergeError:
		if len(mergeErr.errs) != 3 {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0] != mergeErr.errs[1] {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0] != mergeErr.errs[2] {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
		if mergeErr.errs[0] != errNullResult {
			t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
		}
	default:
		t.Errorf("The middleware propagated an unexpected error: %s", err.Error())
	}
}

func Test_incrementalMergeAccumulator_incompleteResponse(t *testing.T) {
	acc := newIncrementalMergeAccumulator(3, combineData)
	acc.Merge(&Response{Data: make(map[string]interface{}, 0), IsComplete: true}, nil)
	acc.Merge(&Response{Data: make(map[string]interface{}, 0), IsComplete: false}, nil)
	acc.Merge(&Response{Data: make(map[string]interface{}, 0), IsComplete: true}, nil)
	res, err := acc.Result()
	if res == nil {
		t.Error("response should not be nil")
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}
	if res.IsComplete {
		t.Error("response should not be completed")
	}
}

func checkRequestParam(t *testing.T, r *Request, k, v string) {
	if r.Params[k] != v {
		t.Errorf("request without the expected set of params: %s - %+v", k, r.Params)
	}
}

func TestNewMergeDataMiddleware_sequential(t *testing.T) {
	timeout := 1000
	endpoint := config.EndpointConfig{
		Backends: []*config.Backend{
			{URLPattern: "/"},
			{URLPattern: "/aaa/{{.Resp0_int}}/{{.Resp0_string}}/{{.Resp0_bool}}/{{.Resp0_float}}/{{.Resp0_struct.foo}}"},
			{URLPattern: "/aaa/{{.Resp0_int}}/{{.Resp0_string}}/{{.Resp0_bool}}/{{.Resp0_float}}/{{.Resp0_struct.foo}}?x={{.Resp1_tupu}}"},
			{URLPattern: "/aaa/{{.Resp0_struct.foo}}/{{.Resp0_struct.struct.foo}}/{{.Resp0_struct.struct.struct.foo}}"},
		},
		Timeout: time.Duration(timeout) * time.Millisecond,
		ExtraConfig: config.ExtraConfig{
			Namespace: map[string]interface{}{
				isSequentialKey: true,
			},
		},
	}
	mw := NewMergeDataMiddleware(&endpoint)
	p := mw(
		dummyProxy(&Response{Data: map[string]interface{}{
			"int":    42,
			"string": "some",
			"bool":   true,
			"float":  3.14,
			"struct": map[string]interface{}{
				"foo": "bar",
				"struct": map[string]interface{}{
					"foo": "bar",
					"struct": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
		}, IsComplete: true}),
		func(ctx context.Context, r *Request) (*Response, error) {
			checkRequestParam(t, r, "Resp0_int", "42")
			checkRequestParam(t, r, "Resp0_string", "some")
			checkRequestParam(t, r, "Resp0_float", "3.14E+00")
			checkRequestParam(t, r, "Resp0_bool", "true")
			checkRequestParam(t, r, "Resp0_struct.foo", "bar")
			return &Response{Data: map[string]interface{}{"tupu": "foo"}, IsComplete: true}, nil
		},
		func(ctx context.Context, r *Request) (*Response, error) {
			checkRequestParam(t, r, "Resp0_int", "42")
			checkRequestParam(t, r, "Resp0_string", "some")
			checkRequestParam(t, r, "Resp0_float", "3.14E+00")
			checkRequestParam(t, r, "Resp0_bool", "true")
			checkRequestParam(t, r, "Resp0_struct.foo", "bar")
			checkRequestParam(t, r, "Resp1_tupu", "foo")
			return &Response{Data: map[string]interface{}{"aaaa": []int{1, 2, 3}}, IsComplete: true}, nil
		},
		func(ctx context.Context, r *Request) (*Response, error) {
			checkRequestParam(t, r, "Resp0_struct.foo", "bar")
			checkRequestParam(t, r, "Resp0_struct.struct.foo", "bar")
			checkRequestParam(t, r, "Resp0_struct.struct.struct.foo", "bar")
			return &Response{Data: map[string]interface{}{"bbbb": []bool{true, false}}, IsComplete: true}, nil
		},
	)
	mustEnd := time.After(time.Duration(2*timeout) * time.Millisecond)
	out, err := p(context.Background(), &Request{Params: map[string]string{}})
	if err != nil {
		t.Errorf("The middleware propagated an unexpected error: %s\n", err.Error())
	}
	if out == nil {
		t.Errorf("The proxy returned a null result\n")
		return
	}
	select {
	case <-mustEnd:
		t.Errorf("We were expecting a response but we got none\n")
	default:
		if len(out.Data) != 8 {
			t.Errorf("We weren't expecting a partial response but we got %v!\n", out)
		}
		if !out.IsComplete {
			t.Errorf("We were expecting a completed response but we got an incompleted one!\n")
		}
	}
}


