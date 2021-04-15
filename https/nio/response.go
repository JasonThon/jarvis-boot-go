package nio

import (
	"bytes"
	"io"
	"net/http"
	"thingworks.net/thingworks/jarvis-boot/utils/bytes2"
)

type BufferedResponseWriter struct {
	writer http.ResponseWriter
	buf    []byte
}

func NewBufferedResponseWriter(writer http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		writer: writer,
	}
}

func (bufWriter *BufferedResponseWriter) Header() http.Header {
	return bufWriter.writer.Header()
}

func (bufWriter *BufferedResponseWriter) Write(data []byte) (int, error) {
	bufWriter.buf = bytes2.NewByteSlice(len(data))
	copy(bufWriter.buf, data)

	return len(data), nil
}

func (bufWriter *BufferedResponseWriter) WriteHeader(statusCode int) {
	bufWriter.writer.WriteHeader(statusCode)
}

func (bufWriter *BufferedResponseWriter) Copy() (int64, error) {
	if bufWriter.buf != nil && len(bufWriter.buf) > 0 {
		reader := bytes.NewBuffer(bufWriter.buf)

		return io.Copy(bufWriter.writer, reader)
	}

	return 0, nil
}
