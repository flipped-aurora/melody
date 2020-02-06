package config

import (
	"errors"
	"regexp"
	"strings"
)

const (
	defaultHttp = `http://`
)

var (
	hostPattern    = regexp.MustCompile(`(https?://)?([a-zA-Z0-9\._\-]+)(:[0-9]{2,6})?/?`)
	errInvalidHost = errors.New("invalid host")
)

//URIParser defines all method that uri needed
type URIParser interface {
	CleanHosts([]string) []string
	CleanHost(string) string
	CleanPath(string) string
	GetEndpointPath(string, []string) string
}

//URI to implement URIParser
type URI int

func (U URI) CleanHosts(hosts []string) []string {
	var cleans []string
	for i := range hosts {
		cleans = append(cleans, U.CleanHost(hosts[i]))
	}

	return cleans
}

func (U URI) CleanHost(host string) string {
	matches := hostPattern.FindAllStringSubmatch(host, -1)
	if len(matches) != 1 {
		panic(errInvalidHost)
	}

	keys := matches[0][1:]
	if keys[0] == "" {
		keys[0] = defaultHttp
	}

	return strings.Join(keys, "")
}

func (U URI) CleanPath(path string) string {
	return "/" + strings.TrimPrefix(path, "/")
}

func (U URI) GetEndpointPath(path string, params []string) string {
	endPoint := path
	if U == BracketsRouterPatternBuilder {

	}

	return endPoint
}
