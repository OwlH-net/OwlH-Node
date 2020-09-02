package analyzer

import (
    // "bufio"
    "encoding/json"
    // "errors"
    "github.com/astaxie/beego/logs"
    // "github.com/hpcloud/tail"
    "io/ioutil"
    // "net"
    // "os"
    // "owlhnode/database"
    // "owlhnode/geolocation"
    "owlhnode/utils"
    // "strconv"
    // "strings"
    // "syscall"
    // "time"
    "regexp"
)

type JA3 struct {
    Enable       bool   `json:"enable"`
    Verbose      bool   `json:"verbose"`
    Stats        bool   `json:"stats"`
    SourceURL    string `json:"sourceurl"`
    Localja3file string `json:"localja3file"`
    // OutputFile             string   `json:"outputfile"`
    // Prefilter              string   `json:"prefilterfile"`
    // Postfilter             string   `json:"postfilterfile"`
    // Tagfile                string   `json:"tagsfile"`
    // Srcfiles               []string `json:"srcfiles"`
    // Feedfiles              []Feedfile
    // Suricatasocket         string `json:"suricatasocket"`
    // SuricatasocketEnabled  bool   `json:"suricatasocketenabled"`
    // WazuhSocketEnabled     bool   `json:"wazuhsocketenabled"`
    // TimebetweenStatusCheck int    `json:"timebetweenstatuscheck"`
    // Timetowaitforfile      int    `json:"timetowaitforfilek"`
    // ChannelWorkers         int    `json:"channelworkers"`
}

type JA3Hash struct {
    Md5            string `json:"md5"`
    Ja3            string `json:"ja3"`
    First_reported string `json:"First_reported"`
    User_agent     string `json:"User-Agent"`
    Last_seen      string `json:"Last_seen"`
    Count          int    `json:"Count"`
}

var ja3hashes = []JA3Hash{}
var ja3Config JA3

func isJA3Enabled() (status bool) {
    return true
    // return ja3Config.Enable
}

func ja3Manage(bline []byte) (nBline []byte) {

    var hasJA3regex = regexp.MustCompile(`ja3`)
    if !hasJA3regex.MatchString(string(bline)) {
        return bline
    }

    jsoninterface := make(map[string]interface{})
    json.Unmarshal([]byte(bline), &jsoninterface)

    logs.Info("JA3 - Line has JA3 Values -> %+v", jsoninterface)

    return bline
}

func loadJA3Config() (err error) {
    ja3CFG, err := utils.GetKeyValueString("analyzer", "ja3conf")
    if err != nil {
        logs.Error("JA3 - Can't open configuration file path from main.conf: " + err.Error())
        return err
    }

    byteValue, err := ioutil.ReadFile(ja3CFG)
    if err != nil {
        logs.Error("JA3 - Error opening analyzer JA3 CFG: " + err.Error())
        return err
    }

    err = json.Unmarshal(byteValue, &ja3Config)
    if err != nil {
        logs.Error(err.Error())
        return err
    }

    return nil
}

func loadJA3Hashes() (err error) {
    ja3LocalFile := ja3Config.Localja3file

    byteValue, err := ioutil.ReadFile(ja3LocalFile)
    if err != nil {
        logs.Error("JA3 - Error opening local ja3 hashes file: %s", err.Error())
        return err
    }

    err = json.Unmarshal(byteValue, &ja3hashes)
    if err != nil {
        logs.Error("JA3 - Error importing ja3 hashes from file to struct: %s ", err.Error())
        return err
    }
    return nil

}

func init() {
    return
}
