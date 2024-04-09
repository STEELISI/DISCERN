package main;

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "os"
    "strings"
    "strconv"
    "context"
    "time"
    "google.golang.org/grpc" 
    "github.com/golang/protobuf/ptypes"

    bridge "FusionBridge/control/bash"
    "DataSorcerers/helpers"
)


var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var err error; 
var client bridge.BashClient;


func main() {

    helpers.SetFileName("Bash")
    helpers.LoadConfig()
    if !helpers.Config.RunBash { return }


    establish_connection()

    for {
        files, err := ioutil.ReadDir("/var/log/ttylog")

        if err != nil {
            helpers.LogInfo(fmt.Sprintf("Cant open /var/log/ttylog: %v", err));
            return
        }
        
        for _, file := range files {
            file_path := "/var/log/ttylog/" + file.Name()

            // Defined in start_ttylog
            arr := strings.Split(file_path, ".")
            if len(arr) < 4 { 
                out := fmt.Sprintf("Improperly formatted file name: %v", file_path)
                helpers.LogInfo(out)
                continue
            }
            host := arr[1]
            user := arr[2]
            cnt, _ := strconv.Atoi(arr[3])

            // Build csv file we'd like to export
            csv, err := build_cmd_csv(file_path, host, user, cnt)

            if err != nil || len(csv) == 0 {
                out := fmt.Sprintf("Error in creating csv: %v", err)
                helpers.LogInfo(out)
                delete_csv(csv)
                continue
            }

            // Send CSV Data to server
            go send_csv_log(csv, host, user, cnt, file_path)

        }
        // wg.Wait()
        time.Sleep(time.Duration(helpers.Config.BashSweepInterval) * time.Second)
    }
}


func establish_connection() {
    conn = helpers.CreateConnection()
    client = bridge.NewBashClient(conn); 
}


func read_cmd_log(file_path string) (string, error) {
     
    code := "/usr/local/src/ttylog/ttylog --read " + file_path

    cmd := exec.Command("bash", "-c", code)

    output, err := cmd.CombinedOutput()
    formatted_output := fmt.Sprintf("%s", output)

    if err != nil {
        out := fmt.Sprintf(`Error in read_cmd_log. 
                File: %v
                Error: %v
                Output: %v`, file_path, err, formatted_output);
        helpers.LogInfo(out)
        return "", err
    }

    return formatted_output, nil
}


func build_cmd_csv(file_path string, user string, host string, cnt int) (string, error) {
    
    csv_filename := fmt.Sprintf("%v-%v-%v.csv", user, host, cnt)

    code := fmt.Sprintf(`/usr/local/src/ttylog/analyze.py %v %v`, 
                            file_path, csv_filename)
    cmd := exec.Command("bash", "-c", code)
    if output, err := cmd.CombinedOutput(); err != nil {

        out := fmt.Sprintf("Error creating csv: %v\n", err)
        helpers.LogInfo(out)
        out = fmt.Sprintf("Output: %v\n", string(output))
        helpers.LogInfo(out)

        delete_csv(csv_filename)
        return "", err
    }
    return csv_filename, nil
}


func send_cmd_log(cmds string, host string, user string, cnt int) error {
 
    // defer wg.Done()

    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    timestamp, err := ptypes.TimestampProto(time.Now())
    if err != nil {
        out := fmt.Sprintf("Error creating timestamp: %v\n", err)
        helpers.LogInfo(out)
        return err
    }

    to_send := bridge.CmdSnapShot{
        SubmissionNumber : 0,
        TimeStamp : timestamp,
        Cmds : cmds, 
        Host: host,
        Count: int32(cnt),
        User: user,
        DevID: helpers.GetID(),
    }

    _, e := client.IngestCmdSnapShot(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in Ingest CMD snapshot: %v", e)
        helpers.LogInfo(out)
        return e
    }
    return nil
} 
    

func empty_cmd_log(file_path string) error {

    // Empty the file
    fd, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        out := fmt.Sprintf("Error opening file: %v", err)
        helpers.LogInfo(out)
        return err
    }
    defer fd.Close()

    // Truncate the file to zero length
    err = fd.Truncate(0)
    if err != nil {
        out := fmt.Sprintf("Error truncating file: %v", err)
        helpers.LogInfo(out)
        return err
    }
    return nil
}

func send_csv_log(csv string, host string, user string, cnt int, file_path string) error {
    
    // defer wg.Done()
    
    cmds, err :=  os.ReadFile(csv)
    if err != nil {
        out := fmt.Sprintf("Error reading CSV logs: %v", err)
        helpers.LogInfo(out)
        delete_csv(csv)
        return err
    }
  
    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    timestamp, err := ptypes.TimestampProto(time.Now())
    if err != nil {
        out := fmt.Sprintf("Error creating timestamp: %v", err)
        helpers.LogInfo(out)
        delete_csv(csv)
        return err
    }

    to_send := bridge.CmdSnapShot{
        SubmissionNumber : 0,
        TimeStamp : timestamp,
        Cmds : string(cmds), 
        Host: host,
        Count: int32(cnt),
        User: user,
        DevID: helpers.GetID(),
    }

    _, e := client.IngestCmdSnapShot(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in Ingest CMD snapshot: %v\n", e)
        helpers.LogInfo(out)
        delete_csv(csv)
        return e
    }

    delete_csv(csv)

    // Empty the logs (as long as the data got properly sent & everything)
    empty_cmd_log(file_path)

    return nil
}

func delete_csv(csv string) {
    err := os.Remove(csv)
    if err != nil { 
        out := fmt.Sprintf("Error deleting csv: %v\n", err)
        helpers.LogInfo(out)
    }
}
