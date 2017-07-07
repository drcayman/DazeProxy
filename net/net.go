package net

import (
	"net"
	"DazeProxy/config"
	"DazeProxy/database"
	"encoding/binary"
	"bytes"
	"errors"
	"DazeProxy/util"
	"log"
	"reflect"
	"unsafe"
	_NET "net"
	"math/rand"
	"time"
	"encoding/json"
)
var Server *net.Listener
type JsonAuth struct{
	Username string
	Password string
}

//F1 01 00 04 31 32 33 34 F2
//F1为头部
//01为命令
//00 04为content长度
//31 32 33 34 为content内容
//F2为尾部
func MakePacket(command byte,content []byte) []byte{
	if content==nil{
		content=[]byte{0x0}
	}
	ContentLen:=uint16(len(content))
	if ContentLen>0xFFFF {
		return nil
	}
	ContentLenBuffer:=bytes.NewBuffer([]byte{})
	buf:=make([]byte,5+len(content))
	buf[0]=0xF1
	buf[1]=command
	binary.Write(ContentLenBuffer,binary.BigEndian,ContentLen)
	copy(buf[2:],ContentLenBuffer.Bytes())
	copy(buf[4:],content)
	buf[len(buf)-1]=0xF2
	return buf
}
func DePacket(buf []byte) (byte,[]byte,error){
	if len(buf)<6 || buf[0]!=0xF1 || buf[len(buf)-1]!=0xF2{
		return 0,nil,errors.New("error1")
	}
	ContentLen:=int(buf[2])*256+int(buf[3])
	if len(buf)-5!=int(ContentLen){
		return 0,nil,errors.New("error2")
	}
	return buf[1],buf[4:4+ContentLen],nil
}
func SendPacket(client *User,data []byte){
	var bufffer *bytes.Buffer
	var AESKey []byte
	if client.IsKeyExchange {
		AESKey = client.AESKey
	}else{
		AESKey= util.GetAESKeyByDay()
	}
	//data,_=util.EncryptAES(data,AESKey)
	dataLen:=len(data)
	header:=[]byte{0xFB,byte(dataLen%0x100),byte(dataLen/0x100),0xFC}
	bufffer=bytes.NewBuffer(header)
	bufffer.Write(data)
	buffferBytes,_:=util.EncryptAES(bufffer.Bytes(),AESKey)
	buffferBytesLen:=len(buffferBytes)
	rand.Seed(time.Now().Unix())
	len1:=rand.Intn(buffferBytesLen)
	if len1<4{
		len1=4
	}
	client.Conn.Write(buffferBytes[:len1])
	client.Conn.Write(buffferBytes[len1:])
}
func SendPacketAndDisconnect(client *User,data []byte){
	SendPacket(client,data)
	client.Conn.Close()
	//DisconnectAndDeleteUser(client)
}
func SaveAESKey(buf []byte,client *User) error{
	Debuf, DeErr := util.DecryptRSA(buf)
	if config.Config.IsDebug{
		log.Println("key解码前长度：", len(buf), "key解码后长度：", len(Debuf))
	}
	if DeErr != nil ||Debuf[len(Debuf)-1]!=0xFF|| len(Debuf) != 33 {
		return errors.New("AESKey Error")
	}
	client.AESKey= Debuf[:len(Debuf)-1]
	client.IsKeyExchange=true
	SendPacket(client, MakePacket(4, nil))
	if config.Config.IsDebug {
		log.Println(client.Conn.RemoteAddr(), " key交换成功")
	}
	return nil
}
func ReadFromClient(client *User) ([]byte,error){
	headerbuf:=make([]byte,4)
	n,err:=client.Conn.Read(headerbuf)
	if n<4 ||err!=nil{
		return nil,errors.New("read header error ")
	}
	AESKey:=client.AESKey
	if AESKey==nil{
		AESKey=util.GetAESKeyByDay()
	}
	header:=headerbuf[:4]
	headerDecode,_:=util.DecryptAES(header,AESKey)
	if headerDecode[0]!=0xFB || headerDecode[3]!=0xFC{
		return nil,errors.New("deheader error")
	}
	buflen:=int(headerDecode[1])+int(headerDecode[2])*256
	buf:=make([]byte,buflen)
	pos:=0
	for{
		n,err:=client.Conn.Read(buf[pos:])
		if err!=nil{
			return nil,errors.New("read body error")
		}
		buflen-=n
		pos+=n
		if buflen<0{
			return nil,errors.New("body len error")
		}
		if buflen==0{
			break
		}
	}
	xbuf:=make([]byte,len(buf)+4)
	copy(xbuf,header)
	copy(xbuf[4:],buf)
	xbuf,_=util.DecryptAES(xbuf,AESKey)
	return xbuf[4:],nil
}
func ServeClient(client *User){
	flag:=0
	defer func(){
		if flag==0{
			client.Conn.Close()
			close(client.AuthHeartBeat)
			close(client.UDPAliveTime)
			if config.Config.IsDebug {
				log.Println(client.Conn.RemoteAddr(), "连接线程关闭")
			}
			return
		}
		if config.Config.IsDebug {
			log.Println(client.Conn.RemoteAddr(), "连接线程已进入代理模式")
		}
	}()
	for{
		buf,err:=ReadFromClient(client)
		if err!=nil{
			return
		}
		if !client.IsKeyExchange{
			if config.Config.IsDebug {
				log.Println(client.Conn.RemoteAddr(), " key交换开始")
			}
			SaveAESErr:=SaveAESKey(buf,client)
			if SaveAESErr!=nil{
				return
			}
			continue
		}
		command,data,derr:=DePacket(buf)
		if derr!=nil{
			continue
		}
		flag=ServeCommand(client,command,data)
		if flag==1{
			return
		}
	}
}
/*
客户端发送到服务端
命令1代表需要公钥
命令2代表设定key
命令3代表用户名密码登录，data里面是json格式的数据，例如{"username":"123","password":"456"}
命令4代表客户端想要用证书登录，data里面是公钥
命令5代表客户端发送过来的随机字符串A，要求服务端对比是否一致
命令A1代表代理IPV4的TCP连接，data里面是IP和端口
命令A2代表代理IPV4的UDP连接，data里面是IP和端口

服务端发送到客户端
命令1代表data里面是key
命令2代表断开性错误（没交换key就执行其他操作）
命令3代表key长度错误
命令4代表接受AESkey
命令5代表利用用户RSA指纹寻找到公钥，并加密了一串随机字符串A
命令7代表客户端发送过来的随机字符串A跟服务端一致（证书登录成功）
命令9代表客户端登录成功

命令A代表用户发送过来的公钥没法加密
命令E1代表IP地址格式错误
命令E2代表连接失败
命令E3代表登录失败
命令E4代表未登录
命令E5代表客户端发送过来的随机字符串A跟服务端不一致（证书登录失败）
命令E6代表用户的公钥在数据库中找不到

命令C1代表成功连接

命令FF代表指令没法识别
*/
func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}
func DecodeConnectPacketAndCheck(data []byte,client *User)(string,error){
	var ips *_NET.IPAddr
	var ResolveErr error=nil
	ip:=util.B2s(data)
	host,port,SplitErr:=_NET.SplitHostPort(ip)
	if SplitErr!=nil{
		return "",errors.New("error1")
	}
	if client.IPv6ResolvePrefer{
		ips,ResolveErr=_NET.ResolveIPAddr("ip6",host)
	}
	if ips==nil{
		ips,ResolveErr=net.ResolveIPAddr("ip",host)
	}
	if ResolveErr!=nil{
		return "",errors.New("error2")
	}
	if ips.IP.IsLoopback(){

	}
	ipstring:=ips.String()
	if len(ipstring)>15{
		ipstring="["+ipstring+"]"
	}
	return ipstring+":"+port,nil
}
func ServeCommand(client *User,command byte,data []byte) int {
		switch command{
		case 0x03:{
			authinfo:=JsonAuth{}
			err:=json.Unmarshal(data,&authinfo)
			if err!=nil||!database.CheckUserPass(authinfo.Username,authinfo.Password){
				SendPacketAndDisconnect(client, MakePacket(0xE3, nil))
				return 0
			}
			client.IsAuth=true
			client.AuthHeartBeat<-1
			SendPacket(client, MakePacket(0x09, nil))
		}
		case 0xA1:fallthrough
		case 0xA2:
			{
				if !client.IsAuth{
					SendPacketAndDisconnect(client, MakePacket(0xE4, nil))
					return 0
				}
				network:="tcp"
				if command==0xA2{
					network="udp"
				}
				address,DecodeConnectPacketErr:=DecodeConnectPacketAndCheck(data,client)
				if DecodeConnectPacketErr!=nil{
					SendPacketAndDisconnect(client, MakePacket(0xE1, nil))
					return 0
				}
				ProxyConn,dailerr:=net.Dial(network,address)
				if config.Config.IsDebug {
					log.Println("客户端", client.Conn.RemoteAddr(), "想要代理", network, address)
				}
				if dailerr!=nil{
					if config.Config.IsDebug {
						log.Println("客户端", client.Conn.RemoteAddr(), "想要代理", network, address, "但连接失败了")
					}
					SendPacketAndDisconnect(client, MakePacket(0xE2, nil))
					return 0
				}
				if config.Config.IsDebug {
					log.Println("客户端", client.Conn.RemoteAddr(), "想要代理", network, address, "，连接成功")
				}
				client.Network=network
				client.IsConnected=true
				SendPacket(client, MakePacket(0xC1, []byte(ProxyConn.RemoteAddr().String())))
				go BridgeClientToRemote(client,ProxyConn)
				go BridgeRemoteToClient(client,ProxyConn)
				return 1
			}
		}
	return 0
}
func CloseChan(client *User) {
	client.Locker.Lock()
	if !client.ChanCloseFlag{
	close(client.AuthHeartBeat)
	close(client.UDPAliveTime)
	}
	client.ChanCloseFlag=true
	client.Locker.Unlock()
}
func BridgeClientToRemote(client *User,Remote net.Conn){
	defer func(){
		client.Conn.Close()
		Remote.Close()
		CloseChan(client)
		if config.Config.IsDebug {
			log.Println(client.Conn.RemoteAddr(), "BCTR退出")
		}
	}()
	for{
		buf,err:=ReadFromClient(client)
		if err!=nil{
			return
		}
		Remote.Write(buf)
	}
}
func BridgeRemoteToClient(client *User,Remote net.Conn){
	defer func(){
		Remote.Close()
		client.Conn.Close()
		CloseChan(client)
		if config.Config.IsDebug {
			log.Println(client.Conn.RemoteAddr(), "BRTC退出")
		}
	}()
	buf:=make([]byte,1024)
	for{
		n,err:=Remote.Read(buf)
		if err!=nil{
			return
		}
		SendPacket(client,buf[:n])
	}
}

func StartServer(targetNet string,port string,ipv6ResolvePrefer bool){
	if port==""{
		return
	}
	targetNet1:="tcp4"
	if targetNet=="ipv6"{
		targetNet1="tcp6"
	}
	l,err:=net.Listen(targetNet1,":"+port)
	if err!=nil{
		log.Fatal(targetNet+"服务端启动失败（原因：",err.Error(),")")
		return
	}
	//Server=&l
	//go StartHeartbeat()
	log.Println(targetNet+"服务端启动成功")
	for {
		conn, AcceptErr := l.Accept()
		//delete(Users,client.RemoteAddr()) //BUG!!!!!
		if AcceptErr!=nil{
			log.Println("客户端接受失败！",AcceptErr.Error())
			continue
		}
		if config.Config.IsDebug {
			log.Println("客户端",conn.RemoteAddr(),"连接")
		}

		//AddUser(client)
		client:=NewUser(conn)
		client.IPv6ResolvePrefer=ipv6ResolvePrefer
		go NewHeartBeatCountDown(client.AuthHeartBeat,5,client,"Auth or Connect")
		SendPacket(client,util.GetPublicKey())
		go ServeClient(client)
	}
}