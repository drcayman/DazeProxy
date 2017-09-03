package obscure

import (
	"net"
	"net/http"

	"bufio"
	"time"
	"math/rand"
	"strconv"
	"errors"
	"bytes"
	"github.com/crabkun/DazeProxy/util"
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
	buffer:=bytes.NewBuffer([]byte("HTTP/1.1 200 OK\r\nServer: nginx\r\nDate: "+
		time.Now().Format("Mon,2 Jan 2006 15:04:05 MST")+
		"\r\nContent-Type: text/html; charset=gbk\r\nContent-Length: "+strconv.Itoa(ContentLength)+"\r\n"+
		"Connection: keep-alive\r\nCache-Control: no-cache\r\n\r\n"))
	buffer.Write([]byte(util.GetRandomString(ContentLength)))
	conn.Write(buffer.Bytes())
	_,err=http.ReadRequest(bufio.NewReader(conn))
	if err!=nil{
		return err
	}
	return nil
}
