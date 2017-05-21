package main


import (
	"./net"
	"./util"
	_ "./config"
)


func main(){
	//log.PrintAlert("aaa","1234567")
	//fmt.Println(runtime.GOOS)
	util.CheckLicense()
	util.CheckKeyAndGen()
	net.StartServer()

}