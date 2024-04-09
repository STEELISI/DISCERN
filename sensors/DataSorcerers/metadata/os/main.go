package main;

import (
    "fmt"
    "os"
    "OS/save"
    "time"
    "DataSorcerers/helpers"
);


func main() {
    helpers.SetFileName("OS")

    helpers.LoadConfig()
    if !helpers.Config.RunOs { return }

    helpers.LogInfo("Hello World from OS metadata scraper")

    build_data_dir()

    // build_sensor_file(filename, commands)

    save.ConnectClient()

    // The sensors timestamp their own data so we send when convenient
    for {
        save.CloseData()
        save.CloseRangeData()
        save.ExecveData()
        save.ExecveAtData()
        save.ForkData()
        save.KillData()
        save.OpenData()
        save.OpenAtData()
        save.OpenAt2Data()
        save.RecvFromData()
        save.RecvMMsgData()
        save.RecvMsgData()
        save.SendMMsgData()
        save.SendMsgData()
        save.SendToData()
        save.SocketData()
        save.SocketPairData()
        save.SysInfoData()
        save.TKillData()
        save.VForkData()

        time.Sleep(time.Duration(helpers.Config.OsInfoDumpInterval) * time.Second)
    }
}

func build_data_dir() {
    if err := os.MkdirAll(helpers.Config.OsDataDir, os.ModePerm); err != nil {
        helpers.LogInfo(fmt.Sprintf("Error creating directory: %v", err))
    }
    // Check if the directory was created or already exists
    if _, err := os.Stat(helpers.Config.OsDataDir); os.IsNotExist(err) {
        helpers.LogInfo("Directory was not created.")
    }
}

func build_sensor_file() {

    // All the different bpftrace files we'd like to record
    sensors := []string{
        "close",
        "close_range",
        "execve",
        "execveat",
        "fork",
        "kill",
        "open",
        "openat",
        "openat2",
        "recvfrom",
        "recvmmsg",
        "recvmsg",
        "sendmmsg",
        "sendmsg",
        "sendto",
        "socket",
        "socketpair",
        "sysinfo",
        "tkill",
        "vfork",
    }

    filename := "./run_sensors.sh"

    err := os.Remove(filename)
    if err != nil {
        helpers.FatalError(fmt.Sprintf("Error opening files: ", err))
    }

    file, err := os.Create(filename)
    defer file.Close()

    if err != nil {
        helpers.FatalError(fmt.Sprintf("Error opening files: ", err))
    }

    file.WriteString("#!/bin/bash\n")

    for _, cmd := range sensors {
        cmdString := "sudo bpftrace -f json ./bpftrace/" + cmd + ".bt >> /tmp/discern/data/os/" + cmd + "-res.txt &\n"
        file.WriteString(cmdString)
    }
}

