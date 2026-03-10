# 🚀 EileCores

> 🔥 增强型文件传输服务器，支持压缩、断点续传和完整性校验

[![License](https://img.shields.io/github/license/404Sec/EileCores)](https://github.com/404Sec/EileCores/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/dl/)
[![GitHub issues](https://img.shields.io/github/issues/404Sec/EileCores)](https://github.com/404Sec/EileCores/issues)
[![GitHub stars](https://img.shields.io/github/stars/404Sec/EileCores)](https://github.com/404Sec/EileCores/stargazers)

---

**🌍 其他语言版本:** [English](README.md) | [中文](README_zh.md)

---

## 📖 目录

- [项目概述](#-项目概述)
- [目录结构](#-目录结构)
- [前提条件](#-前提条件)
- [安装依赖](#-安装依赖)
- [编译](#-编译)
- [使用方法](#-使用方法)
- [功能特性](#-功能特性)
- [日志](#-日志)
- [文件存储](#-文件存储)
- [注意事项](#-注意事项)

---

## 📋 项目概述

**EileCores** 是一个基于 Go 语言开发的增强型文件传输服务器和客户端，支持：
- 📦 文件压缩
- 🔄 断点续传
- 🔐 文件完整性校验

---

## 📁 目录结构

```
EileCores/
├── server/
│   ├── server.go          # 服务器端代码
│   ├── go.mod             # Go 模块定义
│   └── go.sum             # 依赖校验和
├── client/
│   ├── client.go          # 客户端代码
│   ├── go.mod             # Go 模块定义
│   └── go.sum             # 依赖校验和
├── uploads/               # 服务器接收并存储上传文件的目录
├── server.log             # 服务器日志文件
├── README.md              # 英文文档
├── README_zh.md           # 中文文档
└── LICENSE                # MIT 许可证
```

---

## ⚙️ 前提条件

- [Go](https://golang.org/dl/) 1.16 或更高版本

---

## 📦 安装依赖

服务器端和客户端均使用 `github.com/fatih/color` 库来实现终端输出的彩色显示。请在编译前安装该依赖：

```bash
go get github.com/fatih/color
```

---

## 🔨 编译

### 服务器端

进入服务器代码所在目录，运行以下命令编译服务器：

```bash
cd server
go build -o server server.go
```

### 客户端

进入客户端代码所在目录，运行以下命令编译客户端：

```bash
cd ../client
go build -o client client.go
```

---

## 🚀 使用方法

### 启动服务器

```bash
./server -port=59999
```

**参数说明：**
- `-port`: 指定服务器监听的端口，默认为 `59999`

服务器启动后，会在终端显示欢迎信息和当前的传输状态，包括：
- 活跃连接数
- 传输的总字节数
- 当前速度

### 使用客户端传输文件

#### 传输单个文件

```bash
./client -file=/path/to/file -ip=server_ip:59999
```

**参数说明：**
- `-file`: 指定要传输的文件路径
- `-ip`: 指定服务器的 IP 地址和端口，默认为 `localhost:59999`

#### 压缩并传输目录

```bash
./client -path=/path/to/directory -output=output.zip -ip=server_ip:59999
```

**参数说明：**
- `-path`: 指定要压缩的目录路径
- `-output`: 指定压缩后的 ZIP 文件名。如果不指定，默认使用目录名加 `.zip` 后缀
- `-ip`: 指定服务器的 IP 地址和端口，默认为 `localhost:59999`

---

## ✨ 功能特性

### 🔄 断点续传

客户端在传输过程中会自动检测是否需要断点续传。如果传输中断，再次运行相同的传输命令，客户端会从上次中断的位置继续传输。

### 📝 日志

服务器运行过程中会生成 `server.log` 文件，记录详细的日志信息，包括：
- 客户端连接
- 文件传输状态
- 错误信息

### 💾 文件存储

服务器接收到的文件会存储在 `./uploads` 目录下。请确保服务器运行目录下存在该目录，或者在代码中修改 `storageDir` 变量指定其他存储路径。

---

## ⚠️ 注意事项

### 🌐 网络稳定性
为了确保传输的稳定性，建议在网络较为稳定的环境下进行大文件传输。

### 📂 文件权限
确保服务器有权限在指定的存储目录中创建和写入文件。

### 🔒 安全性
当前版本未实现身份验证和加密传输，请在生产环境中考虑增加安全措施。

### ❌ 错误处理
客户端在传输过程中会尝试多次重试，如果所有重试均失败，请检查网络连接和服务器状态。

---

## 📄 许可证

MIT License © 2024-2026 [404Sec](https://github.com/404Sec)

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
