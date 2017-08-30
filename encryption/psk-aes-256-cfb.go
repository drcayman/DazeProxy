package encryption

import (
	"net"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"log"
	"github.com/crabkun/DazeProxy/util"
)

type PskAes256Cfb struct {
	reserved string
}
type PskAesCfb256Tmp struct {
	Key []byte
	Block cipher.Block
}

func (this *PskAes256Cfb) Init(param string,server *interface{})(error){
	key,GenKeyErr:=util.Gen32Md5Key(param)
	if GenKeyErr!=nil{
		return GenKeyErr
	}
	log.Println(len(key))
	t:=PskAesCfb256Tmp{}
	var CipherErr error=nil
	t.Block,CipherErr=aes.NewCipher(key)
	if CipherErr!=nil{
		return CipherErr
	}
	t.Key=key[:t.Block.BlockSize()]
	*server=t
	return nil
}
func (this *PskAes256Cfb)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	return nil
}
func (this *PskAes256Cfb)Encrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	t,flag:=(*server).(PskAesCfb256Tmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Crypter:=cipher.NewCFBEncrypter(t.Block,t.Key)
	Crypter.XORKeyStream(dst,data)
	return dst,nil
}
func (this *PskAes256Cfb)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	t,flag:=(*server).(PskAesCfb256Tmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Decrypter:=cipher.NewCFBDecrypter(t.Block,t.Key)
	Decrypter.XORKeyStream(dst,data)
	return dst,nil
}
