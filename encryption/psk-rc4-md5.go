package encryption

import (
	"net"
	"crypto/rc4"
	"github.com/crabkun/DazeProxy/util"
)

type PskRc4Md5 struct {
	Cipher *rc4.Cipher
}

func (this *PskRc4Md5) Init(param string)(error){
	key,GenKeyErr:=util.Gen16Md5Key(param)
	if GenKeyErr!=nil{
		return GenKeyErr
	}
	var CipherErr error=nil
	this.Cipher,CipherErr=rc4.NewCipher(key)
	return CipherErr
}
func (this *PskRc4Md5)InitUser(conn net.Conn,client *interface{})(error){
	return nil
}
func (this *PskRc4Md5)Encrypt(client *interface{},data []byte)([]byte,error){
	dst:=make([]byte,len(data))
	this.Cipher.Reset()
	this.Cipher.XORKeyStream(dst,data)
	return dst,nil
}
func (this *PskRc4Md5)Decrypt(client *interface{},data []byte)([]byte,error){
	dst:=make([]byte,len(data))
	this.Cipher.Reset()
	this.Cipher.XORKeyStream(dst,data)
	return dst,nil
}
func init(){
	RegisterEncryption("psk-rc4-md5",new(PskRc4Md5))
}