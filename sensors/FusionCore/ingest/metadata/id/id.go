package id;

import "fmt"
import "sync"
import "sort"
import "context"
import "github.com/influxdata/influxdb-client-go/v2"
import "github.com/influxdata/influxdb-client-go/v2/api"

import bridge "FusionBridge/metadata/id"
import Log "FusionCore/log"
import "FusionCore/config"


type IDServer struct {
    bridge.UnimplementedIDServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *IDServer {
    s := &IDServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *IDServer) SaveMyID(ctx context.Context, 
    MSG *bridge.WhoAmI) (*bridge.YouIs, error) {


    for _, interf := range MSG.Interfaces {

        macAddr := interf.MAC
        name := interf.Name

        for _, info := range interf.NetInfo {

            ipAddr := info.IP;
            hostnames := info.HostNames;
            sort.Strings(hostnames)           
            // Write the read data to the writeAPI
            tags := map[string]string{

            }
            fields := map[string]interface{}{
                    "MAC": macAddr,
                    "IP": ipAddr,
                    "HostName": hostnames,
                    "Name":name,
                    "DevID":MSG.DevID,
            }

            point := influxdb2.NewPoint("id", tags, fields, 
                MSG.TimeStamp.AsTime())
            
            err := s.writeAPI.WritePoint(context.Background(), point)
            
            if err != nil {
                Log.LogInfo(fmt.Sprintf("Error in id.SaveMyID: %v", err))
                return &bridge.YouIs{
                    Type: 1, 
                    SubmissionNumber: MSG.SubmissionNumber,
                }, err
            }
        }

    }
    return &bridge.YouIs{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

