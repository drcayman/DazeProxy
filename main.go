package main


import (
	"DazeProxy/net"
	"DazeProxy/util"
	_ "DazeProxy/config"
)


func main(){
	util.CheckLicense()
	util.CheckKeyAndGen()
	net.StartServer()

}
