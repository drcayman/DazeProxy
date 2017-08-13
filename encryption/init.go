package encryption

import (
	"strings"
	"net"
)
type EncryptionAction interface {
	Init(string)(error)
	InitUser(net.Conn,*interface{})(error)
	Encrypt(*interface{},[]byte) ([][]byte,error)
	Decrypt(*interface{},[]byte) ([]byte,error)
}
type regfunc func()(EncryptionAction)
var encryptionMap map[string]regfunc

func GetEncryption(name string) (regfunc,bool){
	name=strings.ToLower(name)
	d,flag:=encryptionMap[name]
	return d,flag
}

func init(){
	encryptionMap=make(map[string]regfunc)

	//这里添加自己开发的伪装模块
	encryptionMap["none"]=func()(EncryptionAction){
		return EncryptionAction(new(none))
	}
}