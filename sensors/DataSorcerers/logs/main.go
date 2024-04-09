package main;

import (
    "errors"
    "fmt";
    "bytes";
    "os/exec";
    "io/ioutil"
    "context"
    "time"
    "google.golang.org/grpc" 
    "github.com/golang/protobuf/ptypes"
    "strings"

    bridge "FusionBridge/log"

    "DataSorcerers/helpers"
)


var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var err error; 
var client bridge.LogClient;

var log_lines_captured int;


// This is necessary cause /var/log/dmesg doesn't always exist
func capture_dmesg() error {

    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    // Try to read dmesg file
    err := send_data("/var/log/dmesg")
    if err == nil { return nil }

    helpers.LogInfo("Capturing dmesg from CLI")
    // If reading dmesg file doesn't work, use CLI
        // and capture output
    cmd := exec.Command("dmesg")
    var outb, errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb
    err = cmd.Run()

    // If there is an error then kill
    if len(errb.String()) != 0 {
        helpers.LogInfo("Error in reading dmesg ftom cli");
        out := fmt.Sprintf("Error: %v", errb.String())
        helpers.LogInfo(out)
        return errors.New("capture_dmesg failed, cant read from CLI")
    }

    output_array := strings.Split(outb.String(), "\n")
    save_array := output_array[(len(output_array) - 1 - log_lines_captured) : (len(output_array) - 1)]

    timestamp, err := ptypes.TimestampProto(time.Now())

    to_send := bridge.LogData{
        SubmissionNumber : 0,
        TimeStamp : timestamp,
        Location: "dmesg-cli", 
        Content: strings.Join(save_array, "\n"),
        DevID: helpers.GetID(),
    }

    _, e := client.SaveLog(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in Save Log Data: %v", e)
        helpers.LogInfo(out)
        return e
    }
    
    out := fmt.Sprintf("Captured dmesg")
    helpers.LogInfo(out)

    return nil
}


// Async this so we can process faster
func send_data(full_path string) error {
 
    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    data, err := ioutil.ReadFile(full_path) 

    if err != nil {
        out := fmt.Sprintf("Error while reading file: %v; %v", full_path, err)
        helpers.LogInfo(out)
        return err;
    }

    lines := strings.Split(string(data), "\n")

    index := len(lines) - 25
    if len(lines) < 25 { index = len(lines) }

    timestamp, err := ptypes.TimestampProto(time.Now())

    to_send := bridge.LogData{
        SubmissionNumber : 0,
        TimeStamp : timestamp,
        Location: full_path, 
        Content: strings.Join(lines[index:], "\n"),
        DevID: helpers.GetID(),
    }

    _, e := client.SaveLog(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in Save Log Data: %v", e)
        helpers.LogInfo(out)
        return e
    }

    return nil
}

func main() {

    helpers.SetFileName("Logs")
    helpers.LoadConfig()
    if !helpers.Config.RunLog { return }

    helpers.LogInfo("Hello world from log scraper");

    log_lines_captured = helpers.Config.LogLinesCaptured;
   
    conn = helpers.CreateConnection()
    client = bridge.NewLogClient(conn); 
    
    
    for {
        // Record dmesg stuff (HAS to be done separately)
        capture_dmesg()
        // Read the log files
        for _, log := range helpers.Config.LogFiles {
            send_data(log)
        }

        helpers.LogInfo("Sending has completed")
        
        time.Sleep(time.Duration(helpers.Config.LogSweepInterval) * time.Second)

    }
}
