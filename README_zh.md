# 🚀 EileCores

> 🔥 高性能文件传输服务器，支持断点续传和完整性校验

[![License](https://img.shields.io/github/license/404Sec/EileCores)](https://github.com/404Sec/EileCores/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/dl/)
[![GitHub issues](https://img.shields.io/github/issues/404Sec/EileCores)](https://github.com/404Sec/EileCores/issues)
[![GitHub stars](https://img.shields.io/github/stars/404Sec/EileCores)](https://github.com/404Sec/EileCores/stargazers)

---

**🌍 其他语言版本:** [English](README.md) | [中文](README_zh.md)

---

## 📖 目录

- [项目概述](#-项目概述)
- [功能特性](#-功能特性)
- [系统架构](#-系统架构)
- [快速开始](#-快速开始)
- [使用指南](#-使用指南)
- [配置说明](#-配置说明)
- [技术细节](#-技术细节)
- [常见问题](#-常见问题)
- [贡献指南](#-贡献指南)
- [许可证](#-许可证)

---

## 📋 项目概述

**EileCores** 是一个使用 Go 语言开发的高性能文件传输服务器和客户端系统。支持文件压缩、断点续传和文件完整性校验，非常适合在网络不稳定的环境下进行大文件可靠传输。

### 💡 应用场景

- 🖥️ **服务器备份** - 将备份文件传输到远程服务器
- 📊 **数据迁移** - 在服务器之间迁移大型数据集
- 🔄 **文件同步** - 跨多个位置同步文件
- 📦 **文件分发** - 向多个客户端分发文件

---

## ✨ 功能特性

| 功能 | 描述 | 状态 |
|------|------|------|
| 🔄 **断点续传** | 从中断位置恢复传输 | ✅ |
| 🔐 **完整性校验** | SHA-256 哈希验证文件完整性 | ✅ |
| 📦 **目录压缩** | 自动将目录压缩为 ZIP 后传输 | ✅ |
| 📊 **实时统计** | 实时显示传输速度和进度 | ✅ |
| 🎨 **彩色输出** | 彩色终端输出，更易阅读 | ✅ |
| 🔁 **自动重试** | 传输失败自动重试（最多 5 次） | ✅ |
| 📝 **详细日志** | 完整的 server.log 日志记录 | ✅ |
| 🚀 **高性能** | 4MB 分块大小，优化吞吐量 | ✅ |
| 🔗 **多客户端** | 支持多个并发客户端连接 | ✅ |
| 📈 **连接追踪** | 监控活动连接和传输状态 | ✅ |

---

## 🏗️ 系统架构

```
┌─────────────────┐         ┌─────────────────┐
│    客户端        │         │    服务器       │
│  (client.go)    │         │  (server.go)    │
│                 │         │                 │
│  - 压缩         │────►    │  - 接收         │
│  - 分块         │         │  - 重组         │
│  - 传输         │         │  - 验证         │
│  - 校验         │◄────    │  - 存储         │
└─────────────────┘         └─────────────────┘
         │                          │
         │      网络 (TCP)          │
         └──────────────────────────┘
```

### 目录结构

```
EileCores/
├── server/
│   ├── server.go          # 服务器实现
│   ├── go.mod             # Go 模块定义
│   └── go.sum             # 依赖校验和
├── client/
│   ├── client.go          # 客户端实现
│   ├── go.mod             # Go 模块定义
│   └── go.sum             # 依赖校验和
├── uploads/               # 服务器存储目录（自动创建）
├── server.log             # 服务器日志文件（自动生成）
├── README.md              # 英文文档
├── README_zh.md           # 中文文档
└── LICENSE                # MIT 许可证
```

---

## 🚀 快速开始

### 环境要求

- [Go](https://golang.org/dl/) 1.16 或更高版本
- Git（用于克隆仓库）

### 安装

```bash
# 克隆仓库
git clone https://github.com/404Sec/EileCores.git
cd EileCores

# 安装依赖
go get github.com/fatih/color
```

### 编译

```bash
# 编译服务器
cd server
go build -o server server.go

# 编译客户端
cd ../client
go build -o client client.go
```

### 运行

**终端 1 - 启动服务器：**
```bash
cd server
./server -port=59999
```

**终端 2 - 传输文件：**
```bash
cd client
./client -file=/path/to/your/file.txt -ip=192.168.1.100:59999
```

---

## 📖 使用指南

### 服务器选项

```bash
./server -port=59999
```

| 参数 | 默认值 | 描述 |
|------|--------|------|
| `-port` | `59999` | 服务器监听端口 |

**服务器输出示例：**
```
╔══════════════════════════════════════════════════╗
║     EileCores 文件传输服务器 v1.4                ║
╚══════════════════════════════════════════════════╝
服务器启动在端口：59999
存储目录：./uploads
等待连接...
```

### 客户端选项

#### 传输单个文件

```bash
./client -file=/path/to/file.zip -ip=192.168.1.100:59999
```

| 参数 | 默认值 | 描述 |
|------|--------|------|
| `-file` | - | 要传输的文件路径 |
| `-ip` | `localhost:59999` | 服务器 IP 和端口 |

#### 压缩并传输目录

```bash
./client -path=/path/to/directory -output=backup.zip -ip=192.168.1.100:59999
```

| 参数 | 默认值 | 描述 |
|------|--------|------|
| `-path` | - | 要压缩的目录路径 |
| `-output` | `<目录名>.zip` | 输出的 ZIP 文件名 |
| `-ip` | `localhost:59999` | 服务器 IP 和端口 |

### 客户端输出示例

```
╔══════════════════════════════════════════════════╗
║      EileCores 文件传输客户端 v1.4               ║
╚══════════════════════════════════════════════════╝
连接到服务器：192.168.1.100:59999
文件：backup.zip (1.2 GB)
分块大小：4 MB
开始传输...

[████████████████████████████████] 100% | 1.2 GB / 1.2 GB
速度：45.6 MB/s | 预计剩余时间：0s
传输完成！
文件哈希：a3f5b891c2d4e6f7...
```

---

## ⚙️ 配置说明

### 服务器配置

服务器可以通过命令行标志进行配置：

```bash
./server -port=8080
```

**存储目录：**
- 默认：`./uploads`
- 修改方法：在 `server.go` 中修改 `storageDir = "/custom/path"`

### 客户端配置

**分块大小：**
- 默认：`4 MB` (4 * 1024 * 1024 字节)
- 修改方法：在 `client.go` 中修改 `const ChunkSize = 4 * 1024 * 1024`

**重试设置：**
- 最大重试次数：`5`
- 重试间隔：`2 秒`
- 修改方法：在 `client.go` 中：
  ```go
  const MaxRetries = 5
  const RetryInterval = 2 * time.Second
  ```

---

## 🔬 技术细节

### 断点续传机制

1. **状态追踪**：服务器使用 `sync.Map` 追踪文件传输进度
2. **分块传输**：文件被分割为 4MB 的数据块
3. **偏移量管理**：记录每个数据块的偏移量
4. **恢复逻辑**：重新连接时，客户端请求最后已知的偏移量

```go
// 服务器端状态管理
var fileState sync.Map
fileState.Store(filename, offset)
```

### 完整性校验

**SHA-256 哈希计算：**

```go
// 客户端计算哈希
hash := sha256.New()
io.Copy(hash, file)
calculatedHash := hex.EncodeToString(hash.Sum(nil))

// 服务器验证哈希
if receivedHash != calculatedHash {
    // 处理验证失败
}
```

### 传输协议

```
1. 客户端连接到服务器
2. 发送文件元数据（名称、大小、哈希）
3. 服务器检查是否存在部分传输
4. 客户端发送带偏移量的数据块
5. 服务器确认每个数据块
6. 客户端验证完成
7. 服务器发送最终哈希进行验证
```

---

## ❓ 常见问题

### Q: 断点续传是如何工作的？

**A:** 服务器会追踪每个文件已接收的字节偏移量。如果传输中断，只需再次运行相同的命令，客户端会向服务器请求最后已知的偏移量以继续传输。

### Q: 如果哈希验证失败会怎样？

**A:** 传输将被标记为失败，文件不会被保存。您可以重试传输，断点续传功能会帮助您从失败的位置继续。

### Q: 可以同时传输多个文件吗？

**A:** 可以，服务器支持多个并发连接。每个客户端连接都被独立追踪。

### Q: 文件存储在哪里？

**A:** 默认情况下，文件存储在服务器运行目录下的 `./uploads` 目录中。您可以在 `server.go` 中修改此路径。

### Q: 传输是加密的吗？

**A:** 不，当前版本使用纯 TCP 传输。对于生产环境，建议通过 SSH 隧道运行或添加 TLS 支持。

### Q: 可以传输多大的文件？

**A:** 没有硬性限制。基于分块的方法允许传输任意大小的文件。已测试传输高达 100GB+ 的文件。

### Q: 如果服务器在传输过程中崩溃怎么办？

**A:** 重启服务器并再次运行客户端命令。断点续传功能会从最后一个成功的数据块继续。

---

## 🤝 贡献指南

欢迎贡献！

### 如何贡献

1. **Fork** 本仓库
2. **创建** 特性分支 (`git checkout -b feature/AmazingFeature`)
3. **提交** 更改 (`git commit -m 'Add some AmazingFeature'`)
4. **推送** 到分支 (`git push origin feature/AmazingFeature`)
5. **开启** Pull Request

### 开发指南

- 遵循 Go 最佳实践
- 为复杂逻辑添加注释
- 充分测试断点续传功能
- 为新功能更新文档

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

## 🔗 相关链接

- [GitHub 仓库](https://github.com/404Sec/EileCores)
- [提交 Issue](https://github.com/404Sec/EileCores/issues)
- [Go 语言官网](https://golang.org/)
- [fatih/color 库](https://github.com/fatih/color)

---

<p align="center">
  <b>如果这个项目对您有帮助，请给个 ⭐ Star！</b><br>
  <sub>由 404Sec 团队用 ❤️ 构建</sub>
</p>
