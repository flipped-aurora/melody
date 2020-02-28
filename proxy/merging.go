package proxy

import (
	"context"
	"fmt"
	"melody/config"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	mergeKey            = "combiner"
	isSequentialKey     = "sequential"
	defaultCombinerName = "default"
)

var responseCombiners = initResponseCombiners()

var sequentialLastParamKeyRegexp = regexp.MustCompile(`\{\{\.Resp(\d+)_([\d\w-_\.]+)\}\}`)

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
		// 链式合并请求
		patterns := make([]string, len(config.Backends))
		for i, v := range config.Backends {
			patterns[i] = v.URLPattern
		}

		return sequentialMerge(patterns, serviceTimeOut, combiner, proxy...)
	}
}

func sequentialMerge(patterns []string, timeout time.Duration, combiner ResponseCombiner, proxy ...Proxy) Proxy {
	return func(ctx context.Context, request *Request) (response *Response, err error) {
		localCtx, cancel := context.WithTimeout(ctx, timeout)

		responses := make([]*Response, len(proxy))
		out := make(chan *Response, 1)
		errChan := make(chan error, 1)

		acc := newIncrementalMergeAccumulator(len(proxy), combiner)

	Loop:
		for i, nextProxy := range proxy {
			if i > 0 {
				for _, match := range sequentialLastParamKeyRegexp.FindAllStringSubmatch(patterns[i], -1) {
					if len(match) > 1 {

						// 第几个backend的下标
						index, err := strconv.Atoi(match[1])
						// 下标不是数字 || 下标大于当前下标 || 下标对应的backend的response为nil
						if err != nil || index >= i || responses[index] == nil {
							continue
						}

						key := "Resp" + match[1] + "_" + match[2]

						var v interface{}
						var ok bool

						data := responses[index].Data
						keys := strings.Split(match[2], ".")
						if len(keys) > 1 {
							for _, k := range keys[:len(keys)-1] {
								v, ok = data[k]
								if !ok {
									break
								}
								switch clean := v.(type) {
								case map[string]interface{}:
									data = clean
								default:
									break
								}
							}
						}
						// 从index对应的backend的response中拿出参数
						v, ok = data[keys[len(keys)-1]]
						if !ok {
							continue
						}

						switch t := v.(type) {
						case string:
							request.Params[key] = t
						case int:
							request.Params[key] = strconv.Itoa(t)
						case float64:
							request.Params[key] = strconv.FormatFloat(t, 'E', -1, 32)
						case bool:
							request.Params[key] = strconv.FormatBool(t)
						default:
							request.Params[key] = fmt.Sprintf("%v", v)
						}
					}
				}
			}
			requestPart(localCtx, nextProxy, request, out, errChan)
			select {
			case err := <-errChan:
				if i == 0 {
					cancel()
					return nil, err
				}
				acc.Merge(nil, err)
				break Loop
			case response := <- out:
				acc.Merge(response, nil)
				if !response.IsComplete {
					break Loop
				}
				responses[i] = response
			}

		}
		result, err := acc.Result()
		cancel()
		return result, err
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
