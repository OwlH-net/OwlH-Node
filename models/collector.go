package models 

import (
    "owlhnode/collector"
    // "owlhnode/changeControl"
)

func PlayCollector() (err error) {
    err = collector.PlayCollector()
    //changecontrol.ChangeControlInsertData(err, "PlayCollector")  
    return err
}
func StopCollector() (err error) {
    err = collector.StopCollector()
    //changecontrol.ChangeControlInsertData(err, "StopCollector")  
    return err
}
func ShowCollector() (data string, err error) {
    data, err = collector.ShowCollector()
    //changecontrol.ChangeControlInsertData(err, "ShowCollector")  
    return data, err
}