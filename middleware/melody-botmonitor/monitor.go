package botmonitor

import (
	"net/http"
	"regexp"

	lru "github.com/hashicorp/golang-lru"
)

// Config 配置定义了检测器的行为
type Config struct {
	Blacklist []string
	Whitelist []string
	Patterns  []string
	CacheSize int
}

// DetectorFunc 是一个func，如果一个请求是由一个机器人发出的，它就会进行chek
type DetectorFunc func(r *http.Request) bool

// New 根据参数返回带有或不带LRU缓存的检测器函数
func New(cfg Config) (DetectorFunc, error) {
	if cfg.CacheSize == 0 {
		d, err := NewDetector(cfg)
		return d.IsBot, err
	}

	d, err := NewLRU(cfg)
	return d.IsBot, err
}

// NewDetector 创建一个检测器
func NewDetector(cfg Config) (*Detector, error) {
	blacklist := make(map[string]struct{}, len(cfg.Blacklist))
	for _, e := range cfg.Blacklist {
		blacklist[e] = struct{}{}
	}
	whitelist := make(map[string]struct{}, len(cfg.Whitelist))
	for _, e := range cfg.Whitelist {
		whitelist[e] = struct{}{}
	}
	patterns := make([]*regexp.Regexp, len(cfg.Patterns))
	for i, p := range cfg.Patterns {
		rp, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		patterns[i] = rp
	}
	return &Detector{
		blacklist: blacklist,
		whitelist: whitelist,
		patterns:  patterns,
	}, nil
}

// Detector (检测器)是一种能够检测机器人发出的请求的结构体
type Detector struct {
	blacklist map[string]struct{}
	whitelist map[string]struct{}
	patterns  []*regexp.Regexp
}

// IsBot : 如果请求是由机器人发出的，则IsBot返回true
func (d *Detector) IsBot(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return false
	}
	if _, ok := d.whitelist[userAgent]; ok {
		return false
	}
	if _, ok := d.blacklist[userAgent]; ok {
		return true
	}
	for _, p := range d.patterns {
		if p.MatchString(userAgent) {
			return true
		}
	}
	return false
}

// NewLRU 创建一个新的LRUDetector
func NewLRU(cfg Config) (*LRUDetector, error) {
	d, err := NewDetector(cfg)
	if err != nil {
		return nil, err
	}

	cache, err := lru.New(cfg.CacheSize)
	if err != nil {
		return nil, err
	}

	return &LRUDetector{
		detectorFunc: d.IsBot,
		cache:        cache,
	}, nil
}

// LRUDetector 是一种能够检测bot发出的请求并缓存结果以供将来重用的结构
type LRUDetector struct {
	detectorFunc DetectorFunc
	cache        *lru.Cache
}

// IsBot 如果请求是由机器人发出的，则IsBot返回true
func (d *LRUDetector) IsBot(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	cached, ok := d.cache.Get(userAgent)
	if ok {
		return cached.(bool)
	}

	res := d.detectorFunc(r)
	d.cache.Add(userAgent, res)

	return res
}
