package encryption

import (
	"strings"
)
type EncryptionAction interface {
	Init(string)(error)
	Encrypt([]byte,int) ([]byte,int,error)
	Decrypt([]byte,int) ([]byte,int,error)
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