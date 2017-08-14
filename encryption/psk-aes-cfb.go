package encryption

import (
	"net"
	"crypto/aes"
	"crypto/cipher"
	"github.com/pkg/errors"
	"strings"
	"DazeClient/util"
)

type PskAesCfb struct {

}

func (this *PskAesCfb) KeyPadding(key string) ([]byte,error){
	//16, 24, or 32
	l:= len(key)
	if l>32{
		return nil,errors.New("key too long")
	}
	if l==0{
		return nil,errors.New("key is null")
	}
	if l<=16{
		p:=16-l
		newkey:=key+strings.Repeat("0",p)
		return util.S2b(&newkey),nil
	}
	if l<=24{
		p:=24-l
		newkey:=key+strings.Repeat("0",p)
		return util.S2b(&newkey),nil
	}
	if l<=32{
		p:=32-l
		newkey:=key+strings.Repeat("0",p)
		return util.S2b(&newkey),nil
	}
	return nil,errors.New("unknown error")
}
func (this *PskAesCfb) Init(arg string,server *interface{})(error){
	key,KeyPaddingErr:=this.KeyPadding(arg)
	if KeyPaddingErr!=nil{
		return KeyPaddingErr
	}
	*server=key
	return nil
}
func (this *PskAesCfb)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	return nil
}
func (this *PskAesCfb)Encrypt(client *interface{},server *interface{},data []byte)([][]byte,error){
	key,flag:=(*server).([]byte)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	block,CipherErr:=aes.NewCipher(key)
	if CipherErr!=nil{
		return nil,CipherErr
	}
	crypter:=cipher.NewCFBEncrypter(block,key)
	crypter.XORKeyStream(dst,data)
	list:=make([][]byte,0)
	list=append(list,dst)
	return list,nil
}
func (this *PskAesCfb)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	key,flag:=(*server).([]byte)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	block,CipherErr:=aes.NewCipher(key)
	if CipherErr!=nil{
		return nil,CipherErr
	}
	decrypter:=cipher.NewCFBDecrypter(block,key)
	decrypter.XORKeyStream(dst,data)
	return dst,nil
}
