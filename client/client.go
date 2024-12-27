//1.4 最终版本
package main

import (
    "archive/zip"
    "crypto/sha256"
    "encoding/binary"
    "encoding/hex"
    "flag"
    "fmt"
    "io"
    "net"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    ChunkSize     = 4 * 1024 * 1024
    MaxRetries    = 5
    RetryInterval = 2 * time.Second
)

func main() {
    zipPath := flag.String("path", "", "指定目录压缩成zip文件")
    output := flag.String("output", "", "指定压缩后的文件名")
    filePath := flag.String("file", "", "指定传输的文件")
    serverAddr := flag.String("ip", "localhost:59999", "指定服务器接收的地址")
    flag.Parse()

    var finalFilePath string

    if *zipPath != "" {
        zipFileName, err := compressDirectory(*zipPath, *output)
        if err != nil {
            fmt.Printf("Failed to compress directory: %v\n", err)
            return
        }
        fmt.Println("Directory compressed to:", zipFileName)
        finalFilePath = zipFileName
    }

    if *filePath != "" {
        finalFilePath = *filePath
    }

    if finalFilePath == "" {
        fmt.Println("No file specified for transfer.")
        return
    }

    err := transferFileWithRetry(*serverAddr, finalFilePath)
    if err != nil {
        fmt.Printf("Failed to transfer file: %v\n", err)
        return
    }

    fmt.Println("File transfer completed successfully.")
}

func compressDirectory(dirPath, outputFileName string) (string, error) {
    if outputFileName == "" {
        outputFileName = filepath.Base(dirPath) + ".zip"
    }
    zipFile, err := os.Create(outputFileName)
    if err != nil {
        return "", err
    }
    defer zipFile.Close()

    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        relPath, err := filepath.Rel(filepath.Dir(dirPath), path)
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        writer, err := zipWriter.Create(relPath)
        if err != nil {
            return err
        }
        _, err = io.Copy(writer, file)
        return err
    })

    if err != nil {
        return "", err
    }
    return outputFileName, nil
}

func transferFileWithRetry(serverAddr, filePath string) error {
    for attempt := 1; attempt <= MaxRetries; attempt++ {
        err := transferFile(serverAddr, filePath)
        if err == nil {
            return nil
        }
        fmt.Printf("Attempt %d/%d failed: %v\n", attempt, MaxRetries, err)
        if attempt < MaxRetries {
            fmt.Println("Retrying...")
            time.Sleep(RetryInterval)
        }
    }
    return fmt.Errorf("all %d attempts failed", MaxRetries)
}

func transferFile(serverAddr, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    fileName := filepath.Base(filePath)
    fileSize, err := getFileSize(filePath)
    if err != nil {
        return fmt.Errorf("failed to get file size: %w", err)
    }

    hash, err := calculateFileHash(filePath)
    if err != nil {
        return fmt.Errorf("failed to calculate file hash: %w", err)
    }

    var offset int64 = 0
    resume := true

    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        fmt.Printf("Connection failed: %v\n", err)
        return fmt.Errorf("error connecting to server: %w", err)
    }
    defer conn.Close()

    fmt.Println("Connection successful.")

    info := fmt.Sprintf("%s|%d|%s|%t", fileName, fileSize, hash, resume)
    infoLength := uint32(len(info))
    lengthBuf := make([]byte, 4)
    binary.BigEndian.PutUint32(lengthBuf, infoLength)

    _, err = conn.Write(lengthBuf)
    if err != nil {
        return fmt.Errorf("failed to send info length: %w", err)
    }
    _, err = conn.Write([]byte(info))
    if err != nil {
        return fmt.Errorf("failed to send file info: %w", err)
    }

    offsetBuf := make([]byte, 256)
    n, err := conn.Read(offsetBuf)
    if err != nil {
        return fmt.Errorf("failed to read resume offset: %w", err)
    }
    offsetStr := strings.TrimSpace(string(offsetBuf[:n]))
    offset, err = strconv.ParseInt(offsetStr, 10, 64)
    if err != nil {
        return fmt.Errorf("invalid resume offset: %w", err)
    }

    if offset > fileSize {
        offset = 0
    }

    _, err = file.Seek(offset, 0)
    if err != nil {
        return fmt.Errorf("failed to seek file: %w", err)
    }

    fmt.Println("Transfer started.")

    buf := make([]byte, ChunkSize)
    for {
        n, err := file.Read(buf)
        if err != nil {
            if err == io.EOF {
                break
            }
            return fmt.Errorf("failed to read from file: %w", err)
        }

        _, err = conn.Write(buf[:n])
        if err != nil {
            return fmt.Errorf("failed to send data: %w", err)
        }
    }

    _, err = conn.Write([]byte(hash))
    if err != nil {
        return fmt.Errorf("failed to send file hash: %w", err)
    }

    return nil
}

func getFileSize(filePath string) (int64, error) {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        return 0, err
    }
    return fileInfo.Size(), nil
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
