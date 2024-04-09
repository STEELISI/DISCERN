package save;

/*
    Sorry to whoever finds this. This should be auto generated since
    they all need slightly different parsing but this works just as 
    well
*/

import (
    "fmt"
    "bufio"
    "log"
    "strconv"
    "os"
    "context"
    "encoding/json"
    "strings"
    "time"
    "google.golang.org/grpc" 
    "github.com/golang/protobuf/ptypes"
    "github.com/golang/protobuf/ptypes/timestamp"

    bridge "FusionBridge/metadata/os"
    helpers "DataSorcerers/helpers"
)




var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var err error; 
var client bridge.OSClient;



func ConnectClient() {
    conn = helpers.CreateConnection()
    client = bridge.NewOSClient(conn); 
}


// CustomError is a custom error type that implements the error interface.
type TimeStampError struct {
    message string
}
// Error returns the error message for the CustomError.
func (e *TimeStampError) Error() string {
    return e.message
}




func check_for_timestamp(value map[string]interface{}) (*timestamp.Timestamp, *TimeStampError) {
    if value["type"] != "time" { return nil, &TimeStampError{"Not a Timestamp"}; }

    // Define the layout format for the time
    layout := "15:04:05"

    // Remove the trailing newline character if present
    str, ok := value["data"].(string)
    if !ok { return nil, &TimeStampError{"Data is not string"}}
    input := string(str)[:len(str)-1]

    // Parse the input string as a time
    t, err := time.Parse(layout, input)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error parsing timestamp: %v\n", err))
        return nil, &TimeStampError{"Error parsing timestamp"}
    }
    timestamp, err := ptypes.TimestampProto(t)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error parsing timestamp: %v\n", err))
        return nil, &TimeStampError{"Error converting timestamp to gRPC format"}
    }
    return timestamp, nil
}


func CloseData() {
    file_loc := helpers.Config.OsDataDir + "/close-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("CloseData Error opening Close Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isFloat64 := val.(float64)
                if !isFloat64 { continue }
                
                // pid,uid
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("CloseData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("CloseData Failed to convert uid %v to int", tmp[1]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("CloseData Failed to convert gid %v to int", tmp[2]))
                }

                MSG := bridge.CloseCall{
                    TimeStamp: timestamp,
                    Pid:   int32(pid),
                    Uid:   int32(uid),
                    Gid:   int32(gid),
                    Count: int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkCloseSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        panic(err)
    }
}


func CloseRangeData() {
    file_loc := helpers.Config.OsDataDir + "/close_range-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Close Range Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isFloat64 := val.(float64)
                if !isFloat64 { continue }
                
                // pid,uid
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("CloseRangeData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("CloseRangeData Failed to convert uid %v to int", tmp[1]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("CloseRangeData Failed to convert gid %v to int", tmp[2]))
                }

                MSG := bridge.CloseRangeCall{
                    TimeStamp: timestamp,
                    Pid:   int32(pid),
                    Uid:   int32(uid),
                    Gid:   int32(gid),
                    Count: int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkCloseRangeSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}



func ExecveData() {
    file_loc := helpers.Config.OsDataDir + "/execve-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Execve Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        // Types can be quite annoying
        for _, value := range data {
            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                arg, isString := val.(string)
                if !isString { continue }
                
                // pid,nsecs,uid,gid,arg#
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveData Failed to convert pid %v to int", tmp[0]))
                }
                nsecs, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveData Failed to convert nsecs %v to int", tmp[1]))
                }
                uid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveData Failed to convert gid %v to int", tmp[3]))
                }
                argNum, err := (strconv.Atoi(string(tmp[4][len(tmp[4]) - 1])))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveData Failed to convert argNum %v to int", tmp[4]))
                }

                // Properly format the timestamps
                nanoseconds := int64(nsecs)

                // Define the Unix epoch as a reference time
                epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

                // Calculate the timestamp by adding nanoseconds to the Unix epoch
                timestamp := epoch.Add(time.Duration(nanoseconds))

                gRPC_timestamp, err := ptypes.TimestampProto(timestamp)
                if err != nil {
                    out := fmt.Sprintf("Error in client.os.Execve. Could not create timestamp: %v", err);
                    helpers.LogInfo(out)
                    return
                }


                MSG := bridge.ExecveCall{
                    TimeStamp: gRPC_timestamp,
                    Pid:    int32(pid),
                    Uid:    int32(uid),
                    Gid:    int32(gid),
                    ArgNum: int32(argNum),
                    Arg:    string(arg),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkExecveSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func ExecveAtData() {
    file_loc := helpers.Config.OsDataDir + "/execveat-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Execveat Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        // Types can be quite annoying
        for _, value := range data {
            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                arg, isString := val.(string)
                if !isString { continue }
                
                // pid,nsecs,uid,gid,arg#
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveAtData Failed to convert pid %v to int\n", tmp[0]))
                }
                nsecs, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveAtData Failed to convert nsecs %v to int\n", tmp[1]))
                }
                uid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveAtData Failed to convert uid %v to int\n", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveAtData Failed to convert gid %v to int\n", tmp[3]))
                }
                argNum, err := (strconv.Atoi(string(tmp[4][len(tmp[4]) - 1])))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ExecveAtData Failed to convert argNum %v to int", tmp[4]))
                }

                // Properly format the timestamps
                nanoseconds := int64(nsecs)

                // Define the Unix epoch as a reference time
                epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

                // Calculate the timestamp by adding nanoseconds to the Unix epoch
                timestamp := epoch.Add(time.Duration(nanoseconds))

                gRPC_timestamp, err := ptypes.TimestampProto(timestamp)
                if err != nil {
                    log.Fatalf("Error in client.os.Execve. Could not create timestamp: %v", err);
                    return
                }


                MSG := bridge.ExecveAtCall{
                    TimeStamp: gRPC_timestamp,
                    Pid:    int32(pid),
                    Uid:    int32(uid),
                    Gid:    int32(gid),
                    ArgNum: int32(argNum),
                    Arg:    string(arg),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkExecveAtSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func ForkData() {
    file_loc := helpers.Config.OsDataDir + "/fork-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Fork Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
    for key := range data {
        keys = append(keys, key)
    }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg#
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ForkData Failed to convert pid %v to int\n", tmp[0]))
                }
                tid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ForkData Failed to convert nsecs %v to int\n", tmp[1]))
                }
                uid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ForkData Failed to convert uid %v to int\n", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("ForkData Failed to convert gid %v to int\n", tmp[3]))
                }

                MSG := bridge.ForkCall{
                    TimeStamp: timestamp,
                    Pid:    int32(pid),
                    Tid:    int32(tid),
                    Uid:    int32(uid),
                    Gid:    int32(gid),
                    Count:  int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkForkSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func KillData() {
    file_loc := helpers.Config.OsDataDir + "/kill-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Kill Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                sig, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg#
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("KillData Failed to convert pid %v to int", tmp[0]))
                }
                argPid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("KillData Failed to convert nsecs %v to int", tmp[1]))
                }
                uid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("KillData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("KillData Failed to convert gid %v to int", tmp[3]))
                }

                MSG := bridge.KillCall{
                    TimeStamp: timestamp,
                    Pid:    int32(pid),
                    ArgPid: int32(argPid),
                    Uid:    int32(uid),
                    Gid:    int32(gid),
                    Sig:    int32(sig),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkKillSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func OpenData() {
    file_loc := helpers.Config.OsDataDir + "/open-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Kill Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenData Failed to convert gid %v to int", tmp[3]))
                }
                filename := tmp[3]

                MSG := bridge.OpenCall{
                    TimeStamp: timestamp,
                    Pid:      int32(pid),
                    Uid:      int32(uid),
                    Gid:      int32(gid),
                    Filename: filename,
                    Count:    int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkOpenSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func OpenAtData() {
    file_loc := helpers.Config.OsDataDir + "/openat-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Kill Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenAtData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenAtData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenAtData Failed to convert gid %v to int", tmp[3]))
                }
                filename := tmp[3]

                MSG := bridge.OpenAtCall{
                    TimeStamp: timestamp,
                    Pid:      int32(pid),
                    Uid:      int32(uid),
                    Gid:      int32(gid),
                    Filename: filename,
                    Count:    int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkOpenAtSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func OpenAt2Data() {
    file_loc := helpers.Config.OsDataDir + "/openat2-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Kill Data file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenAt2Data Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenAt2Data Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("OpenAt2Data Failed to convert gid %v to int", tmp[3]))
                }
                filename := tmp[3]

                MSG := bridge.OpenAt2Call{
                    TimeStamp: timestamp,
                    Pid:      int32(pid),
                    Uid:      int32(uid),
                    Gid:      int32(gid),
                    Filename: filename,
                    Count:    int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkOpenAt2SysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func RecvFromData() {
    file_loc := helpers.Config.OsDataDir + "/recvfrom-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Recv From file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvFromData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvFromData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvFromData Failed to convert gid %v to int", tmp[3]))
                }
                fd, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvFromData Failed to convert gid %v to int", tmp[3]))
                }

                MSG := bridge.RecvFromCall{
                    TimeStamp: timestamp,
                    Pid:      int32(pid),
                    Uid:      int32(uid),
                    Gid:      int32(gid),
                    FileDescriptor: int32(fd),
                    Count:    int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkRecvFromSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func RecvMMsgData() {
    file_loc := helpers.Config.OsDataDir + "/recvmmsg-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Recv MMsg file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMMsgData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMMsgData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMMsgData Failed to convert gid %v to int", tmp[3]))
                }
                fd, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMMsgData Failed to convert gid %v to int", tmp[3]))
                }

                MSG := bridge.RecvMMsgCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    FileDescriptor: int32(fd),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkRecvMMsgSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func RecvMsgData() {
    file_loc := helpers.Config.OsDataDir + "/recvmsg-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Recv Msg file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMsgData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMsgData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMsgData Failed to convert gid %v to int", tmp[3]))
                }
                fd, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("RecvMsgData Failed to convert gid %v to int", tmp[3]))
                }

                MSG := bridge.RecvMsgCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    FileDescriptor: int32(fd),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkRecvMsgSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func SendMMsgData() {
    file_loc := helpers.Config.OsDataDir + "/sendmmsg-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send MMsg file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMMsgData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMMsgData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMMsgData Failed to convert gid %v to int", tmp[3]))
                }
                fd, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMMsgData Failed to convert gid %v to int", tmp[3]))
                }
                length, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMMsgData Failed to convert length %v to int", tmp[3]))
                }

                MSG := bridge.SendMMsgCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    Len:            int32(length),
                    FileDescriptor: int32(fd),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkSendMMsgSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func SendMsgData() {
    file_loc := helpers.Config.OsDataDir + "/sendmsg-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send Msg file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMsgData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMsgData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMsgData Failed to convert gid %v to int", tmp[3]))
                }
                fd, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMsgData Failed to convert gid %v to int", tmp[3]))
                }
                length, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendMsgData Failed to convert length %v to int", tmp[3]))
                }

                MSG := bridge.SendMsgCall{
                    TimeStamp:      timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    Len:            int32(length),
                    FileDescriptor: int32(fd),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkSendMsgSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func SendToData() {
    file_loc := helpers.Config.OsDataDir + "/sendto-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send To file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendToData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendToData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendToData Failed to convert gid %v to int", tmp[3]))
                }
                fd, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendToData Failed to convert gid %v to int", tmp[3]))
                }
                length, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SendToData Failed to convert length %v to int", tmp[3]))
                }

                MSG := bridge.SendMsgCall{
                    TimeStamp:      timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    Len:            int32(length),
                    FileDescriptor: int32(fd),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkSendMsgSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func SocketData() {
    file_loc := helpers.Config.OsDataDir + "/socket-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send To file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketData Failed to convert gid %v to int", tmp[3]))
                }
                family, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketData Failed to convert gid %v to int", tmp[3]))
                }
                kind, err := (strconv.Atoi(tmp[4]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketData Failed to convert length %v to int", tmp[3]))
                }
                protocol, err := (strconv.Atoi(tmp[5]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketData Failed to convert length %v to int", tmp[3]))
                }

                MSG := bridge.SocketCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    Family:         int32(family),
                    Type:           int32(kind),
                    Protocol:       int32(protocol),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkSocketSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func SocketPairData() {
    file_loc := helpers.Config.OsDataDir + "/socketpair-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send To file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketPairData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketPairData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketPairData Failed to convert gid %v to int", tmp[3]))
                }
                family, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketPairData Failed to convert gid %v to int", tmp[3]))
                }
                kind, err := (strconv.Atoi(tmp[4]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketPairData Failed to convert length %v to int", tmp[3]))
                }
                protocol, err := (strconv.Atoi(tmp[5]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SocketPairData Failed to convert length %v to int", tmp[3]))
                }

                MSG := bridge.SocketPairCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    Family:         int32(family),
                    Type:           int32(kind),
                    Protocol:       int32(protocol),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkSocketPairSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func SysInfoData() {
    file_loc := helpers.Config.OsDataDir + "/sysinfo-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send To file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for _, value := range data {
            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                nsecs, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SysInfoData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SysInfoData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("SysInfoData Failed to convert gid %v to int", tmp[3]))
                }

                // Properly format the timestamps
                nanoseconds := int64(nsecs)

                // Define the Unix epoch as a reference time
                epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

                // Calculate the timestamp by adding nanoseconds to the Unix epoch
                timestamp := epoch.Add(time.Duration(nanoseconds))

                gRPC_timestamp, err := ptypes.TimestampProto(timestamp)
                if err != nil {
                    log.Fatalf("Error in client.os.Execve. Could not create timestamp: %v", err);
                    return
                }


                MSG := bridge.SysInfoCall{
                    TimeStamp:      gRPC_timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkSysInfoSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
                }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)

    }
}


func TKillData() {
    file_loc := helpers.Config.OsDataDir + "/tkill-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send To file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }

        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("TKillData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("TKillData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("TKillData Failed to convert gid %v to int", tmp[3]))
                }
                argsPid, err := (strconv.Atoi(tmp[3]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("TKillData Failed to convert uid %v to int", tmp[3]))
                }
                sig, err := (strconv.Atoi(tmp[4]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("TKillData Failed to convert gid %v to int", tmp[4]))
                }

                MSG := bridge.TKillCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    ArgPid:         int32(argsPid),
                    Sig:            int32(sig),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkTKillSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving Close Syscall Event: %v\n", e))
                    return
               }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}


func VForkData() {
    file_loc := helpers.Config.OsDataDir + "/vfork-res.txt"

    file, err := os.OpenFile(file_loc, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        helpers.LogInfo(fmt.Sprintf("Error opening Send To file: %v\n", err))
    }
    defer file.Close()


    var data map[string]interface{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 { continue }
        if err := json.Unmarshal([]byte(line), &data); err != nil {
                helpers.LogInfo(fmt.Sprintf("Error:", err))
                return
        }


        keys := make([]string, 0, len(data))
        for key := range data {
            keys = append(keys, key)
        }

        // Types can be quite annoying
        for i:=0; i< len(data); i++ {
            value := data[keys[i]]

            time_json, ok := value.(map[string]interface{});

            if !ok { continue }

            // Check for timestamp
            timestamp, err := check_for_timestamp(time_json)
            if err != nil { continue } // If its not a timestamp ignore it

            i++;
            value = data[keys[i]]

            mapo, ok := value.(map[string]interface{});

            if !ok { continue }

            atValue, atExists := mapo["@"];

            if !atExists { continue }

            v, isMap := atValue.(map[string]interface{});
            if !isMap { continue }

            for k, val := range v {
                count, isfloat64 := val.(float64)
                if !isfloat64 { continue }
                
                // pid,nsecs,uid,gid,arg
                tmp := strings.Split(k, ",")
                
                pid, err := (strconv.Atoi(tmp[0]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("VForkData Failed to convert pid %v to int", tmp[0]))
                }
                uid, err := (strconv.Atoi(tmp[1]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("VForkData Failed to convert uid %v to int", tmp[2]))
                }
                gid, err := (strconv.Atoi(tmp[2]))
                if err != nil {
                    helpers.LogInfo(fmt.Sprintf("VForkData Failed to convert gid %v to int", tmp[3]))
                }

                MSG := bridge.VForkCall{
                    TimeStamp: timestamp,
                    Pid:            int32(pid),
                    Uid:            int32(uid),
                    Gid:            int32(gid),
                    Count:          int32(count),
                    DevID: helpers.GetID(),
                }
                _, e := client.MarkVForkSysCall(context.Background(), &MSG)
                if e != nil {
                    helpers.LogInfo(fmt.Sprintf("Error while saving VFORK Syscall Event: %v\n", e))
                    return
              }
            }
        }
    }
    err = file.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("truncation error: %v", err)
        helpers.LogInfo(out)
    }
}

