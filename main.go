package main


import (
	"./net"
	"./util"
	_ "./config"
)


func main(){
	util.CheckLicense()
	util.CheckKeyAndGen()
	net.StartServer()
}