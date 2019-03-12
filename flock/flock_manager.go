// package flock
// import (
// 	"encoding/json"
// 	"github.com/astaxie/beego/logs"
// 	"os"
// )



// func Killme()(bool killme){
// 	return killme
// }

// func RegisterFlock()(){
// 	FlockLogger (getPID())
// 	file(GetItem("pidfile"), w).write(string(os.Getpid()))
// }

// func AmiRunning()(){
// 	if _, err := os.Stat(GetItem("pidfile")); err == nil {
// 		// path/to/whatever exists
// 		FlockLogger("I'm running, we don't need two of us, exiting...")
// 		FlockLogger("If you think I'm not running, please check process list and "+GetItem("pidfile"))
// 		return true
// 	}
// 	RegisterFlock()
// }

// func ByeBye()(){
// 	FlockLogger("Time to go home. See you soon")
// 	KillFlockLogger()
// 	os.Remove(conf("pidfile"))
// }

// func CheckPID(pid int)(status bool){
// 	//Check For the existence of a unix pid.
// 	err := os.Kill(pid, 0)
// 	if err!= nil{
// 		return false
// 	}else{
// 		return true
// 	}
// }