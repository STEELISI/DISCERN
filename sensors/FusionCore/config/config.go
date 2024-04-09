package config;

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "gopkg.in/yaml.v2"

    Log "FusionCore/log"
)


var GRPC_Port, ListenIP , ListenProto, InternalDockerPort, 
    ExternalPort, DBurl, DBtoken, BUCKET_NAME, ORG, PASSWORD, 
    USERNAME string;

var StartNewDockerInstance bool;
var ProgramName string = "FusionCore";
var Settings YamlConfig = NewYamlConfig();

type YamlConfig struct {
    // Logging information
    LogFile                        string `default:"/var/log/discern.log"`
    // For gRPC connection
    GRPC_Port                      int64  `default:50051`
    ListenIP                       string `default:"localhost"`
    ListenProto                    string `default:"tcp"`
    // For influxDB connection
    InternalDockerPort             int64  `default:8086`
    ExternalPort                   int64  `default:8086`
    DBurl                          string `default:"http://localhost"`
    DBtoken                        string `default:"BIGElHSa291FOkrliGaBVc7ksnGgQ4vALbkfJzRuH02T2XB8qouH0H3IkYTJACE-XZ-QYV664CH5655LkbQDIQ::"`
    StartNewDockerInstance         bool   `default:false`
    BUCKET_NAME                    string `default:"DISCERN"`
    ORG                            string `default:"ISI"`
    PASSWORD                       string `default:"something"`
    USERNAME                       string `default:"default-user"`
    // For postgres connection
    ExternalPostgresPort           int64  `default:5432`
    PostgresDataMount              string `default:/var/lib/postgresql/data`
    StartNewPostgresInstance       bool   `default: false`
    PostgresUser                   string `default: "postgres"`
    PostgresPassword               string `default: "*Something*"`
    PostgresIP                     string `default: "127.0.0.1"`
}


func NewYamlConfig() YamlConfig {
    return YamlConfig{
        LogFile                    : "/var/log/discern.log",
        GRPC_Port                  : 50051,
        ListenIP                   : "localhost",
        ListenProto                : "tcp",
        InternalDockerPort         : 8086,
        ExternalPort               : 8086,
        DBurl                      : "http://localhost",
        DBtoken                    : "BIGElHSa291FOkrliGaBVc7ksnGgQ4vALbkfJzRuH02T2XB8qouH0H3IkYTJACE-XZ-QYV664CH5655LkbQDIQ::",
        StartNewDockerInstance     : false,
        BUCKET_NAME                : "DISCERN",
        ORG                        : "ISI",
        PASSWORD                   : "something",
        USERNAME                   : "default-user",
        ExternalPostgresPort       : 5432,
        PostgresDataMount          : "/var/lib/postgresql/data",
        StartNewPostgresInstance   : false,
        PostgresUser               : "postgres",
        PostgresPassword           : "*Something*",
        PostgresIP                 : "127.0.0.1",
    }
}

func LoadConfig() {

    yamlFile, err := ioutil.ReadFile("/etc/discern/CoreConfig.yaml")
    if err != nil {
        fmt.Println("Fail one")
        Log.LogInfo(fmt.Sprintf("yamlFile.Get err #%v\n", err))
    } else { // Funny way to do error handling
        err = yaml.Unmarshal(yamlFile, &Settings)
        if err != nil {
            fmt.Println("Fail two")
            Log.LogInfo(fmt.Sprintf("Unmarshal: %v\n", err))
        }
    }

    // Logging information
    Log.SetLogFile(Settings.LogFile)
    // For gRPC connection
    GRPC_Port              = strconv.FormatInt(Settings.GRPC_Port, 10)
    ListenIP               = Settings.ListenIP
    ListenProto            = Settings.ListenProto
    // For influxDB connect
    InternalDockerPort     = strconv.FormatInt(Settings.InternalDockerPort, 10)
    ExternalPort           = strconv.FormatInt(Settings.ExternalPort, 10)
    DBurl                  = Settings.DBurl
    DBtoken                = Settings.DBtoken
    StartNewDockerInstance = Settings.StartNewDockerInstance
    BUCKET_NAME            = Settings.BUCKET_NAME
    ORG                    = Settings.ORG
    PASSWORD               = Settings.PASSWORD
    USERNAME               = Settings.USERNAME
    ProgramName            = Log.ProgramName
}

func FullAddr() string {
    return DBurl + ":" + ExternalPort
}
