# Recon-testing
Arsenal of security tools to perform reconnaissance like port scanning, vulnerability scanning and OSINT 

## Port scanning
This is a simple concurrent port scanner written in Go (Golang). It allows you to scan a range of TCP ports on a given IP address to determine if they are open or closed.

### Features

- Concurrent scanning of multiple ports using goroutines
- Rate limiting to control the rate of connection attempts
- Displays scan results in a table format

### Usage

1. Clone the repository:
`git clone https://github.com/zipponnova/Recon-testing`

2. Navigate to the project directory:
`cd Recon-testing`

3. Build and run the application:
```
go run portscan.go [IP Address] [Start Port] [End Port]
```

Replace `[IP Address]`, `[Start Port]`, and `[End Port]` with the appropriate values for the target IP address and port range.

4. View the scan results:

The application will display the open ports along with their respective states and reasons in a table format.

### Requirements

- Go 1.13 or higher

### Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

### License

This project is licensed under the [MIT License](LICENSE).







