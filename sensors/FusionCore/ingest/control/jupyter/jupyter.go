package jupyter;

import (
    "fmt"
    "sync"
    "context"
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"

    bridge "FusionBridge/control/jupyter"
    config "FusionCore/config"
    Log "FusionCore/log"
    "FusionCore/databases/postgres"
    "FusionCore/helpers"
)


type JupyterServer struct {
    bridge.UnimplementedJupyterServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *JupyterServer {
    s := &JupyterServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *JupyterServer) IngestIPYNB(ctx context.Context, 
    MSG *bridge.IPYNB_Submission) (*bridge.Response, error) {

    // Create a text hash for easy lookup and indexing constraints
    hash := helpers.GetUniqueHash(MSG.FileContents)
    // Store the hash and text as serial
    res, err := postgres.Connection.Query(context.Background(), `
        INSERT INTO JupyterNotebook
            (HASH, CONTENTS)
        VALUES
            ($1, (convert_to($2, 'UTF8')::BYTEA))
        ON CONFLICT (HASH)
        DO UPDATE SET HASH = $1
        `, hash, MSG.FileContents)
    res.Close()
    if err != nil {
        tmp := fmt.Sprintf("Failed to query postgres database: %v", err)
        Log.FatalError(tmp)
    }

    // Write the read data to the writeAPI
    tags := map[string]string{
        "DevID":MSG.DevID,
    }
    fields := map[string]interface{}{
        "ContentHash": hash, // To keep a reference without storing files unnecessarily
        "FileLocation":MSG.FileLocation,
    }

    point := influxdb2.NewPoint("jupyter", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err = s.writeAPI.WritePoint(context.Background(), point)
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in jupyter.IngestIPYBN: %v", err))
        return &bridge.Response{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
    // Could add more error codes for client health and such
    //     Will definitiely add error for write issues
    return &bridge.Response{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


