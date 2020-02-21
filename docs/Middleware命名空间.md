# 关于Melody中Middleware的命名与结构体格式

## melody-xxxxx
- Describe: xxx
- Namespace: `melody_xxxxx`
- Struct:
```
"melody_xxxx": {
    ...
}
```
- Level: [ServiceConfig, Backend, Endpoint]
- Status: xxxx

---


## 1.melody-logstash
- Describe: 包含serviceName的logger
- Namespace: `melody_logstash`
- Struct:
```
"melody_logstash": {
    ... 内容目前不重要，只检察有无该extra config
}
```
- Level: ServiceConfig
- Status: 基本实现

## 2.melody-gologging
- Describe: Base 基础logger middleware
- Namespace: `melody_gologging`
- Struct:
```
"melody_gologging": {
    "level": "DEBUG",
    "prefix": "[GRANT]",
    "syslog": true,
    "stdout": true,
    "format": "default"
}
```
- Level: [ServiceConfig, Backend, Endpoint]
- Status: 基本实现

## 3.melody-viper
- Describe: 基于config中parser实现的viper parser
- Namespace: `melody_parser`
- Struct:
```
"melody_viper": {
    ...
}
```
- Level: [ServiceConfig]
- Status: 基本实现

## 4.melody-gelf
- Describe: 与graylog集成
- Namespace: `melody_gelf`
- Struct:
```
"melody_gelf": {
    "addr": "*:12201",
    "enable_tcp": fasle
}
```
- Level: [ServiceConfig]
- Status: 完成


## 5.melody-metrics
- Describe: 系统的运行数据检测、统计
- Namespace: `melody_metrics`
- Struct:
```
"melody_metrics": {
    proxy_disable     bool
    router_disabled   bool
    backend_disabled  bool
    endpoint_disabled bool
    collection_time   time.Duration
    listen_address       string
}
```
- Level: [ServiceConfig]
- Status: 基本实现


## 6.melody-cors
- Describe: 跨域处理
- Namespace: `melody_cors`
- Struct:
```
"melody_cors": {
    allow_origins     []string
    allow_methods     []string
    allow_headers     []string
    expose_headers    []string
    allow_credentials bool
    max_age           time.Duration
}
```
- Level: [ServiceConfig]
- Status: 基本实现


## 7.melody-httpsecure
- Describe: http 安全相关的一些拦截、过滤、处理
- Namespace: `melody_httpsecure`
- Struct:
```
"melody_httpsecure": {
    "allowed_hosts": [
      "host.known.com:443"
    ],
    "ssl_proxy_headers": {
      "X-Forwarded-Proto": "https"
    },
    "ssl_redirect": true,
    "ssl_host": "ssl.host.domain",
    "ssl_port": "443",
    "ssl_certificate": "/path/to/cert",
    "ssl_private_key": "/path/to/key",
    "sts_seconds": 300,
    "sts_include_subdomains": true,
    "frame_deny": true,
    "custom_frame_options_value": "ALLOW-FROM https://example.com",
    "hpkp_public_key": "pin-sha256=\"base64==\"; max-age=expireTime [; includeSubDomains][; report-uri=\"reportURI\"]",
    "content_type_nosniff": true,
    "browser_xss_filter": true,
    "content_security_policy": "default-src 'self';"
  }
```
- Level: [ServiceConfig]
- Status: 基本实现， 某些字段还不清楚

