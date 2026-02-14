# MQTTZ 配置说明

## 总览

MQTTZ 默认读取当前工作目录下的 `./conf/config.yaml`。可通过命令行参数 `-config` 指定其它路径。

配置文件为 YAML，顶层结构包括：

- `log`：日志
- `server`：服务端（如 HTTP 端口）
- `mqtt_configs`：MQTT 客户端连接与发布/订阅配置列表

## log

| 字段 | 说明 |
|------|------|
| `enable_debug` | 是否开启调试 |
| `level` | 日志级别 |
| `enable_color` | 是否彩色输出 |
| `output_file` | 日志输出文件路径 |
| `max_size` | 单个日志文件最大尺寸（MB） |
| `max_backups` | 最大保留文件数 |
| `max_age` | 最大保留天数 |
| `compress` | 是否压缩历史日志 |

## server

| 字段 | 说明 |
|------|------|
| `port` | 服务端口（可与命令行 `-port` 覆盖） |

## mqtt_configs

每个元素对应一个 MQTT 客户端连接，包含以下字段：

| 字段 | 说明 |
|------|------|
| `broker` | Broker 地址 |
| `port` | Broker 端口 |
| `client_id` | 客户端 ID |
| `username` | 用户名（可选） |
| `password` | 密码（可选） |
| `nickname` | 别名，用作转发等场景下查找的 Key；不填则使用 `{client_id}@{broker}` |
| `pub_configs` | 发布配置列表，见下 |
| `sub_configs` | 订阅配置列表，见下 |

### pub_configs（发布配置）

| 字段 | 说明 |
|------|------|
| `enable_for` | 是否启用该条发布任务 |
| `interval` | 发送间隔（如 `1s`） |
| `source_type` | 数据源类型：`conf`（内联）、`json`、`yaml`（文件路径） |
| `source_path` | 当 `source_type` 为 `json` 或 `yaml` 时的文件路径 |
| `source_data` | 当 `source_type` 为 `conf` 时的内联数据列表，每项含 `topic`、`payload` 等 |
| `is_strong` | 数据模型是否强匹配（可选） |

### sub_configs（订阅配置）

| 字段 | 说明 |
|------|------|
| `topic` | 单个订阅主题（与 `topics` 二选一） |
| `topics` | 多个订阅主题列表 |
| `qos` | QoS（0/1/2） |
| `processors` | 前置处理器列表，见下 |
| `forward_rules` | 转发规则列表，见下 |

### processors（处理器）

用于对订阅到的消息做过滤或拦截，每条处理器包含：

| 字段 | 说明 |
|------|------|
| `type` | 类型：`interceptor`（拦截器）、`filter`（过滤器）、`forwarder`、`extractor`、`generator` |
| `rule` | 规则字符串，格式为 `target:rule` |

**rule 格式**：`target` 与 `rule` 用冒号分隔。`target` 可选：

- `topic`：按主题匹配
- `qos`：按 QoS
- `retain`：按 retain 标志
- `payload`：按负载内容

示例：`topic:contains=2`、`topic:contains=1`（需符合 [validator](https://github.com/go-playground/validator) 语法）。  
拦截器（interceptor）：规则匹配则拦截；过滤器（filter）：规则不匹配则过滤掉。

### forward_rules（转发规则）

| 字段 | 说明 |
|------|------|
| `to_client` | 目标客户端的 `nickname` |
| `processors` | 可选，转发前的处理器列表 |

## 完整示例

以下为与模板和 README 示例一致的配置示例，可直接复制到 `conf/config.yaml` 后按需修改。

```yaml
log:
  enable_debug: true
  enable_color: true

mqtt_configs:
  - broker: 127.0.0.1
    port: 1883
    client_id: mqtt_client_1
    username: ""
    password: ""
    nickname: "MQTTZ_1"
    pub_configs:
      - enable_for: true
        interval: 1s
        source_type: conf
        source_data:
          - topic: "1"
            payload: { "test": 1 }
          - topic: "2"
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
    nickname: "MQTTZ_2"
```
