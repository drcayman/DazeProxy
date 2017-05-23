package net

import (
	"time"
	"net"
	"DazeProxy/util"
	"sync"
)

type User struct {
	conn net.Conn
	LastHeartBeat time.Time
	IsKeyExchange bool
	AESKey []byte
	IsAuth bool
	PublicKeyFlag bool
	IsConnected bool
	RandomData []byte
	RetryTime int
}
func AddUser(conn net.Conn){
	var a sync.Mutex
	a.Lock()
	MapCommandChan<-MapCommand{Command:1,Conn:conn,Locker:&a}
	a.Lock()
	a.Unlock()
}
func DisconnectAndDeleteUser(conn net.Conn){
	conn.Close()
	MapCommandChan<-MapCommand{Command:0,Conn:conn}
}
func SetKeyExchange(conn net.Conn,IsKeyExchange bool){
	MapCommandChan<-MapCommand{Command:3,Conn:conn,Bool:IsKeyExchange}
}
func SetPublicKeyFlag(conn net.Conn,PublicKeyFlag bool){
	MapCommandChan<-MapCommand{Command:5,Conn:conn,Bool:PublicKeyFlag}
}
func SetAESKey(conn net.Conn,key []byte){
	MapCommandChan<-MapCommand{Command:4,Conn:conn,Data:key}
}
func SetAuthed(conn net.Conn,IsAuthed bool){
	MapCommandChan<-MapCommand{Command:6,Conn:conn,Bool:IsAuthed}
}
func RetryTimePlus(conn net.Conn){
	MapCommandChan<-MapCommand{Command:7}
}
func StartOnceHeartbeat(){
	MapCommandChan<-MapCommand{Command:2,Conn:nil}
}
type MapCommand struct{
	Command int
	Conn net.Conn
	Data []byte
	Bool bool
	Locker *sync.Mutex
 }
var Users map[net.Addr]*User
var MapCommandChan chan MapCommand
func MapCommandThread(){
	for c:=range MapCommandChan{
		switch c.Command{
		case 0:
			delete(Users,c.Conn.RemoteAddr())
		case 1:Users[c.Conn.RemoteAddr()]=&User{
			LastHeartBeat:time.Now(),
			conn:c.Conn,
			RandomData:util.GenRandomData(16),
			IsKeyExchange:false,
			}
			(*c.Locker).Unlock()
		case 2:HeartbeatCheck()
		case 3:(Users[c.Conn.RemoteAddr()]).IsKeyExchange=c.Bool
		case 4:	(Users[c.Conn.RemoteAddr()]).AESKey=c.Data
		case 5:	(Users[c.Conn.RemoteAddr()]).PublicKeyFlag=c.Bool
		case 6:	(Users[c.Conn.RemoteAddr()]).IsAuth=c.Bool
		case 7:	(Users[c.Conn.RemoteAddr()]).RetryTime++
		}
	}
}
func init(){
	Users=make(map[net.Addr]*User)
	MapCommandChan=make(chan MapCommand,2048)
	go MapCommandThread()
}
