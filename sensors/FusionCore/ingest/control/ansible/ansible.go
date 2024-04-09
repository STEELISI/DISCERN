package ansible;

import (
    "fmt"
    "sync"
    "context"
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"

    "FusionCore/config"
    bridge "FusionBridge/control/ansible"
    Log "FusionCore/log"
    "FusionCore/databases/postgres"
    "FusionCore/helpers"
)

type AnsibleServer struct {
    bridge.UnimplementedAnsibleServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *AnsibleServer {
    s := &AnsibleServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *AnsibleServer) SaveAnsibleConfig(ctx context.Context, 
    MSG *bridge.ConfigDetails) (*bridge.AnsibleACK, error) {

    // Create a text hash for easy lookup and indexing constraints
    hash := helpers.GetUniqueHash(MSG.Content)
    
    // Save the text. If its a conflict, then get the ID 
    res, err := postgres.Connection.Query(context.Background(), `
        INSERT INTO AnsibleConfig
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

    tags := map[string]string{
            "DevID":MSG.DevID,
            "Location": MSG.Location,
    }
    fields := map[string]interface{}{
            "ContentHash": hash, // To keep a reference without storing files unnecessarily
    }



    point := influxdb2.NewPoint("ansible-config", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err = s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in id.SaveAnsibleConfig: %v", err))
        return &bridge.AnsibleACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }

    return &bridge.AnsibleACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}



func (s *AnsibleServer) SaveAnsiblePlaybook(ctx context.Context, 
    MSG *bridge.PlaybookDetails) (*bridge.AnsibleACK, error) {

    // Create a text hash for easy lookup and indexing constraints
    hash := helpers.GetUniqueHash(MSG.Content)
    // Save the text. If its a conflict, then get the ID 
    res, err := postgres.Connection.Query(context.Background(), `
        INSERT INTO AnsiblePlaybook
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


    tags := map[string]string{
        "Location": MSG.Location,
        "DevID":MSG.DevID,
    }
    fields := map[string]interface{}{
        "ContentHash": hash, // To keep a reference without storing files unnecessarily
        "Verified": MSG.Verified,
    }

    point := influxdb2.NewPoint("ansible-playbooks", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err = s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in id.SaveAnsiblePlaybook: %v", err))
        return &bridge.AnsibleACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }

    return &bridge.AnsibleACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

