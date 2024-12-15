# MQTTZ

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

读取当前路径下的 `config.json` 文件

```json
{
    "broker": "127.0.0.1",
    "port": 1883,
    "client_id": "mqtt_client",
    "username": "",
    "password": "",
    "input_configs": [
        {
            "is_for": false,
            "interval": "1s",
            "mqtt_data": [
                {
                    "topic": "1",
                    "data": {
                        "id": "1"
                    }
                },
                {
                    "topic": "1",
                    "data": {
                        "id": "2"
                    }
                }
            ]
        },
        {
            "is_for": false,
            "interval": "3s",
            "mqtt_data": [
                {
                    "topic": "2",
                    "data": {
                        "id": "3"
                    }
                },
                {
                    "topic": "2",
                    "data": {
                        "id": "4"
                    }
                }
            ]
        },
        {
            "is_for": false,
            "interval": "5s",
            "source": "out.json"
        }
    ],
    "output_config": {
        "sub_topics": [
            "#"
        ],
        "output_file_name": "out"
    }
}
```