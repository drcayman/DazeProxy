package disguise

import "net"

type none struct {

}

func (this *none) Init(arg string)(error){
	return nil
}
func (this *none) Action(conn net.Conn ,arg *interface{}) (error){
	return nil
}