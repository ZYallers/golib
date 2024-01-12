package nets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

func SystemIP() string {
	if netInterfaces, err := net.Interfaces(); err == nil {
		for i := 0; i < len(netInterfaces); i++ {
			if (netInterfaces[i].Flags & net.FlagUp) != 0 {
				addrs, _ := netInterfaces[i].Addrs()
				for _, address := range addrs {
					if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							return ipnet.IP.String()
						}
					}
				}
			}
		}
	}
	return "unknown"
}

func GetIPByPconline(ip string) string {
	var result, url = ip, "http://whois.pconline.com.cn/ipJson.jsp?json=true"
	if ip != "" {
		url += "&ip=" + ip
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result
	}
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	body := mahonia.NewDecoder("GBK").ConvertString(string(b))
	if body == "" {
		return result
	}
	info := struct{ Ip, Addr string }{}
	if json.Unmarshal([]byte(body), &info) != nil {
		return result
	}
	if info.Ip != "" && info.Addr != "" {
		result = fmt.Sprintf("%s %s", info.Ip, strings.ReplaceAll(info.Addr, " ", ""))
	}
	return result
}

func PublicIP() string {
	return GetIPByPconline("")
}

func ClientIP(ip string) string {
	return GetIPByPconline(ip)
}
