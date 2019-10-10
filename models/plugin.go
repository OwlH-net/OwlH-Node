package models

import (
	"owlhnode/plugin"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
	err = plugin.ChangeServiceStatus(anode)
    return err
}

// curl -X PUT \
//   https://52.47.197.22:50002/node/plugin/ChangeMainServiceStatus \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "service": "suricata",
//     "param": "status",
//     "status": "enabled"
// }
func ChangeMainServiceStatus(anode map[string]string)(err error) {
	err = plugin.ChangeMainServiceStatus(anode)
    return err
}

func DeleteService(anode map[string]string)(err error) {
	err = plugin.DeleteService(anode)
    return err
}

func AddPluginService(anode map[string]string) (err error) {
    err = plugin.AddPluginService(anode)
    return err
}

func SaveSuricataInterface(anode map[string]string)(err error) {
	err = plugin.SaveSuricataInterface(anode)
    return err
}

func DeployStapService(anode map[string]string)(err error) {
	err = plugin.DeployStapService(anode)
    return err
}

func StopStapService(anode map[string]string)(err error) {
	err = plugin.StopStapService(anode)
    return err
}

func ModifyStapValues(anode map[string]string)(err error) {
	err = plugin.ModifyStapValues(anode)
    return err
}