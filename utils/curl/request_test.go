package curl

import (
	"fmt"
	"github.com/ZYallers/golib/consts"
	"net/http"
	"testing"
	"time"
)

func TestRequest_Get(t *testing.T) {
	req := NewRequest("https://httpbin.org/get").
		EnableTrace().
		SetHeader("Now-Time", time.Now().Format(consts.TimeFormat)).
		SetQueries(map[string]string{"name": "jack"})
	req.SetTimeOut(3 * time.Second)
	resp, err := req.Get()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("ResponseTotalTime", resp.TotalTime())
	fmt.Print("\r\n")

	//fmt.Println("resp.Body")
	//fmt.Println(resp.Body)

	fmt.Print("\r\n")
	fmt.Println(req.TraceInfo())

	fmt.Print("\r\n")
	fmt.Println(req.DumpAll())
}

func TestRequest_Post(t *testing.T) {
	req := NewRequest("https://httpbin.org/post").
		EnableTrace().
		SetContentType(FormUrlEncodedContentType).
		SetHeader("Now-Time", time.Now().Format(consts.TimeFormat)).
		SetQueries(map[string]string{"age": "12", "from": "baidu"})
	req.SetPostData(map[string]interface{}{"Name": "HaiMei"})
	req.SetTimeOut(3 * time.Second)
	req.SetBody(`{"Name":"Jack"}`)
	//req.SetBody(strings.NewReader(`{"Name":"peter"}`))

	resp, err := req.SetMethod(http.MethodPost).Send()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("ResponseTotalTime", resp.TotalTime())
	fmt.Print("\r\n")

	//fmt.Println("resp.Body")
	//fmt.Println(resp.Body)

	fmt.Print("\r\n")
	fmt.Println(req.TraceInfo())

	fmt.Print("\r\n")
	fmt.Println(req.DumpAll())
}
