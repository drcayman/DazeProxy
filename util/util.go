package util

import (
	"time"
	"math/rand"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GetRandomString(strlen int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	byts := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < strlen; i++ {
		result = append(result, byts[r.Intn(len(byts))])
	}
	return string(result)
}
func StrPadding(str string,count int,char string) string {
	l:=count-len(str)
	newstr:=str+strings.Repeat(char,l)
	return newstr
}
func Gen16Md5Key(key string) ([]byte,error){
	test := md5.New()
	_,err:=test.Write([]byte(key))
	if err!=nil{
		return nil,err
	}
	return test.Sum(nil),nil
}
func Gen32Md5Key(key string) ([]byte,error){
	test := md5.New()
	_,err:=test.Write([]byte(key))
	if err!=nil{
		return nil,err
	}
	md5src:=test.Sum(nil)
	md5dst:=make([]byte,32)
	hex.Encode(md5dst,md5src)
	return md5dst,nil

}
func GetDoubleMd5(str string) string {
	test := md5.New()
	test.Write([]byte(str))
	doublemd5 := fmt.Sprintf("%x", test.Sum(nil))
	test.Reset()
	test.Write([]byte(doublemd5))
	return fmt.Sprintf("%x", test.Sum(nil))
}