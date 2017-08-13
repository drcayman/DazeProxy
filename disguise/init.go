package disguise

import (
	"strings"
	"net"
)
type DisguiseAction interface {
	Init(string)(error)
	Action(net.Conn,*interface{})(error)
}
type regfunc func()(DisguiseAction)
var disguiseMap map[string]regfunc

func GetDisguise(name string) (regfunc,bool){
	name=strings.ToLower(name)
	d,flag:=disguiseMap[name]
	return d,flag
}

func init(){
	disguiseMap=make(map[string]regfunc)

	//这里添加自己开发的伪装模块
	disguiseMap["none"]=func()(DisguiseAction){
		return DisguiseAction(new(none))
	}
}