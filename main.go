package main

import (
	"github.com/crabkun/DazeProxy/common"
	"github.com/crabkun/DazeProxy/helper"
	"github.com/crabkun/DazeProxy/server"
	_ "github.com/crabkun/DazeProxy/database"
	"log"
	"github.com/crabkun/DazeProxy/database"
	"sync"
	"os"
	"math/rand"
	"time"
)

var config common.S_config
var wg sync.WaitGroup

func main(){
	rand.Seed(time.Now().UnixNano())
	//加载配置
	helper.LoadConfig(&config)
	log.Println("DazeProxy V3-201709031")
	//如果处于验证模式则连接数据库
	if !config.NoAuth{
		database.LoadDatabase(config.DatabaseDriver,config.DatabaseConnectionString)
	}else{
		log.Println("警告：服务器处于免验证模式！")
	}
	//启动代理实例
	for _,s:=range config.Proxy{
		s.Config=config
		wg.Add(1)
		go server.StartServer(s,wg)
	}
	wg.Wait()
	os.Exit(-1)
}