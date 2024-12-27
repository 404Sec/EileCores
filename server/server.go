// server.go
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

const ChunkSize = 4 * 1024 * 1024 // 3MB

var (
	// 使用 sync.Map 来安全地在多个 goroutine 中存储和访问文件的偏移量
	fileState             sync.Map
	storageDir            = "./uploads"
	activeConnections     int64
	totalBytesTransferred int64
	serverStartTime       time.Time
	mu                    sync.Mutex
	clients               = make(map[string]*Client)
	clientsMu             sync.Mutex
	completedClients      []*Client
	completedClientsMu    sync.Mutex
)

// Client struct to track each client's transfer status
type Client struct {
	ID             string
	IP             string
	FileName       string
	FileSize       int64
	Received       int64
	Status         string
	Speed          float64
	StartTime      time.Time
	CalculatedHash string
}

// ASCII Art
const asciiArt = `
  ______ _ _        _____                     
 |  ____(_) |      / ____|                    
 | |__   _| | ___ | |     ___  _ __ ___  ___ 
 |  __| | | |/ _ \| |    / _ \| '__/ _ \/ __|
 | |____| | |  __/| |___| (_) | | |  __/\__ \
 |______|_|_|\___| \_____\___/|_|  \___||___/

                            Team：404Sec 
                           Author: WarmBrew
`

func main() {
	port := flag.String("port", "59999", "Port to listen on")
	flag.Parse()

	// Configure logging
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create storage directory
	err = os.MkdirAll(storageDir, os.ModePerm)
	if err != nil {
		log.Println("Failed to create storage directory:", err)
		return
	}

	// Initialize screen
	clearScreen()
	moveCursor(1, 1)

	// Display banner once
	displayBanner()

	// Display initial static information
	fmt.Println("\n") // Add some space after the banner

	// Start listening on IPv4
	listener, err := net.Listen("tcp4", "0.0.0.0:"+*port)
	if err != nil {
		log.Println("Error starting server:", err)
		color.Red("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()
	log.Printf("File server is listening on port %s...\n", *port)
	color.Green("File server is listening on port %s...\n", *port)

	// Initialize server start time
	serverStartTime = time.Now()

	// Start status monitor
	go monitorStatus()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientIP := conn.RemoteAddr().String()
	clientID := fmt.Sprintf("%d", time.Now().UnixNano())

	log.Printf("Client %s connected.\n", clientIP)
	fmt.Printf("Client %s connected.\n", clientIP)

	// Read file info length
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		log.Printf("Client %s: Error reading info length: %v\n", clientIP, err)
		return
	}
	infoLength := binary.BigEndian.Uint32(lengthBuf)

	// Read file info
	infoBuf := make([]byte, infoLength)
	_, err = io.ReadFull(conn, infoBuf)
	if err != nil {
		log.Printf("Client %s: Error reading file info: %v\n", clientIP, err)
		return
	}

	info := strings.Split(string(infoBuf), "|")
	if len(info) < 4 {
		log.Printf("Client %s: Received incomplete file info\n", clientIP)
		return
	}
	fileName := sanitizeFileName(info[0])
	fileSize, err := strconv.ParseInt(info[1], 10, 64)
	if err != nil {
		log.Printf("Client %s: Invalid file size: %v\n", clientIP, err)
		return
	}
	// Remove hash and resume from info, since server will compute hash
	// hash := info[2]
	resume := info[3] == "true"

	log.Printf("Client %s: File Name: %s, File Size: %d, Resume: %t\n", clientIP, fileName, fileSize, resume)

	var offset int64 = 0
	if resume {
		if val, ok := fileState.Load(fileName); ok {
			offset = val.(int64)
			if offset > fileSize {
				offset = 0 // Prevent offset from exceeding file size
			}
		}
		// Send offset back to client
		offsetStr := fmt.Sprintf("%d", offset)
		_, err = conn.Write([]byte(offsetStr))
		if err != nil {
			log.Printf("Client %s: Error sending resume offset: %v\n", clientIP, err)
			return
		}
		log.Printf("Client %s: Sent resume offset: %d\n", clientIP, offset)
	} else {
		// If not resuming, send 0 offset
		_, err = conn.Write([]byte("0"))
		if err != nil {
			log.Printf("Client %s: Error sending initial offset: %v\n", clientIP, err)
			return
		}
		log.Printf("Client %s: Sent initial offset: 0\n", clientIP)
	}

	filePath := filepath.Join(storageDir, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Client %s: Error creating/opening file: %v\n", clientIP, err)
		return
	}
	defer file.Close()

	// Seek to offset
	_, err = file.Seek(offset, 0)
	if err != nil {
		log.Printf("Client %s: Error seeking file: %v\n", clientIP, err)
		return
	}

	// Initialize client status
	client := &Client{
		ID:             clientID,
		IP:             clientIP,
		FileName:       fileName,
		FileSize:       fileSize,
		Received:       offset,
		Status:         "传输中",
		Speed:          0.0,
		StartTime:      time.Now(),
		CalculatedHash: "",
	}

	// Add client to clients map
	clientsMu.Lock()
	clients[clientID] = client
	activeConnections++
	clientsMu.Unlock()

	log.Printf("Client %s: Started transferring file %s (%d bytes)\n", clientIP, fileName, fileSize)
	fmt.Printf("Client %s: Started transferring file %s (%d bytes)\n", clientIP, fileName, fileSize)

	buf := make([]byte, ChunkSize)
	startTime := time.Now()

	for client.Received < client.FileSize {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Client %s: Error reading file chunk: %v\n", clientIP, err)
			client.Status = "传输中断"
			break
		}

		// Write to file
		_, err = file.Write(buf[:n])
		if err != nil {
			log.Printf("Client %s: Error writing to file: %v\n", clientIP, err)
			client.Status = "写入错误"
			break
		}

		client.Received += int64(n)
		mu.Lock()
		totalBytesTransferred += int64(n)
		mu.Unlock()
		fileState.Store(fileName, client.Received)

		// Calculate transfer speed
		elapsed := time.Since(startTime).Seconds()
		if elapsed > 0 {
			client.Speed = float64(n) / elapsed / (1024 * 1024) // MB/s
		}
		startTime = time.Now()
	}

	// Close the file to ensure all data is written
	file.Close()

	// Compute hash of received file
	calculatedHash, err := calculateFileHash(filePath)
	if err != nil {
		log.Printf("Client %s: Error calculating file hash: %v\n", clientIP, err)
		client.Status = "哈希计算错误"
	} else {
		client.CalculatedHash = calculatedHash
		client.Status = "传输完成"
		log.Printf("Client %s: File %s received successfully (%d bytes). Hash: %s\n", clientIP, fileName, client.Received, calculatedHash)
		fmt.Printf("Client %s: File %s received successfully (%d bytes). Hash: %s\n", clientIP, fileName, client.Received, calculatedHash)
	}

	// Move client to completedClients if transfer is completed or encountered an error
	if client.Status == "传输完成" || client.Status != "传输中" {
		completedClientsMu.Lock()
		completedClients = append(completedClients, client)
		completedClientsMu.Unlock()

		// Remove from active clients map
		clientsMu.Lock()
		delete(clients, clientID)
		activeConnections--
		clientsMu.Unlock()
	}

	log.Printf("Client %s: Connection closed.\n", clientIP)
	fmt.Printf("Client %s: Connection closed.\n", clientIP)
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func sanitizeFileName(fileName string) string {
	// Remove path, keep base file name
	baseName := filepath.Base(fileName)
	// Further remove special characters like '..'
	baseName = strings.ReplaceAll(baseName, "..", "")
	return baseName
}

func displayBanner() {
	c := color.New(color.FgCyan).Add(color.Bold)
	c.Println(asciiArt)
	c.Println("Welcome to the Enhanced File Transfer Server!")
}

// ANSI escape codes for terminal control
const (
	esc            = "\033["
	clearScreenSeq = "\033[2J"
	cursorHomeSeq  = "\033[H"
)

// clearScreen clears the entire terminal screen
func clearScreen() {
	fmt.Print(clearScreenSeq)
}

// moveCursor moves the cursor to the specified row and column
func moveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// monitorStatus periodically updates the server status on the terminal
func monitorStatus() {
	ticker := time.NewTicker(500 * time.Millisecond) // 500ms 更新频率
	defer ticker.Stop()

	// Initial position after the banner and initial static information
	// Count the number of lines in asciiArt plus additional lines
	bannerLines := strings.Count(asciiArt, "\n") + 2 // 加上欢迎信息和空行
	statusStartLine := bannerLines + 2               // Adjust based on your layout

	for range ticker.C {
		// Move cursor to status start position
		moveCursor(statusStartLine, 1)

		// Clear from the current line to the end of the screen
		fmt.Print("\033[J") // Clear from cursor to end of screen

		// Collect status information
		mu.Lock()
		conn := activeConnections
		bytesTransferred := totalBytesTransferred
		mu.Unlock()

		// Calculate transfer speed
		elapsed := time.Since(serverStartTime).Seconds()
		var speed float64
		if elapsed > 0 {
			speed = float64(bytesTransferred) / elapsed / (1024 * 1024) // MB/s
		}

		// Build main status string
		mainStatus := fmt.Sprintf("Active Connections: %d | Total Bytes Transferred: %.2f MB | Current Speed: %.2f MB/s",
			conn, float64(bytesTransferred)/(1024*1024), speed)

		fmt.Println(mainStatus)
		fmt.Println("------------------------------------------------------------")

		// Build client status strings
		clientsMu.Lock()
		completedClientsMu.Lock()
		if len(clients) == 0 && len(completedClients) == 0 {
			fmt.Println("No active clients.")
		} else {
			// Display active clients
			for _, client := range clients {
				if client.Status == "传输中" {
					status := fmt.Sprintf("Client %s: %s | File: %s | Size: %s | Received: %s | Speed: %.2f MB/s",
						client.IP, client.Status, client.FileName, formatBytes(client.FileSize), formatBytes(client.Received), client.Speed)
					fmt.Println(status)
				}
			}

			// Display completed clients
			for _, client := range completedClients {
				status := fmt.Sprintf("Client %s: %s | File: %s | Size: %s | Hash: %s",
					client.IP, client.Status, client.FileName, formatBytes(client.FileSize), client.CalculatedHash)
				fmt.Println(status)
			}
		}
		completedClientsMu.Unlock()
		clientsMu.Unlock()

	}
}

// formatBytes formats bytes as human-readable strings
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
