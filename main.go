package main

import (
	_ "DazeProxy/config"
	_ "DazeProxy/database"
	"time"
	"DazeProxy/config"
	"DazeProxy/proxy"
	"DazeProxy/console"
	"DazeProxy/common"
	"github.com/crabkun/go-args"
	"fmt"
	"DazeProxy/database"
	"os"
)
func main(){
	fmt.Println("DazeProxy V2.0-2017082401 Author:螃蟹")
	m:=go_args.ReadArgs()
	if _,consoleFlag:=m.GetArg("-console");consoleFlag {
		console.Start()
		return
	}else if database.GetUserCount()==0{
		fmt.Println("**********注意！**********")
		fmt.Println(" ")
		fmt.Println("检测到用户数为0，请在运行本程序时候带上-console参数来添加用户后再使用！")
		fmt.Println("例如："+os.Args[0]+" -console")
		fmt.Println(" ")
		fmt.Println("**********注意！**********")
	}
	for _,v:=range config.Config.ProxyUnit{
		unit:= common.ProxyUnit{
			Config:v,
		}
		go proxy.StartServer(unit)
	}
	for{
		time.Sleep(time.Second*3600)
	}

}
