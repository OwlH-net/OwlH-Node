package flock

import (
    "github.com/astaxie/beego/logs"
    "os"
    // "os/exec"
    // "strings"
    // "regexp"
	// "owlhnode/utils"
	// "owlhnode/database"
	// "io/ioutil"
	// "errors"
    // "encoding/json"
    // "sync"
	"time"
	"sys"
    // "strconv"
)

func Main()(){
    FlockLogger("init", "INFO")
    if !AmiRunning() {
        for !Killme() {
            InventoryRun()
            time.Sleep(time.Second * 1)
        }
        ByeBye()
    }
}

Main()