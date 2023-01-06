package curl

import (
	"bytes"
	"context"
	"fmt"
	io2 "github.com/ZYallers/golib/funcs/io"
	"github.com/ZYallers/golib/utils/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	JsonContentType           = "application/json"
	JsonUTF8ContentType       = "application/json;charset=utf-8"
	FormUrlEncodedContentType = "application/x-www-form-urlencoded"
)

type Request struct {
	client             *http.Client
	rawRequest         *http.Request
	trace              *clientTrace
	ctx                context.Context
	Response           *Response
	Body               io.Reader
	Method             string
	Url                string
	Timeout            time.Duration
	Headers            map[string]string
	Cookies            map[string]string
	Queries            map[string]string
	PostData           map[string]interface{}
	StartTime          time.Time
	responseReturnTime time.Time
}

// NewRequest new one request
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

// SetHeader set request header
func (r *Request) SetHeader(key, value string) *Request {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	r.Headers[key] = value
	return r
}

// SetHeaders set request headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.Headers = headers
	return r
}

func (r *Request) setHeaders() *Request {
	for k, v := range r.Headers {
		r.rawRequest.Header.Set(k, v)
	}
	return r
}

// SetCookies set request cookies
func (r *Request) SetCookies(cookies map[string]string) *Request {
	r.Cookies = cookies
	return r
}

func (r *Request) setCookies() *Request {
	for k, v := range r.Cookies {
		r.rawRequest.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return r
}

// SetQueries set request query
func (r *Request) SetQueries(queries map[string]string) *Request {
	r.Queries = queries
	return r
}

func (r *Request) setQueries() *Request {
	q := r.rawRequest.URL.Query()
	for k, v := range r.Queries {
		q.Add(k, v)
	}
	r.rawRequest.URL.RawQuery = q.Encode()
	return r
}

// SetPostData set post data
func (r *Request) SetPostData(data map[string]interface{}) *Request {
	if data != nil {
		r.PostData = data
		r.Body = nil
	}
	return r
}

func (r *Request) setPostData() error {
	if ct, ok := r.Headers["Content-Type"]; ok && strings.HasPrefix(ct, JsonContentType) {
		if bts, err := json.Marshal(r.PostData); err != nil {
			return err
		} else {
			r.Body = bytes.NewReader(bts)
			return nil
		}
	}

	// If the Content Type cannot be matched, the default method is application/x-www-form-urlencoded
	postData := url.Values{}
	for k, v := range r.PostData {
		postData.Add(k, fmt.Sprint(v))
	}
	r.Body = strings.NewReader(postData.Encode())

	return nil
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
	r.PostData = nil
	return r
}

// SetTimeOut set the timeout after
func (r *Request) SetTimeOut(timeout time.Duration) *Request {
	if timeout > 0 && timeout <= ClientTimeout {
		r.Timeout = timeout
	}
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

// DisableTrace disables trace.
func (r *Request) DisableTrace() *Request {
	r.trace = nil
	return r
}

// EnableTrace enables trace (http3 currently does not support trace).
func (r *Request) EnableTrace() *Request {
	if r.trace == nil {
		r.trace = &clientTrace{}
	}
	return r
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
			ti.TotalTime = endTime.Sub(r.StartTime)
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

// Get init get request and return response
func (r *Request) Get() (*Response, error) {
	return r.SetMethod(http.MethodGet).Send()
}

// Post init post request and return response
func (r *Request) Post() (*Response, error) {
	return r.SetMethod(http.MethodPost).Send()
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

// Send init request and return response
func (r *Request) Send() (*Response, error) {
	defer func() { r.responseReturnTime = time.Now() }()

	if r.Response == nil {
		r.Response = &Response{Request: r}
	}

	if r.PostData != nil {
		if err := r.setPostData(); err != nil {
			return r.Response, err
		}
	}

	if r.Body != nil {
		var backup io.Reader
		var err error
		if backup, r.Body, err = r.copyBody(); err != nil {
			return r.Response, err
		}
		defer func() { r.Body = backup }()
	}

	if req, err := http.NewRequest(r.Method, r.Url, r.Body); err != nil {
		return r.Response, err
	} else {
		if r.ctx == nil {
			r.ctx = context.Background()
		}
		if r.trace != nil {
			r.ctx = r.trace.createContext(r.ctx)
		}
		if r.Timeout <= 0 || r.Timeout > ClientTimeout {
			r.Timeout = ClientTimeout
		}
		ctx, cancel := context.WithTimeout(r.ctx, r.Timeout)
		defer cancel()
		r.rawRequest = req.WithContext(ctx)
	}

	r.setHeaders().setCookies().setQueries()
	r.StartTime = time.Now()
	httpResp, err := r.client.Do(r.rawRequest)
	if err == nil {
		r.Response.Raw = httpResp
		r.Response.readBody()
		r.Response.setReceivedAt()
	}

	return r.Response, err
}

var reqWriteExcludeHeaderDump = map[string]bool{
	"Host":              true, // not in Header map anyway
	"Transfer-Encoding": true,
	"Trailer":           true,
}

func (r *Request) DumpRequest() string {
	var b bytes.Buffer
	req := r.rawRequest
	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}
	_, _ = fmt.Fprintf(&b, "%s %s HTTP/%d.%d\r\n", valueOrDefault(req.Method, http.MethodGet), reqURI, req.ProtoMajor, req.ProtoMinor)

	absRequestURI := strings.HasPrefix(req.RequestURI, "http://") || strings.HasPrefix(req.RequestURI, "https://")
	if !absRequestURI {
		host := req.Host
		if host == "" && req.URL != nil {
			host = req.URL.Host
		}
		if host != "" {
			_, _ = fmt.Fprintf(&b, "Host: %s\r\n", host)
		}
	}

	if len(req.TransferEncoding) > 0 {
		_, _ = fmt.Fprintf(&b, "Transfer-Encoding: %s\r\n", strings.Join(req.TransferEncoding, ","))
	}
	if req.Close {
		_, _ = fmt.Fprintf(&b, "Connection: close\r\n")
	}

	_ = req.Header.WriteSubset(&b, reqWriteExcludeHeaderDump)
	if r.Body != nil {
		if bte, err := ioutil.ReadAll(r.Body); err == nil && bte != nil {
			_, _ = io.WriteString(&b, "\r\n"+string(bte)+"\r\n")
		}
	}
	return b.String()
}

func (r *Request) DumpResponse() string {
	var b bytes.Buffer
	if r.Response != nil {
		if h := r.Response.HeaderToString(); h != "" {
			_, _ = io.WriteString(&b, h)
		}
		if body := r.Response.Body; body != "" {
			_, _ = io.WriteString(&b, "\r\n"+body)
		}
	}
	return b.String()
}

func (r *Request) DumpAll() string {
	return r.DumpRequest() + "\r\n" + r.DumpResponse()
}

func (r *Request) copyBody() (io.Reader, io.Reader, error) {
	if cp, err := io2.Copy(r.Body); err != nil {
		return nil, r.Body, err
	} else {
		var bf bytes.Buffer
		_, _ = bf.Write(cp)
		return ioutil.NopCloser(&bf), ioutil.NopCloser(bytes.NewReader(cp)), nil
	}
}

func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
