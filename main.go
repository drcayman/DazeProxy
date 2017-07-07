package main


import (
	"DazeProxy/net"
	"DazeProxy/util"
	_ "DazeProxy/config"
	_ "DazeProxy/database"
	"DazeProxy/console"
	"DazeProxy/config"
)


func main(){
	util.CheckLicense()
	util.CheckKeyAndGen()
	go net.StartServer("ipv4",config.Config.IPv4Port,config.Config.IPv6ResolvePrefer)
	go net.StartServer("ipv6",config.Config.IPv6Port,config.Config.IPv6ResolvePrefer)
	console.Start()
}
