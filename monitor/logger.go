package monitor

import (
	"github.com/astaxie/beego/logs"
	// "encoding/json"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	// "owlhnode/database"
	"owlhnode/utils"
)


func Logger() {
	// data,err := ndb.LoadMonitorFiles()
	// if err != nil {logs.Error("Error getting monitor files for logger: "+err.Error())}
	var err error
	//get logger parameters
	loadDataLogger := map[string]map[string]string{}
	loadDataLogger["logs"] = map[string]string{}
	loadDataLogger["logs"]["filename"] = ""
	loadDataLogger["logs"]["maxlines"] = ""
	loadDataLogger["logs"]["maxsize"] = ""
	loadDataLogger["logs"]["daily"] = ""
	loadDataLogger["logs"]["maxdays"] = ""
	loadDataLogger["logs"]["rotate"] = ""
	loadDataLogger["logs"]["level"] = ""
	loadDataLogger, err = utils.GetConf(loadDataLogger)    
	filename := loadDataLogger["logs"]["filename"]
	maxlines := loadDataLogger["logs"]["maxlines"]
	maxsize := loadDataLogger["logs"]["maxsize"]
	daily := loadDataLogger["logs"]["daily"]
	maxdays := loadDataLogger["logs"]["maxdays"]
	rotate := loadDataLogger["logs"]["rotate"]
	level := loadDataLogger["logs"]["level"]
	if err != nil {logs.Error("Error getting data from main.conf for load Logger data: "+err.Error())}

	// //get monitor files
	// jsonPath, err := ioutil.ReadFile("conf/main.conf")
	// if err != nil {logs.Error("Main Error oppening Logger file: "+err.Error())}
	// logFiles := map[string]map[string]string{}
	// json.Unmarshal(jsonPath, &logFiles)

	// exists := false
	// for x,y := range logFiles{
	// 	if x == "monitorfile"{
	// 		for y,_ := range y{
	// 			for param,path := range data {
	// 				for path,_ := range path {
	// 					if data[param][path] == logFiles[x][y]{							
	// 						exists = true
	// 					}
	// 				}
	// 			}
	// 			if !exists {
	// 				//check if file exists
	// 				if _, err := os.Stat(logFiles[x][y]); os.IsNotExist(err) {
	// 					logs.Info("Creating path: "+filepath.Dir(logFiles[x][y]))
	// 					err = os.MkdirAll(filepath.Dir(logFiles[x][y]), os.ModePerm)
	// 					if err != nil {logs.Error("Main Error creating logger file path: "+err.Error())}
	// 					_, err := os.Create(logFiles[x][y])
	// 					logs.Info("Creating file: "+logFiles[x][y])
	// 					if err != nil {logs.Error("Main Error creating logger file: "+err.Error())}
	// 				}
	// 				//insert into db
	// 				uuid := utils.Generate()
	// 				err = ndb.InsertMonitorValue(uuid,"path", logFiles[x][y])
	// 				if err != nil {logs.Error("Main Error inserting logger file into database: "+err.Error())}
	// 			}
	// 			exists = false
	// 		}
	// 	}
	// }
	
	// data,err = ndb.LoadMonitorFiles()
	// if err != nil {logs.Error("Error getting monitor files for logger: "+err.Error())}
	// for id,path := range data {
	// 	for path,_ := range path {
	// 		logs.Notice(data[id][path])
			logs.NewLogger(10000)
			logs.SetLogger(logs.AdapterFile,`{"filename":"`+filename+`", "maxlines":`+maxlines+` ,"maxsize":`+maxsize+`, "daily":`+daily+`, "maxdays":`+maxdays+`, "rotate":`+rotate+`, "level":`+level+`}`)
			// logs.SetLogger(logs.AdapterFile,`{"filename":"`+data[id][path]+`", "maxlines":`+maxlines+` ,"maxsize":`+maxsize+`, "daily":`+daily+`, "maxdays":`+maxdays+`, "rotate":`+rotate+`, "level":`+level+`}`)
	// 	}
	// }
}

func FileRotation(path string)(err error)