package encryption

import (
	"net"
	"crypto/aes"
	"crypto/cipher"
	"github.com/crabkun/DazeProxy/util"
)

type PskAesCfb struct {
	Key []byte
	Block cipher.Block
}

func (this *PskAesCfb) Init(param string)(error){
	key,GenKeyErr:=util.Gen16Md5Key(param)
	if GenKeyErr!=nil{
		return GenKeyErr
	}
	var CipherErr error=nil
	this.Block,CipherErr=aes.NewCipher(key)
	if CipherErr!=nil{
		return CipherErr
	}
	this.Key=key[:this.Block.BlockSize()]
	return nil
}
func (this *PskAesCfb)InitUser(conn net.Conn,client *interface{})(error){
	return nil
}
func (this *PskAesCfb)Encrypt(client *interface{},data []byte)([]byte,error){
	dst:=make([]byte,len(data))
	Crypter:=cipher.NewCFBEncrypter(this.Block,this.Key)
	Crypter.XORKeyStream(dst,data)
	return dst,nil
}
func (this *PskAesCfb)Decrypt(client *interface{},data []byte)([]byte,error){
	dst:=make([]byte,len(data))
	Decrypter:=cipher.NewCFBDecrypter(this.Block,this.Key)
	Decrypter.XORKeyStream(dst,data)
	return dst,nil
}
func init(){
	RegisterEncryption("psk-aes-128-cfb",new(PskAesCfb))
}