package curl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ZYallers/golib/utils/json"
	"github.com/ZYallers/golib/utils/trace"
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

	closeConn          bool
	openTracing        bool
	error              error
	startTime          time.Time
	responseReturnTime time.Time
	body               string
	client             *http.Client
	rawRequest         *http.Request
	clientTrace        *clientTrace
}

// NewRequest new request
func NewRequest(url string) *Request {
	r := &Request{Url: url, client: Client, Timeout: ClientTimeout}
	r.Response = &Response{Request: r}
	return r
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
	if len(headers) == 0 {
		return r
	}
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	for k, v := range headers {
		r.Headers[k] = v
	}
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
	if len(cookies) == 0 {
		return r
	}
	if r.Cookies == nil {
		r.Cookies = map[string]string{}
	}
	for k, v := range cookies {
		r.Cookies[k] = v
	}
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
	if len(queries) == 0 {
		return r
	}
	if r.Queries == nil {
		r.Queries = map[string]string{}
	}
	for k, v := range queries {
		r.Queries[k] = v
	}
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
	if len(data) == 0 {
		return r
	}
	if r.PostData == nil {
		r.PostData = map[string]interface{}{}
	}
	for k, v := range data {
		r.PostData[k] = v
	}
	if ct, ok := r.Headers["Content-Type"]; ok && strings.HasPrefix(ct, JsonContentType) {
		if pb, err := json.Marshal(r.PostData); err != nil {
			r.error = err
		} else {
			r.body = string(pb)
		}
		return r
	}
	posts := url.Values{}
	for k, v := range r.PostData {
		posts.Set(k, fmt.Sprint(v))
	}
	r.body = posts.Encode()
	return r
}

// SetBody set the request body, accepts string, []byte, bytes.Buffer, io.Reader, io.ReadCloser.
func (r *Request) SetBody(body interface{}) *Request {
	if body == nil {
		return r
	}

	switch b := body.(type) {
	case string:
		r.body = b
	case []byte:
		r.body = string(b)
	case bytes.Buffer:
		r.body = b.String()
	case io.Reader:
		if cb, err := ioCopy(b); err == nil {
			r.body = string(cb)
		}
	case io.ReadCloser:
		if cb, err := ioCopy(b); err == nil {
			r.body = string(cb)
		}
	default:
		r.body = fmt.Sprint(body)
	}

	return r
}

// GetBody get the request set body
func (r *Request) GetBody() string {
	return r.body
}

// GetError get the request handle process error
func (r *Request) GetError() error {
	return r.error
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

// CloseConn whether closes the connection after sending this request and reading its response if set to true in HTTP/1.1 and HTTP/2.
// Setting this field prevents re-use of TCP connections between requests to the same hosts event if EnableKeepAlives() were called.
func (r *Request) CloseConn(enable bool) *Request {
	r.closeConn = enable
	return r
}

// OpenTracing whether enable openTracing
func (r *Request) OpenTracing(enable bool) *Request {
	r.openTracing = enable
	return r
}

// EnableTrace enables trace (http3 currently does not support trace).
func (r *Request) EnableTrace() *Request {
	r.clientTrace = &clientTrace{}
	return r
}

// DisableTrace disables trace.
func (r *Request) DisableTrace() *Request {
	r.clientTrace = nil
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

	r.Response = &Response{Request: r}

	ctx := context.Background()
	if r.clientTrace != nil {
		ctx = r.clientTrace.createContext(ctx)
	}
	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	var body io.Reader
	if r.body != "" {
		body = strings.NewReader(r.body)
	}

	if r.rawRequest, r.error = http.NewRequestWithContext(ctx, r.Method, r.Url, body); r.error != nil {
		return r.Response, r.error
	}

	// Set close connection
	if r.closeConn {
		r.rawRequest.Close = true
	}

	// Set header TraceId
	if r.openTracing {
		r.setTraceId()
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
	ct := r.clientTrace

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

// DumpRequest dump request data
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

	if r.body != "" {
		_, _ = fmt.Fprintf(&buf, "\r\n%s\r\n", r.body)
	}

	return buf.String()
}

// DumpResponse dump response data
func (r *Request) DumpResponse() string {
	var buf bytes.Buffer
	if r.Response != nil {
		_, _ = io.WriteString(&buf, fmt.Sprintf("Proto: %s\r\n", r.Response.Proto()))
		_, _ = io.WriteString(&buf, fmt.Sprintf("Status: %s\r\n", r.Response.Status()))
		_, _ = io.WriteString(&buf, fmt.Sprintf("StatusCode: %d\r\n", r.Response.StatusCode()))
		if h := r.Response.HeaderToString(); h != "" {
			_, _ = io.WriteString(&buf, h)
		}
		if r.Response.Body != "" {
			_, _ = io.WriteString(&buf, "\r\n"+r.Response.Body)
		}
	}
	return buf.String()
}

// DumpAll dump request and response data
func (r *Request) DumpAll() string {
	return r.DumpRequest() + "\r\n" + r.DumpResponse()
}

// setTraceId set traceId to request header from current goroutine id
func (r *Request) setTraceId() {
	if traceId := trace.GetGoIdTraceId(); traceId != "" {
		r.SetHeader(trace.IdKey, traceId)
	}
}
