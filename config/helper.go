package config

import (
	"io/ioutil"
	"encoding/json"
	"../log"
)

type ConfigStruct struct{
	ServerPort string
	IsDebug bool
}
var Config ConfigStruct
func init(){
	buf,ReadErr:=ioutil.ReadFile("config.json")
	if ReadErr!=nil{
		log.PrintPanic("配置文件读取错误(config.json)！")
	}
	JsonErr:=json.Unmarshal(buf,&Config)
	if JsonErr!=nil{
		log.PrintPanic("配置文件格式错误！请参考DefaultConfig.json并严格按照JSON格式来修改config.json(",JsonErr.Error(),")")
	}
	log.PrintSuccess("配置文件读取成功：")
	log.PrintSuccess("    端口号：",Config.ServerPort)
}