# MQTTZ 快速启动指南

## 概述

MQTTZ 是一个功能强大的 MQTT 客户端工具，支持多 Broker 连接管理、Mock 数据生成、条件过滤和消息转发等功能。本指南将帮助您快速上手使用 MQTTZ。

## 环境准备

### 1. 准备 MQTT Broker

本指南使用 [EMQX](https://github.com/emqx/emqx) 作为 MQTT Broker。EMQX 提供了便捷的 Docker 镜像，可以快速在本地部署服务。

#### 使用 Docker 启动 EMQX

```bash
docker run -d \
  --name emqx \
  -p 1883:1883 \
  -p 8083:8083 \
  -p 8084:8084 \
  -p 8883:8883 \
  -p 18083:18083 \
  emqx/emqx:latest
```

启动后，您可以通过以下方式访问：
- **MQTT 服务端口**: 1883 (TCP)
- **Web 管理界面**: http://localhost:18083 (用户名/密码: admin/public)

## 获取 MQTTZ

### 方式一：下载预编译版本

访问 [GitHub Releases](https://github.com/RedCrazyGhost/MQTTZ/releases) 页面，下载适合您系统的预编译版本。

### 方式二：从源码构建

1. **克隆仓库**
   ```bash
   git clone https://github.com/RedCrazyGhost/MQTTZ.git
   cd MQTTZ
   ```

2. **构建项目**
   ```bash
   make build
   ```
   
   构建完成后，可执行文件将生成在项目根目录下。

## 配置与运行

### 1. 配置文件准备

MQTTZ 默认读取 `./conf/config.yaml` 配置文件。您可以：

- 复制模板配置文件：
  ```bash
  cp template/conf/config.yaml conf/config.yaml
  ```

- 或创建自定义配置文件：
  ```bash
  mkdir -p conf
  # 编辑 conf/config.yaml 文件
  ```

### 2. 基础配置示例

```yaml
log:
  enable_debug: true
  enable_color: true

mqtt_configs:
  - broker: 127.0.0.1
    port: 1883
    client_id: mqtt_client_1
    username: "emqx"
    password: "emqx2024"
    nickname: "MQTTZ_1"
```

### 3. 启动 MQTTZ

配置完成后，执行以下命令启动：

```bash
./MQTTZ
```

### 常见问题

1. **连接失败**: 确保 EMQX Broker 已正确启动，检查端口是否被占用
2. **配置文件错误**: 检查 YAML 语法，确保缩进正确
3. **权限问题**: 确保对配置文件和日志目录有读写权限

### 获取帮助

- 查看项目 [README](../readme.md) 获取详细文档
- 提交 [Issue](https://github.com/RedCrazyGhost/MQTTZ/issues) 报告问题