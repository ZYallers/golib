package curl

import (
	"fmt"
	"github.com/ZYallers/golib/consts"
	"net/http"
	"testing"
	"time"
)

func TestRequest_Get(t *testing.T) {
	req := NewRequest("https://httpbin.org/get").EnableTrace().
		SetHeader("Now-Time", time.Now().Format(consts.TimeFormat)).
		SetQueries(map[string]string{"test": "abc123"})
	//req.SetTimeOut(time.Second)
	resp, err := req.SetMethod(http.MethodGet).Send()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("resp.Body")
	fmt.Println(resp.Body)

	fmt.Print("\r\n")
	fmt.Println(req.TraceInfo())

	fmt.Print("\r\n")
	fmt.Println(req.DumpAll())
}

func TestRequest_Send_Post(t *testing.T) {
	req := NewRequest("https://httpbin.org/post").EnableTrace().
		SetContentType(FormUrlEncodedContentType).
		SetHeader("Now-Time", time.Now().Format(consts.TimeFormat)).
		SetQueries(map[string]string{"test": "abc123"})
	req.SetPostData(map[string]interface{}{"demo": "def456"})
	//req.SetTimeOut(time.Second)
	//req.SetBody(`{"example":"peter"}`)
	resp, err := req.SetMethod(http.MethodPost).Send()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("resp.Body")
	fmt.Println(resp.Body)

	fmt.Print("\r\n")
	fmt.Println(req.TraceInfo())

	fmt.Print("\r\n")
	fmt.Println(req.DumpAll())
}
