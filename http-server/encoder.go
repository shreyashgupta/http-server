package httpserver

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
)

type Encoding string

const (
	GZIP Encoding = "gzip"
	NONE Encoding = "none"
)

type Encoder interface {
	Encode(s string) (string, error)
	Decode(s string) (string, error)
}

type GzipEncoder struct {
}

func (e *GzipEncoder) Encode(s string) (string, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write([]byte(s))
	if err != nil {
		return "", err
	}

	err = writer.Close() // important to flush and finalize
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
func (e *GzipEncoder) Decode(s string) (string, error) {
	return s, nil
}

func getEcodingFromStr(s string) Encoding {
	requestEncodings := strings.Split(s, ",")

	for _, requestEncoding := range requestEncodings {
		// Return first supported encoding
		if strings.TrimSpace(requestEncoding) == "gzip" {
			return GZIP
		}
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
