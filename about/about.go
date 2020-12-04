package about

import (
    // "bufio"
    // "encoding/json"
    // "errors"
    "github.com/astaxie/beego/logs"
    // "owlhnode/suricata"
    // "github.com/hpcloud/tail"
    // "io/ioutil"
    // "net"
    // "os"
    // "owlhnode/database"
    // "owlhnode/geolocation"
    // "owlhnode/utils"
    // "strconv"
    // "strings"
    // "syscall"
    // "time"
    // "regexp"
)

type AboutST struct {
}

var Version = "Version OwlH Node v0.17.4 build 03122020.2100"

func about() (aboutStruct AboutST) {
    logs.Info("About -> get node details")
    logs.Info("Node Version -> %s", Version)
    logs.Info("Node Name ->")
    logs.Info("Node IP ->")
    // suricataVersion, versionError := suricata.SuricataVersion()
    // logs.Info("Suricata Version -> %+v", suricataVersion)
    return aboutStruct
}

func GetVersion() (version string) {
    return Version
}
