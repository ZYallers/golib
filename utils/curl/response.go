package curl

import (
	"bytes"
	"github.com/ZYallers/golib/funcs/io"
	"net/http"
	"time"
)

type Response struct {
	Raw        *http.Response
	Request    *Request
	Body       string
	receivedAt time.Time
}

// NewResponse new one response
func NewResponse() *Response {
	return &Response{}
}

// Status returns the response status.
func (r *Response) Status() string {
	if r.Raw == nil {
		return ""
	}
	return r.Raw.Status
}

// StatusCode returns the response status code.
func (r *Response) StatusCode() int {
	if r.Raw == nil {
		return 0
	}
	return r.Raw.StatusCode
}

// IsOk if response status code equal 200 return true
func (r *Response) IsOk() bool {
	return r.StatusCode() == http.StatusOK
}

// GetHeader returns the response header value by key.
func (r *Response) GetHeader(key string) string {
	if r.Raw == nil {
		return ""
	}
	return r.Raw.Header.Get(key)
}

// GetHeaderValues returns the response header values by key.
func (r *Response) GetHeaderValues(key string) []string {
	if r.Raw == nil {
		return nil
	}
	return r.Raw.Header.Values(key)
}

// HeaderToString get all header as string.
func (r *Response) HeaderToString() string {
	if r.Raw == nil {
		return ""
	}
	return convertHeaderToString(r.Raw.Header)
}

// GetContentType return the `Content-Type` header value.
func (r *Response) GetContentType() string {
	if r.Raw == nil {
		return ""
	}
	return r.Raw.Header.Get("Content-Type")
}

// TraceInfo returns the TraceInfo from Request.
func (r *Response) TraceInfo() TraceInfo {
	return r.Request.TraceInfo()
}

// TotalTime returns the total time of the request, from request we sent to response we received.
func (r *Response) TotalTime() time.Duration {
	if r.Request.trace != nil {
		return r.Request.TraceInfo().TotalTime
	}
	return r.receivedAt.Sub(r.Request.startTime)
}

// ReceivedAt returns the timestamp that response we received.
func (r *Response) ReceivedAt() time.Time {
	return r.receivedAt
}

// readBody read body to string in response body
func (r *Response) readBody() {
	if r.Raw != nil && r.Raw.Body != nil {
		defer r.Raw.Body.Close()
		if bts, err := io.Copy(r.Raw.Body); err == nil {
			r.Body = string(bts)
		}
	}
}

// setReceivedAt set ReceivedAt time
func (r *Response) setReceivedAt() {
	r.receivedAt = time.Now()
	if r.Request.trace != nil {
		r.Request.trace.endTime = r.receivedAt
	}
}

// convertHeaderToString converts http header to a string.
func convertHeaderToString(h http.Header) string {
	if h == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_ = h.Write(buf)
	return buf.String()
}
