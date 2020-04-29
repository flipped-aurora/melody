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

English | [简体中文](./README-zh_CN.md)

# **Melody API Gateway**

- Build Requirements:  Golang 1.11+
- Melody Version:  1.0.8

## Introduction

Melody is a high-performance open source API gateway to help you to sort out your complex api.

If you are building for the web, mobile, or IoT (Internet of Things) you will likely end up needing common functionality to run your actual software. Melody can help by acting as a gateway for microservices requests while providing load balancing, logging, authentication, rate-limiting, transformations, and more through middlewares.

[Documentation](https://granty1.github.io/melody-docs) | [Installation](https://github.com/flipped-aurora/melody/releases) | [Site](https://granty1.github.io/melody-web/) | [Configer](  https://granty1.github.io/melody-config) | [Test API](https://github.com/granty1/gin-gorm-jwt-quick-start)

## Resource

- [Melody Document Repository](https://github.com/granty1/melody-docs): A Document Repository for Melody.
- [Melody Document](https://granty1.github.io/melody-docs):Detailed documentation for Melody.
- [Melody Data Monitor](https://github.com/granty1/melody-data):Data monitoring system built using Vue.
- [Melody Config Respository](https://github.com/granty1/melody-config):A respository for online visual configuration system.
- [Melody Online Config](https://granty1.github.io/melody-config):Online visual configuration system for Molody.
- [Melody SIte Respository](https://github.com/granty1/melody-web):Melody Web Site Respository.
- [Melody Site](https://granty1.github.io/melody-web/):Melody Web Site.
- [Test API Repository](  https://github.com/granty1/gin-gorm-jwt-quick-start):Using this to quick satrt your golang web application.

## Build

Linux or Mac :

```
make build
```

Windows :

```
go build .
```

## Usage:

  melody [command]


## Available Commands:

| command | *description*                   |
| ------- | ------------------------------- |
| check   | Check that the config           |
| run     | Run the Melody server           |
| help    | Help about any command          |
| graph   | generate graph of melody server |


## Flags:

| flag-abbr | flag-full | description             |
| --------- | --------- | ----------------------- |
| -c        | -config   | Path of the melody.json |
| -d        | -debug    | Enable the Melody debug |
| -h        | -help     | Help for melody         |

## Use example

Clone our project

```shell
$ git clone https://github.com/flipped-aurora/melody.git
$ cd melody
$ go build .
```

Create a file called **melody. Json**

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

Check that the config

```
melody check -c melody.json
```

Use command to run Melody

```
melody run -c melody.json
```

## Features

- **CLI**: Control your Melody API Gateway from the command line.
- **REST API**: Melody can be operated with its RESTful API for maximum flexibility.
- **Transformations**: Add, remove, or manipulate HTTP requests and responses.
- **Rate-limiting**: Block and throttle requests based on many variables.
- **Authentications**: Melody supports JWT and JWK.
- **Monitoring**: Live monitoring provides key load and performance server metrics.
- **Service Discovery**: Melody integrates etcd to support service discovery.
- **Logging**: Log requests and responses to your system.
- **Syslog**: Logging to System log.
- **Security**: Bot Monitor, whitelist/blacklist,SSL, etc...
- **Circuit-Breaker**: Intelligent tracking of unhealthy upstream services.
- **Forward Proxy**: Make Melody connect to intermediary transparent HTTP proxies.
- **metrics**:Monitor and statistic system operation data.
- **influxdb**: Melody integrates with influxDB, writing the data gathered by Metrics to influxDB.
- **CORS**:Support for cross-domain processing.

For more info about middlewares and integrations, you can visit [Documentation](https://granty1.github.io/melody-docs)

## 