package common

type Json_Auth struct{
	Username string
	Password string
	Net string
	Host string
}
type Json_UDP struct {
	Host string
	Data []byte
}
type Json_Ret struct{
	Code int
	Data string
	Spam string
}