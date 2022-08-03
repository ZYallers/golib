package curl

import (
	"github.com/ZYallers/golib/funcs/io"
	"net/http"
)

type Response struct {
	Raw     *http.Response
	Headers map[string]string
	Body    string
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) StatusCode() int {
	if r.Raw == nil {
		return 0
	}
	return r.Raw.StatusCode
}

func (r *Response) IsOk() bool {
	return r.StatusCode() == http.StatusOK
}

func (r *Response) parseHeaders() {
	headers := map[string]string{}
	for k, v := range r.Raw.Header {
		headers[k] = v[0]
	}
	r.Headers = headers
}

func (r *Response) parseBody() error {
	if bts, err := io.Copy(r.Raw.Body); err != nil {
		return err
	} else {
		r.Body = string(bts)
		return nil
	}
}
