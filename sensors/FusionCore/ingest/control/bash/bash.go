package bash;

import "fmt"
import "sync"
import "context"
import "github.com/influxdata/influxdb-client-go/v2"
import "github.com/influxdata/influxdb-client-go/v2/api"

import bridge "FusionBridge/control/bash"
import config "FusionCore/config"
import Log "FusionCore/log"


type BashServer struct {
    bridge.UnimplementedBashServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *BashServer {
    s := &BashServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *BashServer) IngestCmdSnapShot(ctx context.Context, 
    MSG *bridge.CmdSnapShot) (*bridge.BashACK, error) {


    // Write the read data to the writeAPI
    tags := map[string]string{
    }
    fields := map[string]interface{}{
            "Host": MSG.Host,
            "User":MSG.User,
            "Count":MSG.Count,
            "Cmds":MSG.Cmds,
            "DevID":MSG.DevID,
    }

    point := influxdb2.NewPoint("bash", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in bash.IngestCmdSnapShot: %v", err))
        return &bridge.BashACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }

    // Could add more error codes for client health and such
    //     Will definitiely add error for write issues
    return &bridge.BashACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

