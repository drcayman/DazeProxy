package encryption

import "net"

type none struct {
	reserved string
}
func (this *none) Init(arg string,server *interface{})(error){
	return nil
}
func (this *none)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	return nil
}
func (this *none)Encrypt(client *interface{},server *interface{},data []byte)([][]byte,error){
	list:=make([][]byte,0)
	list=append(list,data)
	return list,nil
}
func (this *none)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	return data,nil
}
