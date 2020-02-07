## 关于Melody中Middleware的命名与结构体格式

### melody-xxxxx
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


### 1.melody-logstash
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

### 2.melody-gologging
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

### 3.melody-viper
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
