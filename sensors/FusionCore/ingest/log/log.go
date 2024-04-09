package log;

import (
    "fmt"
    "sync"
    "context"
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"


    bridge "FusionBridge/log"
    config "FusionCore/config"
    Log "FusionCore/log"
    "FusionCore/databases/postgres"
    "FusionCore/helpers"
)


type LogServer struct {
    bridge.UnimplementedLogServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}

func NewServer(client influxdb2.Client) *LogServer {
    s := &LogServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}

func (s *LogServer) SaveLog(ctx context.Context, MSG *bridge.LogData) (*bridge.LogACK, error) {

    // Create a text hash for easy lookup and indexing constraints
    hash := helpers.GetUniqueHash(MSG.Content)
    // Save the text. If its a conflict, then get the ID 
    res, err := postgres.Connection.Query(context.Background(), `
        INSERT INTO Logs
            (HASH, CONTENTS)
        VALUES
            ($1, (convert_to($2, 'UTF8')::BYTEA))
        ON CONFLICT (HASH)
        DO UPDATE SET HASH = $1
        `, hash, MSG.Content)
    res.Close()
    if err != nil {
        tmp := fmt.Sprintf("Failed to query postgres database: %v", err)
        Log.FatalError(tmp)
    }

    // Write the read data to the writeAPI
    tags := map[string]string{
        "Location":MSG.Location,
        "DevID":MSG.DevID,
    }
    // Content is optional but protobufs have auto empty values so ok 
    fields := map[string]interface{}{
        "ContentHash": hash, // To keep a reference without storing files unnecessarily
    }

    point := influxdb2.NewPoint("log", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err = s.writeAPI.WritePoint(context.Background(), point)
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in log.SaveLog: %v", err))
        return &bridge.LogACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }

    // Could add more error codes for client health and such
    //     Will definitiely add error for write issues
    return &bridge.LogACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

