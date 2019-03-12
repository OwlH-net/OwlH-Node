// package flock

// import (
//     "github.com/astaxie/beego/logs"
//     // "os"
//     // "os/exec"
//     // "strings"
//     // "regexp"
// 	// "owlhnode/utils"
// 	// "owlhnode/database"
// 	// "io/ioutil"
// 	// "errors"
//     "encoding/json"
//     // "sync"
// 	// "time"
// 	// "sys"
//     // "strconv"
// )


// //configFile := "/etc/owlh/owlh.conf"


// func LoadInventory()(data map[string]string, err error){

// 	Leer de DB la lista de Servidores y devolverlos.

// 	body, _ := ioutil.ReadAll(configFile)
//     err = json.Unmarshal(body, &conf)
//     if err != nil {
//         return nil,err
// 	}
// 	return conf,nil
// }

// func PrintInventory(owlhs map[string]string)(){
// 	for owlh := range owlhs{
// 		logs.Info("Primitive Inventory: "+owlh)
// 	}
// }

// func InventoryRun()(){
// 	var owlhs = LoadInventory()
// 	for owlh := range owlhs{
// 		//FlockLogger("Check for owlh name -> "+owlh["name"]+" , owlh IP -> "+owlh["ip"])
// 		alive, ssh = CheckOwlhAlive(owlh)
// 		if alive{
// 			//FlockLogger(">>> as owlh name -> "+owlh["name"]+" is alive with check status")
// 			running, status_ok := GetStatusSniffer(owlh, ssh)
// 			//FlockLogger(">>> Running "+running+", Status "+status_ok)
// 			if running {
// 				if !status_ok{
// 					StopSniffer(owlh, ssh)
// 				}
// 			} else if status_ok {
// 				RunSniffer(owlh, ssh)
// 			}
// 			GetFileList(owlh, ssh)
// 			ssh.Close()
// 		}
// 	}
// }