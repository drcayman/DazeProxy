package encryption

import (
	"net"
	"reflect"
	"errors"
)
type EncryptionAction interface {
	//Init，代理实例初始化时执行的操作
	//param：配置文件里面填写的EncryptionParam
	Init(param string)(error)

	//InitUser，用户连接后进行的初始化操作
	//conn：用户的连接套接字
	//client：此用户对象中为加密模块预留的空间
	InitUser(conn net.Conn,	client *interface{})(error)

	//Encrypt，加密
	//client同上
	//data：源数据
	//返回 加密后的数据 与 一个error(若发生了错误)
	Encrypt(client *interface{},	data []byte) ([]byte,error)

	//Decrypt，解密
	//client同上
	//data：加密数据
	//返回 解密后的数据 与 一个error(若发生了错误)
	Decrypt(client *interface{},	data []byte) ([]byte,error)
}
var encryptionMap map[string]reflect.Type

func GetEncryption(name string) (EncryptionAction,bool){
	if encryptionMap==nil{
		goto FAILED
	}
	if v,ok:=encryptionMap[name];ok{
		return reflect.New(v).Interface().(EncryptionAction),true
	}
FAILED:
	return nil,false
}
func GetEncryptionList()[]string{
	list:=make([]string,0)
	for k,_:=range encryptionMap{
		list=append(list, k)
	}
	return list
}
func RegisterEncryption(name string,action EncryptionAction)(error){
	if encryptionMap==nil{
		encryptionMap=make(map[string]reflect.Type)
	}
	if _,ok:=encryptionMap[name];ok{
		return errors.New("exist")
	}
	Ptype:=reflect.ValueOf(action)
	STtype:=reflect.Indirect(Ptype).Type()
	encryptionMap[name]=STtype
	return nil
}

