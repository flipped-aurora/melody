<div align=center>
<img src="http://picture.zyuhn.top/logo.png" onerror="https://github.com/granty1/melody/blob/master/docs/img/melody.png" width="300px" height="336px"/>
</div>
<div align=center>
<img src="https://img.shields.io/badge/golang-1.12-blue"/>
<img src="https://github.com/granty1/melody/workflows/Go/badge.svg"/>
<img src="https://travis-ci.com/granty1/melody.svg?branch=master"/>
<img src="https://coveralls.io/repos/github/granty1/melody/badge.svg?branch=master"/>
<img src="https://goreportcard.com/badge/github.com/granty1/melody"/>
</div>

[English](./README.md) | 简体中文

# **Melody API 网关**

- 编译环境:  Golang 1.11+
- 当前版本:  1.0.0

## 简单介绍

Melody是一个高性能的开源API网关，可以帮助您整理复杂的API。

如果您正在构建web、移动或物联网(Internet of Things)，您可能最终需要使用通用功能来运行实际的软件。Melody可以作为微服务请求的网关，同时通过中间件提供负载平衡、日志记录、身份验证、速率限制、转换等功能。

[在线文档](https://granty1.github.io/melody-docs) | [下载](https://github.com/flipped-aurora/melody/releases) | [站点](https://granty1.github.io/melody-web/) | [可视化配置](  https://granty1.github.io/melody-config)

## 其他资料

- [Melody Document Repository](https://github.com/granty1/melody-docs):  Melody 的在线文档仓库。
- [Melody Document](https://granty1.github.io/melody-docs):  一个详细的Melody说明文档。
- [Melody Data Monitor](https://github.com/granty1/melody-data):  服务监控可视化系统，接受来自Melody的数据并分析展示。
- [Melody Config Respository](https://github.com/granty1/melody-config):  在线可视化配置系统的代码仓库。
- [Melody Online Config](https://granty1.github.io/melody-config):  在线的可视化配置系统。
- [Melody SIte Respository](https://github.com/granty1/melody-web):  Melody官方网站的仓库。
- [Melody Site](https://granty1.github.io/melody-web/): Melody官方网站，提供最新的可执行文件下载。
- [Test API Repository](  https://github.com/granty1/gin-gorm-jwt-quick-start):  使用此来快速开启您的golang web应用程序。

## 代码编译

在Linux 或者 Mac上你可以使用如下命令编译我们的代码 :

```
make build
```

在Windows操作系统上需要使用 :

```
go build .
```

## 使用方式

通过 melody [command] 的方式使用Melody API 网关。


## 可用命令

| 命令  | 描述                     |
| ----- | ------------------------ |
| check | 检查配置文件格式         |
| run   | 启动Melody服务器         |
| help  | 命令提示帮助             |
| graph | 生成Melody配置文件结构图 |


## 命令参数

| 参数缩写 | 完整参数 | 描述                                    |
| -------- | -------- | --------------------------------------- |
| -c       | -config  | melody.json配置文件的路径               |
| -d       | -debug   | 允许开启debug模式，将列出更加详细的信息 |
| -h       | -help    | 命令提示帮助                            |

## 使用示例

使用Git克隆我们的项目：

```shell
$ git clone https://github.com/flipped-aurora/melody.git
$ cd melody
$ go build .
```

创建一个 **melody. Json** 的配置文件，一般而言endpoints代表暴露给客户端的节点，backends对应微服务后端的节点，示例配置如下：

```JSON
{
	"version": 1,
	"extra_config": {
		"melody_gologging": {
			"level": "DEBUG",
			"prefix": "[Grant]",
			"syslog": false,
			"stdout": true,
			"format": "default"
		}
	},
	"timeout": "3000ms",
	"cache_ttl": "300s",
	"output_encoding": "json",
	"port": 8000,
	"endpoints": [{
			"endpoint": "/findone/{name}",
			"method": "GET",
			"extra_config": {
				"melody_proxy": {
					"sequential": true
				}
			},
			"output_encoding": "json",
			"concurrent_calls": 1,
			"backends": [{
					"url_pattern": "/user/{name}",
					"group": "base_info",
					"encoding": "json",
					"extra_config": {},
					"method": "GET",
					"host": [
						"127.0.0.1:9001"
					]
				},
				{
					"url_pattern": "/role/{resp0_base_info.role_id}",
					"encoding": "json",
					"extra_config": {},
					"method": "GET",
					"group": "role_info",
					"host": [
						"127.0.0.1:9001"
					]
				}
			]
		},
		{
			"endpoint": "/ping",
			"backends": [{
				"url_pattern": "/__debug/bar",
				"host": [
					"127.0.0.1:8000"
				]
			}]
		}
	]
}
```

使用命令检查配置文件：

```
melody check -c melody.json
```

启动Melody服务器：

```
melody run -c melody.json
```

## 主要功能

- **命令行**: 使用命令行命令控制Melody网关。
- **支持 REST API**: 可以使用RESTful API操作Melody。
- **转发代理**: 为客户端和服务端提供基于HTTP的转发代理功能。
- **数据转换**: 可以添加、删除或操作HTTP请求和响应，包括:数据合并、筛选、过滤、分组、映射。
- **内容转码**:  Melody可以使用多种内容编码类型，包括:json、string、noop、rss、yaml、xml。
- **流量限制**: 基于令牌桶算法的速率限制方式。
- **签名验证**: Melody支持 JWT 和 JWK.
- **服务监控**: 实时监控提供关键负载和性能服务器指标，集成metrics和influxdb。
- **服务发现:** Melody集成etcd、consul来支持服务发现。
- **记录**: 记录系统的请求和响应信息。
- **系统日志**: 提供系统日志记录。
- **安全**: 点劫持保护、黑名单/白名单、MIME嗅探预防、XSS防护，同时允许HSTS，未来将集成OAuth2。
- **断路器**: 断路器机制将检测到故障并通过不发送可能会失败的请求来防止对服务器造成压力。
- **CORS跨域**: 支持跨域配置。

有关中间件和集成的更多信息，可以访问 [在线文档](https://granty1.github.io/melody-docs)。