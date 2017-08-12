package main


import (
	"DazeProxy/net"
	"DazeProxy/util"
	_ "DazeProxy/config"
	_ "DazeProxy/database"
	"DazeProxy/console"
	"DazeProxy/config"
	"time"
)


func main(){
	util.CheckLicense()
	util.CheckKeyAndGen()
	go net.StartServer("ipv4",config.Config.IPv4Port,config.Config.IPv6ResolvePrefer)
	go net.StartServer("ipv6",config.Config.IPv6Port,config.Config.IPv6ResolvePrefer)
	go console.Start()
	for{
		time.Sleep(time.Second*3600)
	}

}
