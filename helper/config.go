package helper

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"fmt"
	"github.com/crabkun/DazeProxy/common"
)
var IsDebug bool
func LoadConfig(config *common.S_config){
	buf,ReadErr:=ioutil.ReadFile("config.json")
	if ReadErr!=nil{
		fmt.Println("配置文件(config.json)读取错误："+ReadErr.Error())
		os.Exit(-1)
	}
	JsonErr:=json.Unmarshal(buf,config)
	if JsonErr!=nil{
		fmt.Println("配置文件格式错误！请严格按照JSON格式来修改config.json(",JsonErr.Error(),")")
		os.Exit(-1)
	}
	fmt.Println("配置文件读取成功：")
	fmt.Println("    调试模式：",config.Debug)
	fmt.Printf("一共%d个代理服务单元\n",len(config.Proxy))
	os.Stdout.Sync()
	IsDebug=config.Debug
}