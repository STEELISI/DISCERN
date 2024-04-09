package main;

import "fmt";
import "path/filepath"
import "os"
import "io/ioutil"
import "regexp"
import "context"
import "time"
import "google.golang.org/grpc" 
// import "google.golang.org/grpc/credentials/insecure"
import "github.com/golang/protobuf/ptypes"


import bridge "FusionBridge/control/jupyter"
import "DataSorcerers/helpers"


var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var err error; 
var client bridge.JupyterClient;
// var wg sync.WaitGroup;


// Async this so we can process faster
func send_data(full_path string) {
 
    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()
    
    data, err := ioutil.ReadFile(full_path) 

    if err != nil {
        out := fmt.Sprintf("Error while reading file: %v", full_path)
        helpers.LogInfo(out)
        out =  fmt.Sprintf("Error: %v", err)
        helpers.LogInfo(out)
    }
    timestamp, err := ptypes.TimestampProto(time.Now())
    to_send := bridge.IPYNB_Submission{
        SubmissionNumber : 0,
        TimeStamp : timestamp,
        FileLocation : full_path, 
        FileContents : string(data),
        DevID: helpers.GetID(),
    }

    _, e := client.IngestIPYNB(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in Ingest IPYNB: %v", e)
        helpers.LogInfo(out)
    }
}


func walkFn(path string, fd os.FileInfo, err error) error {
    if err != nil {
        out := fmt.Sprintf("Error Crashed jupyter: %v", err)
        helpers.FatalError(out)
        return err
    }

    if fd.IsDir() { return nil }

    res, err := regexp.MatchString(`.*\.ipynb$`, fd.Name())
    if err != nil {
        out := fmt.Sprintf("Error Crashed jupyter: %v", err)
        helpers.FatalError(out)
        return err
    }
    if res {
        go send_data(path);
    }
    return nil
}



func main() {
    helpers.SetFileName("Jupyter")
    
    helpers.LoadConfig()
    if !helpers.Config.RunJupyter { return }

    helpers.LogInfo("Hello world from jupyter scraper");

    conn = helpers.CreateConnection()
    client = bridge.NewJupyterClient(conn); 


    for {
        for _, path := range helpers.Config.JupyterPaths {
            if err := filepath.Walk(path, walkFn); err != nil {
                out := fmt.Sprintf("Error in Jupyter Scraping Walk function: %v", err);
                helpers.LogInfo(out)
            }
        }

        helpers.LogInfo("File search stopped")

        helpers.LogInfo("Sending has completed")

        time.Sleep(time.Duration(helpers.Config.JupyterSweepInterval) * time.Second)

    }
}

