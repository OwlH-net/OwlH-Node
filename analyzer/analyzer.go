package analyzer

import (
    "encoding/json"
    "os"
    "syscall"
    "io/ioutil"
    "strings"
    "github.com/hpcloud/tail"
    "github.com/astaxie/beego/logs"
    "bufio"
    "time"
    "strconv"
    "owlhnode/utils"
    "owlhnode/database"
    "owlhnode/geolocation"
    // "regexp"
)

type chcounter struct {
    CHfeed          int
    CHmapper        int
    CHprefilter     int
    CHpostfilter    int
    CHtag           int
    CHwriter        int
    CHstartpipeline int
    lines           int
}



type Analyzer struct {
    Enable          bool        `json:"enable"`
    OutputFile      string      `json:"outputfile"`
    Prefilter       string      `json:"prefilterfile"`
    Postfilter      string      `json:"postfilterfile"`
    Tagfile         string      `json:"tags"`

    Srcfiles        []string    `json:"srcfiles"`
    Feedfiles       []Feedfile
}

type Feedfile struct {
    File            string      `json:"feedfile"`
    Workers         int         `json:"workers"`
}


type Tags struct{
    Tags            []Tag       `json:"tags"`
}

type Tag struct {
    Tagname         string      `json:"tagname"`
    Type            string      `json:"type"`
    Fields          []string    `json:"fields"`
    Exp             string      `json:"exp"`
    Action          string      `json:"action"`
    Stop            bool        `json:"stop"`
    Tag             string      `json:"tag"`
    Field           string      `json:"field"`
    Value           string      `json:"value"`
}

type Filters struct{
    Filters            []Filter       `json:"filters"`
}

type Filter struct {
    Filtername      string      `json:"filtername"`
    Type            string      `json:"type"`
    Fields          []string    `json:"fields"`
    Exp             string      `json:"exp"`
    Stop            bool        `json:"stop"`
    Tag             string      `json:"tag"`
    Action          string      `json:"action"`
    Fieldname       string      `json:"field"`
    Fieldvalue      string      `json:"value"`
}

type Event struct {
    Source          string
    Line            string
}

var monitorfiles            = map[string]bool{}

var CHstartpipeline         = make (chan string, 10000)
var CHprefilter             = make (chan string, 10000)
var CHmapper                = make (chan string, 10000)
var CHtag                   = make (chan string, 10000)
var CHfeed                  = make (chan string, 10000)
var CHpostfilter            = make (chan string, 10000)
var CHwriter                = make (chan string, 10000)
var CHdispatcher            = make (chan Event, 10000)


var config Analyzer
var tags Tags
var postfilters Filters
var prefilters Filters
var IoCs = map[string][]string{}
var counters chcounter

func readconf()(err error) {
    analyzerCFG, err := utils.GetKeyValueString("analyzer", "analyzerconf")
    if err != nil {logs.Error("AlertLog Error getting data from main.conf: "+err.Error()); return err}

    confFile, err := os.Open(analyzerCFG)
    if err != nil {logs.Error("Error openning analyzer CFG: "+err.Error()); return err}
    defer confFile.Close()
    byteValue, _ := ioutil.ReadAll(confFile)
    err = json.Unmarshal(byteValue, &config)
    if err != nil {logs.Error(err.Error()); return err}
    return nil
}

func readtags() {

    tagFile, err := os.Open(config.Tagfile)
    if err != nil {logs.Error("Error openning analyzer tag file: "+err.Error()); return}
    
    defer tagFile.Close()

    byteValue, _ := ioutil.ReadAll(tagFile)
    err = json.Unmarshal(byteValue, &tags)
    if err != nil {
        logs.Error("tags to json -> Unmarshal error: %s", err.Error())
    }
    logs.Info("tags Loaded")
}

func readpreexcludes() {

    preFile, err := os.Open(config.Prefilter)
    if err != nil {
        logs.Error("Error openning analyzer prefilters file: "+err.Error())
        return
    }
    defer preFile.Close()

    byteValue, _ := ioutil.ReadAll(preFile)
    err = json.Unmarshal(byteValue, &prefilters)
    if err != nil {
        logs.Error("pre filters to json -> Unmarshal error: %s", err.Error())
    }
    logs.Info("prefilters loaded")
}

func readpostexcludes() {
    postFile, err := os.Open(config.Postfilter)
    if err != nil {
        logs.Error("Error openning analyzer postfilters file: "+err.Error())
        return
    }
    defer postFile.Close()

    byteValue, _ := ioutil.ReadAll(postFile)
    err = json.Unmarshal(byteValue, &postfilters)
    if err != nil {
        logs.Error("post filters to json -> Unmarshal error: %s", err.Error())
    }
    logs.Info("postfilters loaded")
}


func readLines(path string) ([]string, error) {
    var lines []string
    file, err := os.Open(path)
    if err != nil {
        return lines, err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func ToDispatcher(source, line string) {
    var event Event
    event.Source    =   source
    event.Line      =   line
    CHdispatcher <- event
}

func DoFeed (wkrid int){
    for {
        line := <- CHfeed
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)
        feeddone := false
        for iocmap := range IoCs {
            for ioc := range IoCs[iocmap] {
                if strings.Contains(line, IoCs[iocmap][ioc]) {
                    feeddone=true
                    break
                }
            }
            if feeddone {
                break
            }
        }
        bline, err := json.Marshal(jsoninterface)
        if err != nil {
            logs.Error("DoFeed - Marshal before ToDispatcher - Error: %s", err.Error())
        }
        ToDispatcher("CHfeed", string(bline))
    }
}

func DoMapper(wkrid int) {
    logs.Info("Mapper -> %d -> Started",wkrid)

    for {
        line := <- CHmapper 
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        renamefield(jsoninterface, "src_ip", "srcip")
        renamefield(jsoninterface, "src_port", "srcport")
        renamefield(jsoninterface, "dest_port", "dstport")
        renamefield(jsoninterface, "dest_ip", "dstip")
        renamefield(jsoninterface, "id.orig_h", "srcip")
        renamefield(jsoninterface, "id.orig_p", "srcport")
        renamefield(jsoninterface, "id.resp_h", "dstip")
        renamefield(jsoninterface, "id.resp_p", "dstport")
        geoinfo(jsoninterface, "srcip", "geolocation_src")
        geoinfo(jsoninterface, "dstip", "geolocation_dst")

        bline, err := json.Marshal(jsoninterface)
        if err != nil {}
        ToDispatcher("CHmapper",string(bline))
    }
}

func renamefield(vjson map[string]interface{}, oldfield, newfield string){
    _, ok := vjson[oldfield] 
    if ok {
        vjson[newfield] = vjson[oldfield]
        delete(vjson, oldfield)
    }
}

func addtag(vjson map[string]interface{}, tag string){
    var currenttag []interface{}
    _, ok := vjson["tag"] 
    if ok {
        for atag := range vjson["tag"].([]string) {
            currenttag = append(currenttag, atag)
        }
    }
    currenttag = append(currenttag, tag)
    vjson["tag"] = currenttag
}

func insertfield(vjson map[string]interface{}, field, value string){
    vjson[field] = value
}

func geoinfo(vjson map[string]interface{}, srcfield, dstfield string){
    _, ok := vjson[srcfield] 
    if ok {
        geodata := geolocation.GetGeoInfo(vjson[srcfield].(string))
        if len(geodata) != 0 {
            vjson[dstfield] = geodata
        }
    }
}

func fieldExists(vjson map[string]interface{}, srcfield string)(exists bool){
    _, ok := vjson[srcfield] 
    if ok {
        return true
    }
    return false
}

func DoPreFilter(wkr int){
    logs.Info("Prefilter -> %d -> Started",wkr)
    for {
        line := <- CHprefilter 
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        exclude := false

        for filter := range prefilters.Filters {
            switch prefilters.Filters[filter].Type {
            case "string":
                for field := range prefilters.Filters[filter].Fields {
                    if fieldExists(jsoninterface,prefilters.Filters[filter].Fields[field]){
                        if jsoninterface[prefilters.Filters[filter].Fields[field]] == prefilters.Filters[filter].Exp {
                            switch prefilters.Filters[filter].Action {
                            case "exclude":
                                exclude = true
                                break
                            default:
                            }
                        }
                    }
                }
            default:
            }
            if exclude {
                break
            }
        }

        if exclude {
            continue
        }

        bline, err := json.Marshal(jsoninterface)
        if err != nil {}
        ToDispatcher("CHprefilter", string(bline))
    }
}

func DoTag(wkr int){
    logs.Info("Tag -> %d -> Started",wkr)
    for {
        line := <- CHtag
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        istheend := false
        
        for tag := range tags.Tags{
            switch tags.Tags[tag].Type {
            case "string":
                for field := range tags.Tags[tag].Fields {
                    if fieldExists(jsoninterface,tags.Tags[tag].Fields[field]){
                        if jsoninterface[tags.Tags[tag].Fields[field]] == tags.Tags[tag].Exp {
                            switch tags.Tags[tag].Action {
                            case "add":
                                addtag(jsoninterface, tags.Tags[tag].Tag)
                            case "insert":
                                insertfield(jsoninterface, tags.Tags[tag].Field, tags.Tags[tag].Value)
                            default:
                            }
                            if tags.Tags[tag].Stop {
                                istheend = true
                                break
                            }
                        }
                    }
                }
            default:
            }
            if istheend {
                break
            }
            logs.Info("tags done")
        }

        if istheend {
            logs.Info("tags continue")
            continue
        }

        bline, err := json.Marshal(jsoninterface)
        if err != nil {}
        ToDispatcher("CHtag", string(bline))
    }
}

func DoPostFilter(wkr int){
    logs.Info("Postfilter -> %d -> Started",wkr)
    for {
        line := <- CHpostfilter 
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        exclude := false
        
        for filter := range postfilters.Filters {
            switch postfilters.Filters[filter].Type {
            case "string":
                for field := range postfilters.Filters[filter].Fields {
                    if fieldExists(jsoninterface,postfilters.Filters[filter].Fields[field]){
                        if jsoninterface[postfilters.Filters[filter].Fields[field]] == postfilters.Filters[filter].Exp {
                            switch postfilters.Filters[filter].Action {
                            case "exclude":
                                exclude = true
                                break
                            default:
                            }
                        }
                    }
                }
            default:
            }
            if exclude {
                break
            }
        }

        if exclude {
            continue
        }

        bline, err := json.Marshal(jsoninterface)
        if err != nil {}
        ToDispatcher("CHpostfilter", string(bline))
    }
}



func DoWriter(wkrid int) {
    // TODO - verify if analyzer.json has outfile. if not try to find on main.conf
    outputfile, err := utils.GetKeyValueString("node", "alertLog")
    if err != nil {
        logs.Error("AlertLog Error getting data from main.conf: " + err.Error())
        return
    }
    ofile, err := os.OpenFile(outputfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        logs.Error("Analyzer Writer: can't open output file: " + outputfile + " -> " + err.Error())
        return
    }
    logs.Info("WRITER -> Started -> " + outputfile)
    _, err = ofile.WriteString("started\n")
    defer ofile.Close()
    go MonitorFile(outputfile, 1, ofile)
    for {
        line := <- CHwriter 
        // logs.Error("WRITER -> writing line %s", line)
        _, err = ofile.WriteString(line+"\n")
        if err != nil {
            logs.Error("Analyzer Writer: can't write line to file: " + outputfile + " -> " + err.Error())
        }
    }
}


func DoDispatcher(x int) {
    logs.Info ("Dispatcher %d --> doing dispatcher stuff", x)
    for {
        line := <- CHdispatcher
        // logs.Info("dispatcher %d -> dispatch line -> source: %s", x, line.Source)
        switch line.Source {
        case "start":
            // logs.Warn("dispatcher %d -> %s -> to Prefilter", x, line.Source)
            counters.CHprefilter += 1
            CHprefilter <- line.Line
        case "CHprefilter":
            // logs.Warn("dispatcher %d -> %s -> to Mapper", x, line.Source)
            counters.CHmapper += 1
            CHmapper <- line.Line
        case "CHmapper":
            // logs.Warn("dispatcher %d -> %s -> to Feed", x, line.Source)
            counters.CHfeed += 1
            CHfeed <- line.Line
        case "CHfeed":
            // logs.Warn("dispatcher %d -> %s -> to Tag", x, line.Source)
            counters.CHtag += 1
            CHtag <- line.Line
        case "CHtag":
            // logs.Warn("dispatcher %d -> %s -> to Postfilter", x, line.Source)
            counters.CHpostfilter += 1
            CHpostfilter <- line.Line
        case "CHpostfilter":
            // logs.Warn("dispatcher %d -> %s -> to Writer", x, line.Source)
            counters.CHwriter += 1
            CHwriter <- line.Line
        default:
            logs.Error("Source %s: Have no idea what is next with this line %s", line.Source, line.Line)
        }
    }
}


func StartDispatcher(wkr int) {
    logs.Info("starting Dispatcher with %d workers", wkr)
    for x:=0; x < wkr; x++ {
        go DoDispatcher(x)
    }
}

func StartPreFilter(wkr int) {
    logs.Info("starting Prefilter with %d workers", wkr)
    for x:=0; x < wkr; x++ {
        go DoPreFilter(x)
    }
}

func StartMapper(wkr int) {
    logs.Info("starting Mapper with %d workers", wkr)
    for x:=0; x < wkr; x++ {
        go DoMapper(x)
    }
}

func StartTag(wkr int) {
    logs.Info("starting Tag with %d workers", wkr)
    for x:=0; x < wkr; x++ {
        go DoTag(x)
    }
}

func StartFeed(wkr int) {
    logs.Info("starting Feed with %d workers", wkr)
    LoadFeed()
    for x:=0; x < wkr; x++ {
        go DoFeed(x)
    }
}

func StartPostFilter(wkr int) {
    logs.Info("starting Postfilter with %d workers", wkr)
    for x:=0; x < wkr; x++ {
        go DoPostFilter(x)
    }
}

func StartWriter(wkr int) {
    for x:=0; x < wkr; x++ {
        go DoWriter(x)
    }
}

func ControlSource(file, uuid string) {
    logs.Info("start file %s control", file)
    filedet, err := os.Stat(file) 
    if os.IsNotExist(err) {
        logs.Error("file %s doesn't exists. we don't control files like this.")
        return
    }
    stat, _ := filedet.Sys().(*syscall.Stat_t)
    var previousinode uint64
    previousinode = stat.Ino
    logs.Info("file %s inode %d ", file, int(previousinode))

    t,err := utils.GetKeyValueString("loop", "ControlSource")
    if err != nil {logs.Error("Search Error: Cannot load inode information.")}
    tDuration, err := strconv.Atoi(t)
    for {
        time.Sleep(time.Second * time.Duration(tDuration)) 
        if fileinfo, err := os.Stat(file); !os.IsNotExist(err) {
            stat, _ := fileinfo.Sys().(*syscall.Stat_t)
            logs.Info("current inode %d vs %d", stat.Ino, previousinode) 
            if previousinode != stat.Ino {
                logs.Warn("file %s inode changed, restarting tail with new inode", file)
                monitorfiles[uuid] = false
                go StartSource(file)
                return
            }
        }
    }
}

func StartSource(file string) {
    logs.Info("Starting tail of source file: "+file)
    var seekv tail.SeekInfo
    seekv.Offset = 0
    seekv.Whence = os.SEEK_END
    uuid := utils.Generate()
    monitorfiles[uuid] = true
    t,err := utils.GetKeyValueString("loop", "StartSource")
    if err != nil {logs.Error("Search Error: Cannot load node information.")}
    tDuration, err := strconv.Atoi(t)
    for {
        logs.Info("tailing - %s", file)
        if _, err := os.Stat(file); os.IsNotExist(err) {
            time.Sleep(time.Second * time.Duration(tDuration)) 
            continue
        }
        t, err := tail.TailFile(file, tail.Config{Follow: true, Location: &seekv})
        if err != nil {
            logs.Error(">>>>>> Tail over file %s error: %s", file, err.Error())
        }
        go ControlSource(file, uuid)
        for line := range t.Lines {
            if !monitorfiles[uuid] {
                return
            }
            counters.lines += 1
            ToDispatcher("start",line.Text)
        }
        logs.Info("End tailing - %s", file)
    }
}

func LoadFeed() {
    logs.Info("loading Feed")
    for file := range config.Feedfiles {
        logs.Info(config.Feedfiles[file].File)
        IoCs[config.Feedfiles[file].File], _ = readLines(config.Feedfiles[file].File)
    }
}

func LoadSources() {
    logs.Info("loading sources")
    for file := range config.Srcfiles {
        go StartSource(config.Srcfiles[file])
    }
}

func CHstats(){
    logs.Info("Channels Status")
    logs.Info("***************")
    logs.Info("CHprefilter %d items",len(CHprefilter))
    logs.Info("CHmapper %d items",len(CHmapper))
    logs.Info("CHtag %d items",len(CHtag))
    logs.Info("CHfeed %d items",len(CHfeed))
    logs.Info("CHpostfilter %d items",len(CHpostfilter))
    logs.Info("CHwriter %d items",len(CHwriter))
    logs.Info("***************")
}

func CHcounter(){
    logs.Info("Channels counters")
    logs.Info("*****************")
    logs.Info("Lines %d readed",counters.lines)
    logs.Info("CHprefilter %d times",counters.CHprefilter)
    logs.Info("CHmapper %d times",counters.CHmapper)
    logs.Info("CHfeed %d times",counters.CHfeed)
    logs.Info("CHtag %d times",counters.CHtag)
    logs.Info("CHpostfilter %d times",counters.CHpostfilter)
    logs.Info("CHwriter %d times",counters.CHwriter)
    logs.Info("***************")

}

func MonitorFile(file string, size int, ofile *os.File) {
    logs.Info(" >>>>>>>>>>  start file %s monitor", file)
    t,err := utils.GetKeyValueString("loop", "MonitorFile")
    if err != nil {logs.Error("Search Error: Cannot load node information.")}
    tDuration, err := strconv.Atoi(t)
    for {
        filedet, err := os.Stat(file) 
        if os.IsNotExist(err) {
            return
        }
        fsize := filedet.Size()
        if fsize > int64(size*1073741824) {
            logs.Error(">>>>>>> File %s size if greater than %dG, rotating...", file, size)
            ofile.Truncate(0)
            ofile.Seek(0,0)
        } 
        time.Sleep(time.Second * time.Duration(tDuration)) 
    }
}

func CHcontrol(){
    t,err := utils.GetKeyValueString("loop", "CHcontrol")
    if err != nil {logs.Error("Search Error: Cannot load node information.")}
    tDuration, err := strconv.Atoi(t)
    for {
        time.Sleep(time.Second * time.Duration(tDuration)) 
        CHstats()
        CHcounter()
        logs.Warn("monitorfiles size is %d", len(monitorfiles))
    }
}

func InitAnalizer() {
    logs.Info("starting analyzer")
    analyzer,_ := PingAnalyzer()
    if analyzer["status"] == "Disabled"{
        logs.Info("Analyzer is Disabled - Nothing to do")
        return
    }
    readconf()
    readtags()
    readpostexcludes()
    readpreexcludes()
    StartWriter(1)
    StartMapper(4)
    StartFeed(4)
    StartTag(4)
    StartDispatcher(4)
    StartPreFilter(4)
    StartPostFilter(4)

    LoadSources()
    go CHcontrol()

    t,err := utils.GetKeyValueString("loop", "InitAnalizer")
    if err != nil {logs.Error("Search Error: Cannot load node information.")}
    tDuration, err := strconv.Atoi(t)
    for {
        analyzer,_ = PingAnalyzer()
        if analyzer["status"] == "Disabled"{
            break
        }
        time.Sleep(time.Second * time.Duration(tDuration)) 
    }
}

func Init(){
    go InitAnalizer()
}

func PingAnalyzer()(data map[string]string ,err error) {
    filePath, err := utils.GetKeyValueString("node", "alertLog")
    if err != nil {logs.Error("PingAnalyzer Error getting data from main.conf")}

    analyzerData := make(map[string]string)
    analyzerData["status"] = "Disabled"

    analyzerStatus,err := ndb.GetStatusAnalyzer()
    if err != nil { logs.Error("Error getting Analyzer data: "+err.Error()); return analyzerData,err}

    analyzerData["status"] = analyzerStatus
    analyzerData["path"] = filePath
    analyzerData["size"] = "0"

    fi, err := os.Stat(filePath)

    if err != nil { logs.Error("Can't access Analyzer ouput file data: "+err.Error()); return analyzerData,nil}
    size := fi.Size()

    analyzerData["size"] = strconv.FormatInt(size, 10)

    return analyzerData, nil
}

func ChangeAnalyzerStatus(anode map[string]string) (err error) {
    if anode["status"] != "Enabled" || anode["status"] != "Disabled" { logs.Error("ChangeAnalyzerStatus bad analyzer value spected for status"); return nil}
    err = ndb.UpdateAnalyzer("analyzer", "status", anode["status"])
    if err != nil { logs.Error("Error updating Analyzer status: "+err.Error()); return err}
    
    return nil
}

func SyncAnalyzer(file map[string][]byte) (err error) { 
    alertFile, err := utils.GetKeyValueString("analyzer", "analyzerconf")
    if err != nil {logs.Error("SyncAnalyzer Error getting data from main.conf")}

    err = utils.WriteNewDataOnFile(alertFile, file["data"])
    if err != nil { logs.Error("Analyzer/SyncAnalyzer Error updating Analyzer file: "+err.Error()); return err}
    return err
}