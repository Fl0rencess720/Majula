# Majula

Majula 是一个基于 Go 语言的 Eino 框架开发的 AI 服务，专注于提供高效的事实核查（Fact Checking）功能。

## 🚀 功能特性

- **事实核查**: 核心功能，通过 AI Agent 对输入文本进行事实核查，并返回核查结果、来源和理由。
- **高性能**: 基于 Gin 框架，提供快速的 API 响应能力。
- **MCP 集成**: 暴露标准化的 MCP 接口，方便与其他支持 MCP 的模型和工具进行交互。
- **可观测性**: 集成了日志（Zap）、链路追踪（OpenTelemetry）和性能分析（Pyroscope），便于系统监控和故障排查。同时集成 CozeLoop 以支持对于 Agent 工作的监控。


## 🛠️ 技术栈

- **后端**: Go
- **Web 框架**: Gin
- **AI Agent 框架**: CloudWeGo Eino, CozeLoop
- **日志**: Zap
- **链路追踪**: OpenTelemetry
- **性能分析**: Pyroscope
- **配置管理**: Viper
- **MCP**: gin-mcp

## 📦 安装与运行

### 前置条件

- Go 1.18+
- Docker (可选，用于部署)

### 方式一 Docker 部署（推荐）

1. **克隆仓库**:
   ```bash
   git clone https://github.com/Fl0rencess720/Majula.git
   cd Majula
   ```
2. **配置环境变量**:
   复制 `.env.example` 文件并重命名为 `.env`，然后根据您的需求修改其中的配置。
   ```bash
   cp configs/.env.example configs/.env
   ```
   请确保配置了必要的 AI 模型 API 密钥等信息。
3. **部署并运行**
   ```bash
   docker-compose up -d
   ```
4. **访问**

    服务默认运行在 `http://localhost:8080` (具体端口可根据 `configs/config.yaml`和`Dockerfile`及`docker-compose.yml`配置而定)。
### 方式二 直接运行

1. **克隆仓库**:
   ```bash
   git clone https://github.com/Fl0rencess720/Majula.git
   cd Majula
   ```

2. **配置环境变量**:
   复制 `.env.example` 文件并重命名为 `.env`，然后根据您的需求修改其中的配置。
   ```bash
   cp configs/.env.example configs/.env
   ```
   请确保配置了必要的 AI 模型 API 密钥等信息。

3. **安装依赖**:
   ```bash
   go mod tidy
   ```

4. **运行服务**:
   ```bash
   make run
   ```
   服务默认运行在 `http://localhost:8080` (具体端口可根据 `configs/config.yaml` 配置而定)。

## 💡 API 接口

### 事实核查 API

- **URL**: `/api/v1/checking`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "text": "待核查的文本内容"
  }
  ```
- **响应体**:
  ```json
  [
    {
      "result": "核查结果 (例如: 真实, 虚假, 无法核实)",
      "sources": ["来源1", "来源2"],
      "reason": "核查理由"
    }
  ]
  ```

### MCP API

- **URL**: `/mcp`
- **描述**: 暴露 MCP 服务器的元数据和工具定义。您可以通过访问此端点来了解 Majula 提供的 MCP 能力。
