package main


import (
	"DazeProxy/net"
	"DazeProxy/util"
	_ "DazeProxy/config"
)


func main(){
	util.CheckLicense()
	util.CheckKeyAndGen()
	go net.StartServerIP6(true)
	net.StartServerIP4(true)

}
