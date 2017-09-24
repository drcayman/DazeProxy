package obscure

import (
	"net"
	mrand "math/rand"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
	"crypto/rsa"
	"crypto/tls"
	"encoding/pem"
	"bytes"
	"github.com/crabkun/DazeProxy/util"
)
type TlsHandshake struct {
	cert tls.Certificate
}

func (this *TlsHandshake) Init(param string)(error){
	privateKey, err := rsa.GenerateKey(rand.Reader,2048)
	if err!=nil{
		return err
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{util.GetRandomString(mrand.Intn(64))},
		},
		NotBefore: time.Now().UTC(),
		NotAfter:  time.Now().Add(time.Hour*24*1024),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:[]string{util.GetRandomString(mrand.Intn(64))+"."+util.GetRandomString(mrand.Intn(3))},
	}
	CertBuf,err := x509.CreateCertificate(rand.Reader,&template,&template,&privateKey.PublicKey,privateKey)
	certPemBuf:=bytes.NewBuffer(make([]byte,0))
	pem.Encode(certPemBuf, &pem.Block{
		Type: "CERTIFICATE",
		Bytes: CertBuf,
	})
	KeyPemBuf:=bytes.NewBuffer(make([]byte,0))
	pem.Encode(KeyPemBuf, &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return err
	}
	this.cert,err=tls.X509KeyPair(
		certPemBuf.Bytes(),
		KeyPemBuf.Bytes(),
	)
	if err!=nil{
		return err
	}
	return nil
}

func (this *TlsHandshake) Action(conn net.Conn) (error){
	c:=tls.Server(conn,&tls.Config{Certificates:[]tls.Certificate{this.cert}})
	return c.Handshake()
}
func init(){
	RegisterObscure("tls_handshake",new(TlsHandshake))
}

