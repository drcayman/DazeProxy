package net

import (
	"net"
	"DazeProxy/config"
	"os"
	"encoding/binary"
	"bytes"
	"errors"
	"io/ioutil"
	"DazeProxy/util"
	"DazeProxy/log"
	"reflect"
	"unsafe"
	"strings"
	"time"
)
var Server *net.Listener

type Packet struct {
	command byte
	data []byte
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
func SendPacket(client net.Conn,data []byte){
	_,err:=client.Write(data)
	if err!=nil{
		DisconnectAndDeleteUser(client)
	}
}
func SendPacketAndDisconnect(client net.Conn,data []byte){
	client.Write(data)
	DisconnectAndDeleteUser(client)
}
func ServeClient(client net.Conn,c chan Packet){
	defer func(){
		close(c)
		DisconnectAndDeleteUser(client)
		log.DebugPrintAlert(client.RemoteAddr(),"连接线程关闭")
	}()
	buf:=make([]byte,65536)
	for{
		n,err:=client.Read(buf)
		if err!=nil{
			return
		}
		if Users[client.RemoteAddr()].IsKeyExchange{
		   newbuf,err:=util.DecryptAES(buf,Users[client.RemoteAddr()].AESKey)
			if err!=nil{
				continue
			}
			buf=nil
			buf=newbuf
		}
		command,data,derr:=DePacket(buf[0:n])
		if derr!=nil{
			continue
		}
		c<-Packet{command:command,data:data}
	}
}
/*
客户端发送到服务端
命令1代表需要公钥
命令2代表设定key
命令3代表用户名密码登录，data里面是json格式的数据，例如{"username":"123","password":"456"}
命令4代表客户端想要用证书登录，data里面是公钥
命令5代表客户端发送过来的随机字符串A，要求服务端对比是否一致
命令A1代表代理IPV4的TCP连接，data里面是sessionID和IP和端口
命令A2代表代理IPV4的UDP连接，data里面是sessionID和IP和端口
命令A

服务端发送到客户端
命令1代表data里面是key
命令2代表断开性错误（没交换key就执行其他操作）
命令3代表key长度错误
命令4代表接受AESkey
命令5代表利用用户RSA指纹寻找到公钥，并加密了一串随机字符串A
命令6代表用户的公钥在数据库中找不到
命令7代表客户端发送过来的随机字符串A跟服务端一致（证书登录成功）
命令8代表客户端发送过来的随机字符串A跟服务端不一致（证书登录失败）
命令9代表登录次数过多
命令A代表用户发送过来的公钥没法加密
命令E1代表IP地址格式错误
命令E2代表连接失败

命令FF代表指令没法识别
*/
func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}
func DecodeConnectPacketAndCheck(data []byte)(string,error){
	if len(data)>21 || len(data)<9{
		return "",errors.New("error1")
	}
	ip:=String(data)
	p1:=strings.IndexByte(ip,':')
	if p1<7 || p1>15{
		return "",errors.New("error2")
	}
	ips:=net.ParseIP(ip[:p1])
	if ips.String()=="test" {
		return "",errors.New("error3")
	}
	return ip,nil
}
func ServeCommand(client net.Conn,c chan Packet) {
	defer func(){
		client.Close()
		log.DebugPrintAlert(client.RemoteAddr(),"处理线程关闭")
	}()
	for packet:=range c{
		if Users[client.RemoteAddr()].IsKeyExchange==false {
			if packet.command == 0x1 {
				PublicKeyBuf, PublicKeyErr := ioutil.ReadFile("public.pem")
				if PublicKeyErr != nil {
					log.PrintPanic("公钥文件丢失！！系统强制退出")
				}
				SendPacket(client, MakePacket(1, PublicKeyBuf))
				continue
			} else if packet.command == 0x2 {
				Debuf, DeErr := util.DecryptRSA(packet.data)
				if DeErr != nil || len(Debuf) != 32 {
					SendPacketAndDisconnect(client, MakePacket(3, nil))
					return
				}
				SetKeyExchange(client, true)
				SetAuthed(client, true) //debug!!!!!!!!!!!!!!!!!
				SetAESKey(client, Debuf)
				SendPacket(client, MakePacket(4, nil))
				continue
			} else {
				SendPacketAndDisconnect(client, MakePacket(2, nil))
				return
			}
		}
		if Users[client.RemoteAddr()].IsAuth==false{
			RetryTimePlus(client)
			if Users[client.RemoteAddr()].RetryTime>4{
				SendPacketAndDisconnect(client, MakePacket(9, nil))
				return
			}
			switch packet.command {
			case 4:
				key := packet.data
				EncryptRandom, EncryptRandomErr := util.EncryptRSAWithKey(Users[client.RemoteAddr()].RandomData, key)
				if EncryptRandomErr != nil {
					SendPacket(client, MakePacket(0xA, nil))
					continue
				}
				SendPacket(client, MakePacket(5, EncryptRandom))
				SetPublicKeyFlag(client, true)
			case 5:
				if Users[client.RemoteAddr()].PublicKeyFlag == false {
					SendPacket(client, MakePacket(0xff, nil))
					continue
				}
				if bytes.Compare(packet.data, Users[client.RemoteAddr()].RandomData) == 0 {
					SetAuthed(client, true)
					SendPacket(client, MakePacket(7, nil))
					continue
				} else {
					SetAuthed(client, false)
					SendPacket(client, MakePacket(8, nil))
					continue
				}
			default:
				SendPacketAndDisconnect(client, MakePacket(0xff, nil))
				return

			}
		}
		if Users[client.RemoteAddr()].IsConnected==false{
			if packet.command==0xA1{
				address,DecodeConnectPacketErr:=DecodeConnectPacketAndCheck(packet.data)
				if DecodeConnectPacketErr!=nil{
					SendPacketAndDisconnect(client, MakePacket(0xE1, nil))
					return
				}
				proxyconn,dailerr:=net.DialTimeout("tcp",address, time.Second * 3)
				log.DebugPrintNormal("客户端",client.RemoteAddr(),"想要代理",address)
				if dailerr!=nil{
					log.DebugPrintNormal("客户端",client.RemoteAddr(),"想要代理",address,"但连接失败了")
					SendPacketAndDisconnect(client, MakePacket(0xE2, nil))
				}
				SetConnected(client,true)
				SetRemoteConn(client,proxyconn)
			}
		}
	}
}
func StartServer(){
	l,err:=net.Listen("tcp",":"+config.Config.ServerPort)
	if err!=nil{
		log.PrintPanic("服务端启动失败（原因：",err.Error(),")")
		os.Exit(-1)
	}
	Server=&l
	go StartHeartbeat()
	log.PrintSuccess("服务端启动成功")
	for {
		client, _ := l.Accept()
		//delete(Users,client.RemoteAddr()) //BUG!!!!!
		AddUser(client)
		c:=make(chan Packet,10)
		go ServeClient(client,c)
		go ServeCommand(client,c)
	}
}
