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
	RemoteConn net.Conn
}
var lock sync.Mutex
func GetIsKeyExchange(conn net.Conn) bool {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].IsKeyExchange
}
func GetAESKey(conn net.Conn) []byte {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	if Users[conn.RemoteAddr()].AESKey==nil{
		return nil
	}
	return Users[conn.RemoteAddr()].AESKey
}
func GetIsAuth(conn net.Conn) bool {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].IsAuth
}
func GetPublicKeyFlag(conn net.Conn) bool {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].PublicKeyFlag
}
func GetIsConnected(conn net.Conn) bool {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].IsConnected
}
func GetRandomData(conn net.Conn) []byte {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].RandomData
}
func GetRetryTime(conn net.Conn) int {
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].RetryTime
}
func GetRemoteConn(conn net.Conn) net.Conn{
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	return Users[conn.RemoteAddr()].RemoteConn
}
func GetAvailable(conn net.Conn) bool{
	defer func(){
	lock.Unlock()
	}()
	lock.Lock()
	_,b:=Users[conn.RemoteAddr()]
	return b
}
func AddUser(conn net.Conn){
	RunCommand(MapCommand{Command:1,Conn:conn})
}
func SetConnected(conn net.Conn,IsConnected bool){
	RunCommand(MapCommand{Command:8,Conn:conn,Bool:IsConnected})
}
func SetRemoteConn(conn net.Conn,RemoteConn net.Conn){
	RunCommand(MapCommand{Command:9,Conn:conn,RemoteConn:RemoteConn})
}
func DisconnectAndDeleteUser(conn net.Conn){
	RunCommand(MapCommand{Command:0,Conn:conn})
}
func SetKeyExchange(conn net.Conn,IsKeyExchange bool){
	RunCommand(MapCommand{Command:3,Conn:conn,Bool:IsKeyExchange})
}
func SetPublicKeyFlag(conn net.Conn,PublicKeyFlag bool){
	RunCommand(MapCommand{Command:5,Conn:conn,Bool:PublicKeyFlag})
}
func SetAESKey(conn net.Conn,key []byte){
	RunCommand(MapCommand{Command:4,Conn:conn,Data:key})
}
func SetAuthed(conn net.Conn,IsAuthed bool){
	RunCommand(MapCommand{Command:6,Conn:conn,Bool:IsAuthed})
}
func RetryTimePlus(conn net.Conn){
	RunCommand(MapCommand{Command:7})
}
func RunCommand(command MapCommand){
	lock.Lock()
	command.Locker=&lock
	MapCommandChan<-command
}
type MapCommand struct{
	Command int
	Conn net.Conn
	Data []byte
	Bool bool
	RemoteConn net.Conn
	Locker *sync.Mutex
 }
var Users map[net.Addr]*User
var MapCommandChan chan MapCommand
func MapCommandThread(){
	for c:=range MapCommandChan{
		switch c.Command{
		case 0:
			c.Conn.Close()
			if Users[c.Conn.RemoteAddr()]!=nil && Users[c.Conn.RemoteAddr()].RemoteConn!=nil{
				Users[c.Conn.RemoteAddr()].RemoteConn.Close()
			}
			delete(Users,c.Conn.RemoteAddr())
		case 1:
			Users[c.Conn.RemoteAddr()]=&User{
			LastHeartBeat:time.Now(),
			conn:c.Conn,
			RandomData:util.GenRandomData(16),
			IsKeyExchange:false,
			RemoteConn:nil,
			PublicKeyFlag:false,
			RetryTime:0,
			IsConnected:false,
			AESKey:nil,
			}
		case 3:(Users[c.Conn.RemoteAddr()]).IsKeyExchange=c.Bool
		case 4:	(Users[c.Conn.RemoteAddr()]).AESKey=c.Data
		case 5:	(Users[c.Conn.RemoteAddr()]).PublicKeyFlag=c.Bool
		case 6:	(Users[c.Conn.RemoteAddr()]).IsAuth=c.Bool
		case 7:	(Users[c.Conn.RemoteAddr()]).RetryTime++
		case 8:	(Users[c.Conn.RemoteAddr()]).IsConnected=c.Bool
		case 9:	(Users[c.Conn.RemoteAddr()]).RemoteConn=c.RemoteConn
		}
		c.Locker.Unlock()
	}
}
func init(){
	Users=make(map[net.Addr]*User)
	MapCommandChan=make(chan MapCommand)
	go MapCommandThread()
}
