package file;

import (
    "fmt"
    "time"
    "sync"
    "context"
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"

    bridge "FusionBridge/metadata/file"
    config "FusionCore/config"
    Log    "FusionCore/log"
)


type FileServer struct {
    bridge.UnimplementedFileServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *FileServer {
    s := &FileServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *FileServer) SaveFsEvent(ctx context.Context, 
    MSG *bridge.FsEvent) (*bridge.FsEventACK, error) {

    // Write the read data to the writeAPI
    tags := map[string]string{
    }
    // Content is optional but protobufs have auto empty values so ok 
    fields := map[string]interface{}{
            "Op": MSG.Op,
            "Location":MSG.Location,
            "Content":MSG.Content,
            "Permissions":MSG.Permissions,
            "Owner":MSG.Owner,
            "Group":MSG.Group,
            "DevID":MSG.DevID,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()


    point := influxdb2.NewPoint("file", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(ctx, point)
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in file.SaveFsEvent: %v", err))
        return &bridge.FsEventACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }

    // Could add more error codes for client health and such
    //     Will definitiely add error for write issues
    return &bridge.FsEventACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

