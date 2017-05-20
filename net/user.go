package net

import (
	"time"
	"net"
)

type User struct {
	conn net.Conn
	LastHeartBeat time.Time
	IsKeyExchange bool
	AESKey []byte
}
type MapCommand struct{
	Command int
	Conn net.Conn
	Data []byte
 }
var Users map[net.Addr]*User
var MapCommandChan chan MapCommand
func MapCommandThread(){
	for c:=range MapCommandChan{
		switch c.Command{
		case 0:delete(Users,c.Conn.RemoteAddr())
		case 1:Users[c.Conn.RemoteAddr()]=&User{
			LastHeartBeat:time.Now(),
			conn:c.Conn,
			}
		case 2:HeartbeatCheck()
		case 3:(Users[c.Conn.RemoteAddr()]).IsKeyExchange=true
		case 4:	(Users[c.Conn.RemoteAddr()]).AESKey=c.Data
		}
	}
}
func init(){
	Users=make(map[net.Addr]*User)
	MapCommandChan=make(chan MapCommand,2048)
	go MapCommandThread()
}
