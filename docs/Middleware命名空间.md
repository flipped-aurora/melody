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


## 8.melody-http
- Describe: http请求过程中的一些处理
- Namespace: `melody_http`
- Struct:
```
"melody_http": {
    "return_error_details": "backend_a"<string>
}
```
- Example:

config:
```
"backend": [
        {
            "host": ["http://127.0.0.1:8081"],
            "url_pattern": "/foo",
            "extra_config": {
                "melody_http": {
                    "return_error_details": "backend_a"
                }
            }
        },
        {
            "host": ["http://127.0.0.1:8081"],
            "url_pattern": "/bar",
            "extra_config": {
                "melody_http": {
                    "return_error_details": "backend_b"
                }
            }
        }
    ]
```
b服务挂掉了

response:
```
{
    "error_backend_b": {
        "http_status_code": 404,
        "http_body": "404 page not found\\n"
    },
    "foo": 42
}
```
- Level: [BackendConfig]
- Status: 基本实现


## 9.melody-proxy(Endpoint)
- Describe: 代理时的一些配置，针对Endpoint层请求
- Namespace: `melody_proxy`
- Struct:
```
"melody_proxy": {
    // 表示开启链式请求
    "sequential": true
    // 静态数据插入
    "static": {
        "strategy": ["always"/"success"/"errored"/"complete"/"imcomplete"],
        "data": {
            "key": value
        }
    }
}
```
- Level: [Endpoint]
- Status: 基本实现

## 10.melody-etcd
- Describe: 创建一个ectd client，注册到sd中
- Namespace: `melody_etcd`
- Struct:
```
"melody_xxxx": {
    "machines": [
        "http://127.0.0.1:8080",
        "http://127.0.0.1:8081",
        "http://127.0.0.1:8082",
        "http://127.0.0.1:8083"
    ],
    "dial_timeout": "5s",
    "dial_keepalive": "30s",
    "header_timeout": "1s",
    "cert": "/https/cert",
    "key": "/https/privateKey",
    "cacert": "/https/CaCert"
}
```
- Level: [ServiceConfig]
- Status: 还有个TODO


## 11.melody-consul
- Describe: custom的服务发现
- Namespace: `sd`
- Struct:
```
backend: [{
    "sd": "custom"
}}
```
- Level: [Backend]
- Status: 刚开始


## 12.melody-dns srv
- Describe: 从 dns srv 中查找对应的host的ip
- Namespace: `sd`
- Struct:
```
backend: [{
    "sd": "dns"
}}
```
- Level: [Backend]
- Status: 基本完成

## 13.melody-ratelimit juju-router
- Describe: endpoint层流量限制
- Namespace: `melody_ratelimit_router`
- Struct:
```
	...
	"extra_config": {
		...
		""melody_ratelimit_router": {
			"maxRate": 2000,
			"strategy": "header",
			"clientMaxRate": 100,
			"key": "X-Private-Token",
		},
		...
	},
	...

```
- Level: [Endpoint]
- Status: 完成

## 14.melody-ratelimit juju-proxy
- Describe: backend层流量限制
- Namespace: `melody_ratelimit_proxy`
- Struct:
```
	...
	"extra_config": {
		...
		"melody-ratelimit_proxy": {
			"maxRate": 100,
			"capacity": 100
		},
		...
	},
	...

```
- Level: [Backend]
- Status: 将要完成


## 15.melody_proxy(backend)
- Describe: 作用于backend层的proxy，针对单个请求、响应
- Namespace: `melody_proxy`
- Struct:
```
"extra_config": {
    "melody_proxy": {
        // 按照正常情况请求后端，但屏蔽该backend的数据回应
        "shadow": true,
        // 针对单个banckend的response为数组的move、del操作
        "flatmap_filter": [
            {
              "type": "move",
              "args": [
                "data.0.name",
                "data.0.id"
              ]
            },
            {
              "type": "del",
              "args": [
                "data.0.name"
              ]
            }
        ]
    }
},

```
- Level: [Backend]
- Status: 完成

