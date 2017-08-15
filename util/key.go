package util

import (
	"crypto/md5"
	"fmt"
)


func GetDoubleMd5(str string) string {
	test := md5.New()
	test.Write([]byte(str))
	doublemd5 := fmt.Sprintf("%x", test.Sum(nil))
	test.Reset()
	test.Write([]byte(doublemd5))
	return fmt.Sprintf("%x", test.Sum(nil))
}