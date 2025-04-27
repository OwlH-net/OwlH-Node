package collector

import (
	"os/exec"

	"github.com/OwlH-net/OwlH-Node/utils"
	"github.com/astaxie/beego/logs"
)

func PlayCollector() (err error) {
	command, err := utils.GetKeyValueString("execute", "command")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return err
	}
	param, err := utils.GetKeyValueString("execute", "param")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return err
	}
	list, err := utils.GetKeyValueString("execute", "list")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return err
	}

	_, err = exec.Command(command, param, list).Output()
	if err != nil {
		logs.Error("Error executing command in PlayCollector function: " + err.Error())
		return err
	}
	return nil
}

func StopCollector() (err error) {
	command, err := utils.GetKeyValueString("execute", "command")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return err
	}
	param, err := utils.GetKeyValueString("execute", "param")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return err
	}
	list, err := utils.GetKeyValueString("execute", "list")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return err
	}

	_, err = exec.Command(command, param, list).Output()
	if err != nil {
		logs.Error("Error executing command in StopCollector function: " + err.Error())
		return err
	}
	return nil
}

func ShowCollector() (data string, err error) {
	status, err := utils.GetKeyValueString("stapCollector", "status")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return "", err
	}
	param, err := utils.GetKeyValueString("stapCollector", "param")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return "", err
	}
	command, err := utils.GetKeyValueString("stapCollector", "command")
	if err != nil {
		logs.Error("Error loading stap collector data: " + err.Error())
		return "", err
	}

	output, err := exec.Command(command, param, status).Output()
	if err != nil {
		logs.Error("Error executing command in ShowCollector function: " + err.Error())
		return "", err
	}

	return string(output), nil
}
