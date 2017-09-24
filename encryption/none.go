package encryption

import (
	"net"
)
//none-无加密
type none struct {
}
func (this *none)Init(param string)(error){
	return nil
}
func (this *none)InitUser(conn net.Conn,client *interface{})(error){
	return nil
}
func (this *none)Encrypt(client *interface{},data []byte)([]byte,error){
	return data,nil
}
func (this *none)Decrypt(client *interface{},data []byte)([]byte,error){
	return data,nil
}
func init(){
	RegisterEncryption("none",new(none))
}
