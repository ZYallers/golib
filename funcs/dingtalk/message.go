package dingtalk

import (
	"github.com/ZYallers/golib/utils/curl"
	"time"
)

var header = map[string]string{"Content-Type": "application/json;charset=utf-8"}

func SendMessage(token, content string, isAtAll bool, timeout time.Duration) (string, error) {
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	uri := "https://oapi.dingtalk.com/robot/send?access_token=" + token
	resp, err := curl.NewRequest(uri).SetHeaders(header).SetPostData(postData).SetTimeOut(timeout).Post()
	if resp == nil {
		return "", err
	}
	return resp.Body, err
}
