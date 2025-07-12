package main

import (
	"net"
	"os"
	"strconv"

	httpserver "github.com/codecrafters-io/http-server-starter-go/http-server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

// type RequestLineParser struct {
// }

// type HeaderParser struct {
// }

// type Bodyparser struct {
// }

// func (p *RequestLineParser) extractPath(req_line string) string {
// 	// ^\w+ → Matches the HTTP method (e.g., GET, POST)
// 	// \s+ → Space after method
// 	// ([^\s]+) → Captures the path and query string, i.e., everything until next space
// 	// \s+HTTP/\d\.\d → Matches the HTTP version (e.g., HTTP/1.1)

// 	re := regexp.MustCompile(`^\w+\s+([^\s]+)\s+HTTP/\d\.\d`)

// 	matches := re.FindStringSubmatch(req_line)

// 	if len(matches) > 1 {
// 		return matches[1]
// 	} else {
// 		return ""
// 	}
// }

// func (p *HeaderParser) parse(headersStr []string) (Headers, error) {
// 	headers := Headers{}
// 	headers.headers = make(map[string]string)
// 	// fmt.Println(headersStr)
// 	re := regexp.MustCompile(`([^\s]+):\s(.+)`)
// 	for _, header := range headersStr {
// 		if len(header) == 0 {
// 			break
// 		}

// 		matches := re.FindStringSubmatch(header)
// 		if len(matches) != 3 {
// 			return Headers{}, fmt.Errorf("Failed to parse headers")
// 		}
// 		headers.headers[matches[1]] = matches[2]

// 		// fmt.Println(headers.headers)
// 	}
// 	return headers, nil
// }

// func (p *RequestLineParser) extractEchoPath(path string) string {

// 	re := regexp.MustCompile(`/echo/([^\s]+)`)

// 	matches := re.FindStringSubmatch(path)

// 	if len(matches) > 1 {
// 		return matches[1]
// 	} else {
// 		return ""
// 	}
// }

// func (p *RequestLineParser) extractUserAgent(path string) string {

// 	re := regexp.MustCompile(`/echo/([^\s]+)`)

// 	matches := re.FindStringSubmatch(path)

// 	if len(matches) > 1 {
// 		return matches[1]
// 	} else {
// 		return ""
// 	}
// }

// type RequestParser struct {
// 	reqLineParser RequestLineParser
// 	headerParser  HeaderParser
// 	bodyparser    Bodyparser
// 	requestLines  []string
// }

// type RequestLine struct {
// 	path string
// }

// func (r *RequestLine) getEchoPathIfEcho() (string, bool) {

// 	re := regexp.MustCompile(`/echo/([^\s]+)`)

// 	matches := re.FindStringSubmatch(r.path)

// 	if len(matches) > 1 {
// 		return matches[1], true
// 	} else {
// 		return "", false
// 	}
// }

// func (r *RequestLine) getFilePathIfFile() (string, bool) {

// 	re := regexp.MustCompile(`/files/([^\s]+)`)

// 	matches := re.FindStringSubmatch(r.path)

// 	if len(matches) > 1 {
// 		return matches[1], true
// 	} else {
// 		return "", false
// 	}
// }

// func (r *RequestLine) isUserAgent() bool {

// 	re := regexp.MustCompile(`/user-agent`)

// 	match := re.Find([]byte(r.path))

// 	return len(match) > 0
// }

// type Headers struct {
// 	headers map[string]string
// }

// type Request struct {
// 	requestLine RequestLine
// 	headers     Headers
// 	body        string
// }

// func NewRequestParser(conn net.Conn) *RequestParser {
// 	reader := bufio.NewReader(conn)
// 	var req []string
// 	for {
// 		line, err := reader.ReadString('\n')
// 		if err != nil {
// 			// Handle error (e.g., io.EOF for connection close)
// 			break
// 		}
// 		req = append(req, strings.TrimSpace(line))
// 		if strings.TrimSpace(line) == "" {
// 			break // End of headers
// 		}
// 	}
// 	return &RequestParser{requestLines: req}
// }

// func (p *RequestParser) Parse() (Request, error) {
// 	reqLinePath := p.reqLineParser.extractPath(p.requestLines[0])
// 	headers, err := p.headerParser.parse(p.requestLines[1:])
// 	if err != nil {
// 		return Request{}, err
// 	}
// 	return Request{requestLine: RequestLine{path: reqLinePath}, headers: headers}, nil

// }

// type ResponseSelector struct {
// }

// func (s *ResponseSelector) getResponse(r Request) Response {
// 	echoPath, isEcho := r.requestLine.getEchoPathIfEcho()
// 	if isEcho {
// 		return Response{code: 200,
// 			codeDesc: "OK",
// 			headers: map[string]string{
// 				"Content-Length": strconv.Itoa(len(echoPath)),
// 				"Content-Type":   "text/plain",
// 			},
// 			body: echoPath,
// 		}
// 	}
// 	isUserAgent := r.requestLine.isUserAgent()
// 	if isUserAgent {
// 		return Response{code: 200,
// 			codeDesc: "OK",
// 			headers: map[string]string{
// 				"Content-Length": strconv.Itoa(len(r.headers.headers["User-Agent"])),
// 				"Content-Type":   "text/plain",
// 			},
// 			body: r.headers.headers["User-Agent"],
// 		}
// 	}
// 	filePath, isFileReq := r.requestLine.getFilePathIfFile()
// 	if isFileReq {
// 		absFilePath := directory + filePath
// 		fileContent, err := readFile(absFilePath)
// 		if err != nil {
// 			return Response{code: 404, codeDesc: "Not Found"}
// 		}
// 		return Response{code: 200,
// 			codeDesc: "OK",
// 			headers: map[string]string{
// 				"Content-Length": strconv.Itoa(len(fileContent)),
// 				"Content-Type":   "application/octet-stream",
// 			},
// 			body: fileContent,
// 		}
// 	}
// 	if len(r.requestLine.path) == 1 {
// 		return Response{code: 200, codeDesc: "OK"}
// 	} else {
// 		return Response{code: 404, codeDesc: "Not Found"}
// 	}

// }

// type Response struct {
// 	code     int
// 	codeDesc string
// 	headers  map[string]string
// 	body     string
// }

// func (resp *Response) getResponse() []byte {
// 	rspLine := fmt.Sprintf("HTTP/1.1 %d %s", resp.code, resp.codeDesc)
// 	headers := ""
// 	for _key, _val := range resp.headers {
// 		headers += fmt.Sprintf("%s: %s\r\n", _key, _val)
// 	}
// 	respStr := rspLine + "\r\n" + headers + "\r\n" + resp.body
// 	return []byte(respStr)
// }

// func handleConnection(conn net.Conn) {
// 	fmt.Println("Handling new connection")
// 	requestParser := NewRequestParser(conn)
// 	req, err := requestParser.parse()
// 	if err != nil {
// 		fmt.Println("Error parsing request: ", err.Error())
// 		os.Exit(1)
// 	}
// 	responseSelector := ResponseSelector{}
// 	response := responseSelector.getResponse(req)
// 	conn.Write(response.getResponse())
// 	conn.Close()
// }

var directory = ""

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path) // Go 1.16+
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	args := os.Args[1:] // skip program name

	if len(args) == 2 {
		directory = args[1]
	}

	// l, err := net.Listen("tcp", "0.0.0.0:4221")
	// if err != nil {
	// 	fmt.Println("Failed to bind to port 4221")
	// 	os.Exit(1)
	// }
	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection: ", err.Error())
	// 		os.Exit(1)
	// 	}
	// 	go handleConnection(conn)
	// }
	mux := httpserver.NewMux()
	mux.Get("/", defaultHandler)
	mux.Get("/echo/{echo_str}", echoHandler)
	mux.Get("/user-agent", userAgentHandler)
	mux.Get("/files/{file_name}", fileHandler)
	mux.Post("/files/{file_name}", filePostHandler)

	server := httpserver.NewHttpServer(mux, "0.0.0.0:4221")

	server.Serve()
}

func defaultHandler(conn net.Conn, captures map[string]string, _ httpserver.Request) {
	resp := httpserver.Response{Code: 200, CodeDesc: "OK"}
	conn.Write(resp.GetResponseStr())
	conn.Close()
}

func echoHandler(conn net.Conn, captures map[string]string, _ httpserver.Request) {
	resp := httpserver.Response{Code: 200,
		CodeDesc: "OK",
		Headers: map[string]string{
			"Content-Length": strconv.Itoa(len(captures["echo_str"])),
			"Content-Type":   "text/plain",
		},
		Body: captures["echo_str"],
	}
	conn.Write(resp.GetResponseStr())
	conn.Close()
}

func userAgentHandler(conn net.Conn, captures map[string]string, req httpserver.Request) {
	resp := httpserver.Response{Code: 200,
		CodeDesc: "OK",
		Headers: map[string]string{
			"Content-Length": strconv.Itoa(len(req.GetHeader("User-Agent"))),
			"Content-Type":   "text/plain",
		},
		Body: req.GetHeader("User-Agent"),
	}
	conn.Write(resp.GetResponseStr())
	conn.Close()
}

func fileHandler(conn net.Conn, captures map[string]string, req httpserver.Request) {
	absFilePath := directory + captures["file_name"]
	fileContent, err := readFile(absFilePath)
	if err != nil {
		resp := httpserver.Response{Code: 500, CodeDesc: "Not Found"}
		conn.Write(resp.GetResponseStr())
		conn.Close()
	}

	resp := httpserver.Response{Code: 200,
		CodeDesc: "OK",
		Headers: map[string]string{
			"Content-Length": strconv.Itoa(len(fileContent)),
			"Content-Type":   "text/plain",
		},
		Body: fileContent,
	}
	conn.Write(resp.GetResponseStr())
	conn.Close()
}

func writeToFile(path string, data []byte) error {
	// Create file (or truncate if it already exists)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to file
	_, err = file.Write(data)
	return err
}

func filePostHandler(conn net.Conn, captures map[string]string, req httpserver.Request) {
	absFilePath := directory + captures["file_name"]
	err := writeToFile(absFilePath, req.GetBodyData())
	if err != nil {
		resp := httpserver.Response{Code: 500, CodeDesc: "failed to write"}
		conn.Write(resp.GetResponseStr())
		conn.Close()
	}

	resp := httpserver.Response{Code: 201,
		CodeDesc: "Created",
	}
	conn.Write(resp.GetResponseStr())
	conn.Close()
}
