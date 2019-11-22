package analyzer

import (
    "encoding/json"
    "os"
    "io/ioutil"
    // "strings"
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


type iocAlert struct {
    Data            Data        `json:"data"`
    Full_log        string      `json:"full_log"`
}

type Data struct {
    Dstport         string      `json:"dstport"`
    Srcport         string      `json:"srcport"`
    Dstip           string      `json:"dstip"`
    Srcip           string      `json:"srcip"`
    IoC             string      `json:"ioc"`
    IoCsource       string      `json:"iocsource"`
    Signature       Signature   `json:"alert"`
}

type Signature struct {
    Signature       string      `json:"signature"`
    Signature_id    string      `json:"signature_id"`
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

var CHstartpipeline         = make (chan string)
var CHprefilter             = make (chan string)
var CHmapper                = make (chan string)
var CHtag                   = make (chan string)
var CHfeed                  = make (chan string)
var CHpostfilter            = make (chan string)
var CHwriter                = make (chan string)
var CHdispatcher            = make (chan Event)

var config Analyzer
var tags Tags
var postfilters Filters
var prefilters Filters

func readconf()(err error) {

    cfg := map[string]map[string]string{}
    cfg["analyzer"] = map[string]string{}
    cfg["analyzer"]["analyzerconf"] = ""
    cfg,err = utils.GetConf(cfg)
    analyzerCFG := cfg["analyzer"]["analyzerconf"]
    if err != nil {
        logs.Error("AlertLog Error getting data from main.conf: "+err.Error())
        return
    }

    confFile, err := os.Open(analyzerCFG)
    if err != nil {
        logs.Error("Error openning analyzer CFG: "+err.Error())
        return err
    }
    defer confFile.Close()
    byteValue, _ := ioutil.ReadAll(confFile)
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
        logs.Error("Error openning analyzer tag file: "+err.Error())
        return
    }
    defer tagFile.Close()

    byteValue, _ := ioutil.ReadAll(tagFile)
    err = json.Unmarshal(byteValue, &tags)
    if err != nil {
        logs.Error("tags to json -> Unmarshal error: %s", err.Error())
    }
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
    logs.Warn(prefilters)
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
}


func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var lines []string
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
        // for ioc := range IoCs {
        //     if strings.Contains(line, IoCs[ioc]) {
        //         logs.Info("Match -> "+ line +" IoC found -> " + IoCs[ioc]  + " wkrid -> " + strconv.Itoa(wkrid))
        //         //ioc
        //         IoCtoAlert(line, IoCs[ioc], iocsrc)
        //     }
        // }

        bline, err := json.Marshal(jsoninterface)
        if err != nil {}
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
        }

        if istheend {
            continue
        }

        bline, err := json.Marshal(jsoninterface)
        if err != nil {}
        ToDispatcher("CHpostfilter", string(bline))
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
    AlertLog := map[string]map[string]string{}
    AlertLog["node"] = map[string]string{}
    AlertLog["node"]["alertLog"] = ""
    AlertLog,err := utils.GetConf(AlertLog)
    outputfile := AlertLog["node"]["alertLog"]
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
            CHprefilter <- line.Line
        case "CHprefilter":
            // logs.Warn("dispatcher %d -> %s -> to Mapper", x, line.Source)
            CHmapper <- line.Line
        case "CHmapper":
            // logs.Warn("dispatcher %d -> %s -> to Feed", x, line.Source)
            CHfeed <- line.Line
        case "CHfeed":
            // logs.Warn("dispatcher %d -> %s -> to Postfilter", x, line.Source)
            CHpostfilter <- line.Line
        case "CHpostfilter":
            // logs.Warn("dispatcher %d -> %s -> to Writer", x, line.Source)
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

func StartFeed(wkr int) {
    logs.Info("starting Feed with %d workers", wkr)
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

func StartSource(file string) {
    logs.Info("Starting tail of source file: "+file)
    var seekv tail.SeekInfo
    seekv.Offset = 0
    seekv.Whence = os.SEEK_END
    for {
        t, err := tail.TailFile(file, tail.Config{Follow: true, Location: &seekv})
        if err != nil {
            logs.Error(">>>>>> Tail over file %s error: %s", file, err.Error())
        }
        for line := range t.Lines {
            ToDispatcher("start",line.Text)
        }
    }
}

func LoadFeed() {
    logs.Info("loading Feed")
    for file := range config.Feedfiles {
        // go StartFeed(config.Feedfiles[file].File, config.Feedfiles[file].Workers)
        go StartFeed(config.Feedfiles[file].Workers)
    }
}


func LoadMapper() {
    logs.Info("loading Mappers")
    go StartMapper(4)
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
    logs.Info("CHstartpipeline %d items",len(CHstartpipeline))
    logs.Info("CHprefilter %d items",len(CHprefilter))
    logs.Info("CHmapper %d items",len(CHmapper))
    logs.Info("CHfeed %d items",len(CHfeed))
    logs.Info("CHpostfilter %d items",len(CHpostfilter))
    logs.Info("CHwriter %d items",len(CHwriter))
    logs.Info("***************")
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
    StartDispatcher(4)
    StartPreFilter(4)
    StartPostFilter(4)

    LoadMapper()
    LoadFeed()
    LoadSources()
    for {
        analyzer,_ = PingAnalyzer()
        if analyzer["status"] == "Disabled"{
            break
        }
        time.Sleep(time.Second * 3)
        CHstats()
    }
}

func Init(){
    go InitAnalizer()
}

func PingAnalyzer()(data map[string]string ,err error) {
    alertFile := map[string]map[string]string{}
    alertFile["node"] = map[string]string{}
    alertFile["node"]["alertLog"] = ""
    alertFile,err = utils.GetConf(alertFile)    
    filePath := alertFile["node"]["alertLog"]
    if err != nil {logs.Error("PingAnalyzer Error getting data from main.conf")}

    analyzerData := make(map[string]string)
    analyzerData["status"] = "Disabled"

    analyzerStatus,err := ndb.GetStatusAnalyzer()
    if err != nil { logs.Error("Error getting Analyzer data: "+err.Error()); return analyzerData,err}

    analyzerData["status"] = analyzerStatus
    analyzerData["path"] = filePath

    fi, err := os.Stat(filePath)
    //logs.Info("analyzer outputfile stats -->")
    //logs.Info(fi)
    //logs.Info("fileinfo.Sys() = %#v\n", fi.Sys())
    if err != nil { logs.Error("Can't access Analyzer ouput file data: "+err.Error()); return analyzerData,err}
    size := fi.Size()

    analyzerData["size"] = strconv.FormatInt(size, 10)

    return analyzerData, nil
}

func ChangeAnalyzerStatus(anode map[string]string) (err error) {
    logs.Emergency("ANALYZER STATUS - NEW STATUS - "+anode["status"])
    err = ndb.UpdateAnalyzer("analyzer", "status", anode["status"])
    if err != nil { logs.Error("Error updating Analyzer status: "+err.Error()); return err}
    
    return nil
}

func SyncAnalyzer(file map[string][]byte) (err error) {
    err = utils.WriteNewDataOnFile("conf/analyzer.json", file["data"])
    if err != nil { logs.Error("Analyzer/SyncAnalyzer Error updating Analyzer file: "+err.Error()); return err}
    return err
}