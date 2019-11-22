package models

import (
    "owlhnode/changeControl"
)

// curl -X GET \
//   https://52.47.197.22:50002/node/changecontrol \
// }
func GetChangeControlNode()(data map[string]map[string]string ,err error) {
    data, err = changecontrol.GetChangeControlNode()  
    changecontrol.ChangeControlInsertData(err, "GetChangeControlNode")  
    return data, err
}