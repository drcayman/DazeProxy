package main


import (
	"./net"
	"./util"
	_ "./config"
)


func main(){
	util.GenAESKey(1024)
	util.CheckLicense()
	util.CheckKeyAndGen()
	net.StartServer()
}