package httpsecure

import (
	"melody/config"

	"github.com/unrolled/secure"
)

// Namespace key of http secure
const Namespace = "melody_httpsecure"

// GetConfig return a httpsecure config interface
func GetConfig(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}

	t, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	c := secure.Options{}

	setInt64(t, "sts_seconds", &c.STSSeconds)

	setStrings(t, "allowed_hosts", &c.AllowedHosts)
	setStrings(t, "host_proxy_headers", &c.HostsProxyHeaders)

	setString(t, "custom_frame_option_value", &c.CustomFrameOptionsValue)
	setString(t, "content_security_policy", &c.ContentSecurityPolicy)
	setString(t, "public_key", &c.PublicKey)
	setString(t, "ssl_host", &c.SSLHost)
	setString(t, "referrer_policy", &c.ReferrerPolicy)

	setBool(t, "content_type_nosniff", &c.ContentTypeNosniff)
	setBool(t, "browser_xss_filter", &c.BrowserXssFilter)
	setBool(t, "is_development", &c.IsDevelopment)
	setBool(t, "sts_include_subdomains", &c.STSIncludeSubdomains)
	setBool(t, "frame_deny", &c.FrameDeny)
	setBool(t, "ssl_redirect", &c.SSLRedirect)

	return c
}

func setStrings(t map[string]interface{}, key string, s *[]string) {
	if v, ok := t[key]; ok {
		var result []string
		for _, s := range v.([]interface{}) {
			if str, ok := s.(string); ok {
				result = append(result, str)
			}
		}
		*s = result
	}
}

func setString(t map[string]interface{}, key string, s *string) {
	if v, ok := t[key]; ok {
		if str, ok := v.(string); ok {
			*s = str
		}
	}
}

func setBool(t map[string]interface{}, key string, b *bool) {
	if v, ok := t[key]; ok {
		if str, ok := v.(bool); ok {
			*b = str
		}
	}
}

func setInt64(t map[string]interface{}, key string, i *int64) {
	if v, ok := t[key]; ok {
		switch a := v.(type) {
		case int64:
			*i = a
		case int:
			*i = int64(a)
		case float64:
			*i = int64(a)
		}
	}
}
