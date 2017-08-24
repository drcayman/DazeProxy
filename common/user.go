package common

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
	Network string
	Locker sync.Mutex
	ChanCloseFlag bool
	IPv6ResolvePrefer bool
	ProxyUnit *ProxyUnit
	EncReserved interface{}
	DsgReserved interface{}
}
func NewUser(conn net.Conn) *User{
	return &User{
		Conn:conn,
		AuthHeartBeat:make(chan int),
	}
}
