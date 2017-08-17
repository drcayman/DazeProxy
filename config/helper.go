package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStruct struct{
	IsDebug bool
	ProxyUnit []ProxyUnitStruct
}
type ProxyUnitStruct struct{
	Port string
	Disguise string
	DisguiseParam string
	Encryption string
	EncryptionParam string
	IPv6ResolvePrefer bool
	NoAuth bool
}
var Config ConfigStruct
func init(){
	buf,ReadErr:=ioutil.ReadFile("config.json")
	if ReadErr!=nil{
		fmt.Println("配置文件读取错误(config.json)！")
	}
	JsonErr:=json.Unmarshal(buf,&Config)
	if JsonErr!=nil{
		fmt.Println("配置文件格式错误！请参考DefaultConfig.json并严格按照JSON格式来修改config.json(",JsonErr.Error(),")")
		os.Exit(-1)
	}
	fmt.Println("配置文件读取成功：")
	fmt.Println("    调试模式：",Config.IsDebug)
	fmt.Printf("一共%d个代理服务单元\n",len(Config.ProxyUnit))
}