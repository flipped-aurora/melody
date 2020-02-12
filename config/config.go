package config

import (
	"errors"
	"fmt"
	"melody/encoding"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	BracketsRouterPatternBuilder = iota
	ColonRouterPatternBuilder
	defaultPort                = 8080
	defaultMaxIdleConnsPreHost = 250
	defaultTimeout             = 2 * time.Second
	CurrVersion                = 2
)

var (
	RoutingPattern          = ColonRouterPatternBuilder
	debugPattern            = "^[^/]|/__debug(/.*)?$"
	sequentialParamsPattern = regexp.MustCompile(`^resp[\d]+_.*$`)
	simpleURLKeysPattern    = regexp.MustCompile(`\{([a-zA-Z\-_0-9]+)\}`)
	errInvalidNoOpEncoding  = errors.New("can not use NoOp encoding with more than one backends connected to the same endpoint")
)

//ServiceConfig contains all config in melody server.
type ServiceConfig struct {
	ExtraConfig    ExtraConfig       `mapstructure:"extra_config"`
	Port           int               `mapstructure:"port"`
	Timeout        time.Duration     `mapstructure:"timeout"`
	Host           []string          `mapstructure:"host"`
	Endpoints      []*EndpointConfig `mapstructure:"endpoints"`
	Version        int               `mapstructure:"version"`
	OutputEncoding string            `mapstructure:"output_encoding"`
	CacheTTL       time.Duration     `mapstructure:"cache_ttl"`

	MaxIdleConnsPerHost int `mapstructure:"max_idle_connections_per_host"`

	DisableStrictREST bool `mapstructure:"disable_rest"`
	//melody is in debug model
	Debug     bool
	uriParser URIParser
}

type EndpointConfig struct {
	// 对外暴露的url
	Endpoint string `mapstructure:"endpoint"`
	// HTTP method of the endpoint (GET, POST, PUT, etc)
	Method string `mapstructure:"method"`
	// 此端点连接的后端的集合
	Backends []*Backend `mapstructure:"backend"`
	// 此端点的并发调用数
	ConcurrentCalls int `mapstructure:"concurrent_calls"`
	// 此端点的超时时间
	Timeout time.Duration `mapstructure:"timeout"`
	// 缓存头的持续时间
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
	// 要从URI提取的查询字符串参数列表
	QueryString []string `mapstructure:"querystring_params"`
	// extra config
	ExtraConfig ExtraConfig `mapstructure:"extra_config"`
	// 可以传递的请求头列表
	HeadersToPass []string `mapstructure:"headers_to_pass"`
	// 响应输出时的编码
	OutputEncoding string `mapstructure:"output_encoding"`
}

type Backend struct {
	// the name of the group the response should be moved to. If empty, the response is
	// not changed
	Group string `mapstructure:"group"`
	// HTTP method of the request to send to the backend
	Method string `mapstructure:"method"`
	// Set of hosts of the API
	Host []string `mapstructure:"host"`
	// False if the hostname should be sanitized
	HostSanitizationDisabled bool `mapstructure:"disable_host_sanitize"`
	// URL pattern to use to locate the resource to be consumed
	URLPattern string `mapstructure:"url_pattern"`
	// set of response fields to remove. If empty, the filter id not used
	Blacklist []string `mapstructure:"blacklist"`
	// set of response fields to allow. If empty, the filter id not used
	Whitelist []string `mapstructure:"whitelist"`
	// map of response fields to be renamed and their new names
	Mapping map[string]string `mapstructure:"mapping"`
	// the encoding format
	Encoding string `mapstructure:"encoding"`
	// the response to process is a collection
	IsCollection bool `mapstructure:"is_collection"`
	// name of the field to extract to the root. If empty, the formater will do nothing
	Target string `mapstructure:"target"`
	// name of the service discovery driver to use
	//SD string `mapstructure:"sd"`

	// list of keys to be replaced in the URLPattern
	URLKeys []string
	// number of concurrent calls this endpoint must send to the API
	ConcurrentCalls int
	// timeout of this backend
	Timeout time.Duration
	// decoder to use in order to parse the received response from the API
	Decoder encoding.Decoder `json:"-"`
	// Backend Extra configuration for customized behaviours
	ExtraConfig ExtraConfig `mapstructure:"extra_config"`
}

//Extra config for melody
type ExtraConfig map[string]interface{}

type EndpointMatchError struct {
	Path   string
	Method string
	Err    error
}

type NoBackendsError struct {
	Path   string
	Method string
}

func (n *NoBackendsError) Error() string {
	return fmt.Sprintf("ERROR: path:%s, method:%s has 0 backends", n.Path, n.Method)
}

func (e *EndpointMatchError) Error() string {
	return fmt.Sprintf("ERROR: parsing endpoint error : url:%s, method:%s, error:%s", e.Path, e.Method, e.Err)
}

// 该方法作用等于clean一下整个struct
func (e *ExtraConfig) sanitize() {
	for module, extra := range *e {
		switch extra := extra.(type) {
		case map[interface{}]interface{}:
			sanitized := map[string]interface{}{}
			for k, v := range extra {
				sanitized[fmt.Sprintf("%v", k)] = v
			}
			(*e)[module] = sanitized
		}
	}
}

func (s *ServiceConfig) Init() error {
	// 初始化URIParser
	s.uriParser = NewURIParser()

	// 判断版本号
	if s.Version != CurrVersion {
		return &UnsupportedVersionError{
			Have: s.Version,
			Want: CurrVersion,
		}
	}
	// 初始化全局参数
	s.initGlobalParams()

	// 初始化Endpoints
	return s.initEndpoints()
}

func (s *ServiceConfig) initGlobalParams() {
	if s.Port == 0 {
		s.Port = defaultPort
	}

	if s.MaxIdleConnsPerHost == 0 {
		s.MaxIdleConnsPerHost = defaultMaxIdleConnsPreHost
	}

	if s.Timeout == 0 {
		s.Timeout = defaultTimeout
	}

	s.Host = s.uriParser.CleanHosts(s.Host)
	s.ExtraConfig.sanitize()
}

func (s *ServiceConfig) initEndpoints() error {
	for i, e := range s.Endpoints {
		e.Endpoint = s.uriParser.CleanPath(e.Endpoint)

		if err := e.validate(); err != nil {
			return err
		}
		// 支Rest风格的前提下
		// 从Endpoint url中提取参数列表
		// 类似 -> /debug/{id}/{name}
		inputUrlParams := s.getPlaceHoldersFromEndpointUrl(e.Endpoint, s.paramExtractionPattern())
		inputParamsSet := map[string]interface{}{}
		for _, v := range inputUrlParams {
			inputParamsSet[v] = nil
		}
		// gin中rest风格与其他有别，所以将路由
		// /debug/{id}/{name} -> /debug/:id/:name
		e.Endpoint = s.uriParser.GetEndpointPath(e.Endpoint, inputUrlParams)
		// 初始化一些全局默认值
		s.initDefaultEndpoints(i)

		//判断encode为NOOP时，backends是否大于1
		if e.OutputEncoding == encoding.NOOP && len(e.Backends) > 1 {
			return errInvalidNoOpEncoding
		}

		e.ExtraConfig.sanitize()

		for j, b := range e.Backends {
			s.initDefaultBackends(i, j)

			// 初始化、解析Backends的url以及对应的参数 （*）
			s.initBackendsURLMappings(i, j, inputParamsSet)

			b.ExtraConfig.sanitize()
		}
	}

	return nil
}

// 提取参数的模式
// 1. 严格模式
// 2. 非严格模式
func (s *ServiceConfig) paramExtractionPattern() *regexp.Regexp {
	if s.DisableStrictREST {
		return simpleURLKeysPattern
	}

	return endpointURLKeysPattern
}

func (s *ServiceConfig) getPlaceHoldersFromEndpointUrl(endpoint string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringSubmatch(endpoint, -1)
	params := make([]string, len(matches))

	for i, v := range matches {
		params[i] = v[1]
	}

	return params
}

func (s *ServiceConfig) initDefaultEndpoints(i int) {
	cur := s.Endpoints[i]

	if cur.Method == "" {
		cur.Method = "GET"
	}

	if s.CacheTTL != 0 && cur.CacheTTL == 0 {
		cur.CacheTTL = s.CacheTTL
	}

	if s.Timeout != 0 && cur.Timeout == 0 {
		cur.Timeout = s.Timeout
	}

	if cur.ConcurrentCalls == 0 {
		cur.ConcurrentCalls = 1
	}

	if cur.OutputEncoding == "" {
		if s.OutputEncoding != "" {
			cur.OutputEncoding = s.OutputEncoding
		} else {
			cur.OutputEncoding = encoding.JSON
		}
	}
}

func (s *ServiceConfig) initDefaultBackends(e, b int) {
	endpoint := s.Endpoints[e]
	backend := endpoint.Backends[b]
	if len(backend.Host) == 0 {
		backend.Host = s.Host
	}

	if !backend.HostSanitizationDisabled {
		backend.Host = s.uriParser.CleanHosts(backend.Host)
	}

	if backend.Method == "" {
		backend.Method = endpoint.Method
	}

	backend.Timeout = endpoint.Timeout
	backend.ConcurrentCalls = endpoint.ConcurrentCalls

	//根据配置的encoding， 加载对应的Decoder
	backend.Decoder = encoding.Get(strings.ToLower(backend.Encoding))(backend.IsCollection)
}

func (s *ServiceConfig) initBackendsURLMappings(e, b int, inputSet map[string]interface{}) error {
	backend := s.Endpoints[e].Backends[b]
	backend.URLPattern = s.uriParser.CleanPath(backend.URLPattern)

	outputParams, outputParamsSize := uniqueOutput(s.getPlaceHoldersFromEndpointUrl(backend.URLPattern, simpleURLKeysPattern))

	inputParams := convertToSlice(inputSet)

	if outputParamsSize > len(inputParams) {
		return &WrongNumberOfParamsError{
			Endpoint:     s.Endpoints[e].Endpoint,
			Method:       s.Endpoints[e].Method,
			Backend:      b,
			InputParams:  inputParams,
			OutputParams: outputParams,
		}
	}

	backend.URLKeys = []string{}
	for _, param := range outputParams {
		if !sequentialParamsPattern.MatchString(param) {
			if _, ok := inputSet[param]; !ok {
				return &UndefinedOutputParamError{
					Endpoint:     s.Endpoints[e].Endpoint,
					Method:       s.Endpoints[e].Method,
					Backend:      b,
					InputParams:  inputParams,
					OutputParams: outputParams,
					Param:        param,
				}
			}
		}
		key := strings.Title(param)
		backend.URLPattern = strings.Replace(backend.URLPattern, "{"+param+"}", "{{."+key+"}}", -1)
		backend.URLKeys = append(backend.URLKeys, key)
	}

	return nil
}

func convertToSlice(inputSet map[string]interface{}) []string {
	var inputParams []string
	for key := range inputSet {
		inputParams = append(inputParams, key)
	}

	sort.Strings(inputParams)
	return inputParams
}

func (e *EndpointConfig) validate() error {
	matched, err := regexp.MatchString(debugPattern, e.Endpoint)
	if err != nil {
		return &EndpointMatchError{
			Path:   e.Endpoint,
			Method: e.Method,
			Err:    err,
		}
	}

	if matched {
		return &EndpointMatchError{
			Path:   e.Endpoint,
			Method: e.Method,
		}
	}

	if len(e.Backends) == 0 {
		return &NoBackendsError{
			Path:   e.Endpoint,
			Method: e.Method,
		}
	}

	return nil
}

func uniqueOutput(output []string) ([]string, int) {
	sort.Strings(output)
	j := 0
	outputSetSize := 0
	for i := 1; i < len(output); i++ {
		if output[j] == output[i] {
			continue
		}
		if !sequentialParamsPattern.MatchString(output[j]) {
			outputSetSize++
		}
		j++
		output[j] = output[i]
	}
	if j == len(output) {
		return output, outputSetSize
	}
	return output[:j+1], outputSetSize
}

type WrongNumberOfParamsError struct {
	Endpoint     string
	Method       string
	Backend      int
	InputParams  []string
	OutputParams []string
}

// Error returns a string representation of the WrongNumberOfParamsError
func (w *WrongNumberOfParamsError) Error() string {
	return fmt.Sprintf(
		"input and output params do not match. endpoint: %s %s, backend: %d. input: %v, output: %v",
		w.Method,
		w.Endpoint,
		w.Backend,
		w.InputParams,
		w.OutputParams,
	)
}

type UndefinedOutputParamError struct {
	Endpoint     string
	Method       string
	Backend      int
	InputParams  []string
	OutputParams []string
	Param        string
}

// Error returns a string representation of the UndefinedOutputParamError
func (u *UndefinedOutputParamError) Error() string {
	return fmt.Sprintf(
		"Undefined output param '%s'! endpoint: %s %s, backend: %d. input: %v, output: %v",
		u.Param,
		u.Method,
		u.Endpoint,
		u.Backend,
		u.InputParams,
		u.OutputParams,
	)
}

type UnsupportedVersionError struct {
	Have int
	Want int
}

func (u *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("Unsupported version: %d (want: %d)", u.Have, u.Want)
}
