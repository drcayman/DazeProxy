package database
import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
	"DazeProxy/util"
	"errors"
)
var engine *xorm.Engine
type User struct {
	Id int64
	Username string
	Password string `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
}
func GetAllUser() (a []User){
	engine.Find(&a)
	return
}
func GetUserById(id int)(User,bool){
	var user User
	b,_:=engine.Where("id = ?", id).Get(&user)
	return user,b
}
func DeleteById(id int) bool{
	v,b:=GetUserById(id)
	if b==false{
		return b
	}
	_,err:=engine.Delete(&v)
	return err==nil
}
func GetUserByUserName(username string) (User,bool){
	var user User
	b,_:=engine.Where("username = ?", username).Get(&user)
	return user,b
}
func AddUser(username string,password string) error{
	if _,b:=GetUserByUserName(username);b{
		return errors.New("exist user")
	}
	user:=User{Username:username,Password:util.GetDoubleMd5(password)}
	_,err:=engine.Insert(&user)
	return err
}
func CheckUserPass(username string,password string) bool{
	var user User
	b,_:=engine.Where("username = ?",username).And("password = ?",util.GetDoubleMd5(password)).Get(&user)
	return b
}
func GetUserCount() int64 {
	var user User
	b,_:=engine.Count(&user)
	return b
}
func init(){
	var err error
	engine,err=xorm.NewEngine("sqlite3","daze.db")
	if err!=nil{
		log.Fatal("数据库连接失败！")
	}
	err=engine.Sync2(new(User))
	if err!=nil{
		log.Fatal("数据库加载失败！")
	}
	log.Print("数据库连接成功！用户数：",GetUserCount())
}