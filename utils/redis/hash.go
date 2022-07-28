package redis

import (
	"errors"
	"github.com/ZYallers/golib/funcs/arrays"
	"github.com/ZYallers/golib/utils/json"
	"strings"
)

const hashAllFieldKey = "all"

func (r *Redis) HGetAll(key string) (result []interface{}) {
	all := r.Client().HGet(key, hashAllFieldKey).Val()
	if all == "" {
		return
	}
	keys := arrays.RemoveDuplicateWithString(strings.Split(all, ","))
	if len(keys) == 0 {
		return
	}
	result = r.Client().HMGet(key, keys...).Val()
	return
}

func (r *Redis) HMSet(key string, data map[string]interface{}) error {
	fields := make([]string, 0)
	fieldValues := make(map[string]interface{}, 0)
	for k, v := range data {
		if k == "" || v == nil {
			continue
		}
		if b, err := json.Marshal(v); err == nil {
			fieldValues[k] = string(b)
			fields = append(fields, k)
		}
	}

	if len(fields) == 0 {
		return errors.New("the data that can be saved is empty")
	}

	if val := r.Client().HGet(key, hashAllFieldKey).Val(); val != "" {
		fields = append(fields, strings.Split(val, ",")...)
	}

	var allFieldValue string
	if len(fields) > 0 {
		allFieldValue = strings.Join(arrays.RemoveDuplicateWithString(fields), ",")
	}
	fieldValues[hashAllFieldKey] = allFieldValue
	return r.Client().HMSet(key, fieldValues).Err()
}

func (r *Redis) HMDel(key string, fields ...string) error {
	newFields := make([]string, 0)
	if val := r.Client().HGet(key, hashAllFieldKey).Val(); val != "" {
		newFields = append(newFields, strings.Split(val, ",")...)
	}
	if len(newFields) > 0 {
		for _, field := range fields {
			newFields = arrays.RemoveWithString(newFields, field)
		}
	}

	var allFieldValue string
	if len(newFields) > 0 {
		allFieldValue = strings.Join(newFields, ",")
	}

	pl := r.Client().Pipeline()
	pl.HDel(key, fields...)
	pl.HSet(key, hashAllFieldKey, allFieldValue)
	_, err := pl.Exec()
	return err
}
