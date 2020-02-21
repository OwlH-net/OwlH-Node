package monitor

import (
	"github.com/astaxie/beego/logs"
	"os"
	"owlhnode/utils"
	"owlhnode/database"
	"strconv"
	"time"
	"path/filepath"
	"bufio"
	"strings"
	// "io/ioutil"
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
			//Check if file exists.
			_, err := os.Stat(rotate[x]["path"]);
			if err==nil && rotate[x]["rotate"] == "Enabled"{
				currentTime := time.Now()
				fileDateModified, err := os.Stat(rotate[x]["path"]); if err != nil {logs.Error("FileRotation ERROR Checking file modification date: "+err.Error())}
				modifiedtime := fileDateModified.ModTime()

				//delete files older than maxDays
				err = filepath.Walk(filepath.Dir(rotate[x]["path"]),
					func(fileSearch string, info os.FileInfo, err error) error {
					if err != nil {return err}
					if !info.IsDir() {
						if strings.Contains(fileSearch, filepath.Base(rotate[x]["path"])+"-"){
							maxDays,_ := strconv.Atoi(rotate[x]["maxDays"])
							oldFile := info.ModTime().AddDate(0, 0, maxDays)
							if oldFile.Format("2006-01-02") < currentTime.Format("2006-01-02") {
								err = os.Remove(filepath.Dir(rotate[x]["path"])+"/"+info.Name())
								if err != nil {logs.Error("FileRotation Error deleting file: "+filepath.Dir(rotate[x]["path"])+"/"+info.Name())}								
							}
						}
					}
					return nil
				})
				if err != nil {logs.Error("FileRotation Error filepath walk finish: "+err.Error())}

				file, err := os.OpenFile(rotate[x]["path"], os.O_RDWR, 0755); if err != nil {logs.Error("FileRotation ERROR readding file: "+err.Error())}
				defer file.Close()
				fileInfo, err := file.Stat()
	
				//get number of lines for check maxLines
				lines := 0
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					lines++
				}

				fileLines,_ := strconv.Atoi(rotate[x]["maxLines"])
				fileSize,_ := strconv.ParseInt(rotate[x]["maxSize"], 10, 64)	
				if lines > fileLines{
					//CHECK MAX LINES
					err = utils.BackupFullPath(rotate[x]["path"]); if err != nil {logs.Error("FileRotation ERROR creating backup by maxLines: "+err.Error())}						
					err = file.Truncate(0); if err != nil {logs.Error("FileRotation ERROR: "+err.Error())}
					_,err = file.Seek(0,0); if err != nil {logs.Error("FileRotation ERROR2: "+err.Error())}
				}else if fileInfo.Size() > fileSize{
					//CHECK FILE SIZE			
					err = utils.BackupFullPath(rotate[x]["path"]); if err != nil {logs.Error("FileRotation ERROR creating backup by maxSize: "+err.Error())}
					err = file.Truncate(0); if err != nil {logs.Error("FileRotation ERROR: "+err.Error())}
					_,err = file.Seek(0,0); if err != nil {logs.Error("FileRotation ERROR2: "+err.Error())}
				}else if currentTime.Format("2006-01-02") >  modifiedtime.Format("2006-01-02"){
					//CHECK FILE MODIFICATION DATE				
					err = utils.BackupFullPath(rotate[x]["path"]); if err != nil {logs.Error("FileRotation ERROR creating backup by maxSize: "+err.Error())}
					err = file.Truncate(0); if err != nil {logs.Error("FileRotation ERROR: "+err.Error())}
					_,err = file.Seek(0,0); if err != nil {logs.Error("FileRotation ERROR2: "+err.Error())}
				}			
			}
		}
		logs.Info("Monitor files rotated!")
        time.Sleep(time.Minute*1)
	}
}

func EditRotation(anode map[string]string)(err error){
	err = ndb.UpdateMonitorFileValue(anode["file"], "path", anode["path"]); if err != nil {logs.Error("EditRotation monitor files edit path Error: "+err.Error()); return err}
	err = ndb.UpdateMonitorFileValue(anode["file"], "maxSize", anode["size"]); if err != nil {logs.Error("EditRotation monitor files edit maxSize Error: "+err.Error()); return err}
	err = ndb.UpdateMonitorFileValue(anode["file"], "maxLines", anode["lines"]); if err != nil {logs.Error("EditRotation monitor files edit maxLines Error: "+err.Error()); return err}
	err = ndb.UpdateMonitorFileValue(anode["file"], "maxDays", anode["days"]); if err != nil {logs.Error("EditRotation monitor files edit maxDays Error: "+err.Error()); return err}
	return nil
}