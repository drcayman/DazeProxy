package main


import (
	"DazeProxy/net"
	"DazeProxy/util"
	_ "DazeProxy/config"
	_ "DazeProxy/database"
	"DazeProxy/console"
)


func main(){
	util.CheckLicense()
	util.CheckKeyAndGen()
	go net.StartServer("ipv4",true)
	go net.StartServer("ipv6",true)
	console.Start()
}
