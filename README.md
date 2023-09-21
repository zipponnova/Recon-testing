# TCP Scanner in Go

This is a simple, concurrent TCP port scanner written in Go. The scanner allows you to scan a range of ports for multiple IP addresses. It provides features like setting a timeout for scans, retrying on failure, saving results to a log file, and displaying a progress bar during the scan.

## Features

- **Concurrency**: Scans ports concurrently for faster results.
- **Timeouts**: Set custom timeouts for each scan attempt.
- **Retries**: Retry scanning a port if it fails.
- **Logging**: Save scan results to a log file.
- **Progress Bar**: Visual progress bar to track scanning progress.
- **Quiet Mode**: Option to only display IPs with open ports.
- **Colorful Output**: Improved readability with colored output for IP addresses and open ports.

## Prerequisites

1. Go installed on your machine.
2. Required Go packages: 
    - `github.com/cheggaaa/pb/v3` for the progress bar.
    - `github.com/fatih/color` for colored output.

You can install the required packages using:
```bash
go get github.com/cheggaaa/pb/v3
go get github.com/fatih/color
```

## Features

1. Clone the repository or download the tcp_scanner.go script.
2. Navigate to the directory containing the script.
3. Compile the script:
```
go build port_scanner.go
```
4. Run the compiled binary:
```
./port_scanner
```
5. Follow the on-screen prompts to input a single IP, multiple IPs from a file, specify a port range, and other preferences.
