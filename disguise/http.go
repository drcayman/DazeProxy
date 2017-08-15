package disguise

import (
	"net"
	"net/http"
	"bufio"
	"time"
)
type HTTP struct {
	reserved string
}

func (this *HTTP) Init(arg string,server *interface{})(error){
	return nil
}
func (this *HTTP) Action(conn net.Conn ,client *interface{}, server *interface{}) (error){
	_,err:=http.ReadRequest(bufio.NewReader(conn))
	if err!=nil{
		return err
	}
	rsp:=http.Response{
		ProtoMajor:1,
		ProtoMinor:1,
		StatusCode:200,
		Body:nil,
	}
	rsp.Header.Add("Content-Type","text/html; charset=gbk")
	rsp.Header.Add("Connection","keep-alive")
	rsp.Header.Add("Date",time.Now().Format("Mon,2 Jan 2006 15:04:05 MST"))
	rsp.Write(conn)
	conn.Read(make([]byte,1))
	return nil
}
