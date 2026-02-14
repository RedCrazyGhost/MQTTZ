# MQTTZ

<img align="right" width="274px" src="./docs/logo.png">

- 多 Brock 连接管理
- 支持循环间隔发送生成 Mock 数据的 topic
- 支持条件过滤、拦截订阅的 topic
- 支持多 Brock 之间转发 topic

如果这个项目对您有帮助，请给一个 ⭐ Star 支持！

## 环境要求

- **Go 1.26+**
- 从源码构建时需 [wire](https://github.com/google/wire)（`make build` 会自动检查并安装）

## 构建与运行

- **构建**：`make build`，产物为 `./bin/MQTTZ`
- **运行**：
  - `./bin/MQTTZ`：默认读取 `./conf/config.yaml`
  - `make run`：会带 `-config ./bin/conf/config.yaml`（与上者默认路径不同，可按需复制配置或自行指定 `-config`）
- **命令行参数**：`-config`（配置文件路径）、`-port`（服务端口）、`-debug`（调试模式）

## 项目结构

- `cmd/`：程序入口与 Wire 依赖注入
- `pkg/`：conf（配置）、logger（日志）、mqtt（连接管理）、processor（订阅端处理器）
- `model/`：配置与处理器模型
- `utils/`：通用工具
- `template/conf/`：配置模板

## 文档

- [快速开始](docs/quick_start.md)
- [配置说明](docs/configuration.md)

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
      nickname: "MQTTZ_1" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key
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
            forward_rules:
                - to_client: MQTTZ_2
          - topics:
                - "2"
                - "1"
            processors:
                - type: interceptor
                  rule: "topic:contains=2"
                - type: filter
                  rule: "topic:contains=1"


    - broker: 127.0.0.1
      port: 1883
      client_id: mqtt_client_2
      username: ""
      password: ""
      nickname: "MQTTZ_2" # 别名用作查找的 Key，如果不存在则是用 {client_id}@{broker} 作为 Key

```