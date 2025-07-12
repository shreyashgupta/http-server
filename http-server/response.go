package httpserver

import "fmt"

type Response struct {
	Code     int
	CodeDesc string
	Headers  map[string]string
	Body     string
}

func (resp *Response) GetResponseStr() []byte {
	rspLine := fmt.Sprintf("HTTP/1.1 %d %s", resp.Code, resp.CodeDesc)
	headers := ""
	for _key, _val := range resp.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", _key, _val)
	}
	respStr := rspLine + "\r\n" + headers + "\r\n" + resp.Body
	return []byte(respStr)
}
