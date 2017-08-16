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
)
func main(){
	fmt.Println("DazeProxy V2.0-2017081601 Author:螃蟹")
	m:=go_args.ReadArgs()
	if _,consoleFlag:=m.GetArg("-console");consoleFlag{
		console.Start()
		return
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
