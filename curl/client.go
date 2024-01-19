package curl

import (
	"net/http"
	"time"
)

const (
	ClientTimeout             = 15 * time.Second
	JsonContentType           = "application/json"
	JsonUTF8ContentType       = "application/json;charset=utf-8"
	FormUrlEncodedContentType = "application/x-www-form-urlencoded"
)

// Client Customize http.Client to optimize the performance based on the original http.DefaultTransport
// @see https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance
var Client = (func() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	return &http.Client{Transport: t, Timeout: ClientTimeout}
})()
