package helper

import "log"

func DebugPrintln(msg string){
	if IsDebug{
		log.Println(msg)
	}
}
func Println(msg string){
		log.Println(msg)
}