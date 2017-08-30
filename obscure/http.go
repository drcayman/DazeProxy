package obscure

import (
	"net"
	"net/http"

	"bufio"
	"time"
	"math/rand"
	"strconv"
	"errors"
)
type Http struct {
	RegArg string
}
func (this *Http) Init(param string,server *interface{})(error) {
	return nil
}

func (this *Http) Action(conn net.Conn , server *interface{}) (error){
	var err error
	_,err=http.ReadRequest(bufio.NewReader(conn))
	if err!=nil{
		return err
	}
	ContentLength:=1+rand.Intn(256)
	conn.Write([]byte("HTTP/1.1 200 OK\r\nServer: nginx\r\nDate: "+
		time.Now().Format("Mon,2 Jan 2006 15:04:05 MST")+
		"\r\nContent-Type: text/html; charset=gbk\r\nContent-Length: "+strconv.Itoa(ContentLength)+"\r\n"+
		"Connection: keep-alive\r\nCache-Control: no-cache\r\n\r\n"))
	//conn.Write([]byte(util.GetRandomString(ContentLength)))
	SafeRead(conn,ContentLength)
	return nil
}
func SafeRead(conn net.Conn,length int) ([]byte,error) {
	buf:=make([]byte,length)
	for pos:=0;pos<length;{
		n,err:=conn.Read(buf[pos:])
		if err!=nil {
			return nil,errors.New("根据Content-Length读取负载错误")
		}
		pos+=n
	}
	return buf,nil
}