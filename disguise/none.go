package disguise

type none struct {

}

func (this *none) Init(arg string)(error){
	return nil
}
func (this *none) Action() (error){
	return nil
}