package obscure

import(
	"net"
)
//none-无伪装
type none struct {
}

func (this *none) Init(param string)(error){
	return nil
}
func (this *none) Action(conn net.Conn) (error){
	return nil
}
func init(){
	RegisterObscure("none",new(none))
}
