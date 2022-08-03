package php

import (
	"github.com/syyongx/php2go"
	"github.com/techoner/gophp/serialize"
	"strings"
)

func Unserialize(str string) map[string]interface{} {
	vars := make(map[string]interface{}, 10)
	offset := 0
	sl := php2go.Strlen(str)
	for offset < sl {
		if index := strings.Index(php2go.Substr(str, uint(offset), -1), "|"); index < 0 {
			break
		}

		pos := php2go.Strpos(str, "|", offset)
		num := pos - offset

		varname := php2go.Substr(str, uint(offset), num)
		offset += num + 1
		data, _ := serialize.UnMarshal([]byte(php2go.Substr(str, uint(offset), -1)))
		vars[varname] = data

		jsonbyte, _ := serialize.Marshal(data)
		offset += php2go.Strlen(string(jsonbyte))
	}
	return vars
}

func Serialize(vars map[string]interface{}) (str string) {
	for k, v := range vars {
		sa, _ := serialize.Marshal(v)
		str += k + "|" + string(sa)
	}
	return
}
