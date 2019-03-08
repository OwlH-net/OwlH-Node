package stap

import (
    // "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
  	//"owlhnode/utils"
  	//"owlhnode/database"
	  // "io/ioutil"
	  //"errors"
	  //"encoding/json"
)

func run()(status bool) {
    // get Servers list
    // for each server run a controller or collector
    // keep running until stop is requested. Stop all servers collection before end.
    return true
}

func collector(server string)(status bool) {
    return true
}

func controller(server string)(status bool) {
    //check if server stap is enabled in our config
    //check if server is reacheble
    //    if disabled - be sure the stap if off
    //    check if server status (CPU, RAM, STORAGE) is ok
    //    check stap status - stop or start as needed.
    //collect pcap files
    return true
}
