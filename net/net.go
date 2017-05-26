package net

import (
	"net"
	"DazeProxy/config"
	"os"
	"encoding/binary"
	"bytes"
	"errors"
	"DazeProxy/util"
	"DazeProxy/log"
	"reflect"
	"unsafe"
	_NET "net"
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
	var bufffer *bytes.Buffer
	var buf []byte
	var AESKey []byte
	if GetAvailable(client)==false{
		return
	}
	if GetIsKeyExchange(client) {
		AESKey = GetAESKey(client)
	}else{
		AESKey= util.GetAESKeyByDay()
	}
	data,_=util.EncryptAES(data,AESKey)
	datelen:=len(data)
	buf,_=util.EncryptAES([]byte{0xFB,byte(datelen%0x100),byte(datelen/0x100),0xFC},AESKey)
	bufffer=bytes.NewBuffer(buf)
	bufffer.Write(data)
	//fmt.Println("发送了",len(bufffer.Bytes()))
	client.Write(bufffer.Bytes())
}
func SendRawPacket(client net.Conn,data []byte){
	client.Write(data)
}
func SendPacketAndDisconnect(client net.Conn,data []byte){
	SendPacket(client,data)
	client.Close()
	//DisconnectAndDeleteUser(client)
}
func SaveAESKey(buf []byte,client net.Conn) error{
	Debuf, DeErr := util.DecryptRSA(buf)
	log.DebugPrintSuccess("key解码前长度：",len(buf),"key解码后长度：",len(Debuf))
	if DeErr != nil ||Debuf[len(Debuf)-1]!=0xFF|| len(Debuf) != 33 {
		client.Close()
		return errors.New("AESKey Error")
	}
	SetAuthed(client, true) //debug!!!!!!!!!!!!!!!!!
	SetAESKey(client, Debuf[:len(Debuf)-1])
	SetKeyExchange(client, true)
	SendPacket(client, MakePacket(4, nil))
	log.DebugPrintSuccess(client.RemoteAddr()," key交换成功")
	return nil
}
func ReadFromClient(where net.Conn) ([]byte,error){
	headerbuf:=make([]byte,4)
	n,err:=where.Read(headerbuf)
	if n<4 ||err!=nil{
		return nil,errors.New("read header error ")
	}
	AESKey:=GetAESKey(where)
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
		n,err:=where.Read(buf[pos:])
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
	decodeBuf,_:=util.DecryptAES(buf,AESKey)
	return decodeBuf,nil
}
func ServeClient(client net.Conn,c chan Packet){
	defer func(){
		close(c)
		if GetIsConnected(client){
			GetRemoteConn(client).Close()
		}
		DisconnectAndDeleteUser(client)
		log.DebugPrintAlert(client.RemoteAddr(),"连接线程关闭")
	}()
	for{
		buf,err:=ReadFromClient(client)
		if err!=nil{
			return
		}
		if !GetIsKeyExchange(client){
			log.DebugPrintSuccess(client.RemoteAddr()," key交换开始")
			SaveAESKey(buf,client)
			continue
		}
		if GetIsConnected(client){
			c<-Packet{command:0,data:buf}
			continue
		}
		command,data,derr:=DePacket(buf)
		if derr!=nil{
			continue
			}
		c<-Packet{command:command,data:data}
		//fmt.Println(command)
		//newbuf,err:=util.DecryptAES(buf[:n],GetAESKey(client))
		//if err!=nil{
		//	continue
		//}
		//if GetIsConnected(client){
		//	GetRemoteConn(client).Write(newbuf)
		//}else{
		//	command,data,derr:=DePacket(newbuf)
		//	if derr!=nil{
		//		continue
		//	}
		//	c<-Packet{command:command,data:data}
		//}

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
func DecodeConnectPacketAndCheck(data []byte)(string,error){
	ip:=util.B2s(data)
	host,port,SplitErr:=_NET.SplitHostPort(ip)
	if SplitErr!=nil{
		return "",errors.New("error1")
	}
	ips,ResolveErr:=net.ResolveIPAddr("ip4",host)
	if ResolveErr!=nil{
		return "",errors.New("error2")
	}
	if ips.IP.IsLoopback(){

	}
	return ips.String()+":"+port,nil
}
func ServeCommand(client net.Conn,c chan Packet) {
	defer func(){
		client.Close()
		log.DebugPrintAlert(client.RemoteAddr(),"处理线程关闭")
	}()
	for packet:=range c{
		//if GetIsAuth(client)==false{
		//	RetryTimePlus(client)
		//	if GetRetryTime(client)>4{
		//		SendPacketAndDisconnect(client, MakePacket(9, nil))
		//		return
		//	}
		//	switch packet.command {
		//	case 4:
		//		key := packet.data
		//		EncryptRandom, EncryptRandomErr := util.EncryptRSAWithKey(GetRandomData(client), key)
		//		if EncryptRandomErr != nil {
		//			SendPacket(client, MakePacket(0xA, nil))
		//			continue
		//		}
		//		SendPacket(client, MakePacket(5, EncryptRandom))
		//		SetPublicKeyFlag(client, true)
		//	case 5:
		//		if GetPublicKeyFlag(client) == false {
		//			SendPacket(client, MakePacket(0xff, nil))
		//			continue
		//		}
		//		if bytes.Compare(packet.data, GetRandomData(client)) == 0 {
		//			SetAuthed(client, true)
		//			SendPacket(client, MakePacket(7, nil))
		//			continue
		//		} else {
		//			SetAuthed(client, false)
		//			SendPacket(client, MakePacket(8, nil))
		//			continue
		//		}
		//	default:
		//		SendPacketAndDisconnect(client, MakePacket(0xff, nil))
		//		return
		//
		//	}
		//}
		if GetIsConnected(client)==false{
			if packet.command==0xA1{
				address,DecodeConnectPacketErr:=DecodeConnectPacketAndCheck(packet.data)
				if DecodeConnectPacketErr!=nil{
					SendPacketAndDisconnect(client, MakePacket(0xE1, nil))
					return
				}
				ProxyConn,dailerr:=net.Dial("tcp",address)
				log.DebugPrintNormal("客户端",client.RemoteAddr(),"想要代理",address)
				if dailerr!=nil{
					log.DebugPrintNormal("客户端",client.RemoteAddr(),"想要代理",address,"但连接失败了")
					SendPacketAndDisconnect(client, MakePacket(0xE2, nil))
					return
				}
				log.DebugPrintNormal("客户端",client.RemoteAddr(),"想要代理",address,"，连接成功")
				SetConnected(client,true)
				SetRemoteConn(client,ProxyConn)
				SendPacket(client, MakePacket(0xC1, nil))
				go ProxyRecvHandle(c,ProxyConn,client)
				//SendPacket(client, MakePacket(0xC1, nil))
				ProxySendHandle(c,ProxyConn)
				//return
			}
		}
	}
}
//从远端接受发给用户
func ProxyRecvHandle(c chan Packet,remote net.Conn,client net.Conn){
	defer func(){
		client.Close()
		log.DebugPrintAlert(client.RemoteAddr(),"代理线程关闭")
	}()
	buf:=make([]byte,1024)
	for{
		n,err:=remote.Read(buf)
		if err!=nil{
			return
		}
		SendPacket(client,buf[:n])
	}
}
func ProxySendHandle(c chan Packet,Remote net.Conn){
	for packet:=range c{
		switch packet.command {
			case 0:
				SendRawPacket(Remote,packet.data)
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
	//go StartHeartbeat()
	log.PrintSuccess("服务端启动成功")
	for {
		client, _ := l.Accept()
		//delete(Users,client.RemoteAddr()) //BUG!!!!!
		AddUser(client)
		SendPacket(client,util.GetPublicKey())
		c:=make(chan Packet,128)
		go ServeClient(client,c)
		go ServeCommand(client,c)
	}
}
