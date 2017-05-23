package net

import (
	"time"
	"DazeProxy/config"
	"DazeProxy/log"
)

func HeartbeatCheck(){
	for _,v:=range Users{
		if v.IsConnected==false{
			d:=time.Now().Sub(v.LastHeartBeat)
			if d.Seconds()>10{
				if config.Config.IsDebug{
					log.PrintAlert("用户",v.conn.RemoteAddr(),"验证或者连接超时，断开")
				}
				DisconnectAndDeleteUser(v.conn)
			}
		}
	}
}
func StartHeartbeat(){
	for{
		StartOnceHeartbeat()
		time.Sleep(time.Second)
	}
}