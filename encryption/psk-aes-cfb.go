package encryption

import (
	"net"
	"crypto/aes"
	"crypto/cipher"
	"github.com/pkg/errors"
	"crypto/md5"
)

type PskAesCfb struct {
	reserved string
}
type PskAesCfbTmp struct {
	Key []byte
	Block cipher.Block
}
func (this *PskAesCfb) GenKey(key string) ([]byte,error){
	test := md5.New()
	_,err:=test.Write([]byte(key))
	if err!=nil{
		return nil,err
	}
	return test.Sum(nil),nil
}
func (this *PskAesCfb) Init(arg string,server *interface{})(error){
	key,GenKeyErr:=this.GenKey(arg)
	if GenKeyErr!=nil{
		return GenKeyErr
	}
	t:=PskAesCfbTmp{}
	var CipherErr error=nil
	t.Block,CipherErr=aes.NewCipher(key)
	if CipherErr!=nil{
		return CipherErr
	}
	t.Key=key[:t.Block.BlockSize()]
	*server=t
	return nil
}
func (this *PskAesCfb)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	return nil
}
func (this *PskAesCfb)Encrypt(client *interface{},server *interface{},data []byte)([][]byte,error){
	t,flag:=(*server).(PskAesCfbTmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Crypter:=cipher.NewCFBEncrypter(t.Block,t.Key)
	Crypter.XORKeyStream(dst,data)
	list:=make([][]byte,0)
	list=append(list,dst)
	return list,nil
}
func (this *PskAesCfb)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	t,flag:=(*server).(PskAesCfbTmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Decrypter:=cipher.NewCFBDecrypter(t.Block,t.Key)
	Decrypter.XORKeyStream(dst,data)
	return dst,nil
}
