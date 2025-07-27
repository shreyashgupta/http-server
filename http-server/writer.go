package httpserver

import (
	"fmt"
	"net"
	"strconv"
)

type StatusCode int

const (
	HTTP_OK          StatusCode = 200
	HTTP_CREATED     StatusCode = 201
	HTTP_NOT_FOUND   StatusCode = 404
	HTTP_BAD_REQUEST StatusCode = 500
)

var StatusDescMap = map[StatusCode]string{
	HTTP_OK:          "OK",
	HTTP_CREATED:     "Created",
	HTTP_BAD_REQUEST: "Bad Request",
	HTTP_NOT_FOUND:   "Not Found",
}

type RespStatus struct {
	code StatusCode
	desc string
}

type Writer struct {
	conn           net.Conn
	headers        Headers
	status         RespStatus
	body           string
	requestHeaders *Headers
}

func NewWriter(conn net.Conn, requestHeaders *Headers) *Writer {
	return &Writer{
		conn:           conn,
		requestHeaders: requestHeaders,
		headers:        Headers{headers: make(map[string]string)},
		status:         RespStatus{HTTP_OK, StatusDescMap[HTTP_OK]},
	}
}

func (w *Writer) SetStatusCode(code StatusCode) {
	desc, ok := StatusDescMap[code]
	if ok {
		w.status = RespStatus{code: code, desc: desc}
	}
}

func (w *Writer) SetHeader(key string, value string) {
	w.headers.headers[key] = value
}

func (w *Writer) SetContent(data string) {
	w.body = data
}

func (w *Writer) encodeBody(body string) (string, Encoding) {
	encoding := getEcodingFromStr(w.requestHeaders.headers["Accept-Encoding"])
	encoder, err := GetEncoder(encoding)
	if err != nil {
		return body, NONE
	}
	encodedBody, err := encoder.Encode(body)
	if err != nil {
		// add a log here
		return body, NONE
	}
	return encodedBody, encoding
}

func (w *Writer) Write() {
	encodedBody, encoding := w.encodeBody(w.body)
	if encoding != NONE {
		w.SetHeader("Content-Encoding", string(encoding))
	}
	headerConnection, ok := w.requestHeaders.headers["Connection"]

	if ok && headerConnection == "close" {
		w.SetHeader("Connection", "close")
	}
	w.SetHeader("Content-Length", strconv.Itoa(len(encodedBody)))
	rspLine := fmt.Sprintf("HTTP/1.1 %d %s", w.status.code, w.status.desc)
	headers := ""
	for _key, _val := range w.headers.headers {
		headers += fmt.Sprintf("%s: %s\r\n", _key, _val)
	}
	respStr := rspLine + "\r\n" + headers + "\r\n" + encodedBody
	w.conn.Write([]byte(respStr))
	// w.conn.Close()
}
