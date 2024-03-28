package dingtalk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func SendMessage(token string, content string, isAtAll bool, timeout time.Duration) (string, error) {
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content + "\n"},
		"at":      map[string]interface{}{"isAtAll": isAtAll},
	}
	uri := "https://oapi.dingtalk.com/robot/send?access_token=" + token
	b, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), err
}
