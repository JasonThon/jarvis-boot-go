package bytes2

import "bytes"

func NewByteSlice(size int) []byte {
	return make([]byte, size)
}

func NewByteBuffer() *bytes.Buffer {
	return bytes.NewBuffer(NewByteSlice(0))
}
