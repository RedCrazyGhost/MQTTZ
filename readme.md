# MQTTZ

## 功能目标

- 提供对多 MQTT 服务处理的能力
    - 添加多个服务的连接配置
    - 提供数据转发能力
    - 提供多种数据源的处理（json、yaml、MQTTClient）

```mermaid
input 处理器（过滤、序列化...）-> dataCh -> output 处理器（反序列化、过滤...）
```

## 前置条件

使用 [EMQX](https://github.com/emqx/emqx) 的 docker 镜像本地部署服务

```shell
docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx:latest
```

使用 MQTTX 作为 MQTT-Client 监看

```shell
go get github.com/eclipse/paho.mqtt.golang
```

## 支持功能

- 基础功能
    - 配置 MQTT 协议配置，并且提供默认值
- 提供发送 topic 数据的能力
    - 提供循环发送
    - 提供发送间隔顺序的能力
- 提供录制接收到的 topic 数据的能力
    - 录制成文件形式，并且可以直接给发送 topic 模块使用

## 使用参数

默认读取 `./conf/config.yaml` 文件

```yaml
mqtt_configs:
    - broker: 127.0.0.1
      port: 1883
      client_id: mqtt_client_1
      username: ""
      password: ""
      nickname: "abc" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key
    - broker: 127.0.0.1
      port: 1883
      client_id: mqtt_client_2
      username: ""
      password: ""
      nickname: "" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key

input_configs:
```