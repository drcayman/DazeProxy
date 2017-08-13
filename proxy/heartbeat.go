package proxy

import (
	"time"
	"log"
	"DazeProxy/config"
	"DazeProxy/common"
)

func NewHeartBeatCountDown(c chan int,second int,client *common.User,msg string){
	select {
	case <- time.After(time.Duration(second)* time.Second):
		if config.Config.IsDebug{
			log.Println(client.Conn.RemoteAddr(),"time out",msg)
		}
		client.Conn.Close()
		return
	case <- c:
		return
	}
}
