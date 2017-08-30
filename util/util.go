package util

import (
	"time"
	"math/rand"
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

