package httpserver

type RequestLine struct {
	reqType RequestType
	path    string
}
type Headers struct {
	headers map[string]string
}

type BodyType int

const (
	JSON BodyType = iota
	OCTET_STREAM
	UNKNOWN
)

type Body struct {
	bodyData []byte
	bodyType BodyType
}

func (body *Body) format() any {
	if body.bodyType == OCTET_STREAM {
		return body.bodyData
	}
	if body.bodyType == JSON {
		return string(body.bodyData)
	}
	return ""
}

func getBodyType(bodyType string) BodyType {
	if bodyType == "application/octet-stream" {
		return OCTET_STREAM
	}
	if bodyType == "application/json" {
		return JSON
	}
	return UNKNOWN
}

type Request struct {
	requestLine RequestLine
	headers     Headers
	body        Body
	captures    map[string]string
}

func (r *Request) GetCapture(key string) string {
	return r.captures[key]
}

func (r *Request) GetHeader(key string) (string, bool) {
	val, ok := r.headers.headers[key]
	return val, ok
}

func (r *Request) GetBodyData() []byte {
	return r.body.bodyData
}
