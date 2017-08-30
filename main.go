package main

import (
	"github.com/crabkun/DazeProxy/common"
	"github.com/crabkun/DazeProxy/helper"
	"github.com/crabkun/DazeProxy/server"
	"time"
	"log"
)

var config common.S_config

func main(){
	log.Println("DazeProxy V3-201708301")
	//加载配置
	helper.LoadConfig(&config)
	//启动代理实例
	for _,s:=range config.Proxy{
		go server.StartServer(s)
	}
	for{
		time.Sleep(time.Second*10)
	}
}