package httpserver

type RequestLine struct {
	reqType RequestType
	path    string
}
type Headers struct {
	headers map[string]string
}

type Request struct {
	requestLine RequestLine
	headers     Headers
	body        string
}

func (r *Request) GetHeader(key string) string {
	return r.headers.headers[key]
}
