# PetWell Backend (Go Version)

PetWell App 的后端服务仓库。本项目计划使用 **Go (Golang)** 重构，以提供更高性能、并发更强的后端服务。目前作为 iOS 客户端的静态数据源和未来业务逻辑的承载体。

## 📁 规划目录结构

```
Backend/
├── main.go           # Go 程序入口 (HTTP Server)
├── go.mod            # Go 模块定义
├── handlers/         # 路由处理逻辑 (Controllers)
├── models/           # 数据模型结构体
├── vaccines.json     # 静态数据源 (暂存)
└── README.md         # 说明文档
```

## 🚀 快速开始

### 1. 环境准备
- 请确保已安装 [Go](https://go.dev/dl/) (建议 1.21+ 版本)

### 2. 初始化项目 (如果是首次运行)
在终端中进入 `Backend` 目录：

```bash
# 初始化 Go module
go mod init petwell-backend

# 如果使用 Gin 框架 (推荐)，安装依赖
go get -u github.com/gin-gonic/gin
```

### 3. 启动服务
编写好 `main.go` 后，运行：

```bash
go run main.go
```

服务默认建议运行在 `http://localhost:8080` (Go 常用端口)。

### 4. API 接口说明

#### 获取疫苗列表
- **Endpoint**: `/vaccines`
- **Method**: `GET`
- **Response**: JSON 数组，结构与 iOS 端 `Vaccine` 模型一致。

**测试命令**:
```bash
curl http://localhost:8080/vaccines
```

## 📱 前端适配指南

由于 Go 服务通常默认使用 8080 端口（Python 默认 8000），请注意更新 iOS 前端配置。

### 1. 疫苗服务 (VaccineService)
前端文件: `Services/VaccineService.swift`

修改 `urlString`：
```swift
// 注意端口变化：8000 -> 8080
private let urlString = "http://localhost:8080/vaccines" 
// 真机调试请使用局域网 IP
```

### 2. AI 服务 (AIService) - *规划中*
前端文件: `Services/AIService.swift`

计划在 Go 后端实现 AI 代理路由：
- **目标**: 隐藏 OpenAI/DeepSeek API Key。
- **实现**: 使用 Go 的高并发特性处理流式响应 (Server-Sent Events) 转发给客户端。
- **路径**: `/ai/chat` (POST)

### 3. 保险服务 (InsuranceService) - *规划中*
前端文件: `Services/InsuranceService.swift`

- 将 SQLite 数据迁移至服务端数据库 (PostgreSQL/MySQL)。
- Go 后端提供 RESTful API 查询保险产品。

## 🛠 开发建议

1.  **Web 框架**: 推荐使用 **Gin** 或 **Echo**，它们轻量且性能极佳。
2.  **配置管理**: 使用 **Viper** 处理配置文件 (加载 API Keys 等)。
3.  **数据库**: 推荐使用 **GORM** 或 **sqlc** 进行数据库操作。


1. **修改数据**: 直接编辑 `vaccines.json` 即可更新疫苗列表，无需重启服务（每次请求都会重新读取文件）。
2. **扩展功能**: 建议未来引入 Flask 或 FastAPI 框架来替代当前的 `http.server`，以支持更复杂的路由和数据库连接。
