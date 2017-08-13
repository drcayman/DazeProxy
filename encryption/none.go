package encryption

type none struct {

}
func (this *none) Init(arg string)(error){
	return nil
}
func (this *none)Encrypt(data []byte,len int)([]byte,int,error){
	return data,len,nil
}
func (this *none)Decrypt(data []byte,len int)([]byte,int,error){
	return data,len,nil
}
