package main


import (
	_ "DazeProxy/config"
	_ "DazeProxy/database"
	"time"
	"DazeProxy/config"
	"DazeProxy/proxy"
	//"DazeProxy/console"
	"DazeProxy/common"
)


func main(){

	//go console.Start()
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
