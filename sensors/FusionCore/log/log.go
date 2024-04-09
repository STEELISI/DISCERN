package logging;

import (
    "log"
    "os"
    "time"
)

// Given predefined value so we can log if config doesn't load correctly
var LogFile string = "/var/log/discern.log"
var ProgramName string = "FusionCore"; // Used by config on boot

func SetLogFile(path string) {
    LogFile = path
}

func LogInfo(msg string) {
    fd, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer fd.Close()
    log.SetOutput(fd)

    currentTime := time.Now().Format("2006-01-02 15:04:05")
    log.Printf("[%s] %s: %s\n", currentTime, ProgramName, msg)
}

func FatalError(msg string) {
    fd, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer fd.Close()
    log.SetOutput(fd)
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    log.Fatalf("[%s] %s: %s\n", currentTime, ProgramName, msg)
}

