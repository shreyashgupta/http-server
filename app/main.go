package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type RequestLineParser struct {
}

type HeaderParser struct {
}

type Bodyparser struct {
}

func (p *RequestLineParser) extractPath(req_line string) string {
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

func (p *HeaderParser) parse(headersStr []string) (Headers, error) {
	headers := Headers{}
	headers.headers = make(map[string]string)
	// fmt.Println(headersStr)
	re := regexp.MustCompile(`([^\s]+):\s(.+)`)
	for _, header := range headersStr {
		if len(header) == 0 {
			break
		}

		matches := re.FindStringSubmatch(header)
		if len(matches) != 3 {
			return Headers{}, fmt.Errorf("Failed to parse headers")
		}
		headers.headers[matches[1]] = matches[2]

		// fmt.Println(headers.headers)
	}
	return headers, nil
}

func (p *RequestLineParser) extractEchoPath(path string) string {

	re := regexp.MustCompile(`/echo/([^\s]+)`)

	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		return matches[1]
	} else {
		return ""
	}
}

func (p *RequestLineParser) extractUserAgent(path string) string {

	re := regexp.MustCompile(`/echo/([^\s]+)`)

	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		return matches[1]
	} else {
		return ""
	}
}

type RequestParser struct {
	reqLineParser RequestLineParser
	headerParser  HeaderParser
	bodyparser    Bodyparser
	requestLines  []string
}

type RequestLine struct {
	path string
}

func (r *RequestLine) getEchoPathIfEcho() (string, bool) {

	re := regexp.MustCompile(`/echo/([^\s]+)`)

	matches := re.FindStringSubmatch(r.path)

	if len(matches) > 1 {
		return matches[1], true
	} else {
		return "", false
	}
}

func (r *RequestLine) isUserAgent() bool {

	re := regexp.MustCompile(`/user-agent`)

	match := re.Find([]byte(r.path))

	return len(match) > 0
}

type Headers struct {
	headers map[string]string
}

type Request struct {
	requestLine RequestLine
	headers     Headers
	body        string
}

func NewRequestParser(conn net.Conn) *RequestParser {
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
	return &RequestParser{requestLines: req}
}

func (p *RequestParser) parse() (Request, error) {
	reqLinePath := p.reqLineParser.extractPath(p.requestLines[0])
	headers, err := p.headerParser.parse(p.requestLines[1:])
	if err != nil {
		return Request{}, err
	}
	return Request{requestLine: RequestLine{path: reqLinePath}, headers: headers}, nil

}

func (p *RequestParser) getEchoPath() (string, error) {
	path := p.reqLineParser.extractPath(p.requestLines[0])
	echoPath := p.reqLineParser.extractEchoPath(path)

	if echoPath != "" {
		return echoPath, nil
	}
	return "", fmt.Errorf("Not an echo request")
}

func (p *RequestParser) getUserAgent() (string, error) {
	path := p.reqLineParser.extractPath(p.requestLines[0])
	echoPath := p.reqLineParser.extractEchoPath(path)

	if echoPath != "" {
		return echoPath, nil
	}
	return "", fmt.Errorf("Not an echo request")
}

type ResponseSelector struct {
}

func (s *ResponseSelector) getResponse(r Request) Response {
	echoPath, isEcho := r.requestLine.getEchoPathIfEcho()
	if isEcho {
		return Response{code: 200,
			codeDesc: "OK",
			headers: map[string]string{
				"Content-Length": strconv.Itoa(len(echoPath)),
				"Content-Type":   "text/plain",
			},
			body: echoPath,
		}
	}
	isUserAgent := r.requestLine.isUserAgent()
	if isUserAgent {
		return Response{code: 200,
			codeDesc: "OK",
			headers: map[string]string{
				"Content-Length": strconv.Itoa(len(r.headers.headers["User-Agent"])),
				"Content-Type":   "text/plain",
			},
			body: r.headers.headers["User-Agent"],
		}
	}

	if len(r.requestLine.path) == 1 {
		return Response{code: 200, codeDesc: "OK"}
	} else {
		return Response{code: 404, codeDesc: "Not Found"}
	}

}

type Response struct {
	code     int
	codeDesc string
	headers  map[string]string
	body     string
}

func (resp *Response) getResponse() []byte {
	rspLine := fmt.Sprintf("HTTP/1.1 %d %s", resp.code, resp.codeDesc)
	headers := ""
	for _key, _val := range resp.headers {
		headers += fmt.Sprintf("%s: %s\r\n", _key, _val)
	}
	respStr := rspLine + "\r\n" + headers + "\r\n" + resp.body
	return []byte(respStr)
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
	requestParser := NewRequestParser(conn)

	req, err := requestParser.parse()
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}

	responseSelector := ResponseSelector{}

	response := responseSelector.getResponse(req)
	conn.Write(response.getResponse())

	// reqLineParser := RequestLineParser{req[0]}
	// path := reqLineParser.extractPath(req[0])
	// if len(path) == 1 {
	// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	// } else {
	// 	echoPath := reqLineParser.extractEchoPath(path)
	// 	if len(echoPath) == 0 {
	// 		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	// 	} else {
	// 		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoPath), echoPath)))
	// 	}
	// }

}
