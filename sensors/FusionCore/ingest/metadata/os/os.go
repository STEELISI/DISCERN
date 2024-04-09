package os;

import "fmt"
import "sync"
import "context"
import "github.com/influxdata/influxdb-client-go/v2"
import "github.com/influxdata/influxdb-client-go/v2/api"

import bridge "FusionBridge/metadata/os"
import Log "FusionCore/log"
import "FusionCore/config"


type OSServer struct {
    bridge.UnimplementedOSServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *OSServer {
    s := &OSServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *OSServer) MarkCloseSysCall(ctx context.Context, 
    MSG *bridge.CloseCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,           
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Count":MSG.Count,
           "Range":false, // Place Close and Close range into same index
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("SysClose", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in network.LogNetworkActivity: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkCloseRangeSysCall(ctx context.Context, 
    MSG *bridge.CloseRangeCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,           
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Count":MSG.Count,
           "Range":true, // Place Close and Close range into same index
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("SysClose", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.MarkCloseRangeSysCall: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkExecveSysCall(ctx context.Context, 
    MSG *bridge.ExecveCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,           
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Arg":MSG.Arg,
           "ArgNum":MSG.ArgNum,
           "At":false, // Place Close and Close range into same index
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Execve", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.MarkExecveSysCall: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkExecveAtSysCall(ctx context.Context, 
    MSG *bridge.ExecveAtCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,           
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Arg":MSG.Arg,
           "ArgNum":MSG.ArgNum,
           "At":true, // Place Close and Close range into same index
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Execve", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.MarkExecveAtSysCall: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkForkSysCall(ctx context.Context, 
    MSG *bridge.ForkCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,           
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Count":MSG.Count,
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Fork", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.Fork: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkKillSysCall(ctx context.Context, 
    MSG *bridge.KillCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,           
           "ArgPid":MSG.ArgPid,           
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Sig":MSG.Sig,
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Kill", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.Kill: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkOpenSysCall(ctx context.Context, 
    MSG *bridge.OpenCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Filename":MSG.Filename,
           "Count":MSG.Count,
           "Version":"",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Open", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.Open: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkOpenAtSysCall(ctx context.Context, 
    MSG *bridge.OpenAtCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Filename":MSG.Filename,
           "Count":MSG.Count,
           "Version":"At",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Open", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.OpenAt: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkOpenAt2SysCall(ctx context.Context, 
    MSG *bridge.OpenAt2Call) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Filename":MSG.Filename,
           "Count":MSG.Count,
           "Version":"At2",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Open", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.OpenAt2: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkRecvFromSysCall(ctx context.Context, 
    MSG *bridge.RecvFromCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "FileDescriptor":MSG.FileDescriptor,
           "Count":MSG.Count,
           "Version":"from",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Recv", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.RecvFrom: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkRecvMMsgSysCall(ctx context.Context, 
    MSG *bridge.RecvMMsgCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "FileDescriptor":MSG.FileDescriptor,
           "Count":MSG.Count,
           "Version":"mmsg",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Recv", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.RecvMMsg: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkRecvMsgSysCall(ctx context.Context, 
    MSG *bridge.RecvMsgCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "FileDescriptor":MSG.FileDescriptor,
           "Count":MSG.Count,
           "Version":"msg",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Recv", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.RecvMMsg: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkSendMMsgSysCall(ctx context.Context, 
    MSG *bridge.SendMMsgCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "FileDescriptor":MSG.FileDescriptor,
           "Len":MSG.Len,
           "Count":MSG.Count,
           "Version":"mmsg",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Send", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.SendMMsg: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkSendMsgSysCall(ctx context.Context, 
    MSG *bridge.SendMsgCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "FileDescriptor":MSG.FileDescriptor,
           "Len":MSG.Len,
           "Count":MSG.Count,
           "Version":"msg",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Send", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.SendMsg: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkSendToSysCall(ctx context.Context, 
    MSG *bridge.SendToCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "FileDescriptor":MSG.FileDescriptor,
           "Len":MSG.Len,
           "Count":MSG.Count,
           "Version":"to",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Send", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.SendTo: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkSocketSysCall(ctx context.Context, 
    MSG *bridge.SocketCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Family":MSG.Family,
           "Type":MSG.Type,
           "Protocol":MSG.Protocol,
           "Count":MSG.Count,
           "Version":"",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Socket", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.Socket: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkSocketPairSysCall(ctx context.Context, 
    MSG *bridge.SocketPairCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Family":MSG.Family,
           "Type":MSG.Type,
           "Protocol":MSG.Protocol,
           "Count":MSG.Count,
           "Version":"pair",
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("Socket", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.Socket: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkSysInfoSysCall(ctx context.Context, 
    MSG *bridge.SysInfoCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("SysInfo", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.SysInfo: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkTKillSysCall(ctx context.Context, 
    MSG *bridge.TKillCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "ArgPid":MSG.ArgPid,
           "Sig":MSG.Sig,
           "Count":MSG.Count,
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("TKill", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.TKill: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}


func (s *OSServer) MarkVForkSysCall(ctx context.Context, 
    MSG *bridge.VForkCall) (*bridge.OS_ACK, error) {

    tags := map[string]string{

    }
    fields := map[string]interface{}{
           "Pid":MSG.Pid,
           "Uid":MSG.Uid,
           "Gid":MSG.Gid,
           "Count":MSG.Count,
           "DevID":MSG.DevID,
    }


    point := influxdb2.NewPoint("vFork", tags, fields, 
        MSG.TimeStamp.AsTime())
    
    err := s.writeAPI.WritePoint(context.Background(), point)
    
    if err != nil {
        Log.LogInfo(fmt.Sprintf("Error in os.TKill: %v", err))
        return &bridge.OS_ACK{
            Type: 1, 
            SubmissionNumber: MSG.SubmissionNumber,
        }, err
    }
 
    return &bridge.OS_ACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

