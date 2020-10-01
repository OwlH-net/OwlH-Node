package incidents

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
)

func GetIncidentsNode() (anode map[string]map[string]string, err error) {
    incidents, err := ndb.GetIncidentsNode()    
    if err != nil {logs.Error("GetIncidentsNode error getting incidents from database %s", err.Error()); return nil,err}
    return incidents,nil
}

func PutIncidentNode(anode map[string]string) (err error) {
    err = ndb.PutIncidentNode(anode["uuid"], anode["param"], anode["value"])    
    if err != nil {logs.Error("PutIncidentNode error putting incidents into database: %s", err.Error()); return err}
    return nil
}