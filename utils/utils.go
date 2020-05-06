package utils

import (
    "archive/tar"
    "compress/gzip"
    "crypto/md5"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/astaxie/beego/logs"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

func UpdateBPFFile(path string, file string, bpf string) (err error) {
    //delete file content
    err = os.Truncate(path+file, 0)
    if err != nil {
        logs.Error("Error truncate BPF file: " + err.Error())
        return err
    }

    //write new bpf content
    bpfByteArray := []byte(bpf)
    err = WriteNewDataOnFile(path+file, bpfByteArray)
    if err != nil {
        logs.Error("Error writing new BPF data into file: " + err.Error())
        return err
    }
    return nil
}

//create a BPF backup
func BackupFullPath(path string) (err error) {
    copy, err := GetKeyValueString("execute", "copy")
    if err != nil {
        logs.Error("Error getting data from main.conf: " + err.Error())
    }

    t := time.Now()
    destFolder := path + "-" + strconv.FormatInt(t.Unix(), 10)
    cpCmd := exec.Command(copy, path, destFolder)
    err = cpCmd.Run()
    if err != nil {
        logs.Error("utils.BackupFullPath Error exec cmd command: " + err.Error())
        return err
    }

    return nil
}

func BackupFile(path string, fileName string) (err error) {
    backupFolder, err := GetKeyValueString("node", "backupFolder")
    if err != nil {
        logs.Error("utils.BackupFile Error getting backup path: " + err.Error())
        return err
    }
    copy, err := GetKeyValueString("execute", "copy")
    if err != nil {
        logs.Error("Error getting data from main.conf: " + err.Error())
    }

    // check if folder exists
    if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
        err = os.MkdirAll(backupFolder, 0755)
        if err != nil {
            logs.Error("utils.BackupFile Error creating main backup folder: " + err.Error())
            return err
        }
    }
    // check if file exists
    if _, err := os.Stat(path + fileName); os.IsNotExist(err) {
        logs.Error("utils.BackupFile file to backup doesn't exists. Can't backup the file.")
        return nil
    }
    //return nil

    //get older backup file
    listOfFiles, err := FilelistPathByFile(backupFolder, fileName)
    if err != nil {
        logs.Error("utils.BackupFile Error walking through backup folder: " + err.Error())
        return err
    }
    count := 0
    previousBck := ""
    for x := range listOfFiles {
        count++
        if previousBck == "" {
            previousBck = x
            continue
        } else if previousBck > x {
            previousBck = x
        }
    }

    //delete older bck file if there are 5 bck files
    if count == 5 {
        err = os.Remove(backupFolder + previousBck)
        if err != nil {
            logs.Error("utils.BackupFile Error deleting older backup file: " + err.Error())
        }
    }

    //create backup
    t := time.Now()
    newFile := fileName + "-" + strconv.FormatInt(t.Unix(), 10)
    srcFolder := path + fileName
    destFolder := backupFolder + newFile

    //check if file exist
    if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
        return errors.New("utils.BackupFile error: Source file doesn't exists")
    } else {
        cpCmd := exec.Command(copy, srcFolder, destFolder)
        err = cpCmd.Run()
        if err != nil {
            logs.Error("utils.BackupFile Error exec cmd command: " + err.Error())
            return err
        }
    }
    return nil
}

//write data on a file
func WriteNewDataOnFile(path string, data []byte) (err error) {
    logs.Notice(path)
    logs.Warn(string(data))
    err = ioutil.WriteFile(path, data, 0644)
    if err != nil {
        logs.Error("Error WriteNewData")
        return err
    }

    return nil
}

//Read files
func GetConfFiles() (loadDataReturn map[string]string, err error) {
    confFilePath := "./conf/main.conf"
    JSONconf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConfFiles -> can't open Conf file: " + confFilePath)
        return nil, err
    }
    var anode map[string]map[string]string
    json.Unmarshal(JSONconf, &anode)
    return anode["files"], nil
}

//Generate a 16 bytes unique id
func Generate() (uuid string) {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        logs.Error(err)
    }
    uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return uuid
}

func LoadDefaultServerData(fileName string) (json map[string]string, err error) {
    //Get full path
    file, err := GetKeyValueString("files", fileName)
    if err != nil {
        logs.Error("LoadDefaultServerData Error getting data from main.conf: " + err.Error())
        return nil, err
    }
    fileContent := make(map[string]string)
    rawData, err := ioutil.ReadFile(file)
    if err != nil {
        logs.Error("LoadDefaultServerData Error reading file: " + err.Error())
        return nil, err
    }
    fileContent["fileContent"] = string(rawData)
    return fileContent, nil
}

func CopyFile(dstfolder string, srcfolder string, file string, BUFFERSIZE int64) (err error) {
    if BUFFERSIZE == 0 {
        BUFFERSIZE = 1000
    }
    sourceFileStat, err := os.Stat(srcfolder + file)
    if err != nil {
        logs.Error("Error checking file at CopyFile function" + err.Error())
        return err
    }
    if !sourceFileStat.Mode().IsRegular() {
        logs.Error("%s is not a regular file.", sourceFileStat)
        return errors.New(sourceFileStat.Name() + " is not a regular file.")
    }
    source, err := os.Open(srcfolder + file)
    if err != nil {
        return err
    }
    defer source.Close()
    _, err = os.Stat(dstfolder + file)
    if err == nil {
        return errors.New("File " + dstfolder + file + " already exists.")
    }
    destination, err := os.Create(dstfolder + file)
    if err != nil {
        logs.Error("Error Create =-> " + err.Error())
        return err
    }
    defer destination.Close()
    logs.Info("copy file -> " + srcfolder + file)
    logs.Info("to file -> " + dstfolder + file)
    buf := make([]byte, BUFFERSIZE)
    for {
        n, err := source.Read(buf)
        if err != nil && err != io.EOF {
            logs.Error("Error no EOF=-> " + err.Error())
            return err
        }
        if n == 0 {
            break
        }
        if _, err := destination.Write(buf[:n]); err != nil {
            logs.Error("Error Writing File: " + err.Error())
            return err
        }
    }
    return err
}

func RemoveFile(path string, file string) (err error) {
    err = os.Remove(path + file)
    if err != nil {
        logs.Error("Error deleting file " + path + file + ": " + err.Error())
        return err
    }
    return nil
}

func RunCommand(cmdtxt string, params string) (err error) {
    cmd := exec.Command(cmdtxt, params)
    err = cmd.Run()
    if err != nil {
        logs.Error("utils run command -> " + err.Error())
        return err
    }
    return err
}

func StartCommand(cmdtxt string, params string) (err error) {
    cmd := exec.Command(cmdtxt, params)
    err = cmd.Start()
    if err != nil {
        logs.Error("utils run command -> " + err.Error())
        return err
    }
    return err
}

func FilelistPathByFile(path string, fileToSearch string) (files map[string][]byte, err error) {
    pathMap := make(map[string][]byte)
    err = filepath.Walk(path,
        func(file string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }

            if !info.IsDir() {
                pathSplit := strings.Split(file, "/")
                if strings.Contains(pathSplit[len(pathSplit)-1], fileToSearch) {
                    content, err := ioutil.ReadFile(file)
                    if err != nil {
                        logs.Error("Error filepath walk: " + err.Error())
                        return err
                    }
                    pathMap[pathSplit[len(pathSplit)-1]] = content
                }
            }
            return nil
        })
    if err != nil {
        logs.Error("Error filepath walk finish: " + err.Error())
        return nil, err
    }

    return pathMap, nil
}

//extract tar.gz files
func ExtractFile(tarGzFile string, pathDownloads string) (err error) {
    base := filepath.Base(tarGzFile)
    fileType := strings.Split(base, ".")

    wget, err := GetKeyValueString("execute", "command")
    if err != nil {
        logs.Error("ExtractFile Error getting data from main.conf")
        return err
    }

    if fileType[len(fileType)-1] == "rules" {
        cmd := exec.Command(wget, tarGzFile, "-O", pathDownloads)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Run()
    } else {
        logs.Info("file to open -> %s", tarGzFile)
        file, err := os.Open(tarGzFile)
        defer file.Close()
        if err != nil {
            logs.Error("something went wrong opening file %s", tarGzFile)
            return err
        }

        Untar(file, pathDownloads)
        return nil
        uncompressedStream, err := gzip.NewReader(file)
        if err != nil {
            return err
        }

        tarReader := tar.NewReader(uncompressedStream)
        for true {
            header, err := tarReader.Next()
            if err == io.EOF {
                break
            }
            if err != nil {
                return err
            }

            switch header.Typeflag {
            case tar.TypeDir:
                err := os.MkdirAll(pathDownloads+"/"+header.Name, 0755)
                if err != nil {
                    logs.Error("TypeDir: " + err.Error())
                    return err
                }
            case tar.TypeReg:
                outFile, err := os.Create(pathDownloads + "/" + header.Name)
                _, err = io.Copy(outFile, tarReader)
                if err != nil {
                    logs.Error("TypeReg: " + err.Error())
                    return err
                }
            default:
                logs.Error(
                    "ExtractTarGz: uknown type: %s in %s",
                    header.Typeflag,
                    header.Name)
            }
        }
    }

    return nil
}

func CalculateMD5(path string) (md5Data string, err error) {
    file, err := os.Open(path)
    if err != nil {
        logs.Error("Error calculating md5: %s", err.Error())
        return "", err
    }
    defer file.Close()
    hash := md5.New()
    _, err = io.Copy(hash, file)
    if err != nil {
        logs.Error("Error copying md5: %s", err.Error())
        return "", err
    }
    hashInBytes := hash.Sum(nil)[:16]
    returnMD5String := hex.EncodeToString(hashInBytes)

    return returnMD5String, nil
}

func Untar(r io.Reader, dir string) error {
    return untar(r, dir)
}

func untar(r io.Reader, dir string) (err error) {
    t0 := time.Now()
    nFiles := 0
    madeDir := map[string]bool{}
    // defer func() {
    //     td := time.Since(t0)
    //     if err == nil {
    //         logs.Info("extracted tarball into %s: %d files, %d dirs (%v)", dir, nFiles, len(madeDir), td)
    //     } else {
    //         logs.Info("error extracting tarball into %s after %d files, %d dirs, %v: %v", dir, nFiles, len(madeDir), td, err)
    //     }
    // }()
    zr, err := gzip.NewReader(r)
    if err != nil {
        return fmt.Errorf("requires gzip-compressed body: %v", err)
    }
    logs.Debug(".- gz done, lets tar -.")
    tr := tar.NewReader(zr)
    loggedChtimesError := false
    logs.Debug(".- taring -.")

    for {
        logs.Debug(".- taring for -.")
        f, err := tr.Next()
        if err == io.EOF {
            logs.Error(".- EOF -.")
            break
        }
        if err != nil {
            logs.Error("tar reading error: %v", err)
            return fmt.Errorf("tar error: %v", err)
        }
        // if !validRelPath(f.Name) {
        //     logs.Error(".- Error, not valid Path %s -.", f.Name)
        //     return fmt.Errorf("tar contained invalid name error %q", f.Name)
        // }
        logs.Debug("let's write %s", f.Name)
        rel := filepath.FromSlash(f.Name)
        logs.Debug("rel is -> %s", rel)
        abs := filepath.Join(dir, rel)
        logs.Debug("abs is -> %s", abs)

        fi := f.FileInfo()
        mode := fi.Mode()
        switch {
        case mode.IsRegular():
            // Make the directory. This is redundant because it should
            // already be made by a directory entry in the tar
            // beforehand. Thus, don't check for errors; the next
            // write will fail with the same error.
            dir := filepath.Dir(abs)
            if !madeDir[dir] {
                if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
                    return err
                }
                madeDir[dir] = true
            }
            wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
            if err != nil {
                return err
            }
            n, err := io.Copy(wf, tr)
            if closeErr := wf.Close(); closeErr != nil && err == nil {
                err = closeErr
            }
            if err != nil {
                return fmt.Errorf("error writing to %s: %v", abs, err)
            }
            if n != f.Size {
                return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, f.Size)
            }
            modTime := f.ModTime
            if modTime.After(t0) {
                // Clamp modtimes at system time. See
                // golang.org/issue/19062 when clock on
                // buildlet was behind the gitmirror server
                // doing the git-archive.
                modTime = t0
            }
            if !modTime.IsZero() {
                if err := os.Chtimes(abs, modTime, modTime); err != nil && !loggedChtimesError {
                    // benign error. Gerrit doesn't even set the
                    // modtime in these, and we don't end up relying
                    // on it anywhere (the gomote push command relies
                    // on digests only), so this is a little pointless
                    // for now.
                    logs.Info("error changing modtime: %v (further Chtimes errors suppressed)", err)
                    loggedChtimesError = true // once is enough
                }
            }
            nFiles++
        case mode.IsDir():
            logs.Debug("let's create a folder -> %s", abs)
            if err := os.MkdirAll(abs, 0755); err != nil {
                logs.Debug("error creting directory -> %s", abs)
                return err
            }
            madeDir[abs] = true
        default:
            logs.Debug("tar file entry %s contained unsupported file type %v", f.Name, mode)
            return fmt.Errorf("tar file entry %s contained unsupported file type %v", f.Name, mode)
        }
    }
    return nil
}

func validRelativeDir(dir string) bool {
    if strings.Contains(dir, `\`) || path.IsAbs(dir) {
        return false
    }
    dir = path.Clean(dir)
    if strings.HasPrefix(dir, "../") || strings.HasSuffix(dir, "/..") || dir == ".." {
        return false
    }
    return true
}

func validRelPath(p string) bool {
    if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
        return false
    }
    return true
}
