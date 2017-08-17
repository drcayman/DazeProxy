package console

import (
	"fmt"
	"bufio"
	"os"
	"DazeProxy/util"
	"strings"
	"DazeProxy/database"
)

func ShowMenu(){
	fmt.Println("**********命令列表**********")
	fmt.Println("help 显示此帮助")
	fmt.Println("users 显示所有用户")
	fmt.Println("add 增加一个新用户（比如add test 1234意思是增加一个用户名为test，密码为1234的用户名）")
	fmt.Println("del 删除一个用户（比如del 4就是删除掉ID为4的用户）")
	fmt.Println("exit 退出用户管理")
	fmt.Println("****************************")
}
func Start(){
	fmt.Println("**********注意！**********")
	fmt.Println(" ")
	fmt.Println("你目前在用户管理模式，代理功能没有生效，仅供管理用户！")
	fmt.Println(" ")
	fmt.Println("**********注意！**********")
	ShowMenu()
	r:=bufio.NewReader(os.Stdin)
	command:=""
	for{
		fmt.Print(">>>>>>")
		buf,_,err:=r.ReadLine()
		if err!=nil{
			return
		}
		bufstr:=util.B2s(buf)
		fmt.Sscanf(bufstr,"%s",&command)
		switch strings.ToLower(command) {
		case "help":
			ShowMenu()
		case "users":
			users:=database.GetAllUser()
			fmt.Printf("一共有%d个用户\n",database.GetUserCount())
			for _,v:=range users{
				fmt.Printf("ID:%d   用户名：%s\n",v.Id,v.Username)
			}
		case "add":
			var username,password string
			n,_:=fmt.Sscanf(bufstr,"%s%s%s",&command,&username,&password)
			if n!=3{
				fmt.Println("命令格式错误")
				continue
			}
			err:=database.AddUser(username,password)
			if err!=nil{
				fmt.Println("添加用户失败，原因：",err.Error())
				continue
			}
			fmt.Printf("添加用户成功！\n用户名:%s\n密码:%s\n",username,password)
		case "del":
			var id int
			n,_:=fmt.Sscanf(bufstr,"%s%d",&command,&id)
			if n!=2{
				fmt.Println("命令格式错误")
				continue
			}
			if database.DeleteById(id){
				fmt.Println("删除用户成功")
			}else{
				fmt.Println("删除用户失败，或许是ID错误了？")
			}
		case "exit":
			return
		default:
			fmt.Println("命令格式错误，请输入help来查看帮助")
		}
	}
}