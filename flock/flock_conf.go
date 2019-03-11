package flock
import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

var conf map[string]string
configFile := "/etc/owlh/owlh.conf"
LoadConf()


func LoadConf()(data map[string]string, err error){
	body, _ := ioutil.ReadAll(configFile)
    err = json.Unmarshal(body, &conf)
    if err != nil {
        return nil,err
	}
	return conf,nil
}

func GetItem(item string)(data map[string]string, err error){
	return conf[item], nil
}

func PrintConf()(){
	for item := range conf{
		logs.Info(item+" --> "+conf[item])
	}
}

