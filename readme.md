# MQTTX

<img align="right" width="274px" src="./docs/logo.png">

## 功能
- 多 Brock 连接管理
- 支持循环间隔发送生成 Mock 数据的 topic
- 支持条件过滤、拦截订阅的 topic
- 支持多 Brock 之间转发 topic

## 依赖库

- [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) - MQTT 客户端库
- [zap](https://github.com/uber-go/zap) - 高性能日志库
- [sonic](https://github.com/bytedance/sonic) - 高性能 JSON 处理
- [wire](https://github.com/google/wire) - 依赖注入
- [gofakeit](https://github.com/brianvoe/gofakeit) - 数据生成工具

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
      pub_configs:
          - enable_for: true
            interval: 1s
            source_type: conf
            source_data:
                - topic: 1
                  payload: { "test": 1 }
                - topic: 2
                  payload: { "test": 2 }
          - enable_for: true
            interval: 1s
            source_type: json
            source_path: ./bin/conf/data.json
      sub_configs:
          - topic: json/1
            qos: 0
          - topics:
                - "2"
                - "1"
    - broker: 127.0.0.1
      port: 1883
      client_id: mqtt_client_2
      username: ""
      password: ""
      nickname: "" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key

```