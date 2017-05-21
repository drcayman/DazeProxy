package net

import (
	"time"
	"../config"
	"../log"
)

func HeartbeatCheck(){
	for _,v:=range Users{
		d:=time.Now().Sub(v.LastHeartBeat)
		if d.Seconds()>10{
			if config.Config.IsDebug{
				log.PrintAlert("用户",v.conn.RemoteAddr(),"心跳超时，断开")
			}
			v.conn.Close()
		}
	}
}
func StartHeartbeat(){
	for{
		StartOnceHeartbeat()
		time.Sleep(time.Second)
	}
}