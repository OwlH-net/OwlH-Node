package stap

import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
	"owlhnode/utils"
	"owlhnode/database"
	// "io/ioutil"
	"errors"
)

func AddServer(elem map[string]string) (err error){

	jsonDefaultData,err := utils.LoadDefaultServerData("defaults.json")
	logs.Info(jsonDefaultData)

	nodeName:= elem["nodeName"]
	nodeIP:= elem["nodeIP"]
	uuid:= elem["uuid"]

	uuidServer := utils.Generate()
	if ndb.Sdb == nil {
        logs.Error("AddServer stap -- Can't access to database")
        return errors.New("AddServer stap -- Can't access to database")
	} 
	logs.Info("stap/AddServer INSERT server-node")
	insertServer, err := ndb.Sdb.Prepare("insert into server_node (server_uniqueid, node_uniqueid) values (?,?);")
	_, err = insertServer.Exec(&uuidServer, &uuid)  
	defer insertServer.Close()
	if err != nil {
        return err
    }
		// for k,v in elem {
	logs.Info("stap/AddServer INSERT new server ip")
	insertServerIP, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
	_, err = insertServerIP.Exec(&uuidServer, "ip", &nodeIP)  
	defer insertServerIP.Close()
		// }
	if err != nil {
		return err
	}

	logs.Info("stap/AddServer INSERT new server name")
	insertServerName, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
	_, err = insertServerName.Exec(&uuidServer, "name", &nodeName)  
	defer insertServerName.Close()
    if err != nil {
        return err
	}
	
	return nil
}