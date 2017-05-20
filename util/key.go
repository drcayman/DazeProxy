package util

import (
	"encoding/pem"
	"crypto/x509"
	"os"
	"crypto/rsa"
	"crypto/rand"
	"fmt"
)

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
func CheckKeyAndGen()  {
	fmt.Println("[!]开始检查密钥文件")
	_,privateKeyErr:=os.Stat("private.pem")
	_,publicKeyErr:=os.Stat("private.pem")
	if privateKeyErr!=nil || publicKeyErr!=nil{
		fmt.Println("[!]密钥文件不存在，开始生成密钥文件")
		os.Remove("private.pem")
		os.Remove("public.pem")
		GenRsaKey(1024)
		fmt.Println("[√]密钥文件生成成功")
	}
	fmt.Println("[√]密钥文件检查完毕")
}
func CheckLicense(){

}