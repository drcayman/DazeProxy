package encryption

import (
	"net"
	"crypto/rsa"
	"crypto/rand"
	"errors"
	"crypto/aes"
	"crypto/cipher"
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
	this.privateKey, err = rsa.GenerateKey(rand.Reader,8*(128+mrand.Intn(127)))
	return err
}
func (this *KeypairAes)InitUser(conn net.Conn,client *interface{})(error){
	var err error
	keylen:=this.privateKey.PublicKey.N.BitLen()/8
	keyBuf:=make([]byte,keylen+1)
	keyBuf[0]=byte(keylen)
	copy(keyBuf[1:],this.privateKey.PublicKey.N.Bytes())
	conn.Write(keyBuf)
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