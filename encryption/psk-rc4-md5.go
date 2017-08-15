package encryption

import (
	"net"
	"crypto/rc4"
	"github.com/pkg/errors"
	"crypto/md5"
)

type PskRc4Md5 struct {
	reserved string
}
type PskRc4Md5Tmp struct {
	Cipher *rc4.Cipher
}
func (this *PskRc4Md5) GenKey(key string) ([]byte,error){
	test := md5.New()
	_,err:=test.Write([]byte(key))
	if err!=nil{
		return nil,err
	}
	return test.Sum(nil),nil
}
func (this *PskRc4Md5) Init(arg string,server *interface{})(error){
	key,GenKeyErr:=this.GenKey(arg)
	if GenKeyErr!=nil{
		return GenKeyErr
	}
	t:=PskRc4Md5Tmp{}
	var CipherErr error=nil
	t.Cipher,CipherErr=rc4.NewCipher(key)
	if CipherErr!=nil{
		return CipherErr
	}
	*server=t
	return nil
}
func (this *PskRc4Md5)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	return nil
}
func (this *PskRc4Md5)Encrypt(client *interface{},server *interface{},data []byte)([][]byte,error){
	t,flag:=(*server).(PskRc4Md5Tmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	t.Cipher.Reset()
	t.Cipher.XORKeyStream(dst,data)
	list:=make([][]byte,0)
	list=append(list,dst)
	return list,nil
}
func (this *PskRc4Md5)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	t,flag:=(*server).(PskRc4Md5Tmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	t.Cipher.Reset()
	t.Cipher.XORKeyStream(dst,data)
	return dst,nil
}
