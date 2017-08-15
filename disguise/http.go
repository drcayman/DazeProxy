package disguise

import (
	"net"
	"net/http"
	"bufio"
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
		Proto:"HTTP/1.1",
		StatusCode:200,
		Body:nil,
	}
	rsp.Write(conn)
	conn.Read(make([]byte,1))
	return nil
}
