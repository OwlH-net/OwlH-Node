package analyzer

import (
    "bufio"
    "encoding/json"
    "errors"
    "github.com/astaxie/beego/logs"
    "github.com/hpcloud/tail"
    "io/ioutil"
    "net"
    "os"
    "owlhnode/database"
    "owlhnode/geolocation"
    "owlhnode/utils"
    "strconv"
    "strings"
    "syscall"
    "time"
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
    Enable                 bool     `json:"enable"`
    Force                  bool     `json:"force"`
    Verbose                bool     `json:"verbose"`
    Stats                  bool     `json:"stats"`
    OutputFile             string   `json:"outputfile"`
    Prefilter              string   `json:"prefilterfile"`
    Postfilter             string   `json:"postfilterfile"`
    Tagfile                string   `json:"tagsfile"`
    Srcfiles               []string `json:"srcfiles"`
    Feedfiles              []Feedfile
    Suricatasocket         string `json:"suricatasocket"`
    SuricatasocketEnabled  bool   `json:"suricatasocketenabled"`
    WazuhSocketEnabled     bool   `json:"wazuhsocketenabled"`
    TimebetweenStatusCheck int    `json:"timebetweenstatuscheck"`
    Timetowaitforfile      int    `json:"timetowaitforfilek"`
    ChannelWorkers         int    `json:"channelworkers"`
}

var config Analyzer

type Feedfile struct {
    File    string `json:"feedfile"`
    Workers int    `json:"workers"`
}

type Tags struct {
    Tags []Tag `json:"tags"`
}

type Tag struct {
    Tagname string   `json:"tagname"`
    Type    string   `json:"type"`
    Fields  []string `json:"fields"`
    Exp     string   `json:"exp"`
    Action  string   `json:"action"`
    Stop    bool     `json:"stop"`
    Tag     string   `json:"tag"`
    Field   string   `json:"field"`
    Value   string   `json:"value"`
}

type Filters struct {
    Filters []Filter `json:"filters"`
}

type Filter struct {
    Filtername string   `json:"filtername"`
    Type       string   `json:"type"`
    Fields     []string `json:"fields"`
    Exp        string   `json:"exp"`
    Stop       bool     `json:"stop"`
    Tag        string   `json:"tag"`
    Action     string   `json:"action"`
    Fieldname  string   `json:"field"`
    Fieldvalue string   `json:"value"`
}

type Event struct {
    Source string
    Line   string
}

type monitfile struct {
    File   string
    Status bool
    Error  bool
}

var monitorfiles = map[string]monitfile{}

var CHstartpipeline = make(chan string, 1000000)
var CHprefilter = make(chan string, 1000000)
var CHmapper = make(chan string, 1000000)
var CHtag = make(chan string, 1000000)
var CHfeed = make(chan string, 1000000)
var CHpostfilter = make(chan string, 1000000)
var CHwriter = make(chan string, 1000000)
var CHdispatcher = make(chan Event, 1000000)

var tags Tags
var postfilters Filters
var prefilters Filters
var IoCs = map[string][]string{}
var counters, previous chcounter

func readconf() (err error) {
    analyzerCFG, err := utils.GetKeyValueString("analyzer", "analyzerconf")
    if err != nil {
        logs.Error("AlertLog Error getting data from main.conf: " + err.Error())
        return err
    }

    byteValue, err := ioutil.ReadFile(analyzerCFG)
    if err != nil {
        logs.Error("Error openning analyzer CFG: " + err.Error())
        return err
    }

    err = json.Unmarshal(byteValue, &config)
    if err != nil {
        logs.Error(err.Error())
        return err
    }

    return nil
}

func readtags() {

    tagFile, err := os.Open(config.Tagfile)
    if err != nil {
        logs.Error("Error openning analyzer tag file: " + err.Error())
        return
    }

    defer tagFile.Close()

    byteValue, _ := ioutil.ReadAll(tagFile)
    err = json.Unmarshal(byteValue, &tags)
    if err != nil {
        logs.Error("tags to json -> Unmarshal error: %s", err.Error())
    }
    if config.Verbose {
        logs.Info("tags Loaded")
    }
}

func readpreexcludes() {

    preFile, err := os.Open(config.Prefilter)
    if err != nil {
        logs.Error("Error openning analyzer prefilters file: " + err.Error())
        return
    }
    defer preFile.Close()

    byteValue, _ := ioutil.ReadAll(preFile)
    err = json.Unmarshal(byteValue, &prefilters)
    if err != nil {
        logs.Error("pre filters to json -> Unmarshal error: %s", err.Error())
    }
    if config.Verbose {
        logs.Info("prefilters loaded")
    }

}

func readpostexcludes() {
    postFile, err := os.Open(config.Postfilter)
    if err != nil {
        logs.Error("Error openning analyzer postfilters file: " + err.Error())
        return
    }
    defer postFile.Close()

    byteValue, _ := ioutil.ReadAll(postFile)
    err = json.Unmarshal(byteValue, &postfilters)
    if err != nil {
        logs.Error("post filters to json -> Unmarshal error: %s", err.Error())
    }

    if config.Verbose {
        logs.Info("postfilters loaded")
    }
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
    event.Source = source
    event.Line = line
    CHdispatcher <- event
}

func DoFeed(wkrid int) {
    for {
        line := <-CHfeed
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)
        feeddone := false
        for iocmap := range IoCs {
            for ioc := range IoCs[iocmap] {
                if strings.Contains(line, IoCs[iocmap][ioc]) {
                    feeddone = true
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
        if isJA3Enabled() {
            bline = ja3Manage(bline)
        }
        ToDispatcher("CHfeed", string(bline))
    }
}

func DoMapper(wkrid int) {

    for {
        line := <-CHmapper
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
        if err != nil {
        }
        ToDispatcher("CHmapper", string(bline))
    }
}

func renamefield(vjson map[string]interface{}, oldfield, newfield string) {
    _, ok := vjson[oldfield]
    if ok {
        vjson[newfield] = vjson[oldfield]
        delete(vjson, oldfield)
    }
}

func addtag(vjson map[string]interface{}, tag string) {
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

func insertfield(vjson map[string]interface{}, field, value string) {
    vjson[field] = value
}

func geoinfo(vjson map[string]interface{}, srcfield, dstfield string) {
    _, ok := vjson[srcfield]
    if ok {
        geodata := geolocation.GetGeoInfo(vjson[srcfield].(string))
        if len(geodata) != 0 {
            vjson[dstfield] = geodata
        }
    }
}

func fieldExists(vjson map[string]interface{}, srcfield string) (exists bool) {
    _, ok := vjson[srcfield]
    if ok {
        return true
    }
    return false
}

func DoPreFilter(wkr int) {
    for {
        line := <-CHprefilter
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        exclude := false

        for filter := range prefilters.Filters {
            switch prefilters.Filters[filter].Type {
            case "string":
                for field := range prefilters.Filters[filter].Fields {
                    if fieldExists(jsoninterface, prefilters.Filters[filter].Fields[field]) {
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
        if err != nil {
        }
        ToDispatcher("CHprefilter", string(bline))
    }
}

func DoTag(wkr int) {
    for {
        line := <-CHtag
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        istheend := false

        for tag := range tags.Tags {
            switch tags.Tags[tag].Type {
            case "string":
                for field := range tags.Tags[tag].Fields {
                    if fieldExists(jsoninterface, tags.Tags[tag].Fields[field]) {
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
        }

        if istheend {
            logs.Info("tags continue")
            continue
        }

        bline, err := json.Marshal(jsoninterface)
        if err != nil {
        }
        ToDispatcher("CHtag", string(bline))
    }
}

func DoPostFilter(wkr int) {
    for {
        line := <-CHpostfilter
        jsoninterface := make(map[string]interface{})
        json.Unmarshal([]byte(line), &jsoninterface)

        exclude := false

        for filter := range postfilters.Filters {
            switch postfilters.Filters[filter].Type {
            case "string":
                for field := range postfilters.Filters[filter].Fields {
                    if fieldExists(jsoninterface, postfilters.Filters[filter].Fields[field]) {
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
        if err != nil {
        }
        ToDispatcher("CHpostfilter", string(bline))
    }
}

func DoWriter(wkrid int) {
    outputfile := config.OutputFile
    ofile, err := os.OpenFile(outputfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        logs.Error("Analyzer Writer: can't open or create output file: " + outputfile + " -> " + err.Error())
        return
    }
    if config.Verbose {
        logs.Info("WRITER -> Started -> " + outputfile)
    }
    _, err = ofile.WriteString("started\n")
    defer ofile.Close()
    go MonitorFile(outputfile, 1, ofile)
    for {
        line := <-CHwriter
        // logs.Error("WRITER -> writing line %s", line)
        _, err = ofile.WriteString(line + "\n")
        if err != nil {
            logs.Error("Analyzer Writer: can't write line to file: " + outputfile + " -> " + err.Error())
        }
        dgramWriter(line)
    }
}

func DoDispatcher(x int) {
    for {
        line := <-CHdispatcher
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
    for x := 0; x < wkr; x++ {
        go DoDispatcher(x)
    }
}

func StartPreFilter(wkr int) {
    for x := 0; x < wkr; x++ {
        go DoPreFilter(x)
    }
}

func StartMapper(wkr int) {
    for x := 0; x < wkr; x++ {
        go DoMapper(x)
    }
}

func StartTag(wkr int) {
    for x := 0; x < wkr; x++ {
        go DoTag(x)
    }
}

func StartFeed(wkr int) {
    LoadFeed()
    for x := 0; x < wkr; x++ {
        go DoFeed(x)
    }
}

func StartPostFilter(wkr int) {
    for x := 0; x < wkr; x++ {
        go DoPostFilter(x)
    }
}

func StartWriter(wkr int) {
    for x := 0; x < wkr; x++ {
        go DoWriter(x)
    }
}

func ControlSource(file, uuid string) {

    if config.Verbose {
        logs.Info("start file %s control", file)
    }
    filedet, err := os.Stat(file)
    if os.IsNotExist(err) {
        logs.Error("file %s doesn't exists. we don't control files like this.", file)
        return
    }

    stat, _ := filedet.Sys().(*syscall.Stat_t)

    var previousinode uint64
    previousinode = stat.Ino
    if config.Verbose {
        logs.Info("file %32s == %40s inode %7d ", uuid, file, int(previousinode))
    }

    t, err := utils.GetKeyValueString("loop", "ControlSource")
    if err != nil {
        logs.Error("Search Error: Cannot load controlsource - loop time from main.conf.")
        return
    }
    tDuration, err := strconv.Atoi(t)

    for {
        analyzer, _ := PingAnalyzer()
        if analyzer["status"] == "Disabled" {
            logs.Info("Analyzer is Disabled - Nothing to do")
            return
        }
        time.Sleep(time.Second * time.Duration(tDuration))
        if fileinfo, err := os.Stat(file); !os.IsNotExist(err) {

            stat, _ := fileinfo.Sys().(*syscall.Stat_t)
            if config.Verbose {
                logs.Info("AN - %40s, %32s, inode c- %7d vs %7d", file, uuid, stat.Ino, previousinode)
            }
            if previousinode != stat.Ino {
                if config.Verbose {
                    logs.Warn("AN - file %s inode changed, restarting tail with new inode", file)
                }
                mfile := monitorfiles[uuid]
                mfile.Status = false
                monitorfiles[uuid] = mfile
                //go StartSource(file)
                return
            }
        }
    }
}

func IsFileMonitored(file string) (monitored bool) {
    for uuid := range monitorfiles {
        mfile := monitorfiles[uuid]
        if mfile.File == file && mfile.Status {
            return true
        }
    }
    return false
}

func IsFileInMonitoredFiles(file string) (uuid string, isInList bool) {
    for uuid := range monitorfiles {
        mfile := monitorfiles[uuid]
        if mfile.File == file {
            return uuid, true
        }
    }
    return "", false
}

func StartSource(file, uuid string) {

    logs.Debug("Starting tail of source file: " + file)
    var seekv tail.SeekInfo

    seekv.Offset = 0
    seekv.Whence = os.SEEK_END

    logs.Debug("tailing - %40s - [%32s]", file, uuid)
    if _, err := os.Stat(file); os.IsNotExist(err) {
        delete(monitorfiles, uuid)
        return
    }
    t, err := tail.TailFile(file, tail.Config{Follow: true, Poll: false, Location: &seekv})
    if err != nil {
        logs.Error(">>>>>> Tail over file %s error: %s", file, err.Error())
        return
    }

    for line := range t.Lines {
        if !monitorfiles[uuid].Status {
            //delete(monitorfiles, uuid)
            //go StartSource(file)
            if config.Verbose {
                logs.Info("file %s is not ready anymore, closing tail for uuid %s", file, uuid)
            }
            t.Stop()
            t.Cleanup()
            return
        }
        counters.lines += 1
        ToDispatcher("start", line.Text)
    }

    logs.Debug("End tailing - %s", file)
}

func LoadFeed() {
    if config.Verbose {
        logs.Info("loading Feed")
    }
    for file := range config.Feedfiles {
        logs.Info(config.Feedfiles[file].File)
        IoCs[config.Feedfiles[file].File], _ = readLines(config.Feedfiles[file].File)
    }
}

func LoadSources() {
    if config.Verbose {
        logs.Info("Init - loading sources")
    }
    for {
        analyzer, _ := PingAnalyzer()
        for file := range config.Srcfiles {
            if analyzer["status"] == "Disabled" {
                logs.Info("Analyzer is Disabled - Nothing to do")
                return
            }

            _, err := os.Stat(config.Srcfiles[file])
            if os.IsNotExist(err) {
                continue
            }

            if IsFileMonitored(config.Srcfiles[file]) {
                continue
            }

            uuid, isInMonitoredFiles := IsFileInMonitoredFiles(config.Srcfiles[file])
            if !isInMonitoredFiles {
                uuid = utils.Generate()
            }

            logs.Debug("AN - file %s is not being monitored, start monitoring", config.Srcfiles[file])

            var mfile monitfile
            mfile.Status = true
            mfile.File = config.Srcfiles[file]
            mfile.Error = false
            monitorfiles[uuid] = mfile
            logs.Debug("AN - file added to monitor list - %40s, %32s, %t", mfile.File, uuid, mfile.Status)

            go StartSource(config.Srcfiles[file], uuid)
            go ControlSource(config.Srcfiles[file], uuid)

        }
        time.Sleep(time.Second * time.Duration(config.TimebetweenStatusCheck))
    }
}

func CHmonitor() {
    if config.Verbose {
        logs.Info("AN - Control Monitor")
        logs.Info("AN - ***************")
    }
    for uuid := range monitorfiles {
        mfile := monitorfiles[uuid].File
        mstatus := monitorfiles[uuid].Status

        if config.Verbose {
            logs.Info("AN - Mon file - %40s, %32s, %t", mfile, uuid, mstatus)
        }
    }
}

func CHstats() {
    logs.Info("Channels Status")
    logs.Info("***************")
    logs.Info("CHprefilter %d items", len(CHprefilter))
    logs.Info("CHmapper %d items", len(CHmapper))
    logs.Info("CHtag %d items", len(CHtag))
    logs.Info("CHfeed %d items", len(CHfeed))
    logs.Info("CHpostfilter %d items", len(CHpostfilter))
    logs.Info("CHwriter %d items", len(CHwriter))
    logs.Info("***************")
}

func CHEPS(t int) {

    EPSPrefilter := (counters.CHprefilter - previous.CHprefilter) / t
    previous.CHprefilter = counters.CHprefilter

    EPSMapper := (counters.CHmapper - previous.CHmapper) / t
    previous.CHmapper = counters.CHmapper

    EPSTag := (counters.CHtag - previous.CHtag) / t
    previous.CHtag = counters.CHtag

    EPSFeed := (counters.CHfeed - previous.CHfeed) / t
    previous.CHfeed = counters.CHfeed

    EPSPostfilter := (counters.CHpostfilter - previous.CHpostfilter) / t
    previous.CHpostfilter = counters.CHpostfilter

    EPSWriter := (counters.CHwriter - previous.CHwriter) / t
    previous.CHwriter = counters.CHwriter

    logs.Info("Channels EPS")
    logs.Info("***************")
    logs.Info("CHprefilter %d eps", EPSPrefilter)
    logs.Info("CHmapper %d eps", EPSMapper)
    logs.Info("CHtag %d eps", EPSTag)
    logs.Info("CHfeed %d eps", EPSFeed)
    logs.Info("CHpostfilter %d eps", EPSPostfilter)
    logs.Info("CHwriter %d epss", EPSWriter)
    logs.Info("***************")
}

func CHcounter() {
    logs.Info("Channels counters")
    logs.Info("*****************")
    logs.Info("Lines %d readed", counters.lines)
    logs.Info("CHprefilter %d times", counters.CHprefilter)
    logs.Info("CHmapper %d times", counters.CHmapper)
    logs.Info("CHfeed %d times", counters.CHfeed)
    logs.Info("CHtag %d times", counters.CHtag)
    logs.Info("CHpostfilter %d times", counters.CHpostfilter)
    logs.Info("CHwriter %d times", counters.CHwriter)
    logs.Info("***************")

}

func MonitorFile(file string, size int, ofile *os.File) {
    if config.Verbose {
        logs.Info("AN -  start file %s monitor", file)
    }
    t, err := utils.GetKeyValueString("loop", "MonitorFile")
    if err != nil {
        logs.Error("Search Error: Cannot load node information.")
    }
    tDuration, err := strconv.Atoi(t)
    for {
        filedet, err := os.Stat(file)
        if os.IsNotExist(err) {
            return
        }
        fsize := filedet.Size()
        if fsize > int64(size*1073741824) {
            if config.Verbose {
                logs.Error("AN - File %s size if greater than %dG, rotating...", file, size)
            }
            ofile.Truncate(0)
            ofile.Seek(0, 0)
        }
        time.Sleep(time.Second * time.Duration(tDuration))
    }
}

func FilesControl() {
    if config.Verbose {
        logs.Warn("AN - monitored files slice size is %4d", len(monitorfiles))
    }
    for uuid := range monitorfiles {
        if !monitorfiles[uuid].Status {
            delete(monitorfiles, uuid)
        }
    }
    if config.Verbose {
        logs.Warn("AN - monitored files slice size after cleaning is %4d", len(monitorfiles))
    }

}

func CHcontrol() {
    t, err := utils.GetKeyValueString("loop", "CHcontrol")
    if err != nil {
        logs.Error("Search Error: Cannot load node information.")
    }
    tDuration, err := strconv.Atoi(t)
    for {
        analyzer, _ := PingAnalyzer()
        if analyzer["status"] == "Disabled" {
            logs.Info("Stop channel Control")
            return
        }
        time.Sleep(time.Second * time.Duration(tDuration))
        CHstats()
        CHcounter()
        CHEPS(tDuration)
        CHmonitor()
        FilesControl()
    }
}

func InitAnalizer() {
    logs.Info("starting analyzer")
    analyzer, _ := PingAnalyzer()
    if analyzer["status"] == "Disabled" {
        logs.Info("Analyzer is Disabled - Nothing to do")
        return
    }
    readtags()
    readpostexcludes()
    readpreexcludes()
    StartWriter(1)
    StartMapper(config.ChannelWorkers)
    StartFeed(config.ChannelWorkers)
    StartTag(config.ChannelWorkers)
    StartDispatcher(config.ChannelWorkers)
    StartPreFilter(config.ChannelWorkers)
    StartPostFilter(config.ChannelWorkers)

    go LoadSources()
    logs.Info("AN - AN - AN - config - Verbose - %t", config.Verbose)
    if config.Stats {
        go CHcontrol()
    }

    // t, err := utils.GetKeyValueString("loop", "InitAnalizer")
    // if err != nil {
    //     logs.Error("Search Error: Cannot load node information.")
    // }
    // tDuration, err := strconv.Atoi(t)
    for {
        analyzer, _ = PingAnalyzer()
        if analyzer["status"] == "Disabled" {
            break
        }
        time.Sleep(time.Second * time.Duration(config.TimebetweenStatusCheck))
    }
}

func Init() {
    readconf()
    //forced by main conf
    if config.Force && !config.Enable {
        logs.Info("AN - Forced by main.conf configuration. Analyzer is disabled")
        return
    }
    go InitAnalizer()
    go dgram()
    loadJA3Config()
    loadJA3Hashes()
}

func PingAnalyzer() (data map[string]string, err error) {

    //unmarshal analyzer.conf into data struct
    readconf()

    filePath := config.OutputFile

    analyzerData := make(map[string]string)
    analyzerData["status"] = "Disabled"

    analyzerStatus, err := ndb.GetStatusAnalyzer()
    if err != nil {
        logs.Error("Error getting Analyzer data from DB: %s, defaulting to: %s", err.Error(), config.Enable)
        if config.Enable {
            analyzerStatus = "Enabled"
        }
    }

    analyzerData["status"] = analyzerStatus
    analyzerData["path"] = config.OutputFile
    analyzerData["size"] = "0"

    if analyzerData["status"] == "Disabled" {
        return analyzerData, nil
    }

    fi, err := os.Stat(config.OutputFile)
    if err != nil {
        logs.Error("Can't access Analyzer output file %s Error: %s", filePath, err.Error())
        return analyzerData, err
    }
    size := fi.Size()

    analyzerData["size"] = strconv.FormatInt(size, 10)

    return analyzerData, nil
}

// Analyzer status changed from Master/UI
func ChangeAnalyzerStatus(anode map[string]string) (err error) {
    if config.Verbose {
        logs.Debug(anode)
    }
    if anode["status"] == "Enabled" || anode["status"] == "Disabled" {
        err = ndb.UpdateAnalyzer("analyzer", "status", anode["status"])
        if err != nil {
            logs.Error("Error updating Analyzer status: " + err.Error())
            return err
        }
        return nil
    } else {
        if config.Verbose {
            logs.Error("ChangeAnalyzerStatus bad analyzer value expected for status")
        }
        return errors.New("ChangeAnalyzerStatus bad analyzer value expected for status")
    }

}

// Analyzer configuration file written from Master/UI.
func SyncAnalyzer(file map[string][]byte) (err error) {
    alertFile, err := utils.GetKeyValueString("analyzer", "analyzerconf")
    if err != nil {
        logs.Error("SyncAnalyzer Error getting data from main.conf")
    }

    err = utils.WriteNewDataOnFile(alertFile, file["data"])
    if err != nil {
        logs.Error("Analyzer/SyncAnalyzer Error updating Analyzer file: " + err.Error())
        return err
    }
    return err
}

func dgram() {
    if !config.SuricatasocketEnabled {
        return
    }
    socketPath := config.Suricatasocket

    // unlink it before doing anything
    syscall.Unlink(socketPath)

    // resolve unix address
    laddr, err := net.ResolveUnixAddr("unixgram", socketPath)
    if err != nil {
        logs.Error("Could not resolve unix socket: " + err.Error())
        return
    }

    // listen on the socket
    conn, err := net.ListenUnixgram("unixgram", laddr)
    if err != nil {
        logs.Error("Could not listen on unix socket datagram: " + err.Error())
        return
    }
    // close socket when we finish
    defer conn.Close()

    // scan text
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        //logs.Info("line is - %s", string(scanner.Bytes()))
        ToDispatcher("start", string(scanner.Bytes()))
    }
}

func dgramWriter(line string) {
    if !config.WazuhSocketEnabled {
        return
    }
    c, err := net.Dial("unixgram", "/var/ossec/queue/ossec/queue")
    if err != nil {
        logs.Error("Dial error", err)
        return
    }
    defer c.Close()

    msg := "1:/var/log/owlh/alerts.json/:" + line
    _, err = c.Write([]byte(msg))
    if err != nil {
        logs.Error("Write error:", err)
        return
    }
    logs.Info("Client sent:", msg)
}
