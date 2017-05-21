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
	IsAuth bool
	IsConnected bool
}
func AddUser(conn net.Conn){
	MapCommandChan<-MapCommand{Command:1,Conn:conn}
}
func DisconnectAndDeleteUser(conn net.Conn){
	conn.Close()
	MapCommandChan<-MapCommand{Command:0,Conn:conn}
}
func SetKeyExchange(conn net.Conn,IsKeyExchange bool){
	MapCommandChan<-MapCommand{Command:3,Conn:conn,Bool:IsKeyExchange}
}
func SetAESKey(conn net.Conn,key []byte){
	MapCommandChan<-MapCommand{Command:4,Conn:conn,Data:key}
}
func StartOnceHeartbeat(){
	MapCommandChan<-MapCommand{Command:2,Conn:nil}
}
type MapCommand struct{
	Command int
	Conn net.Conn
	Data []byte
	Bool bool
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
		case 3:(Users[c.Conn.RemoteAddr()]).IsKeyExchange=c.Bool
		case 4:	(Users[c.Conn.RemoteAddr()]).AESKey=c.Data
		case 5:
		}
	}
}
func init(){
	Users=make(map[net.Addr]*User)
	MapCommandChan=make(chan MapCommand,2048)
	go MapCommandThread()
}
