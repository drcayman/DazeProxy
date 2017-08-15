package encryption

import (
	"strings"
	"net"
)
type EncryptionAction interface {
	Init(string,*interface{})(error)
	InitUser(net.Conn,*interface{},*interface{})(error)
	Encrypt(*interface{},*interface{},[]byte) ([][]byte,error)
	Decrypt(*interface{},*interface{},[]byte) ([]byte,error)
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

	//这里添加自己开发的加密模块
	encryptionMap["none"]=func()(EncryptionAction){
		return EncryptionAction(&none{})
	}
	encryptionMap["psk-aes-cfb"]=func()(EncryptionAction){
		return EncryptionAction(&PskAesCfb{})
	}
	encryptionMap["psk-aes-256-cfb"]=func()(EncryptionAction){
		return EncryptionAction(&PskAes256Cfb{})
	}
	encryptionMap["psk-rc4-md5"]=func()(EncryptionAction){
		return EncryptionAction(&PskRc4Md5{})
	}
	encryptionMap["keypair-aes"]=func()(EncryptionAction){
		return EncryptionAction(&KeypairAes{})
	}
}