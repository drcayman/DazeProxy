package util

import (
	"unsafe"
	"reflect"
)

func B2s(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

func S2b(s *string) []byte {
	return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(s))))
}