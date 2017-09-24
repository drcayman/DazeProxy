package obscure

import (
	"net"
	"reflect"
	"errors"
)
type ObscureAction interface {
	//Init，代理实例初始化时执行的操作
	//param：配置文件里面填写的ObscureParam
	Init(param string)(error)

	//Action，用户连接后进行的伪装操作
	//conn：用户的连接套接字
	Action(conn net.Conn)(error)
}
var obscureMap map[string]reflect.Type

func GetObscure(name string)(ObscureAction,bool){
	if obscureMap==nil{
		goto FAILED
	}
	if v,ok:=obscureMap[name];ok{
		return reflect.New(v).Interface().(ObscureAction),true
	}
FAILED:
	return nil,false
}

func GetObscureList()[]string{
	list:=make([]string,0)
	for k:=range obscureMap{
		list=append(list, k)
	}
	return list
}
func RegisterObscure(name string,action ObscureAction)(error){
	if obscureMap==nil{
		obscureMap=make(map[string]reflect.Type)
	}
	if _,ok:=obscureMap[name];ok{
		return errors.New("exist")
	}
	Ptype:=reflect.ValueOf(action)
	STtype:=reflect.Indirect(Ptype).Type()
	obscureMap[name]=STtype
	return nil
}
