package curl

import (
	"bytes"
	"io"
	"sync"
)

var bufferPool = sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 4096)) }}

func ioCopy(r io.Reader) ([]byte, error) {
	dst := bufferPool.Get().(*bytes.Buffer)
	dst.Reset()
	defer func() {
		if dst != nil {
			bufferPool.Put(dst)
			dst = nil
		}
	}()

	if _, err := io.Copy(dst, r); err != nil {
		return nil, err
	}
	bodyB := dst.Bytes()

	bufferPool.Put(dst)
	dst = nil

	return bodyB, nil
}
