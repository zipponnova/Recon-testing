package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"github.com/olekukonko/tablewriter"
)

const (
	concurrency = 1000
)

type ScanResult struct {
	Port   int
	State  string
	Reason string
}

func worker(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, ports chan int, results chan ScanResult, hostname string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case p, ok := <-ports:
			if !ok {
				return
			}
			address := fmt.Sprintf("%s:%d", hostname, p)
			err := limiter.Wait(ctx) // This will block until limiter permits another event
			if err != nil {
				break
			}
			conn, err := net.DialTimeout("tcp", address, 5*time.Second)
			result := ScanResult{Port: p}
			if err != nil {
				result.State = "Closed"
				result.Reason = err.Error()
			} else {
				result.State = "Open"
				conn.Close()
			}
			results <- result
		}
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run portscan.go [IP Address] [Start Port] [End Port]")
		os.Exit(1)
	}

	hostname := os.Args[1]
	start, _ := strconv.Atoi(os.Args[2])
	end, _ := strconv.Atoi(os.Args[3])

	ports := make(chan int, concurrency)
	results := make(chan ScanResult, concurrency)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	limiter := rate.NewLimiter(500, concurrency)

	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go worker(ctx, &wg, limiter, ports, results, hostname)
	}

	go func() {
		for i := start; i <= end; i++ {
			ports <- i
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Port", "State", "Reason"})

	for r := range results {
		if r.State == "Open" {
			table.Append([]string{strconv.Itoa(r.Port), r.State, r.Reason})
		}
	}

	table.Render()
}

