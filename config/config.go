package config

import (
	"fmt"
	"melody/encoding"
	"regexp"
	"time"
)

const (
	BracketsRouterPatternBuilder = iota
	ColonRouterPatternBuilder
	defaultPort                = 8080
	defaultMaxIdleConnsPreHost = 250
	defaultTimeout             = 2 * time.Second
)

var (
	RoutingPattern       = ColonRouterPatternBuilder
	debugPattern         = "^[^/]|/__debug(/.*)?$"
	simpleURLKeysPattern = regexp.MustCompile(`\{([a-zA-Z\-_0-9]+)\}`)
)

//ServiceConfig contains all config in melody server.
type ServiceConfig struct {
	ExtraConfig ExtraConfig       `mapstructure:"extra_config"`
	Port        int               `mapstructure:"port"`
	Timeout     time.Duration     `mapstructure:"timeout"`
	Host        []string          `mapstructure:"host"`
	Endpoints   []*EndpointConfig `mapstructure:"endpoints"`

	OutputEncoding string `mapstructure:"output_encoding"`
	CacheTTL time.Duration `mapstructure:"cache_ttl"`

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

	//TODO 判断版本一致

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

		//TODO NOOP encode (目前不知道noop什么意思)

		e.ExtraConfig.sanitize()

		//TODO 初始化 Endpoints下的Backends （*）
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
