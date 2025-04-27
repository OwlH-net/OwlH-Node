package models

import (
	changecontrol "github.com/OwlH-net/OwlH-Node/changeControl"
)

//	curl -X GET \
//	  https://52.47.197.22:50002/node/changecontrol \
//	}
func GetChangeControlNode(username string) (data map[string]map[string]string, err error) {
	data, err = changecontrol.GetChangeControlNode()
	//changecontrol.ChangeControlInsertData(err, "GetChangeControlNode", username)
	return data, err
}
