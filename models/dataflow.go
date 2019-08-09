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

func SaveNewLocal(anode map[string]string)(err error) {
    err = dataflow.SaveNewLocal(anode)
    return err
}

func SaveVxLAN(anode map[string]string)(err error) {
    err = dataflow.SaveVxLAN(anode)
    return err
}

func SaveSocketToNetworkSelected(anode map[string]string)(err error) {
    err = dataflow.SaveSocketToNetworkSelected(anode)
    return err
}

func DeleteDataFlowValueSelected(anode map[string]string)(err error) {
    err = dataflow.DeleteDataFlowValueSelected(anode)
    return err
}