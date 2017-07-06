package log1

import (
	"fmt"
	"time"
	"os"
	"syscall"
	"runtime"
	"DazeProxy/config"
)
const(
	linux_yellow="\033[01;33m"
	linux_red="\033[22;31m"
	linux_green="\033[22;32m"
	linux_turnoff="\033[0m"
	windows_turnoff=15
	windows_yellow=14
	windows_red=12
	windows_green=10
)
var kernel32 *syscall.DLL
var SetConsoleTextAttribute *syscall.Proc
var StdoutHandle syscall.Handle
func appendArr(src *[]interface{},target *[]interface{}){
	for _,v :=range *src{
		*target=append(*target,v)
	}
}
func DebugPrintAlert(a ...interface{}){
	if config.Config.IsDebug{
		PrintAlert(a)
	}
}
func DebugPanic(a ...interface{}){
	if config.Config.IsDebug{
		PrintPanic(a)
	}
}
func DebugPrintPanicWithoutExit(a ...interface{}){
	if config.Config.IsDebug{
		PrintPanicWithoutExit(a)
	}
}
func DebugPrintSuccess(a ...interface{}){
	if config.Config.IsDebug{
		PrintSuccess(a)
	}
}
func DebugPrintNormal(a ...interface{}){
	if config.Config.IsDebug{
		PrintNormal(a)
	}
}
func PrintAlert(a ...interface{}){
	var tmp []interface{}
	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_yellow)
	}
	tmp=append(tmp, time.Now().Format("2006-01-02 03:04:05.000 | [！]"))
	appendArr(&a,&tmp)

	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_turnoff)
	}else{
		SetConsoleTextAttribute.Call(uintptr(StdoutHandle),uintptr(windows_yellow))
	}
	_Print(tmp)
}
func PrintPanic(a ...interface{}){
	var tmp []interface{}
	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_red)
	}
	tmp=append(tmp, time.Now().Format("2006-01-02 03:04:05.000 | [×]"))
	appendArr(&a,&tmp)

	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_turnoff)
	}else{
		SetConsoleTextAttribute.Call(uintptr(StdoutHandle),uintptr(windows_red))
	}
	_Print(tmp)
	os.Exit(-1)
}
func PrintPanicWithoutExit(a ...interface{}){
	var tmp []interface{}
	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_red)
	}
	tmp=append(tmp, time.Now().Format("2006-01-02 03:04:05.000 | [×]"))
	appendArr(&a,&tmp)
	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_turnoff)
	}else{
		SetConsoleTextAttribute.Call(uintptr(StdoutHandle),uintptr(windows_red))
	}
	_Print(tmp)
}
func PrintSuccess(a ...interface{}){
	var tmp []interface{}
	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_green)
	}
	tmp=append(tmp, time.Now().Format("2006-01-02 03:04:05.000 | [√]"))
	appendArr(&a,&tmp)

	if runtime.GOOS!="windows"{
		tmp=append(tmp,linux_turnoff)
	}else{
		SetConsoleTextAttribute.Call(uintptr(StdoutHandle),uintptr(windows_green))
	}
	_Print(tmp)
}
func PrintNormal(a ...interface{}){
	if runtime.GOOS=="windows"{
		SetConsoleTextAttribute.Call(uintptr(StdoutHandle),uintptr(windows_turnoff))
	}
	var tmp []interface{}
	tmp=append(tmp, time.Now().Format("2006-01-02 03:04:05.000 |"))
	appendArr(&a,&tmp)
	_Print(tmp)
}
func _Print(a []interface{}){
	fmt.Println(a)
}
func init(){
	if runtime.GOOS=="windows" {
		kernel32, _ = syscall.LoadDLL("kernel32.dll")
		SetConsoleTextAttribute, _ = kernel32.FindProc("SetConsoleTextAttribute")
		StdoutHandle, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	}
}