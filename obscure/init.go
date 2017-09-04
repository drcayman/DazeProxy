package obscure

import (
	"strings"
	"net"
)
type Action interface {
	//Init，代理实例初始化时执行的操作
	//param：配置文件里面填写的ObscureParam
	//server：此代理实例中为伪装模块预留的空间
	Init(param string,	server *interface{})(error)

	//Action，用户连接后进行的伪装操作
	//conn：用户的连接套接字
	//server：此代理实例中为伪装模块预留的空间
	Action(conn net.Conn,	server *interface{})(error)
}
type regfunc func()(Action)
var obscureMap map[string]regfunc

func GetObscure(name string) (regfunc,bool){
	name=strings.ToLower(name)
	d,flag:=obscureMap[name]
	return d,flag
}

func init(){
	obscureMap=make(map[string]regfunc)
	//自己开发的伪装模块必需在此注册

	obscureMap["none"]=func()(Action){
		return Action(&none{})
	}
	obscureMap["tls_handshake"]=func()(Action){
		return Action(&TlsHandshake{})
	}
	obscureMap["http"]=func()(Action){
		return Action(&Http{})
	}
}

