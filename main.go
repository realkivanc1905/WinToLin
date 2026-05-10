package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// progressTracker wraps an io.Reader to print progress to the terminal
type progressTracker struct {
	io.Reader
	totalBytes int64
	readBytes  int64
	lastUpdate time.Time
}

func (pt *progressTracker) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.readBytes += int64(n)

	now := time.Now()
	// Update progress every 100ms or when finished
	if now.Sub(pt.lastUpdate) > 100*time.Millisecond || pt.readBytes == pt.totalBytes {
		percent := float64(pt.readBytes) / float64(pt.totalBytes) * 100
		fmt.Printf("\rProgress: %.2f%% (%d / %d bytes)", percent, pt.readBytes, pt.totalBytes)
		pt.lastUpdate = now
	}

	return n, err
}

// getPublicIP fetches the external IP address of the machine
func getPublicIP() string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

// getBestLocalIP finds the primary local IP used for internet traffic
func getBestLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// Fallback: search for first non-loopback IPv4 address
		addrs, _ := net.InterfaceAddrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func sendFile(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		fmt.Println("Folder transfer is not supported yet. Please zip the folder before sending.")
		os.Exit(1)
	}

	// Listen on a random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Printf("Failed to open port: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	localIP := getBestLocalIP()
	publicIP := getPublicIP()

	fmt.Println("=====================================================")
	fmt.Println("Ready to Send! Run the command below on the receiver:")
	fmt.Println("-----------------------------------------------------")
	if publicIP != "" && publicIP != localIP {
		fmt.Printf("Over Internet:  wintolin receive %s:%d\n", publicIP, port)
		fmt.Printf("Over Local/Wifi: wintolin receive %s:%d\n", localIP, port)
		fmt.Printf("\n(Note: You might need to forward port %d on your router for internet transfers)\n", port)
	} else {
		fmt.Printf("wintolin receive %s:%d\n", localIP, port)
	}
	fmt.Println("=====================================================")
	fmt.Println("Waiting for receiver...")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Receiver connected. Starting transfer...")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Could not open file: %v\n", err)
		return
	}
	defer file.Close()

	// Sanitize filename for protocol safety
	fileName := filepath.Base(filePath)
	fileName = strings.ReplaceAll(fileName, "\n", "")
	fileSize := fileInfo.Size()

	// Send metadata
	fmt.Fprintf(conn, "%s\n", fileName)
	fmt.Fprintf(conn, "%d\n", fileSize)

	// Start file transfer with progress tracking
	tracker := &progressTracker{
		Reader:     file,
		totalBytes: fileSize,
		lastUpdate: time.Now(),
	}

	_, err = io.Copy(conn, tracker)
	if err != nil {
		fmt.Printf("\nTransfer error: %v\n", err)
		return
	}

	fmt.Println("\nTransfer completed successfully!")
}

func receiveFile(address string) {
	fmt.Printf("Connecting to: %s...\n", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read file metadata
	fileName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read filename: %v\n", err)
		return
	}
	fileName = strings.TrimSpace(fileName)

	fileSizeStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read file size: %v\n", err)
		return
	}
	fileSizeStr = strings.TrimSpace(fileSizeStr)
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid file size received: %v\n", err)
		return
	}

	fmt.Printf("Incoming File: %s (%d bytes)\n", fileName, fileSize)

	// Create local file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Could not create file: %v\n", err)
		return
	}
	defer file.Close()

	tracker := &progressTracker{
		Reader:     reader,
		totalBytes: fileSize,
		lastUpdate: time.Now(),
	}

	// Read exactly fileSize bytes
	_, err = io.CopyN(file, tracker, fileSize)
	if err != nil && err != io.EOF {
		fmt.Printf("\nTransfer error: %v\n", err)
		return
	}

	fmt.Println("\nTransfer completed successfully! Saved as:", fileName)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  To Send:    wintolin send <file_path>")
		fmt.Println("  To Receive: wintolin receive <ip:port>")
		os.Exit(1)
	}

	command := os.Args[1]
	target := os.Args[2]

	switch command {
	case "send":
		sendFile(target)
	case "receive":
		receiveFile(target)
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands: send, receive")
		os.Exit(1)
	}
}
