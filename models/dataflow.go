package models 

import (
    "owlhnode/dataflow"
)

func ChangeDataflowValues(anode map[string]string) (err error) {
    err = dataflow.ChangeDataflowValues(anode)
    return err
}

func SaveSocketToNetwork(anode map[string]string) (err error) {
	err = dataflow.SaveSocketToNetwork(anode)
	return err
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
    data, err = dataflow.LoadDataflowValues()
    return data, err
}

// func LoadSocketToNetwork() (data map[string]map[string]string, err error) {
//     err = dataflow.LoadSocketToNetwork(anode)
//     return err
// }