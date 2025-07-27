package httpserver

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type RequestLineParser struct {
}

type HeaderParser struct {
}

type Bodyparser struct {
}

func (p *RequestLineParser) extractPathAndType(req_line string) (string, string) {
	// ^\w+ → Matches the HTTP method (e.g., GET, POST)
	// \s+ → Space after method
	// ([^\s]+) → Captures the path and query string, i.e., everything until next space
	// \s+HTTP/\d\.\d → Matches the HTTP version (e.g., HTTP/1.1)

	re := regexp.MustCompile(`(^\w+)\s+([^\s]+)\s+HTTP/\d\.\d`)

	matches := re.FindStringSubmatch(req_line)

	if len(matches) > 2 {
		return matches[1], matches[2]
	} else {
		return "", ""
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

func (p *Bodyparser) parse(headers Headers, reader *bufio.Reader) ([]byte, error) {
	var body []byte
	if lenStr, ok := headers.headers["Content-Length"]; ok {
		contentLength, err := strconv.Atoi(lenStr)
		if err != nil {
			return []byte{}, fmt.Errorf("invalid Content-Length: %w", err)
		}

		body = make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return []byte{}, fmt.Errorf("error reading body: %w", err)
		}
	}
	return body, nil
}

type RequestParser struct {
	reqLineParser RequestLineParser
	headerParser  HeaderParser
	bodyparser    Bodyparser
}

func (p *RequestParser) Parse(reader *bufio.Reader) (Request, error) {
	var requestLines []string
	_, err := reader.Peek(1)
	if err != nil {
		// No data available
		if err == io.EOF {
			return Request{}, io.EOF // connection closed
		}
		// Can also return nil error to allow graceful wait
		return Request{}, err
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Handle error (e.g., io.EOF for connection close)
			break
		}
		requestLines = append(requestLines, strings.TrimSpace(line))
		if strings.TrimSpace(line) == "" {
			break // End of headers
		}
	}
	reqType, reqLinePath := p.reqLineParser.extractPathAndType(requestLines[0])

	headers, err := p.headerParser.parse(requestLines[1:])
	if err != nil {
		return Request{}, err
	}

	body, err := p.bodyparser.parse(headers, reader)
	if err != nil {
		return Request{}, err
	}
	return Request{
		requestLine: RequestLine{reqType: ToRequestType(reqType), path: reqLinePath},
		headers:     headers,
		body:        Body{bodyType: getBodyType(headers.headers["Content-Type"]), bodyData: body}}, nil
}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}
