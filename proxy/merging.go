package proxy

import (
	"context"
	"melody/config"
	"strings"
	"time"
)

const (
	mergeKey            = "combiner"
	isSequentialKey     = "sequential"
	defaultCombinerName = "default"
)

var responseCombiners = initResponseCombiners()



// ResponseCombiner 将多个proxy.Response 整合成一个
type ResponseCombiner func(int, []*Response) *Response

type incrementalMergeAccumulator struct {
	pending  int
	data     *Response
	combiner ResponseCombiner
	errs     []error
}

func initResponseCombiners() *combinerRegister {
	return newCombinerRegister(map[string]ResponseCombiner{defaultCombinerName: combineData}, combineData)
}

// NewMergeDataMiddleware 为多个backends的endpoint包裹一层middleware去合并response
func NewMergeDataMiddleware(config *config.EndpointConfig) Middleware {
	totalBackends := len(config.Backends)
	if totalBackends == 0 {
		panic(ErrNoBackends)
	}
	if totalBackends == 1 {
		return EmptyMiddleware
	}

	serviceTimeOut := time.Duration(85*config.Timeout.Nanoseconds()/100) * time.Nanosecond
	combiner := getResponseCombiner(config.ExtraConfig)
	return func(proxy ...Proxy) Proxy {
		if len(proxy) != totalBackends {
			panic(ErrNotEnoughProxies)
		}

		if !shouldRunSequentialMerger(config.ExtraConfig) {
			// 并行合并请求
			return parallelMerge(serviceTimeOut, combiner, proxy...)
		}
		// TODO **链式合并请求
		return nil
	}
}

func parallelMerge(timeout time.Duration, rc ResponseCombiner, next ...Proxy) Proxy {
	return func(ctx context.Context, request *Request) (response *Response, e error) {
		localCtx, cancel := context.WithTimeout(ctx, timeout)
		responses := make(chan *Response, len(next))
		failed := make(chan error, len(next))
		for _, v := range next {
			go requestPart(localCtx, v, request, responses, failed)
		}

		acc := newIncrementalMergeAccumulator(len(next), rc)
		for i := 0 ; i < len(next) ; i++ {
			select {
			case resp := <-responses:
				acc.Merge(resp, nil)
			case err := <-failed:
				acc.Merge(nil, err)
			}
		}
		res, err := acc.Result()
		cancel()

		return res, err
	}
}

func (i *incrementalMergeAccumulator) Merge(res *Response, err error) {
	i.pending--
	if err != nil {
		i.errs = append(i.errs, err)
		if i.data != nil {
			i.data.IsComplete = false
		}
		return
	}
	if res == nil {
		i.errs = append(i.errs, errNullResult)
		return
	}

	if i.data == nil {
		i.data = res
		return
	}

	i.data = i.combiner(2, []*Response{i.data, res})
}

func (i *incrementalMergeAccumulator) Result() (*Response, error) {
	if i.data == nil {
		return &Response{
			Data: map[string]interface{}{},
			IsComplete: false,
		}, newMergeError(i.errs)
	}

	if i.pending != 0 || len(i.errs) != 0 {
		i.data.IsComplete = false
	}

	return i.data, newMergeError(i.errs)
}

func newIncrementalMergeAccumulator(total int, combiner ResponseCombiner) *incrementalMergeAccumulator {
	return &incrementalMergeAccumulator{
		pending:  total,
		combiner: combiner,
		errs:     []error{},
	}
}

// 是否开启链式请求
func shouldRunSequentialMerger(config config.ExtraConfig) bool {
	if v, ok := config[Namespace]; ok {
		if temp, ok := v.(map[string]interface{}); ok {
			if str, ok := temp[isSequentialKey]; ok {
				r , ok := str.(bool)
				return r && ok
			}
		}
	}
	return false
}

func getResponseCombiner(extra config.ExtraConfig) ResponseCombiner {
	combiner, _ := responseCombiners.GetResponseCombiner(defaultCombinerName)
	if v, ok := extra[Namespace]; ok {
		if temp, ok := v.(map[string]interface{}); ok {
			if s, ok := temp[mergeKey]; ok {
				if c, ok := responseCombiners.GetResponseCombiner(s.(string)); ok {
					combiner = c
				}
			}
		}
	}
	return combiner
}

func combineData(count int, responses []*Response) *Response {
	isComplete := len(responses) == count
	var resp *Response
	for _, v := range responses {
		if v == nil || v.Data == nil {
			isComplete = false
			continue
		}
		isComplete = isComplete && v.IsComplete
		if resp == nil {
			resp = v
			continue
		}
		for k, v := range v.Data {
			resp.Data[k] = v
		}
	}
	if nil == resp {
		return  &Response{
			Data:       map[string]interface{}{},
			IsComplete: isComplete,
		}
	}
	resp.IsComplete = isComplete
	return resp
}

//TODO 实现自定义combineData，可以根据mergeKey读取


func requestPart(ctx context.Context, next Proxy, request *Request, out chan<- *Response, failed chan<- error) {
	localCtx, cancel := context.WithCancel(ctx)

	resp, err := next(localCtx, request)

	if err != nil {
		failed <- err
		cancel()
		return
	}

	if resp == nil {
		failed <- errNullResult
		cancel()
		return
	}

	select {
	case out <- resp:
	case <- ctx.Done():
		failed <- ctx.Err()
		cancel()
	}

	cancel()
}


func newMergeError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return mergeError{errs}
}

type mergeError struct {
	errs []error
}

func (m mergeError) Error() string {
	msg := make([]string, len(m.errs))
	for i, err := range m.errs {
		msg[i] = err.Error()
	}
	return strings.Join(msg, "\n")
}
