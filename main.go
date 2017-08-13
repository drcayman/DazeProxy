package main


import (
	_ "DazeProxy/config"
	_ "DazeProxy/database"
	"time"
	"DazeProxy/config"
	"DazeProxy/proxy"
)


func main(){
	//util.CheckLicense()
	//util.CheckKeyAndGen()
	//go net.StartServer("ipv4",config.Config.IPv4Port,config.Config.IPv6ResolvePrefer)
	//go net.StartServer("ipv6",config.Config.IPv6Port,config.Config.IPv6ResolvePrefer)
	//go console.Start()
	for _,v:=range config.Config.ProxyUnit{
		unit:= proxy.ProxyUnit{
			Config:v,
		}
		go proxy.StartServer(unit)
	}
	for{
		time.Sleep(time.Second*3600)
	}

}
