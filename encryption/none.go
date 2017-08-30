package encryption

import (
	"net"
)
//none-无加密
type none struct {
	RegArg string
}
func (this *none)Init(param string,server *interface{})(error){
	return nil
}
func (this *none)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	return nil
}
func (this *none)Encrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	return data,nil
}
func (this *none)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	return data,nil
}
