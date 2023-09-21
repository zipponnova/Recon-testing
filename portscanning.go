package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
)

func tcpConnectScan(target string, port int, timeout time.Duration, retries int, results chan<- int) {
	address := fmt.Sprintf("%s:%d", target, port)
	for i := 0; i < retries; i++ {
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			results <- port
			conn.Close()
			return
		}
	}
}

func main() {
	var ipFile, singleIP, portRange, logFile string
	var startPort, endPort, retries int
	var timeoutDuration time.Duration
	var quietMode, saveToLog bool

	ipColor := color.New(color.FgCyan).SprintFunc()
	openPortsColor := color.New(color.FgGreen).SprintFunc()

	fmt.Print("Enter single IP (or leave blank if using IP file): ")
	fmt.Scanln(&singleIP)

	fmt.Print("Enter IP file path (or leave blank if using single IP): ")
	fmt.Scanln(&ipFile)

	fmt.Print("Enter port range (e.g., 1-1024): ")
	fmt.Scanln(&portRange)
	_, err := fmt.Sscanf(portRange, "%d-%d", &startPort, &endPort)
	if err != nil {
		fmt.Println("Invalid port range format. Exiting.")
		return
	}

	fmt.Print("Enter timeout in seconds (e.g., 1): ")
	var timeoutSec int
	fmt.Scan(&timeoutSec)
	timeoutDuration = time.Duration(timeoutSec) * time.Second

	fmt.Print("Number of retries per port (e.g., 1): ")
	fmt.Scan(&retries)

	fmt.Print("Quiet mode (only display IPs with open ports)? (yes/no): ")
	var quietInput string
	fmt.Scan(&quietInput)
	quietMode = strings.ToLower(quietInput) == "yes"

	fmt.Print("Save results to log file? (yes/no): ")
	var logInput string
	fmt.Scan(&logInput)
	saveToLog = strings.ToLower(logInput) == "yes"
	if saveToLog {
		fmt.Print("Enter log file path (e.g., scan_results.log): ")
		fmt.Scan(&logFile)
	}

	var targets []string

	if singleIP != "" {
		targets = append(targets, singleIP)
	} else if ipFile != "" {
		file, err := os.Open(ipFile)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			targets = append(targets, strings.TrimSpace(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading lines:", err)
			return
		}
	} else {
		fmt.Println("Either a single IP or IP file must be provided.")
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50) // Limit concurrency to 50 goroutines
	totalScans := len(targets) * (endPort - startPort + 1)
	progressBar := pb.StartNew(totalScans)

	for _, ip := range targets {
		ipResults := make(chan int, endPort-startPort+1)
		startTime := time.Now()

		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			go func(ip string, port int) {
				defer wg.Done()
				semaphore <- struct{}{}
				tcpConnectScan(ip, port, timeoutDuration, retries, ipResults)
				<-semaphore
				progressBar.Increment()
			}(ip, port)
		}

		go func() {
			wg.Wait()
			close(ipResults)
		}()

		ports := []int{}
		for port := range ipResults {
			ports = append(ports, port)
		}

		sort.Ints(ports)

		if len(ports) > 0 || !quietMode {
			fmt.Println("\nScanning IP:", ipColor(ip))
			if len(ports) > 0 {
				portStrings := []string{}
				for _, p := range ports {
					portStrings = append(portStrings, strconv.Itoa(p))
				}
				scanResult := fmt.Sprintf("Open Ports: %s", openPortsColor(strings.Join(portStrings, ", ")))
				fmt.Println(scanResult)
				if saveToLog {
					saveScanResultToLog(logFile, ip, scanResult)
				}
			} else {
				fmt.Println("No open ports found.")
			}
			fmt.Printf("Scan Duration: %v\n", time.Since(startTime))
		}
	}

	progressBar.Finish()
}

func saveScanResultToLog(logFile, ip, result string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s: %s\n", ip, result))
	if err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
