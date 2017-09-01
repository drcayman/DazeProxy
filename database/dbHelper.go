package database
import (
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)
var engine *xorm.Engine
type User struct {
	Id int64
	Username string
	Password string `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
	Expired time.Time
	Group string

}
func (User) TableName() string {
	return "DF_User"
}
func CheckUserPass(username string,password string) (bool,time.Time,string){
	var user User
	b,_:=engine.Where("username = ?",username).And("password = ?",password).Get(&user)
	return b,user.Expired,user.Group
}
func GetUserCount() int64 {
	var user User
	b,_:=engine.Count(&user)
	return b
}
func LoadDatabase(driver string,connectionString string){
	var err error
	engine,err=xorm.NewEngine(driver,connectionString)
	if err!=nil{
		log.Fatal("数据库连接失败！原因：",err)
	}
	err=engine.Sync2(new(User))
	if err!=nil{
		log.Fatal("数据库加载失败！原因：",err)
	}
	count:=GetUserCount()
	log.Println("数据库连接成功！用户数：",count)
	if count==0{
		log.Println("检测到用户数为0，请自行通过数据库管理工具添加用户！")
	}
}