// package flock

// import (
//     "github.com/astaxie/beego/logs"
//     "os"
//     // "os/exec"
//     // "strings"
//     // "regexp"
// 	// "owlhnode/utils"
// 	// "owlhnode/database"
// 	// "io/ioutil"
// 	// "errors"
//     // "encoding/json"
//     // "sync"
// 	"time"
// 	"sys"
//     // "strconv"
// )

// func InitFlockLogger()(){
//     var logfile string
//     bufSize := 0
//     logfile,err := file(GetItem("logfile"),'a',bufSize)
//     if err!=nil {
//         logs.Err("InitFlockLogger error --> "+err.Error())
//     }
// }

// func KillFlockLogger()(){
//     FlockLogger("Closing log output")
// }

// func FlockLogger(text, "Info", "flock", getPID())(){
//     logfile.write(datetime.utcnow().strftime("%a %d %b %Y %H:%M:%S.%f") + " [" + proc + "] (" + str(id) +") [" + level+ "]: " + text + "\n")
// }

// logFile=""
// InitFlockLogger()