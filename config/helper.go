package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStruct struct{
	IsDebug bool
	IPv4Port string
	IPv6Port string
	IPv6ResolvePrefer bool
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
	fmt.Println("    IPv4端口号：",Config.IPv4Port)
	fmt.Println("    IPv6端口号：",Config.IPv6Port)
	fmt.Println("    优先解析域名的IPv6地址：",Config.IPv6ResolvePrefer)
}