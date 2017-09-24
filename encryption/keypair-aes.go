package encryption

import (
	"net"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"time"
	"strconv"
	"crypto/aes"
	"crypto/cipher"
	"github.com/crabkun/DazeProxy/util"
	mrand "math/rand"
	"io"
)

type KeypairAes struct {
	privateKey *rsa.PrivateKey
}
type KeypairAesTmp struct {
	Key []byte
	Block cipher.Block
}
func (this *KeypairAes) Init(param string)(error){
	var err error
	this.privateKey, err = rsa.GenerateKey(rand.Reader,8*(128+mrand.Intn(64)))
	if err!=nil{
		return err
	}
	return nil
}
func (this *KeypairAes)InitUser(conn net.Conn,client *interface{})(error){
	var err error
	utc:=time.Now().UTC()
	s,err:=time.ParseDuration(utc.Format("-15h04m05s"))
	if err!=nil{
		return err
	}
	utc=utc.Add(s)
	UTCunix:=utc.Unix()
	UTCunixStr:=strconv.FormatInt(UTCunix,10)
	UTCunixStrPadded:=util.StrPadding(UTCunixStr,16,"0")

	aesKey,err:=util.Gen16Md5Key(UTCunixStrPadded)
	if err!=nil{
		return err
	}
	Cipher,err:=aes.NewCipher(aesKey)
	if err!=nil{
		return err
	}
	enc:=cipher.NewCFBEncrypter(Cipher,aesKey[:Cipher.BlockSize()])
	pubkey,err:=x509.MarshalPKIXPublicKey(&this.privateKey.PublicKey)
	if err!=nil{
		return errors.New("unknown error")
	}
	keyEncoded:=make([]byte,len(pubkey))
	enc.XORKeyStream(keyEncoded,pubkey)
	keyEncodedBuf:=make([]byte,len(keyEncoded)+1)
	keyEncodedBuf[0]=byte(len(keyEncoded))
	copy(keyEncodedBuf[1:],keyEncoded)
	conn.Write(keyEncodedBuf)
	buf,err:=this.SafeRead(conn,this.privateKey.N.BitLen()/8)
	if err!=nil{
		return errors.New("无法接收客户端的密钥")
	}
	DecryptBuf,DecryptErr:=rsa.DecryptPKCS1v15(rand.Reader,this.privateKey,buf)
	if DecryptErr!=nil{
		return errors.New("无法解密客户端发送过来的数据")
	}

	t:=KeypairAesTmp{}
	t.Block,err=aes.NewCipher(DecryptBuf)
	if err!=nil{
		return err
	}
	t.Key=DecryptBuf[:t.Block.BlockSize()]
	*client=t
	return nil
}
func (this *KeypairAes)Encrypt(client *interface{},data []byte)([]byte,error){
	t,flag:=(*client).(KeypairAesTmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Crypter:=cipher.NewCFBEncrypter(t.Block,t.Key)
	Crypter.XORKeyStream(dst,data)
	return dst,nil
}
func (this *KeypairAes)Decrypt(client *interface{},data []byte)([]byte,error){
	t,flag:=(*client).(KeypairAesTmp)
	if !flag{
		return nil,errors.New("unknown error")
	}
	dst:=make([]byte,len(data))
	Decrypter:=cipher.NewCFBDecrypter(t.Block,t.Key)
	Decrypter.XORKeyStream(dst,data)
	return dst,nil
}
func (this *KeypairAes)SafeRead(conn net.Conn,length int)([]byte,error){
	buf:=make([]byte,length)
	_,err:=io.ReadFull(conn,buf)
	return buf,err
}
func init(){
	RegisterEncryption("keypair-aes",new(KeypairAes))
}