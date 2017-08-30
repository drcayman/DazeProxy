package obscure
import(
	"net"
)
//none-无伪装
type none struct {
	RegArg string
}

func (this *none) Init(param string,server *interface{})(error){
	return nil
}
func (this *none) Action(conn net.Conn , server *interface{}) (error){
	return nil
}
