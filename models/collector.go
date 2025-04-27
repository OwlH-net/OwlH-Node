package models

import (
	"github.com/OwlH-net/OwlH-Node/collector"
	// "github.com/OwlH-net/OwlH-Node/changeControl"
)

func PlayCollector(username string) (err error) {
	err = collector.PlayCollector()
	//changecontrol.ChangeControlInsertData(err, "PlayCollector", username)
	return err
}
func StopCollector(username string) (err error) {
	err = collector.StopCollector()
	//changecontrol.ChangeControlInsertData(err, "StopCollector", username)
	return err
}
func ShowCollector(username string) (data string, err error) {
	data, err = collector.ShowCollector()
	//changecontrol.ChangeControlInsertData(err, "ShowCollector", username)
	return data, err
}
