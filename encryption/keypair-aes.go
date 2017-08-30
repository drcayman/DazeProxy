package encryption

import (
	"net"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"strings"
	"time"
	"strconv"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
)

type KeypairAes struct {
	reserved string
}
type KeypairAesTmp struct {
	Key []byte
	Block cipher.Block
}
func (this *KeypairAes) Init(param string,server *interface{})(error){
	privateKey, err := rsa.GenerateKey(rand.Reader,2048)
	if err!=nil{
		return err
	}
	*server=privateKey
	return nil
}
func (this *KeypairAes)InitUser(conn net.Conn,client *interface{},server *interface{})(error){
	key,flag:=(*server).(*rsa.PrivateKey)
	if !flag{
		return errors.New("unknown error")
	}

	utc:=time.Now().UTC()
	s,ParseDuration:=time.ParseDuration(utc.Format("-15h04m05s"))
	if ParseDuration!=nil{
		return ParseDuration
	}
	utc=utc.Add(s)
	UTCunix:=utc.Unix()
	UTCunixStr:=strconv.FormatInt(UTCunix,10)
	UTCunixStrPadded:=this.StrPadding(UTCunixStr)

	aesKey,GenMd5Err:=this.GenMd5Key(UTCunixStrPadded)
	if GenMd5Err!=nil{
		return GenMd5Err
	}

	Cipher,CipherErr:=aes.NewCipher(aesKey)
	if CipherErr!=nil{
		return CipherErr
	}
	enc:=cipher.NewCFBEncrypter(Cipher,aesKey[:Cipher.BlockSize()])
	pubkey,GenPublicKeyErr:=x509.MarshalPKIXPublicKey(&key.PublicKey)
	if GenPublicKeyErr!=nil{
		return errors.New("unknown error")
	}
	keyEncoded:=make([]byte,len(pubkey))
	enc.XORKeyStream(keyEncoded,pubkey)
	conn.Write(keyEncoded)
	pos:=0
	buf:=make([]byte,256)
	for pos<256{
		n,err:=conn.Read(buf[pos:])
		if err!=nil{
			return errors.New("客户端在握手期间断开连接"+err.Error())
		}
		pos+=n
	}
	DecryptBuf,DecryptErr:=rsa.DecryptPKCS1v15(rand.Reader,key,buf)
	if DecryptErr!=nil{
		return errors.New("无法解密客户端发送过来的数据")
	}

	t:=KeypairAesTmp{}
	t.Block,CipherErr=aes.NewCipher(DecryptBuf)
	if CipherErr!=nil{
		return CipherErr
	}
	t.Key=DecryptBuf[:t.Block.BlockSize()]
	*client=t
	return nil
}
func (this *KeypairAes)Encrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	t,flag:=(*client).(KeypairAesTmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Crypter:=cipher.NewCFBEncrypter(t.Block,t.Key)
	Crypter.XORKeyStream(dst,data)
	return dst,nil
}
func (this *KeypairAes)Decrypt(client *interface{},server *interface{},data []byte)([]byte,error){
	t,flag:=(*client).(KeypairAesTmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Decrypter:=cipher.NewCFBDecrypter(t.Block,t.Key)
	Decrypter.XORKeyStream(dst,data)
	return dst,nil
}
func (this *KeypairAes)StrPadding(str string) string {
	l:=16-len(str)
	newstr:=str+strings.Repeat("0",l)
	return newstr
}
func (this *KeypairAes) GenMd5Key(key string) ([]byte,error){
	test := md5.New()
	_,err:=test.Write([]byte(key))
	if err!=nil{
		return nil,err
	}
	return test.Sum(nil),nil
}