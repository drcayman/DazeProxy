package main

import (
	"github.com/crabkun/DazeProxy/common"
	"github.com/crabkun/DazeProxy/helper"
	"github.com/crabkun/DazeProxy/server"
	"log"
	"github.com/crabkun/DazeProxy/database"
	"sync"
	"os"
	"math/rand"
	"time"
	"github.com/crabkun/go-args"
)

var config common.S_config
var wg sync.WaitGroup

func main(){
	var logflag bool
	rand.Seed(time.Now().UnixNano())
	//处理启动参数
	args:=go_args.ReadArgs()
	if _,logflag=args.GetArg("-log");logflag{
		logfile,err:=os.OpenFile("DazeProxy.log",os.O_CREATE|os.O_WRONLY, 0666)
		if err!=nil{
			log.Println("日志文件创建失败！原因：",err.Error())
		}else{
			helper.IsDebug=false
			log.Println("进入日志模式，所有控制台输出已重定向到DazeProxy.log")
			log.SetOutput(logfile)
			log.Println("进入日志模式，已关闭调试模式")
		}
	}
	//加载配置
	helper.LoadConfig(&config,logflag)
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