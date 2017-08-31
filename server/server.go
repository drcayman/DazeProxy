package server

import(
	"github.com/crabkun/DazeProxy/common"
	"github.com/crabkun/DazeProxy/encryption"
	"github.com/crabkun/DazeProxy/obscure"
	"github.com/crabkun/DazeProxy/helper"
	"log"
	"net"
	"fmt"
	"sync"
)
//启动代理服务单元
func StartServer(s common.S_proxy,wg sync.WaitGroup){
	defer func(){
		if err := recover(); err != nil {
			log.Printf("代理实例（端口：%s）启动失败（原因：%s）\n",s.Port,err)
		}
		wg.Done()
	}()
	if s.Port==""{
		panic("端口不能为空")
	}

	//加载加密模块和初始化
	E,ExistFlag:=encryption.GetEncryption(s.Encryption)
	if !ExistFlag{
		panic("加密方式"+s.Encryption+"不存在")
	}
	s.E=E()
	EInitErr:=s.E.Init(s.EncryptionParam,&s.EReserved)
	if EInitErr!=nil{
		panic("加密方式"+s.Encryption+"加载错误！原因："+EInitErr.Error())
	}

	//加载伪装模块和初始化
	Ob,ExistFlag:=obscure.GetObscure(s.Obscure)
	if !ExistFlag{
		panic("伪装方式"+s.Obscure+"不存在")
	}
	s.Ob=Ob()
	ObInitErr:=s.Ob.Init(s.ObscureParam,&s.ObReserved)
	if ObInitErr!=nil{
		panic("伪装方式"+s.Obscure+"加载错误！原因："+ObInitErr.Error())
	}

	//开始监听
	l,err:=net.Listen("tcp",":"+s.Port)
	if err!=nil{
		panic(err.Error())
	}
	log.Printf("\n代理实例启动成功\n（\n端口：%s\n加密方式：%s\n加密参数：%s\n伪装方式：%s\n伪装参数：%s\n）\n",
		s.Port,
		s.Encryption,
		s.EncryptionParam,
		s.Obscure,
		s.ObscureParam)
	for {
		conn, AcceptErr := l.Accept()
		if AcceptErr!=nil{
			log.Printf("代理实例（端口：%s）接受客户端失败！（原因：%s）\n",s.Port,AcceptErr.Error())
			if err,ok:=AcceptErr.(net.Error);ok && !err.Temporary(){
				return
			}
			continue
		}
		helper.DebugPrintln(fmt.Sprintf("代理实例（端口：%s）接受客户端（%s）",s.Port,conn.RemoteAddr()))
		client:=PackNewUser(conn,&s)
		go NewClientComing(client)
	}

}
