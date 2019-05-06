package collector

import (
	"github.com/astaxie/beego/logs"
	"owlhnode/utils"
	"os/exec"
)

func PlayCollector()(err error) {   
	// stapCollector := map[string]map[string]string{}
    // stapCollector["stapCollector"] = map[string]string{}
    // stapCollector["stapCollector"]["start"] = ""
    // stapCollector["stapCollector"]["param"] = ""
    // stapCollector["stapCollector"]["command"] = ""
	// stapCollector,err = utils.GetConf(stapCollector)
	// start := stapCollector["stapCollector"]["start"]
	// param := stapCollector["stapCollector"]["param"]
	// command := stapCollector["stapCollector"]["command"]

	// _, err = exec.Command(command, param, start).Output()
    // if err != nil{
	// 	logs.Error("Error executing command in PlayCollector function: "+err.Error())
    //     return err    
	// }
	_, err = exec.Command("bash","-c","ls -la").Output()
    if err != nil{
		logs.Error("Error executing command in PlayCollector function: "+err.Error())
        return err    
	}
	return nil
}

func StopCollector()(err error) {   
	// stapCollector := map[string]map[string]string{}
    // stapCollector["stapCollector"] = map[string]string{}
    // stapCollector["stapCollector"]["stop"] = ""
    // stapCollector["stapCollector"]["param"] = ""
    // stapCollector["stapCollector"]["command"] = ""
	// stapCollector,err = utils.GetConf(stapCollector)
	// stop := stapCollector["stapCollector"]["stop"]
	// param := stapCollector["stapCollector"]["param"]
	// command := stapCollector["stapCollector"]["command"]

	// _, err = exec.Command(command, param, stop).Output()
    // if err != nil{
	// 	logs.Error("Error executing command in StopCollector function: "+err.Error())
    //     return err    
	// }
	_, err = exec.Command("bash","-c","ls -la").Output()
    if err != nil{
		logs.Error("Error executing command in StopCollector function: "+err.Error())
        return err    
	}
	return nil
}

func ShowCollector() (data string, err error) {
	stapCollector := map[string]map[string]string{}
    stapCollector["stapCollector"] = map[string]string{}
    stapCollector["stapCollector"]["status"] = ""
    stapCollector["stapCollector"]["param"] = ""
    stapCollector["stapCollector"]["command"] = ""
	stapCollector,err = utils.GetConf(stapCollector)
	status := stapCollector["stapCollector"]["status"]
	param := stapCollector["stapCollector"]["param"]
	command := stapCollector["stapCollector"]["command"]

	output, err := exec.Command(command, param, status).Output()
	logs.Debug(output)
	logs.Debug(output)
	logs.Debug(output)
	logs.Debug(output)
	logs.Debug(output)
	logs.Debug(output)
	logs.Debug(output)
    if err != nil{
		logs.Error("Error executing command in ShowCollector function: "+err.Error())
        return "",err    
	}
	return string(output),nil
}