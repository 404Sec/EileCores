# 🚀 EileCores

> 🔥 High-performance file transfer server with resume capability and integrity verification

[![License](https://img.shields.io/github/license/404Sec/EileCores)](https://github.com/404Sec/EileCores/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)](https://golang.org/dl/)
[![GitHub issues](https://img.shields.io/github/issues/404Sec/EileCores)](https://github.com/404Sec/EileCores/issues)
[![GitHub stars](https://img.shields.io/github/stars/404Sec/EileCores)](https://github.com/404Sec/EileCores/stargazers)

---

**🌍 Read this in other languages:** [English](README.md) | [中文](README_zh.md)

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [Usage Guide](#-usage-guide)
- [Configuration](#-configuration)
- [Technical Details](#-technical-details)
- [FAQ](#-faq)
- [Contributing](#-contributing)
- [License](#-license)

---

## 📋 Overview

**EileCores** is an enhanced file transfer server and client system developed in Go. It supports file compression, breakpoint resume, and file integrity verification, making it ideal for reliable large file transfers over unstable networks.

### 💡 Use Cases

- 🖥️ **Server Backup** - Transfer backup files to remote servers
- 📊 **Data Migration** - Migrate large datasets between servers
- 🔄 **File Synchronization** - Sync files across multiple locations
- 📦 **Distribution** - Distribute files to multiple clients

---

## ✨ Features

| Feature | Description | Status |
|---------|-------------|--------|
| 🔄 **Breakpoint Resume** | Resume interrupted transfers from last position | ✅ |
| 🔐 **Integrity Verification** | SHA-256 hash verification for file integrity | ✅ |
| 📦 **Directory Compression** | Auto-compress directories to ZIP before transfer | ✅ |
| 📊 **Real-time Statistics** | Live transfer speed and progress display | ✅ |
| 🎨 **Colorful Output** | Color-coded terminal output for better readability | ✅ |
| 🔁 **Auto Retry** | Automatic retry on transfer failures (up to 5 times) | ✅ |
| 📝 **Detailed Logging** | Comprehensive logging to server.log | ✅ |
| 🚀 **High Performance** | 4MB chunk size for optimal throughput | ✅ |
| 🔗 **Multi-client** | Support multiple concurrent client connections | ✅ |
| 📈 **Connection Tracking** | Monitor active connections and transfer status | ✅ |

---

## 🏗️ Architecture

```
┌─────────────────┐         ┌─────────────────┐
│    Client       │         │    Server       │
│  (client.go)    │         │  (server.go)    │
│                 │         │                 │
│  - Compress     │────►    │  - Receive      │
│  - Split        │         │  - Reassemble   │
│  - Transfer     │         │  - Verify       │
│  - Verify       │◄────    │  - Store        │
└─────────────────┘         └─────────────────┘
         │                          │
         │      Network (TCP)       │
         └──────────────────────────┘
```

### Directory Structure

```
EileCores/
├── server/
│   ├── server.go          # Server implementation
│   ├── go.mod             # Go module definition
│   └── go.sum             # Dependency checksums
├── client/
│   ├── client.go          # Client implementation
│   ├── go.mod             # Go module definition
│   └── go.sum             # Dependency checksums
├── uploads/               # Server storage directory (auto-created)
├── server.log             # Server log file (auto-generated)
├── README.md              # English documentation (default)
├── README_zh.md           # Chinese documentation
├── readme.md              # Original Chinese documentation (legacy)
└── LICENSE                # MIT License
```

---

## 🚀 Quick Start

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher
- Git (for cloning the repository)

### Installation

```bash
# Clone the repository
git clone https://github.com/404Sec/EileCores.git
cd EileCores

# Install dependencies
go get github.com/fatih/color
```

### Build

```bash
# Build server
cd server
go build -o server server.go

# Build client
cd ../client
go build -o client client.go
```

### Run

**Terminal 1 - Start Server:**
```bash
cd server
./server -port=59999
```

**Terminal 2 - Transfer File:**
```bash
cd client
./client -file=/path/to/your/file.txt -ip=192.168.1.100:59999
```

---

## 📖 Usage Guide

### Server Options

```bash
./server -port=59999
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `-port` | `59999` | Server listening port |

**Server Output Example:**
```
╔══════════════════════════════════════════════════╗
║     EileCores File Transfer Server v1.4          ║
╚══════════════════════════════════════════════════╝
Server started on port: 59999
Storage directory: ./uploads
Waiting for connections...
```

### Client Options

#### Transfer Single File

```bash
./client -file=/path/to/file.zip -ip=192.168.1.100:59999
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `-file` | - | File path to transfer |
| `-ip` | `localhost:59999` | Server IP and port |

#### Compress and Transfer Directory

```bash
./client -path=/path/to/directory -output=backup.zip -ip=192.168.1.100:59999
```

| Parameter | Default | Description |
|-----------|---------|-------------|
| `-path` | - | Directory path to compress |
| `-output` | `<dirname>.zip` | Output ZIP filename |
| `-ip` | `localhost:59999` | Server IP and port |

### Client Output Example

```
╔══════════════════════════════════════════════════╗
║      EileCores File Transfer Client v1.4         ║
╚══════════════════════════════════════════════════╝
Connecting to server: 192.168.1.100:59999
File: backup.zip (1.2 GB)
Chunk Size: 4 MB
Starting transfer...

[████████████████████████████████] 100% | 1.2 GB / 1.2 GB
Speed: 45.6 MB/s | ETA: 0s
Transfer completed successfully!
File Hash: a3f5b891c2d4e6f7...
```

---

## ⚙️ Configuration

### Server Configuration

The server can be configured via command-line flags:

```bash
./server -port=8080
```

**Storage Directory:**
- Default: `./uploads`
- Modify in `server.go`: `storageDir = "/custom/path"`

### Client Configuration

**Chunk Size:**
- Default: `4 MB` (4 * 1024 * 1024 bytes)
- Modify in `client.go`: `const ChunkSize = 4 * 1024 * 1024`

**Retry Settings:**
- Max Retries: `5`
- Retry Interval: `2 seconds`
- Modify in `client.go`:
  ```go
  const MaxRetries = 5
  const RetryInterval = 2 * time.Second
  ```

---

## 🔬 Technical Details

### Breakpoint Resume Mechanism

1. **State Tracking**: Server uses `sync.Map` to track file transfer progress
2. **Chunk-based Transfer**: Files are split into 4MB chunks
3. **Offset Management**: Each chunk's offset is recorded
4. **Resume Logic**: On reconnection, client requests last known offset from server

```go
// Server-side state management
var fileState sync.Map
fileState.Store(filename, offset)
```

### Integrity Verification

**SHA-256 Hash Calculation:**

```go
// Client calculates hash
hash := sha256.New()
io.Copy(hash, file)
calculatedHash := hex.EncodeToString(hash.Sum(nil))

// Server verifies hash
if receivedHash != calculatedHash {
    // Handle verification failure
}
```

### Transfer Protocol

```
1. Client connects to server
2. Send file metadata (name, size, hash)
3. Server checks for existing partial transfer
4. Client sends chunks with offset
5. Server acknowledges each chunk
6. Client verifies completion
7. Server sends final hash for verification
```

---

## ❓ FAQ

### Q: How does breakpoint resume work?

**A:** The server tracks the received byte offset for each file. If the transfer is interrupted, simply run the same command again, and the client will request the last known offset from the server to resume.

### Q: What happens if the hash verification fails?

**A:** The transfer will be marked as failed, and the file will not be saved. You can retry the transfer, and the breakpoint resume feature will help you continue from where it failed.

### Q: Can I transfer multiple files simultaneously?

**A:** Yes, the server supports multiple concurrent connections. Each client connection is tracked independently.

### Q: Where are the files stored?

**A:** By default, files are stored in the `./uploads` directory relative to where the server is running. You can modify this in `server.go`.

### Q: Is the transfer encrypted?

**A:** No, the current version uses plain TCP. For production use, consider running over SSH tunnel or adding TLS support.

### Q: How large files can I transfer?

**A:** There's no hard limit. The chunk-based approach allows transferring files of any size. Tested with files up to 100GB+.

### Q: What if the server crashes during transfer?

**A:** Restart the server and run the client command again. The breakpoint resume feature will continue from the last successful chunk.

---

## 🤝 Contributing

Contributions are welcome!

### How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/AmazingFeature`)
3. **Commit** your changes (`git commit -m 'Add some AmazingFeature'`)
4. **Push** to the branch (`git push origin feature/AmazingFeature`)
5. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices
- Add comments for complex logic
- Test breakpoint resume functionality thoroughly
- Update documentation for new features

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

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
