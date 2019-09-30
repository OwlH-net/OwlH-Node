package models

import (
    "owlhnode/analyzer"
)

func PingAnalyzer()(data string ,err error) {
	data, err = analyzer.PingAnalyzer()	
	return data, err
}

func ChangeAnalyzerStatus(uuid map[string]string) (err error) {
	err = analyzer.ChangeAnalyzerStatus(uuid)
	return err
}