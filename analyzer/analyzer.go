package analyzer

import (
    "encoding/json"
    "os"
    // "github.com/google/uuid"
    "io/ioutil"
    "strings"
	"github.com/hpcloud/tail"
	"github.com/astaxie/beego/logs"
	"bufio"
	"time"
	"strconv"
	"owlhnode/utils"
	"owlhnode/database"
)

type iocAlert struct {
    Data      Data     `json:"data"`
    Full_log  string   `json:"full_log"`
}

type Data struct {
    Dstport     string    `json:"dstport"`
    Srcport     string    `json:"srcport"`
    Dstip       string    `json:"dstip"`
    Srcip       string    `json:"srcip"`
    IoC         string    `json:"ioc"`
    IoCsource   string    `json:"iocsource"`
    Signature   Signature `json:"alert"`
}

type Signature struct {
    Signature       string `json:"signature"`
    Signature_id    string `json:"signature_id"`
}


type Analyzer struct {
    Enable 	bool 		`json:"enable"`
    Srcfiles 	[]string 	`json:"srcfiles"`
    Feedfiles 	[]Feedfile 
}

type Feedfile struct {
    File		string		`json:"feedfile"`
    Workers		int			`json:"workers"`
}

var dispatcher = make(map[string]chan string)

var config Analyzer
// var analyzerconf = "conf/analyzer.json"

func readconf()(err error) {

	cfg := map[string]map[string]string{}
	cfg["analyzer"] = map[string]string{}
	cfg["analyzer"]["analyzerconf"] = ""
	cfg,err = utils.GetConf(cfg)
	analyzerCFG := cfg["analyzer"]["analyzerconf"]
	if err != nil {
		logs.Error("AlertLog Error getting data from main.conf: "+err.Error())
		return
	}

    confFile, err := os.Open(analyzerCFG)
    if err != nil {
        logs.Error(err.Error())
        return err
    }
    defer confFile.Close()
    byteValue, _ := ioutil.ReadAll(confFile)
    err = json.Unmarshal(byteValue, &config)
    if err != nil {
        logs.Error(err.Error())
        return err
    }
    return nil
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func Registerchannel(uuid string) {
    dispatcher[uuid] = make(chan string)
}

func Domystuff(IoCs []string, uuid string, wkrid int, iocsrc string) {
    for {
        line := <- dispatcher[uuid] 
        for ioc := range IoCs {
            if strings.Contains(line, IoCs[ioc]) {
				logs.Info("Match -> "+ line +" IoC found -> " + IoCs[ioc]  + " wkrid -> " + strconv.Itoa(wkrid))
				//ioc
				IoCtoAlert(line, IoCs[ioc], iocsrc)
            }
        }
    }
}

func Startanalyzer(file string, wkr int) {
    newuuid := utils.Generate()
    logs.Info(newuuid + ": starting analyzer with feed: "+file + " with " + strconv.Itoa(wkr) + " workers")
    Registerchannel(newuuid)
    IoCs, _ := readLines(file)
    for x:=0; x < wkr; x++ {
        go Domystuff(IoCs, newuuid, x, file)
    }
}


func Starttail(file string) {
    t, _ := tail.TailFile(file, tail.Config{Follow: true})
    for line := range t.Lines {
        dispatch(line.Text)
    }
}

func LoadAnalyzers() {
	logs.Info("loading analyzers")
    for file := range config.Feedfiles {
        go Startanalyzer(config.Feedfiles[file].File, config.Feedfiles[file].Workers)
    }
}

func LoadSources() {
    logs.Info("loading sources")
    for file := range config.Srcfiles {
        go Starttail(config.Srcfiles[file])
    }
}

func dispatch(line string) {
    for channel := range dispatcher {
        dispatcher[channel] <- line
    }
}


func IoCtoAlert(line, ioc, iocsrc string) {
	var err error
	AlertLog := map[string]map[string]string{}
	AlertLog["Node"] = map[string]string{}
	AlertLog["Node"]["AlertLog"] = ""
	AlertLog,err = utils.GetConf(AlertLog)
	AlertLogJson := AlertLog["Node"]["AlertLog"]
	if err != nil {
		logs.Error("AlertLog Error getting data from main.conf: "+err.Error())
		return
	}

	alert     := iocAlert{}
	data      := Data{}
	signature := Signature{}

	signature.Signature = "OwlH IoC found - "+ioc
	signature.Signature_id = "8000101"

	// data.Dstport = dstport
	// data.Dstip = dstip
	// data.Srcip = srcip
	// data.Srcport = srcport
	data.Signature = signature
	data.IoC = ioc
	data.IoCsource = iocsrc

	alert.Data = data
	alert.Full_log = line
	alertOutput, _ := json.Marshal(alert)

	err = utils.WriteNewDataOnFile(AlertLogJson, alertOutput)
	if err != nil {
		logs.Error("Error saving data IoCtoAlert: %s", err.Error())
	}

}

func InitAnalizer() {
    logs.Info("starting analyzer")
    readconf()
    LoadAnalyzers()
    LoadSources()
    for {
		status,_ := PingAnalyzer()
		if status == "Disabled"{
			break
		}
		time.Sleep(time.Second * 60)
		//check if is active at DB
		//sleep 1 min
    }
}

func Init(){
	go InitAnalizer()
}

func PingAnalyzer()(data string ,err error) {
	analyzerData,err := ndb.GetStatusAnalyzer()
	if err != nil { logs.Error("Error getting Analyzer data: "+err.Error()); return "",err}

	return analyzerData	, nil
}

func ChangeAnalyzerStatus(anode map[string]string) (err error) {
	err = ndb.UpdateAnalyzer(anode["uuid"], "status", anode["status"])
	if err != nil { logs.Error("Error updating Analyzer status: "+err.Error()); return err}

	return nil
}