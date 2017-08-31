package helper

import "log"

func DebugPrintln(msg string){
	if IsDebug{
		log.Println(msg)
	}
}
