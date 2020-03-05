package collector

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "os/exec"
)

func PlayCollector()(err error) {   
    _, err = exec.Command("bash","-c","ls -la").Output()
    if err != nil{
        logs.Error("Error executing command in PlayCollector function: "+err.Error())
        return err    
    }
    return nil
}

func StopCollector()(err error) {   
    _, err = exec.Command("bash","-c","ls -la").Output()
    if err != nil{
        logs.Error("Error executing command in StopCollector function: "+err.Error())
        return err    
    }
    return nil
}

func ShowCollector() (data string, err error) {
    status, err := utils.GetKeyValueString("stapCollector", "status")
    if err != nil{logs.Error("Error loading stap collector data: "+err.Error()); return "",err}
    param, err := utils.GetKeyValueString("stapCollector", "param")
    if err != nil{logs.Error("Error loading stap collector data: "+err.Error()); return "",err}
    command, err := utils.GetKeyValueString("stapCollector", "command")
    if err != nil{logs.Error("Error loading stap collector data: "+err.Error()); return "",err}
    
    output, err := exec.Command(command, param, status).Output()
    if err != nil{logs.Error("Error executing command in ShowCollector function: "+err.Error()); return "",err}

    return string(output),nil
}