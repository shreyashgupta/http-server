package httpserver

import "fmt"

type Encoding string

const (
	GZIP Encoding = "gzip"
	NONE Encoding = "none"
)

type Encoder interface {
	Encode(s string) string
	Decode(s string) string
}

type GzipEncoder struct {
}

func (e *GzipEncoder) Encode(s string) string {
	return s
}
func (e *GzipEncoder) Decode(s string) string {
	return s
}

func getEcodingFromStr(s string) Encoding {
	if s == "gzip" {
		return GZIP
	}
	return NONE
}

func GetEncoder(encoding Encoding) (Encoder, error) {

	switch encoding {
	case GZIP:
		return &GzipEncoder{}, nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}
