package helper

import (
	"time"
)

const robotUriPrefix = "https://oapi.dingtalk.com/robot/send?access_token="

var header = map[string]string{"Content-Type": "application/json;charset=utf-8"}

func SendMessage(token, content string, isAtAll bool, timeout time.Duration) (*Response, error) {
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	uri := robotUriPrefix + token
	return NewRequest(uri).SetHeaders(header).SetPostData(postData).SetTimeOut(timeout).Post()
}
