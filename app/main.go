package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func extractPath(req_line string) string {
	// ^\w+ → Matches the HTTP method (e.g., GET, POST)
	// \s+ → Space after method
	// ([^\s]+) → Captures the path and query string, i.e., everything until next space
	// \s+HTTP/\d\.\d → Matches the HTTP version (e.g., HTTP/1.1)

	re := regexp.MustCompile(`^\w+\s+([^\s]+)\s+HTTP/\d\.\d`)

	matches := re.FindStringSubmatch(req_line)

	if len(matches) > 1 {
		return matches[1]
	} else {
		return ""
	}
}

func extractEchoPath(path string) string {

	re := regexp.MustCompile(`/echo/([^\s]+)`)

	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		return matches[1]
	} else {
		return ""
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// resp := "HTTP/1.1 200 OK\r\n\r\n"
	// _, err = conn.Write([]byte(resp))
	reader := bufio.NewReader(conn)
	var req []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Handle error (e.g., io.EOF for connection close)
			break
		}
		req = append(req, strings.TrimSpace(line))
		if strings.TrimSpace(line) == "" {
			break // End of headers
		}
	}
	// for _, line := range req {
	// 	fmt.Println(line)
	// }
	path := extractPath(req[0])
	if len(path) == 1 {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		echoPath := extractEchoPath(path)
		if len(echoPath) == 0 {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoPath), echoPath)))
		}
	}

}
