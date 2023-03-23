package curl

import (
	"bytes"
	"context"
	"fmt"
	libIo "github.com/ZYallers/golib/funcs/io"
	"github.com/ZYallers/golib/utils/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Url      string
	Method   string
	Timeout  time.Duration
	Headers  map[string]string
	Cookies  map[string]string
	Queries  map[string]string
	PostData map[string]interface{}
	Body     io.Reader
	Response *Response

	close              bool
	error              error
	startTime          time.Time
	responseReturnTime time.Time
	ctx                context.Context
	client             *http.Client
	rawRequest         *http.Request
	trace              *clientTrace
}

// NewRequest new request
func NewRequest(url string) *Request {
	req := &Request{Url: url, client: Client, Response: &Response{}, Timeout: ClientTimeout}
	req.Response.Request = req
	return req
}

// SetMethod set request method
func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

// SetUrl set request url
func (r *Request) SetUrl(url string) *Request {
	r.Url = url
	return r
}

// SetHeaders set request headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.Headers = headers
	return r
}

// SetHeader set request header
func (r *Request) SetHeader(key, value string) *Request {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	r.Headers[key] = value
	return r
}

// SetCookies set request cookies
func (r *Request) SetCookies(cookies map[string]string) *Request {
	r.Cookies = cookies
	return r
}

// SetCookie set request cookie
func (r *Request) SetCookie(key, value string) *Request {
	if r.Cookies == nil {
		r.Cookies = map[string]string{}
	}
	r.Cookies[key] = value
	return r
}

// SetQueries set request query
func (r *Request) SetQueries(queries map[string]string) *Request {
	r.Queries = queries
	return r
}

// SetQuery set request query
func (r *Request) SetQuery(key, value string) *Request {
	if r.Queries == nil {
		r.Queries = map[string]string{}
	}
	r.Queries[key] = value
	return r
}

// SetPostData set post data
func (r *Request) SetPostData(data map[string]interface{}) *Request {
	if data == nil {
		return r
	}

	r.PostData = data

	if ct, ok := r.Headers["Content-Type"]; ok && strings.HasPrefix(ct, JsonContentType) {
		if postDataBytes, err := json.Marshal(r.PostData); err != nil {
			r.error = err
		} else {
			r.Body = bytes.NewReader(postDataBytes)
		}
		return r
	}

	// If the Content Type cannot be matched, the default method is application/x-www-form-urlencoded
	posts := url.Values{}
	for k, v := range r.PostData {
		posts.Set(k, fmt.Sprint(v))
	}
	r.Body = strings.NewReader(posts.Encode())

	return r
}

// SetBody set the request Body, accepts string, []byte, io.Reader, io.ReadCloser.
func (r *Request) SetBody(body interface{}) *Request {
	if body == nil {
		return r
	}

	switch b := body.(type) {
	case io.ReadCloser:
		r.Body = b
	case io.Reader:
		r.Body = b
	case []byte:
		r.Body = bytes.NewReader(b)
	case string:
		r.Body = strings.NewReader(b)
	default:
		r.Body = strings.NewReader(fmt.Sprint(body))
	}

	return r
}

// SetTimeOut set request timeout after
func (r *Request) SetTimeOut(timeout time.Duration) *Request {
	r.Timeout = timeout
	return r
}

// SetContentType set the `Content-Type` for the request.
func (r *Request) SetContentType(contentType string) *Request {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	r.Headers["Content-Type"] = contentType
	return r
}

// SetContext method sets the context.Context for current Request. It allows
// to interrupt the request execution if ctx.Done() channel is closed.
// See https://blog.golang.org/context article and the "context" package
// documentation.
func (r *Request) SetContext(ctx context.Context) *Request {
	if ctx != nil {
		r.ctx = ctx
	}
	return r
}

// EnableCloseConn closes the connection after sending this request and reading its response if set to true in HTTP/1.1 and HTTP/2.
// Setting this field prevents re-use of TCP connections between requests to the same hosts event if EnableKeepAlives() were called.
func (r *Request) EnableCloseConn() *Request {
	r.close = true
	return r
}

// DisableCloseConn disable close connection
func (r *Request) DisableCloseConn() *Request {
	r.close = false
	return r
}

// EnableTrace enables trace (http3 currently does not support trace).
func (r *Request) EnableTrace() *Request {
	if r.trace == nil {
		r.trace = &clientTrace{}
	}
	return r
}

// DisableTrace disables trace.
func (r *Request) DisableTrace() *Request {
	r.trace = nil
	return r
}

// Get init get request and return response
func (r *Request) Get() (*Response, error) {
	return r.SetMethod(http.MethodGet).Send()
}

// Post init post request and return response
func (r *Request) Post() (*Response, error) {
	return r.SetMethod(http.MethodPost).Send()
}

// Send init request and return response
func (r *Request) Send() (*Response, error) {
	defer func() { r.responseReturnTime = time.Now() }()

	if r.Response == nil {
		r.Response = &Response{Request: r}
	}

	if r.error != nil {
		return r.Response, r.error
	}

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	if r.trace != nil {
		r.ctx = r.trace.createContext(r.ctx)
	}

	if r.Timeout > 0 {
		var cancel context.CancelFunc
		r.ctx, cancel = context.WithTimeout(r.ctx, r.Timeout)
		defer cancel()
	}

	if r.Body != nil {
		if bodyBytes, err := libIo.Copy(r.Body); err != nil {
			r.error = err
			return r.Response, r.error
		} else {
			defer func() { r.Body = bytes.NewReader(bodyBytes) }()
			r.Body = bytes.NewReader(bodyBytes)
		}
	}

	if r.rawRequest, r.error = http.NewRequestWithContext(r.ctx, r.Method, r.Url, r.Body); r.error != nil {
		return r.Response, r.error
	}

	// Set close connection
	if r.close {
		r.rawRequest.Close = true
	}

	// Set header
	if len(r.Headers) > 0 {
		for k, v := range r.Headers {
			r.rawRequest.Header.Set(k, v)
		}
	}

	// Set cookies
	if len(r.Cookies) > 0 {
		for k, v := range r.Cookies {
			r.rawRequest.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	// Set query
	if len(r.Queries) > 0 {
		q := r.rawRequest.URL.Query()
		for k, v := range r.Queries {
			q.Add(k, v)
		}
		r.rawRequest.URL.RawQuery = q.Encode()
	}

	r.startTime = time.Now()
	r.Response.Raw, r.error = r.client.Do(r.rawRequest)
	if r.error == nil {
		r.Response.setReceivedAt()
		r.Response.readBody()
	}

	return r.Response, r.error
}

// TraceInfo returns the trace information, only available if trace is enabled
func (r *Request) TraceInfo() TraceInfo {
	ct := r.trace

	if ct == nil {
		return TraceInfo{}
	}

	ti := TraceInfo{
		IsConnReused:  ct.gotConnInfo.Reused,
		IsConnWasIdle: ct.gotConnInfo.WasIdle,
		ConnIdleTime:  ct.gotConnInfo.IdleTime,
	}

	endTime := ct.endTime
	if endTime.IsZero() { // in case timeout
		endTime = r.responseReturnTime
	}

	if !ct.tlsHandshakeStart.IsZero() {
		if !ct.tlsHandshakeDone.IsZero() {
			ti.TLSHandshakeTime = ct.tlsHandshakeDone.Sub(ct.tlsHandshakeStart)
		} else {
			ti.TLSHandshakeTime = endTime.Sub(ct.tlsHandshakeStart)
		}
	}

	if ct.gotConnInfo.Reused {
		ti.TotalTime = endTime.Sub(ct.getConn)
	} else {
		if ct.dnsStart.IsZero() {
			ti.TotalTime = endTime.Sub(r.startTime)
		} else {
			ti.TotalTime = endTime.Sub(ct.dnsStart)
		}
	}

	dnsDone := ct.dnsDone
	if dnsDone.IsZero() {
		dnsDone = endTime
	}

	if !ct.dnsStart.IsZero() {
		ti.DNSLookupTime = dnsDone.Sub(ct.dnsStart)
	}

	// Only calculate on successful connections
	if !ct.connectDone.IsZero() {
		ti.TCPConnectTime = ct.connectDone.Sub(dnsDone)
	}

	// Only calculate on successful connections
	if !ct.gotConn.IsZero() {
		ti.ConnectTime = ct.gotConn.Sub(ct.getConn)
	}

	// Only calculate on successful connections
	if !ct.gotFirstResponseByte.IsZero() {
		ti.FirstResponseTime = ct.gotFirstResponseByte.Sub(ct.gotConn)
		ti.ResponseTime = endTime.Sub(ct.gotFirstResponseByte)
	}

	// Capture remote address info when connection is non-nil
	if ct.gotConnInfo.Conn != nil {
		ti.RemoteAddr = ct.gotConnInfo.Conn.RemoteAddr()
		ti.LocalAddr = ct.gotConnInfo.Conn.LocalAddr()
	}

	return ti
}

func (r *Request) DumpRequest() string {
	var buf bytes.Buffer

	req := r.rawRequest
	_, _ = fmt.Fprintf(&buf, "Proto: %s\r\nMethod: %s\r\n", req.Proto, req.Method)
	if req.URL != nil {
		_, _ = fmt.Fprintf(&buf, "Scheme: %s\r\nHost: %s\r\nPath: %s\r\nQuery: %s\r\nURL: %s\r\n",
			req.URL.Scheme, req.URL.Host, req.URL.Path, req.URL.RawQuery, req.URL.String())
	}

	if len(req.TransferEncoding) > 0 {
		_, _ = fmt.Fprintf(&buf, "Transfer-Encoding: %s\r\n", strings.Join(req.TransferEncoding, ","))
	}

	if req.Close {
		_, _ = fmt.Fprintf(&buf, "Connection: close\r\n")
	}

	_ = req.Header.WriteSubset(&buf, map[string]bool{"Host": true, "Transfer-Encoding": true, "Trailer": true})

	if r.Body != nil {
		if bodyBytes, err := libIo.Copy(r.Body); err == nil {
			defer func() { r.Body = bytes.NewReader(bodyBytes) }()
			_, _ = fmt.Fprintf(&buf, "\r\n%s\r\n", string(bodyBytes))
		}
	}

	return buf.String()
}

func (r *Request) DumpResponse() string {
	var buf bytes.Buffer
	if r.Response != nil {
		_, _ = io.WriteString(&buf, fmt.Sprintf("Proto: %s\r\n", r.Response.Proto()))
		_, _ = io.WriteString(&buf, fmt.Sprintf("Status: %s\r\n", r.Response.Status()))
		_, _ = io.WriteString(&buf, fmt.Sprintf("StatusCode: %d\r\n", r.Response.StatusCode()))
		if h := r.Response.HeaderToString(); h != "" {
			_, _ = io.WriteString(&buf, h)
		}
		if body := r.Response.Body; body != "" {
			_, _ = io.WriteString(&buf, "\r\n"+body)
		}
	}
	return buf.String()
}

func (r *Request) DumpAll() string {
	return r.DumpRequest() + "\r\n" + r.DumpResponse()
}
