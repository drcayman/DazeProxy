package proxy

import (
	"net"
	"DazeProxy/config"
	"DazeProxy/database"
	"encoding/binary"
	"bytes"
	"errors"
	"log"
	"encoding/json"
	"DazeProxy/encryption"
	"DazeProxy/disguise"
	. "DazeProxy/common"
)
var Server *net.Listener
type JsonAuth struct{
	Username string
	Password string
	Net string
	Host string
	Port string
}

//生成控制数据包
//[头部：F1][命令][内容长度][内容][尾部：F2]
//头部尾部均为1字节
//命令的长度为1字节
//内容长度的长度为2字节
//内容的长度不限
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
	binary.Write(ContentLenBuffer,binary.LittleEndian,ContentLen)
	copy(buf[2:],ContentLenBuffer.Bytes())
	copy(buf[4:],content)
	buf[len(buf)-1]=0xF2
	return buf
}
//解析控制数据包，解析看上面
func DePacket(buf []byte) (byte,[]byte,error){
	if len(buf)<6 || buf[0]!=0xF1 || buf[len(buf)-1]!=0xF2{
		return 0,nil,errors.New("error1")
	}
	ContentLen:=int(buf[2])+int(buf[3])*256
	if len(buf)-5!=int(ContentLen){
		return 0,nil,errors.New("error2")
	}
	return buf[1],buf[4:4+ContentLen],nil
}

//打包并发送数据包给客户端
//[头部][内容]
//头部和内容分开加密
//头部为4字节,[FB][内容长度][FC]
//内容无限长
func SendPacket(client *User,data []byte){
	var bufffer *bytes.Buffer
	packets,encErr:=client.ProxyUnit.Encryption.Encrypt(&client.EncReserved,&client.ProxyUnit.EncReserved,data)
	if encErr!=nil{
		return
	}
	for _,pkt:=range packets{
		dataLen:=len(pkt)
		header:=[]byte{0xFB,byte(dataLen%0x100),byte(dataLen/0x100),0xFC}
		headerEncoded,headerEncodedErr:=client.ProxyUnit.Encryption.Encrypt(&client.EncReserved,&client.ProxyUnit.EncReserved,header)
		if headerEncodedErr!=nil{
			return
		}
		bufffer=bytes.NewBuffer(headerEncoded[0])
		bufffer.Write(pkt)
		client.Conn.Write(bufffer.Bytes())
	}

}

//发送数据包然后断开客户端
func SendPacketAndDisconnect(client *User,data []byte){
	SendPacket(client,data)
	client.Conn.Close()
}

//解析客户端发过来的数据包，解析看上面
func ReadFromClient(client *User) ([]byte,error){
	HeaderBuf:=make([]byte,4)
	n,err:=client.Conn.Read(HeaderBuf)
	if n<4 ||err!=nil{
		return nil,errors.New("read header error ")
	}
	headerDecode,err:=client.ProxyUnit.Encryption.Decrypt(&client.EncReserved,&client.ProxyUnit.EncReserved,HeaderBuf)
	if err!=nil || headerDecode[0]!=0xFB || headerDecode[3]!=0xFC{
		return nil,errors.New("decode header error")
	}
	PacketLen:=int(headerDecode[1])+int(headerDecode[2])*256
	buf:=make([]byte,PacketLen)
	pos:=0
	for{
		n,err:=client.Conn.Read(buf[pos:])
		if err!=nil{
			return nil,errors.New("read body error")
		}
		PacketLen-=n
		pos+=n
		if PacketLen<0{
			return nil,errors.New("body len error")
		}
		if PacketLen==0{
			break
		}
	}
	//copy(buf,HeaderBuf)
	DecodeBuf,err:=client.ProxyUnit.Encryption.Decrypt(&client.EncReserved,&client.ProxyUnit.EncReserved,buf)
	if err!=nil{
		return nil,errors.New("decode body error")
	}
	return DecodeBuf,nil
}

//接待客户端
func ServeClient(client *User){
	flag:=0
	defer func(){
		if flag==0{
			DisconnectUser(client)
			if config.Config.IsDebug {
				log.Println(client.Conn.RemoteAddr(), "连接线程关闭")
			}
		}else if config.Config.IsDebug {
				log.Println(client.Conn.RemoteAddr(), "连接线程已进入代理模式")
		}
	}()
	for{
		buf,err:=ReadFromClient(client)
		if err!=nil{
			return
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

//解析客户端发送过来的域名
func DecodeAddressAndCheck(host string,port string,IPv6ResolvePrefer bool)(string,error){
	var ips *net.IPAddr
	var ResolveErr error=nil
	if IPv6ResolvePrefer{
		ips,ResolveErr=net.ResolveIPAddr("ip6",host)
	}
	if ips==nil{
		ips,ResolveErr=net.ResolveIPAddr("ip",host)
	}
	if ResolveErr!=nil{
		return "",errors.New("error2")
	}
	if ips.IP.IsLoopback(){
		return "",errors.New("error3")
	}
	ipstring:=ips.String()
	if len(ipstring)>15{
		ipstring="["+ipstring+"]"
	}
	return ipstring+":"+port,nil
}

/*
客户端发送到服务端：
命令1代表读取公告
命令2代表登录并尝试代理，data里面为JSON格式，例如
{
"username":"abc",
"password":"456",
"net":"tcp",
"host":"www.baidu.com",
"port":"80"
}

服务端发送到客户端：
命令E1代表IP地址格式错误
命令E2代表连接失败
命令E3代表登录失败
命令E4代表协议错误
命令C1代表成功连接
命令FF代表指令没法识别

*/

//处理客户端发送过来的控制数据包
func ServeCommand(client *User,command byte,data []byte) int {
		switch command{
		case 0x01:{

		}
		case 0x02:{
			authinfo:=JsonAuth{}
			err:=json.Unmarshal(data,&authinfo)
			if err!=nil || !database.CheckUserPass(authinfo.Username,authinfo.Password){
				if config.Config.IsDebug {
					log.Printf("客户端%s想要代理[%s]%s:%s但验证失败了！\n",client.Conn.RemoteAddr(),authinfo.Net,authinfo.Host,authinfo.Port)
				}
				SendPacketAndDisconnect(client, MakePacket(0xE3, nil))
				return 0
			}
			client.IsAuth=true
			client.AuthHeartBeat<-1
			if authinfo.Net!="tcp" && authinfo.Net!="udp"{
				SendPacketAndDisconnect(client, MakePacket(0xE4, nil))
				return 0
			}
			address,DecodeConnectPacketErr:=DecodeAddressAndCheck(authinfo.Host,authinfo.Port,client.IPv6ResolvePrefer)
			if DecodeConnectPacketErr!=nil{
				SendPacketAndDisconnect(client, MakePacket(0xE1, nil))
				return 0
			}
			ProxyConn,dailerr:=net.Dial(authinfo.Net,address)
			if dailerr!=nil{
				if config.Config.IsDebug {
					log.Printf("客户端%s想要代理[%s]%s:%s但连接失败了！\n",client.Conn.RemoteAddr(),authinfo.Net,authinfo.Host,authinfo.Port)
				}
				SendPacketAndDisconnect(client, MakePacket(0xE2, nil))
				return 0
			}
			if config.Config.IsDebug {
				log.Printf("客户端%s成功代理[%s]%s:%s\n",client.Conn.RemoteAddr(),authinfo.Net,authinfo.Host,authinfo.Port)
			}
			client.RemoteConn=ProxyConn
			client.Network=authinfo.Net
			client.IsConnected=true
			SendPacket(client, MakePacket(0xC1, []byte(ProxyConn.RemoteAddr().String())))
			go BridgeClientToRemote(client,ProxyConn)
			go BridgeRemoteToClient(client,ProxyConn)
			return 1
			}
		default:{
			SendPacketAndDisconnect(client, MakePacket(0xFF, nil))
			return 0
			}
		}
	return 0
}

//关闭客户端的心跳chan
func CloseChan(client *User) {
	client.Locker.Lock()
	if !client.ChanCloseFlag{
	close(client.AuthHeartBeat)
	close(client.UDPAliveTime)
	}
	client.ChanCloseFlag=true
	client.Locker.Unlock()
}

//IO桥：客户端到目标服务器
func BridgeClientToRemote(client *User,Remote net.Conn){
	defer func(){
		DisconnectUserAndRemoteConn(client)
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

//IO桥：目标服务器到客户端
func BridgeRemoteToClient(client *User,Remote net.Conn){
	defer func(){
		DisconnectUserAndRemoteConn(client)
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

//启动代理服务单元
func StartServer(unit ProxyUnit){
	defer func(){
		if err := recover(); err != nil {
			log.Printf("代理服务单元（端口：%s）启动失败（原因：%s）\n",unit.Config.Port,err)
		}
	}()
	if unit.Config.Port==""{
		panic("端口不能为空")
	}

	//加载加密模块和初始化
	enc,encflag:=encryption.GetEncryption(unit.Config.Encryption)
	if !encflag{
		panic("加密方式"+unit.Config.Encryption+"不存在")
	}
	unit.Encryption=enc()
	encInitErr:=unit.Encryption.Init(unit.Config.EncryptionParam,&unit.EncReserved)
	if encInitErr!=nil{
		panic("加密方式"+unit.Config.Encryption+"加载错误！原因："+encInitErr.Error())
	}
	//加载伪装模块和初始化
	dsg,dsgflag:=disguise.GetDisguise(unit.Config.Disguise)
	if !dsgflag{
		panic("伪装方式"+unit.Config.Disguise+"不存在")
	}
	unit.Disguise=dsg()
	dsgInitErr:=unit.Disguise.Init(unit.Config.DisguiseParam,&unit.DsgReserved)
	if dsgInitErr!=nil{
		panic("伪装方式"+unit.Config.Disguise+"加载错误！原因："+dsgInitErr.Error())
	}

	l,err:=net.Listen("tcp",":"+unit.Config.Port)
	if err!=nil{
		panic(err.Error())
	}
	log.Printf("代理服务单元启动成功（端口：%s，加密方式：%s，加密参数：%s，伪装方式：%s，伪装参数：%s，优先解析IPV6：%v）\n",
		unit.Config.Port,
		unit.Config.Encryption,
		unit.Config.EncryptionParam,
		unit.Config.Disguise,
		unit.Config.DisguiseParam,
		unit.Config.IPv6ResolvePrefer)
	for {
		conn, AcceptErr := l.Accept()
		if AcceptErr!=nil{
			log.Printf("代理服务单元（端口：%s）接受客户端失败！（原因：%s）\n",unit.Config.Port,AcceptErr.Error())
			continue
		}
		if config.Config.IsDebug {
			log.Printf("代理服务单元（端口：%s）接受客户端（%s）\n",unit.Config.Port,conn.RemoteAddr())
		}
		client:=NewUser(conn)
		client.IPv6ResolvePrefer=unit.Config.IPv6ResolvePrefer
		client.ProxyUnit=&unit
		go NewClientComing(client)
	}
}

//断开用户
func DisconnectUser(client *User){
	CloseChan(client)
	client.Conn.Close()
}
//断开用户和目标服务器链接
func DisconnectUserAndRemoteConn(client *User){
	CloseChan(client)
	client.Conn.Close()
	client.RemoteConn.Close()
}

//新客户端到来时的准备工作
func NewClientComing(client *User){
	defer func(){
		if err := recover(); err != nil{
			if config.Config.IsDebug{
				log.Printf("客户端(%s)处理失败（原因：%s）\n",client.Conn.RemoteAddr(),err)
				DisconnectUser(client)
			}
		}
	}()
	go NewHeartBeatCountDown(client.AuthHeartBeat,5,client,"Auth or Connect")
	dsgErr:=client.ProxyUnit.Disguise.Action(client.Conn,&client.DsgReserved,&client.ProxyUnit.DsgReserved)
	if dsgErr!=nil{
		panic("伪装时出现错误："+dsgErr.Error())
	}
	encErr:=client.ProxyUnit.Encryption.InitUser(client.Conn,&client.EncReserved,&client.ProxyUnit.EncReserved)
	if encErr!=nil{
		panic("为用户初始化加密方式时出现错误："+encErr.Error())
	}
	go ServeClient(client)
}