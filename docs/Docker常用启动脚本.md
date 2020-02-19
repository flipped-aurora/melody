
# 1. 与Graylog集成

> 由于环境依赖 Mongo和Elasticsearch，所以建议Docker方式在本地启动


### 配置docker-compose 文件

找个地方
```
vim graylog.yml
```
直接拷贝如

```yaml
version: '2'
services:
  # MongoDB: https://hub.docker.com/_/mongo/
  mongodb:
    image: mongo:3
  # Elasticsearch: https://www.elastic.co/guide/en/elasticsearch/reference/5.6/docker.html
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:5.6.3
    environment:
      - http.host=0.0.0.0
      - transport.host=localhost
      - network.host=0.0.0.0
      # Disable X-Pack security: https://www.elastic.co/guide/en/elasticsearch/reference/5.6/security-settings.html#general-security-settings
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    mem_limit: 1g
  # Graylog: https://hub.docker.com/r/graylog/graylog/
  graylog:
    image: graylog/graylog:2.4.0-1
    environment:
      # CHANGE ME!
      - GRAYLOG_PASSWORD_SECRET=somepasswordpepper
      # Password: admin
      - GRAYLOG_ROOT_PASSWORD_SHA2=8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918
      - GRAYLOG_WEB_ENDPOINT_URI=http://127.0.0.1:9000/api
    links:
      - mongodb:mongo
      - elasticsearch
    depends_on:
      - mongodb
      - elasticsearch
    ports:
      # Graylog web interface and REST API
      - 9000:9000
      # Syslog TCP
      - 514:514
      # Syslog UDP
      - 514:514/udp
      # GELF TCP
      - 12201:12201
      # GELF UDP
      - 12201:12201/udp
```

### 拉取、启动镜像

```shell script
sudo docker-compose -f graylog.yml up
```

***由于镜像较多，过程可能比较慢，1-2h***

### 访问控制台

查看容器启动状况，确定mongo, elasticsearch和graylog启动成功，并且graylog开放了tcp端口12201和udp端口12201
```shell script
docker ps 
```

### 配置

- 点击 “System>Inputs”, 选择GDELF UDP之后点击Launch new input
    - Node 只有一个可选
    - Title 自定义
    - 其他默认
    - Save 保存在下方
    
### 测试

启动日志测试
```shell script
docker run -d \
           --log-driver=gelf \
           --log-opt gelf-address=udp://localhost:12201 \
           --log-opt tag="{{.ImageName}}/{{.Name}}/{{.ID}}" \
           busybox sh -c 'while true; do echo "Hello, this is A"; sleep 10; done;'
```


# 2. metrics模块与InfluxDB集成，并在Grafana展示

### 2.1. 启动influxDB container
```shell
> docker run -d -p 8086:8086 \
	  -e INFLUXDB_DB=melody \
	  -e INFLUXDB_USER=melody -e INFLUXDB_USER_PASSWORD=melody \
	  -e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=admin \
	  -it --name=influx \
	  influxdb
```
### 2.2. 启动Grafana
```shell
> docker run \
  -d \
  -p 3000:3000 \
  --name=grafana \
  grafana/grafana
```

grafana admin : localhost:3000

在连接influxDB时，ip地址不是localhost或者127，应该是influxDB所在docker的ip地址， 查看方式
```shell
> docker exec -it influx influx
> cat /etc/hosts
```
