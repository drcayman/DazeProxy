package proxy

import (
	"net"
	"sync"
)

type User struct {
	Conn net.Conn
	AuthHeartBeat chan int
	IsKeyExchange bool
	AESKey []byte
	IsAuth bool
	PublicKeyFlag bool
	IsConnected bool
	RandomData []byte
	RetryTime int
	RemoteConn net.Conn
	UDPAliveTime chan int
	Network string
	Locker sync.Mutex
	ChanCloseFlag bool
	IPv6ResolvePrefer bool
}
func NewUser(conn net.Conn) *User{
	return &User{
		Conn:conn,
		AuthHeartBeat:make(chan int),
		UDPAliveTime:make(chan int),
	}
}
