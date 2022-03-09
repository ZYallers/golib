package helper

import (
	"reflect"
	"unsafe"
)

//  String2Bytes ...
//  @author Cloud|2021-12-12 18:24:23
//  @param s string ...
//  @return []byte ...
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//  Bytes2String ...
//  @author Cloud|2021-12-12 18:24:20
//  @param b []byte ...
//  @return string ...
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
