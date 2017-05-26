package net

import (
	"time"
	"DazeProxy/log"
)

func HeartbeatCheck(){
	for _,v:=range Users{
		if v.IsConnected==false{
			d:=time.Now().Sub(v.LastHeartBeat)
			if d.Seconds()>10{
				log.DebugPrintAlert("用户",v.conn.RemoteAddr(),"验证或者连接超时，断开")
				DisconnectAndDeleteUser(v.conn)
			}
		}
	}
}
func StartHeartbeat(){
	for{
		//HeartbeatCheck()
		time.Sleep(time.Second)
	}
}