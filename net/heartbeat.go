package net

import (
	"time"
	"DazeProxy/log"
)

func NewHeartBeatCountDown(c chan int,second int,client *User,msg string){
	select {
	case <- time.After(time.Duration(second)* time.Second):
		log.DebugPrintNormal(client.Conn.RemoteAddr(),"time out",msg)
		client.Conn.Close()
		return
	case <- c:
		return
	}
}
