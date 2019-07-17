package models 

import (
    "owlhnode/dataflow"
)

func ChangeDataflowValues(anode map[string]string) (err error) {
    err = dataflow.ChangeDataflowValues(anode)
    return err
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
    data, err = dataflow.LoadDataflowValues()
    return data, err
}