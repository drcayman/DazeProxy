package net

import (
	"net"
	"../config"
	"fmt"
	"os"
	"encoding/binary"
	"bytes"
	"errors"
	"io/ioutil"
	"../util"
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

//客户端发送到服务端
//命令1代表需要公钥
//命令2代表设定key

//服务端发送到客户端
//命令1代表data里面是key
//命令2代表断开性错误（没交换key就执行其他操作）
//命令3代表key长度错误
//命令4代表接受key
func MakePacket(command byte,content []byte) []byte{
	if content==nil{
		content=[]byte{0x0}
	}
	//取内容长度
	ContentLen:=uint16(len(content))
	if ContentLen>0xFFFF {
		return nil
	}
	ContentLenBuffer:=bytes.NewBuffer([]byte{})
	//缓冲区，5=头部+命令+16位长度+尾部
	buf:=make([]byte,5+len(content))
	//填充头部和命令
	buf[0]=0xF1
	buf[1]=command
	//填充内容长度
	binary.Write(ContentLenBuffer,binary.BigEndian,ContentLen)
	copy(buf[2:],ContentLenBuffer.Bytes())
	//填充内容
	copy(buf[4:],content)
	//填充尾部
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
		client.Close()
	}
}
func SendPacketAndDisconnect(client net.Conn,data []byte){
	client.Write(data)
	client.Close()
}
func ServeClient(client net.Conn,c chan Packet){
	defer func(){
		close(c)
		DisconnectAndDeleteUser(client)
		if config.Config.IsDebug{
			fmt.Println("[!]",client.RemoteAddr(),"连接线程关闭")
		}
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
func ServeCommand(client net.Conn,c chan Packet) {
	defer func(){
		client.Close()
		if config.Config.IsDebug{
			fmt.Println("[!]",client.RemoteAddr(),"处理线程关闭")
		}
	}()
	for packet:=range c{
		if packet.command==0x1{
			PublicKeyBuf,PublicKeyErr:=ioutil.ReadFile("public.pem")
			if PublicKeyErr!=nil{
				fmt.Println("[×]公钥文件丢失！！系统强制退出")
				os.Exit(-1)
			}
			go SendPacket(client, MakePacket(1,PublicKeyBuf))
			continue
		}else if packet.command==0x2{
			Debuf,DeErr:=util.DecryptRSA(packet.data)
			if DeErr!=nil || len(Debuf)!=32{
				SendPacketAndDisconnect(client,MakePacket(3,nil))
				return
			}
			SetKeyExchange(client,true)
			SetAESKey(client,Debuf)
			go SendPacket(client, MakePacket(4,nil))
			continue
		}else{
			if Users[client.RemoteAddr()].IsKeyExchange==false{
				SendPacketAndDisconnect(client,MakePacket(2,nil))
				return
			}
		}

	}
}
func StartServer(){
	l,err:=net.Listen("tcp",":"+config.Config.ServerPort)
	if err!=nil{
		fmt.Println("[×]服务端启动失败（原因：",err.Error(),")")
		os.Exit(-1)
	}
	Server=&l
	//go StartHeartbeat()
	fmt.Println("[√]服务端启动成功")
	for {
		client, _ := l.Accept()
		//delete(Users,client.RemoteAddr()) //BUG!!!!!
		AddUser(client)
		c:=make(chan Packet,10)
		go ServeClient(client,c)
		go ServeCommand(client,c)
	}
}
