package monitor

import (
    "bufio"
    "github.com/astaxie/beego/logs"
    "os"
    "owlhnode/database"
    "owlhnode/utils"
    "path/filepath"
    "strconv"
    "strings"
    // "syscall"
    "time"
    // "io/ioutil"
)

type Targetfile struct {
    lastrotated time.Time
}

var Targets = map[string]Targetfile{}

func Logger() {
    var err error
    //get logger parameters
    filepath, err := utils.GetKeyValueString("logs", "filepath")
    if err != nil {
        logs.Error("Main Error getting data from main.conf for load Logger data: " + err.Error())
    }
    filename, err := utils.GetKeyValueString("logs", "filename")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    maxlines, err := utils.GetKeyValueString("logs", "maxlines")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    maxsize, err := utils.GetKeyValueString("logs", "maxsize")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    daily, err := utils.GetKeyValueString("logs", "daily")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    maxdays, err := utils.GetKeyValueString("logs", "maxdays")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    rotate, err := utils.GetKeyValueString("logs", "rotate")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    level, err := utils.GetKeyValueString("logs", "level")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }
    maxfiles, err := utils.GetKeyValueInt("logs", "maxfiles")
    if err != nil {
        logs.Error("Error getting data from main.conf for load Logger data: " + err.Error())
    }

    //transform maxsize to bytes
    newMaxSize,_ := utils.GetBytesFromSizeType(maxsize)
    pattern := "owlhnode-api[.]\\d{4}[-]\\d{2}[-]\\d{2}[.]\\d{3}.log"
    err = utils.ClearOlderLogFiles(filepath, filename+"." , maxfiles, pattern)
    if err != nil {logs.Error(err.Error())}
    
    logs.NewLogger(10000)
    logs.SetLogger(logs.AdapterFile, `{"filename":"`+filepath+filename+`", "maxlines":`+maxlines+` ,"maxsize":`+newMaxSize+`, "daily":`+daily+`, "maxdays":`+maxdays+`, "rotate":`+rotate+`, "level":`+level+`}`)

}

func FileRotation() {
    for {
        var err error
        rotate, err := ndb.LoadMonitorFiles()
        if err != nil {
            logs.Error("FileRotation ERROR readding rotation files from DB: " + err.Error())
        }

        for x := range rotate {
            //transform maxsize to bytes
            newMaxSize,_ := utils.GetBytesFromSizeType(rotate[x]["maxSize"])
            
            fileNumber, _ := strconv.Atoi(rotate[x]["maxFiles"])
            pattern := filepath.Base(rotate[x]["path"])+"-\\d{10}"
            err = utils.ClearOlderLogFiles(filepath.Dir(rotate[x]["path"])+"/", filepath.Base(rotate[x]["path"]), fileNumber, pattern)
            if err != nil {logs.Error(err.Error())}

            //Check if file exists.
            _, err := os.Stat(rotate[x]["path"])
            if err == nil && rotate[x]["rotate"] == "Enabled" {
                currentTime := time.Now()
                if _, ok := Targets[rotate[x]["path"]]; !ok {
                    targetLast := Targets[rotate[x]["path"]]
                    targetLast.lastrotated = currentTime
                    Targets[rotate[x]["path"]] = targetLast
                }
                // fileDateModified, err := os.Stat(rotate[x]["path"])
                if err != nil {
                    logs.Error("FileRotation ERROR Checking file modification date: " + err.Error())
                }
                // modifiedtime := fileDateModified.ModTime()
                // stat := fileDateModified.Sys().(*syscall.Stat_t)
                // ctime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
                err = filepath.Walk(filepath.Dir(rotate[x]["path"]),
                    func(fileSearch string, info os.FileInfo, err error) error {
                        if err != nil {
                            return err
                        }
                        if !info.IsDir() {
                            if strings.Contains(fileSearch, filepath.Base(rotate[x]["path"])+"-") {

                                maxFiles, _ := strconv.Atoi(rotate[x]["maxFiles"])
                                oldFile := info.ModTime().AddDate(0, 0, maxFiles)
                                if oldFile.Format("2006-01-02") < currentTime.Format("2006-01-02") {
                                    err = os.Remove(filepath.Dir(rotate[x]["path"]) + "/" + info.Name())
                                    if err != nil {
                                        logs.Error("FileRotation Error deleting file: " + filepath.Dir(rotate[x]["path"]) + "/" + info.Name())
                                    }
                                }
                            }
                        }
                        return nil
                    })
                if err != nil {
                    logs.Error("FileRotation Error filepath walk finish: " + err.Error())
                }

                file, err := os.OpenFile(rotate[x]["path"], os.O_RDWR, 0644)
                if err != nil {
                    logs.Error("FileRotation ERROR readding file: " + err.Error())
                }
                defer file.Close()
                fileInfo, err := file.Stat()

                //get number of lines for check maxLines
                lines := 0
                scanner := bufio.NewScanner(file)
                for scanner.Scan() {
                    lines++
                }

                fileLines, _ := strconv.Atoi(rotate[x]["maxLines"])
                fileSize, _ := strconv.ParseInt(newMaxSize, 10, 64)
                if lines > fileLines {
                    //CHECK MAX LINES
                    err = utils.BackupFullPath(rotate[x]["path"])
                    if err != nil {
                        logs.Error("FileRotation ERROR creating backup by maxLines: " + err.Error())
                    }
                    err = file.Truncate(0)
                    if err != nil {
                        logs.Error("FileRotation ERROR: " + err.Error())
                    }
                    _, err = file.Seek(0, 0)
                    if err != nil {
                        logs.Error("FileRotation ERROR2: " + err.Error())
                    }
                } else if fileInfo.Size() > fileSize {
                    //CHECK FILE SIZE
                    err = utils.BackupFullPath(rotate[x]["path"])
                    if err != nil {
                        logs.Error("FileRotation ERROR creating backup by maxSize: " + err.Error())
                    }
                    err = file.Truncate(0)
                    if err != nil {
                        logs.Error("FileRotation ERROR: " + err.Error())
                    }
                    _, err = file.Seek(0, 0)
                    if err != nil {
                        logs.Error("FileRotation ERROR2: " + err.Error())
                    }
                    // } else if currentTime.For<mat("2006-01-02") > modifiedtime.Format("2006-01-02") {
                } else if Targets[rotate[x]["path"]].lastrotated.Format("2006-01-02") > currentTime.Format("2006-01-02") {
                    //CHECK FILE MODIFICATION DATE
                    err = utils.BackupFullPath(rotate[x]["path"])
                    if err != nil {
                        logs.Error("FileRotation ERROR creating backup by maxSize: " + err.Error())
                    }
                    err = file.Truncate(0)
                    if err != nil {
                        logs.Error("FileRotation ERROR: " + err.Error())
                    }
                    _, err = file.Seek(0, 0)
                    if err != nil {
                        logs.Error("FileRotation ERROR2: " + err.Error())
                    }
                    targetlast := Targets[rotate[x]["path"]]
                    targetlast.lastrotated = currentTime
                    Targets[rotate[x]["path"]] = targetlast
                }
            }
        }
        logs.Info("Monitor files rotated!")
        t, err := utils.GetKeyValueString("loop", "FileRotation")
        if err != nil {
            logs.Error("Search Error: Cannot load node information.")
        }
        tDuration, err := strconv.Atoi(t)
        time.Sleep(time.Minute * time.Duration(tDuration))
    }
}

func EditRotation(anode map[string]string) (err error) {
    logs.Debug(anode)
    err = ndb.UpdateMonitorFileValue(anode["file"], "path", anode["path"])
    if err != nil {
        logs.Error("EditRotation monitor files edit path Error: " + err.Error())
        return err
    }
    err = ndb.UpdateMonitorFileValue(anode["file"], "maxSize", anode["size"])
    if err != nil {
        logs.Error("EditRotation monitor files edit maxSize Error: " + err.Error())
        return err
    }
    err = ndb.UpdateMonitorFileValue(anode["file"], "maxLines", anode["lines"])
    if err != nil {
        logs.Error("EditRotation monitor files edit maxLines Error: " + err.Error())
        return err
    }
    err = ndb.UpdateMonitorFileValue(anode["file"], "maxFiles", anode["files"])
    if err != nil {
        logs.Error("EditRotation monitor files edit maxFiles Error: " + err.Error())
        return err
    }
    err = ndb.UpdateMonitorFileValue(anode["file"], "maxDays", anode["days"])
    if err != nil {
        logs.Error("EditRotation monitor files edit maxDays Error: " + err.Error())
        return err
    }
    return nil
}
