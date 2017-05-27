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
				v.conn.Close()
			}
		}
		if v.Network=="udp"{
			d:=time.Now().Sub(v.UDPAliveTime)
			if d.Seconds()>20{
				log.DebugPrintAlert("用户",v.conn.RemoteAddr(),"UDP心跳死了哦，断开")
				v.conn.Close()
			}
		}
	}
}
func StartHeartbeat(){
	for{
		StartOnceHBCheck()
		time.Sleep(time.Second)
	}
}