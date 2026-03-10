# 🚀 EileCores

> 🔥 Enhanced file transfer server with compression, breakpoint resume, and integrity verification

[![License](https://img.shields.io/github/license/404Sec/EileCores)](https://github.com/404Sec/EileCores/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/dl/)
[![GitHub issues](https://img.shields.io/github/issues/404Sec/EileCores)](https://github.com/404Sec/EileCores/issues)
[![GitHub stars](https://img.shields.io/github/stars/404Sec/EileCores)](https://github.com/404Sec/EileCores/stargazers)

---

**🌍 Read this in other languages:** [English](README.md) | [中文](README_zh.md)

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Directory Structure](#-directory-structure)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Build](#-build)
- [Usage](#-usage)
- [Features](#-features)
- [Logging](#-logging)
- [File Storage](#-file-storage)
- [Important Notes](#-important-notes)

---

## 📋 Overview

**EileCores** is an enhanced file transfer server and client developed in Go. It supports:
- 📦 File compression
- 🔄 Breakpoint resume
- 🔐 File integrity verification

---

## 📁 Directory Structure

```
EileCores/
├── server/
│   ├── server.go          # Server-side code
│   ├── go.mod             # Go module definition
│   └── go.sum             # Dependency checksums
├── client/
│   ├── client.go          # Client-side code
│   ├── go.mod             # Go module definition
│   └── go.sum             # Dependency checksums
├── uploads/               # Directory for storing uploaded files
├── server.log             # Server log file
├── README.md              # English documentation
├── README_zh.md           # Chinese documentation
└── LICENSE                # MIT License
```

---

## ⚙️ Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher

---

## 📦 Installation

Both server and client use the `github.com/fatih/color` library for colored terminal output. Install the dependency before building:

```bash
go get github.com/fatih/color
```

---

## 🔨 Build

### Server

Navigate to the server directory and build:

```bash
cd server
go build -o server server.go
```

### Client

Navigate to the client directory and build:

```bash
cd ../client
go build -o client client.go
```

---

## 🚀 Usage

### Start Server

```bash
./server -port=59999
```

**Parameters:**
- `-port`: Server listening port (default: `59999`)

After starting, the server will display welcome information and current transfer status, including:
- Active connections
- Total bytes transferred
- Current transfer speed

### Transfer Files with Client

#### Transfer Single File

```bash
./client -file=/path/to/file -ip=server_ip:59999
```

**Parameters:**
- `-file`: Path to the file to transfer
- `-ip`: Server IP and port (default: `localhost:59999`)

#### Compress and Transfer Directory

```bash
./client -path=/path/to/directory -output=output.zip -ip=server_ip:59999
```

**Parameters:**
- `-path`: Directory path to compress
- `-output`: Output ZIP filename (default: `<directory_name>.zip`)
- `-ip`: Server IP and port (default: `localhost:59999`)

---

## ✨ Features

### 🔄 Breakpoint Resume

The client automatically detects if breakpoint resume is needed. If the transfer is interrupted, running the same transfer command again will resume from the last interrupted position.

### 📝 Logging

The server generates a `server.log` file with detailed log information including:
- Client connections
- File transfer status
- Error messages

### 💾 File Storage

Files received by the server are stored in the `./uploads` directory. Ensure this directory exists in the server's running directory, or modify the `storageDir` variable in the code to specify a different storage path.

---

## ⚠️ Important Notes

### 🌐 Network Stability
For stable transfers, it's recommended to transfer large files in a stable network environment.

### 📂 File Permissions
Ensure the server has permission to create and write files in the specified storage directory.

### 🔒 Security
The current version does not implement authentication or encrypted transmission. Consider adding security measures for production environments.

### ❌ Error Handling
The client will attempt multiple retries during transfer. If all retries fail, please check the network connection and server status.

---

## 📄 License

MIT License © 2024-2026 [404Sec](https://github.com/404Sec)

---

## 🔗 Related Links

- [GitHub Repository](https://github.com/404Sec/EileCores)
- [Submit Issue](https://github.com/404Sec/EileCores/issues)
- [Go Language](https://golang.org/)
- [fatih/color Library](https://github.com/fatih/color)

---

<p align="center">
  <b>If this project helped you, please give it a ⭐ Star!</b><br>
  <sub>Built with ❤️ by 404Sec Team</sub>
</p>
