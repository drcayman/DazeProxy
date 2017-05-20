package net

import (
	"time"
	"net"
)

type User struct {
	conn net.Conn
	LastHeartBeat time.Time
	IsKeyExchange bool
}
type MapCommand struct{
	Command int
	Conn net.Conn
 }
var Users map[net.Addr]User
var MapCommandChan chan MapCommand
func MapCommandThread(){
	for c:=range MapCommandChan{
		switch c.Command{
		case 0:delete(Users,c.Conn.RemoteAddr());break
		case 1:Users[c.Conn.RemoteAddr()]=User{
			LastHeartBeat:time.Now(),
			conn:c.Conn,
			};break
		case 2:HeartbeatCheck();break
		}
	}
}
func init(){
	Users=make(map[net.Addr]User)
	MapCommandChan=make(chan MapCommand,2048)
	go MapCommandThread()
}
