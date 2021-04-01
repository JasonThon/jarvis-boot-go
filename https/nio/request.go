package nio

import (
	"bytes"
	"io"
	"net/http"
)

func NewBufferedRequest(request *http.Request) (*http.Request, error) {
	defer request.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0))

	_, err := io.Copy(buf, request.Body)

	if err != nil {
		return nil, err
	}

	return http.NewRequest(request.Method, request.URL.RequestURI(), buf)
}
