package encryption

import "net"

type none struct {

}
func (this *none) Init(arg string)(error){
	return nil
}
func (this *none)InitUser(conn net.Conn,arg *interface{})(error){
	return nil
}
func (this *none)Encrypt(arg *interface{},data []byte)([][]byte,error){
	list:=make([][]byte,0)
	list=append(list,data)
	return list,nil
}
func (this *none)Decrypt(arg *interface{},data []byte)([]byte,error){
	return data,nil
}
