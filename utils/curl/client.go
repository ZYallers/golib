package curl

import (
	"net/http"
	"time"
)

const clientTimeout = 15 * time.Second

// 自定义http.Client，在原来http.DefaultTransport的基础上进行性能优化
// @see https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance
var Client = (func() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	return &http.Client{Transport: t, Timeout: clientTimeout}
})()
