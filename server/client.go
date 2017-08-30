package server

import(
	"fmt"
	"time"
	"net"
	"encoding/json"
	"github.com/crabkun/DazeProxy/util"
	"github.com/crabkun/DazeProxy/common"
	"github.com/crabkun/DazeProxy/helper"
)
type S_Client struct {
	//代理用户的套接字
	UserConn net.Conn
	//代理目标服务器TCP套接字
	RemoteTCPConn net.Conn
	//代理目标服务器UDP套接字
	RemoteUDPConn *net.UDPConn
	//是否已登录
	Authed bool
	//是否已连接
	Connected bool
	//代理目标网络协议
	Network string
	//预留给加密
	EReserved interface{}
	//预留给混淆
	ObReserved interface{}
	//所属代理实例
	Proxy *common.S_proxy

}
func (client *S_Client)Decode(data []byte) []byte{
	buf,err:=client.Proxy.E.Decrypt(&client.EReserved,&client.Proxy.EReserved,data)
	if err!=nil{
		panic(err.Error())
	}
	return buf
}
func (client *S_Client)Encode(data []byte) []byte{
	buf,err:=client.Proxy.E.Encrypt(&client.EReserved,&client.Proxy.EReserved,data)
	if err!=nil{
		panic(err.Error())
	}
	return buf
}
func (client *S_Client)Disconnect(){
	if client.Connected{
		if client.Network=="tcp"{
			client.RemoteTCPConn.Close()
		}else{
			client.RemoteUDPConn.Close()
		}
		client.Connected=false
	}
	client.UserConn.Close()

}
func (client *S_Client)Read() []byte {
	//读取头部
	headerEncoded:=client.SafeRead(client.UserConn,4)
	//解码头部
	header:=client.Decode(headerEncoded)
	if header[0]!=0xF1 && header[3]!=0xF2{
		panic("头部不匹配")
	}
	//读取负载
	length:=int(header[1])+int(header[2])*256
	if length==0{
		panic("长度错误")
	}
	//解码负载
	bodyEncoded:=client.SafeRead(client.UserConn,length)
	return client.Decode(bodyEncoded)
}
func (client *S_Client)Write(data []byte){
	length:=len(data)
	if data==nil || length==0 || length>65535{
		panic("数据长度不正确(1-65535)")
	}
	header:=[]byte{0xF1,byte(length%0x100),byte(length/0x100),0xF2}
	client.SafeSend(client.Encode(header),client.UserConn)
	client.SafeSend(client.Encode(data),client.UserConn)
}
func (client *S_Client)SafeSend(data []byte,conn net.Conn){
	length:=len(data)
	for pos:=0;pos<length;{
		n,err:=conn.Write(data[pos:])
		if err!=nil {
			panic("连接正常断开")
		}
		pos+=n
	}
}
func (client *S_Client)SafeRead(conn net.Conn,length int) ([]byte) {
	buf:=make([]byte,length)
	for pos:=0;pos<length;{
		n,err:=conn.Read(buf[pos:])
		if err!=nil {
			panic("连接正常断开")
		}
		pos+=n
	}
	return buf
}
func (client *S_Client)Serve(){
	var err error
	authjson:=client.Read()
	authinfo:=new(common.Json_Auth)
	err=json.Unmarshal(authjson,authinfo)
	if err!=nil{
		client.WriteJsonRet(-1,"")
		panic("登录数据解码错误:"+err.Error())
	}
	if authinfo.Net!="tcp" && authinfo.Net!="udp"{
		client.WriteJsonRet(-2,"")
		panic("网络协议有误")
	}
	//验证成功，去除验证超时
	client.Authed=true
	client.Network=authinfo.Net
	//连接对面
	if client.Network=="tcp"{
		host,_,err:=net.SplitHostPort(authinfo.Host)
		if err!=nil{
			client.WriteJsonRet(-3,"")
			panic("ip地址有误")
		}
		if ip:=net.ParseIP(host);ip!=nil && ip.IsMulticast(){
			client.WriteJsonRet(-3,"")
			panic("ip地址有误")
		}
		client.RemoteTCPConn,err=net.Dial(authinfo.Net,authinfo.Host)
		if err!=nil{
			client.WriteJsonRet(-4,"")
			panic("无法连接指定地址"+err.Error())
		}
		client.WriteJsonRet(1,client.RemoteTCPConn.RemoteAddr().String())
		go client.BridgeTCPClientToRemote()
		go client.BridgeTCPRemoteToClient()
		helper.DebugPrintln(fmt.Sprintf("客户端(%s)连接了[%s]%s",client.UserConn.RemoteAddr(),authinfo.Net,authinfo.Host))
	}else{
		client.RemoteUDPConn,err=net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
		if err!=nil{
			client.WriteJsonRet(-4,"")
			panic("无法连接指定地址"+err.Error())
		}
		if client.Network=="udp"{
			client.RemoteUDPConn.SetReadDeadline(time.Now().Add(time.Second*5))
			client.RemoteUDPConn.SetWriteDeadline(time.Now().Add(time.Second*5))
		}
		client.WriteJsonRet(1,"")
		go client.BridgeUDPClientToRemote()
		go client.BridgeUDPRemoteToClient()
		helper.DebugPrintln(fmt.Sprintf("客户端(%s)建立了UDP通道",client.UserConn.RemoteAddr()))
	}
	client.UserConn.SetDeadline(time.Time{})
	client.Connected=true
}
func (client *S_Client)BridgeTCPClientToRemote(){
	defer func(){
		recover()
		client.RemoteTCPConn.Close()
	}()
	for{
		client.SafeSend(client.Read(),client.RemoteTCPConn)
	}
}
func (client *S_Client)BridgeTCPRemoteToClient(){
	defer func(){
		recover()
		client.UserConn.Close()
	}()
	buf:=make([]byte,4096)
	for{
		n,err:=client.RemoteTCPConn.Read(buf)
		if err!=nil{
			panic(err)
		}
		client.Write(buf[:n])
	}
}
func (client *S_Client) BridgeUDPClientToRemote(){
	defer func(){
		recover()
		client.RemoteUDPConn.Close()
	}()
	var UDP common.Json_UDP
	var ADDR *net.UDPAddr
	var LastAddr string
	var err error
	for{
		buf:=client.Read()
		if err=json.Unmarshal(buf,&UDP);err!=nil{
			return
		}
		if LastAddr!=UDP.Host{
			ADDR,err=net.ResolveUDPAddr("udp",UDP.Host)
			if err!=nil{
				return
			}
			LastAddr=UDP.Host
		}
		_,err=client.RemoteUDPConn.WriteToUDP(UDP.Data,ADDR)
		if err!=nil{
			return
		}

		client.RemoteUDPConn.SetWriteDeadline(time.Now().Add(time.Second*5))
	}
}
func (client *S_Client) BridgeUDPRemoteToClient(){
	defer func(){
		recover()
		client.UserConn.Close()
	}()
	buf:=make([]byte,65507)
	var err error
	var n int
	var addr net.Addr
	var UDP common.Json_UDP
	for{
		n,addr,err=client.RemoteUDPConn.ReadFrom(buf)
		if err!=nil{
			return
		}
		UDP.Host=addr.String()
		UDP.Data=buf[:n]
		ret,err:=json.Marshal(UDP)
		if err!=nil{
			return
		}
		client.Write(ret)
		client.RemoteUDPConn.SetReadDeadline(time.Now().Add(time.Second*5))
	}
}
func (client *S_Client)WriteJsonRet(code int,data string){
	authretbuf,_:=json.Marshal(common.Json_Ret{
		Code:code,
		Data:data,
		Spam:util.GetRandomString(256),
	})
	client.Write(authretbuf)
}
func PackNewUser(conn net.Conn,s *common.S_proxy) *S_Client{
	return &S_Client{
		Proxy:s,
		UserConn:conn,
	}
}
func NewClientComing(client *S_Client){
	defer func(){
		if err := recover(); err != nil{
				helper.DebugPrintln(fmt.Sprintf("客户端(%s)处理结束（原因：%s）",client.UserConn.RemoteAddr(),err))
				client.Disconnect()
		}
	}()
	//设置验证超时时间
	client.UserConn.SetDeadline(time.Now().Add(time.Second*5))
	//开始混淆
	obErr:=client.Proxy.Ob.Action(client.UserConn,&client.Proxy.ObReserved)
	if obErr!=nil{
		panic("伪装时出现错误："+obErr.Error())
	}
	//为用户初始化加密方式
	eErr:=client.Proxy.E.InitUser(client.UserConn,&client.EReserved,&client.Proxy.EReserved)
	if eErr!=nil{
		panic("为用户初始化加密方式时出现错误："+eErr.Error())
	}
	client.Serve()
}