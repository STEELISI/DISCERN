package helpers;

import (
    "fmt"
    "os"
    "io/ioutil"
    "strconv"
    "net/http"
    "encoding/json"
    "bytes"
    "time"
    "gopkg.in/yaml.v2"
    "log"
    "google.golang.org/grpc" 
    "google.golang.org/grpc/credentials/insecure"
)

var Config YamlConfig = *NewYamlConfig();
var FileName string;
var LoadIDFromRemote bool = true;
var logFile *os.File;

type YamlConfig struct {
    // Logging information
    LogFile                        string `default:"/var/log/discern.log"`
    // For gRPC connection
    GRPC_Port                      int64  `default:50051`
    GRPC_IP                        string `default:"localhost"`
    GRPC_PROTO                     string `default:"tcp"`
    MaxRecvMsgSize                 int  `default:10485760`
    // For TLS connection
    TLS                            bool   `default:false`
    CertFile                       string `default:""`
    KeyFile                        string `default:""`
    // For ID section
    IdApiUrl                       string `default:"localhost:50051/api/id"`
    IdApiBody                      string `default:""`
    RunIdApiServer                 bool `default:true`
    // Specifying which service should run
    RunLog                         bool `default:true`
    RunAnsible                     bool `default:true`
    RunBash                        bool `default:true`
    RunJupyter                     bool `default:true`
    RunFile                        bool `default:true`
    RunId                          bool `default:true`
    RunIdServer                    bool `default:true`
    RunNetwork                     bool `default:true`
    RunOs                          bool `default:true`
    // For metadata/os
    OsInfoDumpInterval             int `default:3001`
    OsDataDir                      string `default:"/tmp/discern/data/os"`
    // For metadata/file
    BlackListDirs                  []string `default:[]string{"/home/.*/\.config","/home/.*/\..*"}`
    StartingDirs                   []string `default:[]string{"/home","/tmp"}`
    // For metadata/network
    NetworkSliceLength             uint `default:5`
    // For Log info
    LogSweepInterval               int `default:5`
    LogLinesCaptured               int `default:5`
    LogFiles                       []string `default:[]string{"/var/log/messages", "/var/log/syslog", "/var/log/auth.log", "/var/log/daemon.log", "/var/log/kern.log"}`
    // For jupyter data
    JupyterSweepInterval           int `default:5`
    JupyterPaths                   []string `default:[]string{"/home"}`
    // For ansible data
    SaveAnsibleConfigs             bool `default:false`
    SaveAnsiblePlaybooks           bool `default:false`
    AnsibleConfigSweepInterval     int `default:5`
    AnsiblePlaybookSweepInterval   int `default:5`
    // For bash logging data
    BashSweepInterval              int `default:5`
    // For ID information
    ScrapeForUsers                 int `default:30000`
    ScrapeInterfaceInfo            int `default:3000`
}


func NewYamlConfig() *YamlConfig {
    return &YamlConfig{
        LogFile:                      "/var/log/discern.log",

        GRPC_Port:                    50051,
        GRPC_IP:                      "localhost",
        GRPC_PROTO:                   "tcp",
        MaxRecvMsgSize:               10485760,

        TLS:                          false,
        CertFile:                     "",
        KeyFile:                      "",

        IdApiUrl:                     "http://localhost:50052/api/id",
        IdApiBody:                    "",
        RunIdApiServer:               true,

        RunLog:                       true,
        RunAnsible:                   true,
        RunBash:                      true,
        RunJupyter:                   true,
        RunFile:                      true,
        RunId:                        true,
        RunIdServer:                  true,
        RunNetwork:                   true,
        RunOs:                        true,

        OsInfoDumpInterval:           3001,
        OsDataDir:                    "/tmp/discern/data/os",

        BlackListDirs:                []string{`/home/.*/\.config`, `/home/.*/\..*`},
        StartingDirs:                 []string{"/home", "/tmp"},

        NetworkSliceLength:           5,

        LogSweepInterval:             5,
        LogLinesCaptured:             5,
        LogFiles:                     []string{"/var/log/messages", "/var/log/syslog", "/var/log/auth.log", "/var/log/daemon.log", "/var/log/kern.log"},

        JupyterSweepInterval:         5,
        JupyterPaths:                 []string{"/home"},

        SaveAnsibleConfigs:           true,
        SaveAnsiblePlaybooks:         true,
        AnsibleConfigSweepInterval:   5,
        AnsiblePlaybookSweepInterval: 5,

        BashSweepInterval:            5,

        ScrapeForUsers:               30000,
        ScrapeInterfaceInfo:          3000,
    }
}


func (in *YamlConfig) ServerAddr() string {
    return in.GRPC_IP + ":" + strconv.FormatInt(in.GRPC_Port, 10)
}


func LoadConfig() {
    yamlFile, err := ioutil.ReadFile("/etc/discern/SorcererConfig.yaml")
    if err != nil {
        LogInfo(fmt.Sprintf("yamlFile.Get err #%v", err))
    } else { // Funny way to do error handling
        err = yaml.Unmarshal(yamlFile, &Config)
        if err != nil {
            LogInfo(fmt.Sprintf("Unmarshal: %v", err))
        }
    }
    if LoadIDFromRemote {
        LoadID()
    }
    // This should be synchronized but whatever. Good enough
}


var ID_STRING string = "";

func GetID() string {
    return ID_STRING;
}


func LoadID() {
    b := new(bytes.Buffer)
    
    err := json.NewEncoder(b).Encode(Config.IdApiBody)
    if err != nil {
        LogInfo(fmt.Sprintf("Error loading the ID"))
        return
    }
    
    client := http.Client{
        Timeout: 5 * time.Second,
    }
    resp, err := client.Post(Config.IdApiUrl, "application/json", b)
    if err != nil {
        LogInfo(fmt.Sprintf("%v\n", err))
        return
    }

    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        // Read the response body
        responseBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            out := fmt.Sprintln("Error reading the response body:", err)
            LogInfo(out)
            return
        }

        // Print the response as a string
        responseString := string(responseBytes)
        ID_STRING = responseString
    } else {
        out := fmt.Sprintf("Request failed with status code: %d\n", resp.StatusCode)
        LogInfo(out)
    }
    return
}

func CreateConnection() *grpc.ClientConn {

    var opts []grpc.DialOption;
    opts = append(opts, 
        grpc.WithTimeout(20 * time.Second))
    opts = append(opts, 
        grpc.WithTransportCredentials(insecure.NewCredentials()));
   
    conn, err := grpc.Dial(Config.ServerAddr(), opts...);
    if err != nil {
        out := fmt.Sprintf("Failed to dial: %v", err)
        FatalError(out)
    }

    return conn
}

func LogInfo(msg string) {
    fd, err := os.OpenFile(Config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer fd.Close()
    log.SetOutput(fd)

    currentTime := time.Now().Format("2006-01-02 15:04:05")
    log.Printf("[%s] %s: %s\n", currentTime, FileName, msg)
}

func FatalError(msg string) {
    fd, err := os.OpenFile(Config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer fd.Close()
    log.SetOutput(fd)
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    log.Fatalf("[%s] %s: %s\n", currentTime, FileName, msg)
}

func SetFileName(name string) {
    FileName = name
}
