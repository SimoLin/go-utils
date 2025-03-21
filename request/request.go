package request

import (
	"bytes"
	"compress/gzip"
	"io"
	"slices"
)

func DecodeGzip(header map[string][]string, body []byte) (result []byte, err error) {
	header_content_encoding, ok := header["Content-Encoding"]
	if !ok {
		return body, nil
	}

	if slices.Contains(header_content_encoding, "gzip") {
		return body, nil
	}

	reader := bytes.NewReader(body)

	reader_gzip, err := gzip.NewReader(reader)
	if err != nil {
		return
	}
	defer reader_gzip.Close()

	result, err = io.ReadAll(reader_gzip)
	if err != nil {
		return
	}
	return
}
