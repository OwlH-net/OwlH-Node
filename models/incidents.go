package models

import (
    "owlhnode/incidents"
//    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

// curl -X GET \
//   https://52.47.197.22:50002/node/incidents \
// }
func GetIncidentsNode()(data map[string]map[string]string ,err error) {
	data, err = incidents.GetIncidentsNode()	
	return data, err
}

// curl -X POST \
//   https://52.47.197.22:50002/node/incidents \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "nodeuuid": "d",
//     "uuid": "v",
//     "param": "v",
//     "value": "v",
// }
func PutIncidentNode(anode map[string]string)(err error){
    cc := anode
    logs.Info("============")
    logs.Info("INCIDENT - PutIncidentNode")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = incidents.PutIncidentNode(anode)
    return err
}