package helper

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/crabkun/DazeProxy/common"
	"log"
)
var IsDebug bool
func LoadConfig(config *common.S_config,logflag bool){
	buf,ReadErr:=ioutil.ReadFile("config.json")
	if ReadErr!=nil{
		log.Println("配置文件(config.json)读取错误："+ReadErr.Error())
		os.Exit(-1)
	}
	JsonErr:=json.Unmarshal(buf,config)
	if JsonErr!=nil{
		log.Println("配置文件格式错误！请严格按照JSON格式来修改config.json(",JsonErr.Error(),")")
		os.Exit(-1)
	}
	log.Println("配置文件读取成功：")
	log.Println("    调试模式：",!logflag && config.Debug)
	log.Printf("一共%d个代理服务单元\n",len(config.Proxy))
	if !logflag{
		IsDebug=config.Debug
	}
}