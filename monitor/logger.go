package monitor

import (
	"github.com/astaxie/beego/logs"
	// "encoding/json"
	// "io/ioutil"
	"os"
	// "path/filepath"
	"owlhnode/utils"
	"owlhnode/database"
	"strconv"
	"time"
	"bufio"
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

func FileRotation()(){	
	for{
		var err error
		rotate, err := ndb.LoadMonitorFiles()

		if err != nil {logs.Error("FileRotation ERROR readding rotation files from DB: "+err.Error())}
		for x := range rotate {
			if rotate[x]["rotate"] == "Enabled"{

				//delete more than 7 days
				currentTime := time.Now()
				fileDateModified, err := os.Stat(rotate[x]["path"]); if err != nil {logs.Error("FileRotation ERROR Checking file modification date: "+err.Error())}
				modifiedtime := fileDateModified.ModTime()
				// minusDays,_ := strconv.Atoi(rotate[x]["maxDays"])
				// lastDays = currentTime.AddDate(0, 0, -minusDays)

				// //delete file if is older than max days
				// if modifiedtime.Format("2006-01-02") < lastDays.Format("2006-01-02"){
				// 	err = os.Remove(rotate[x]["path"]); if err != nil {logs.Error("FileRotation ERROR deleting older files than max days: "+err.Error())}
				// 	err = ndb.DeleteMonitorFile(x); if err != nil {logs.Error("FileRotation ERROR deleting older files than max days at database: "+err.Error())}
				// }else{
					// file, err := os.Open(rotate[x]["path"])
					file, err := os.OpenFile(rotate[x]["path"], os.O_RDWR, 0755); if err != nil {logs.Error("FileRotation ERROR readding file: "+err.Error())}
					defer file.Close()
	
					fileInfo, err := file.Stat()
		
					lines := 0
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						lines++
					}
		
					//CHECK MAX LINES
					fileLines,_ := strconv.Atoi(rotate[x]["maxLines"])
					fileSize,_ := strconv.ParseInt(rotate[x]["maxSize"], 10, 64)	
					if lines > fileLines{
						err = utils.BackupFullPath(rotate[x]["path"])
						if err != nil {logs.Error("FileRotation ERROR creating backup by maxLines: "+err.Error())}
						err = file.Truncate(0); if err != nil {logs.Error("FileRotation ERROR: "+err.Error())}
						_,err = file.Seek(0,0); if err != nil {logs.Error("FileRotation ERROR2: "+err.Error())}
					}
					//CHECK FILE SIZE			
					if fileInfo.Size() > fileSize{
						err = utils.BackupFullPath(rotate[x]["path"])
						if err != nil {logs.Error("FileRotation ERROR creating backup by maxSize: "+err.Error())}
						err = file.Truncate(0); if err != nil {logs.Error("FileRotation ERROR: "+err.Error())}
						_,err = file.Seek(0,0); if err != nil {logs.Error("FileRotation ERROR2: "+err.Error())}
					}
					//CHECK FILE MODIFICATION DATE				
					if currentTime.Format("2006-01-02") >  modifiedtime.Format("2006-01-02"){
						logs.Notice("DAILY")
						err = utils.BackupFullPath(rotate[x]["path"])
						if err != nil {logs.Error("FileRotation ERROR creating backup by maxSize: "+err.Error())}
						err = file.Truncate(0); if err != nil {logs.Error("FileRotation ERROR: "+err.Error())}
						_,err = file.Seek(0,0); if err != nil {logs.Error("FileRotation ERROR2: "+err.Error())}
					}
				// }				
			}
		}
		logs.Info("Monitor files rotated!")
        time.Sleep(time.Minute*1)
	}
}