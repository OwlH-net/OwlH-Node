package models 

import (
    "owlhnode/collector"
)

func PlayCollector() (err error) {
    err = collector.PlayCollector()
    return err
}
func StopCollector() (err error) {
	err = collector.StopCollector()
    return err
}
func ShowCollector() (data string, err error) {
    data, err = collector.ShowCollector()
    return data, err
}