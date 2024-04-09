package main;

import (
    "fmt"
    "os"
    "context"
    "time"
    "net"
    "google.golang.org/grpc" 
    "github.com/golang/protobuf/ptypes"

    bridge "FusionBridge/metadata/id"

    "DataSorcerers/helpers"
)


var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var client bridge.IDClient;
var err error; 
var UID string;


// Async this so we can process faster
func send_data(interfaces []*bridge.InterfaceEntry) error {

    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    timestamp, err := ptypes.TimestampProto(time.Now())
    if err != nil {
        out := fmt.Sprintf("Error in client.id.send_data. Could not create timestamp: %v", err);
        helpers.FatalError(out)
        return err
    }
    to_send := bridge.WhoAmI{
        SubmissionNumber : 0,
        TimeStamp : timestamp,
        Interfaces : interfaces,
        DevID: helpers.GetID(),
    }

    _, e := client.SaveMyID(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in client.id.send_data: %v", e)
        helpers.LogInfo(out)
        return e
    }
    return nil
}

func scrapeIDInfo() ([]*bridge.InterfaceEntry, error) {
    // Get list of network interfaces
    interfaces, err := net.Interfaces()
    if err != nil { return nil, err }

    // var ret []InterfaceEntry;
    var ret []*bridge.InterfaceEntry;
    
    for _, iface := range interfaces {

        entry := &bridge.InterfaceEntry{};

        entry.MAC = iface.HardwareAddr.String()
        entry.Name = iface.Name

        addrs, err := iface.Addrs()
        if err != nil {
            out := fmt.Sprintf("Error in client.id.scrapeIDInfo iface.Addrs(): %v", err)
            helpers.FatalError(out)
            return nil, err
        }
        
        for _, addr := range addrs {
            ip, _, _ := net.ParseCIDR(addr.String())
            hostnames, err := net.LookupAddr(ip.String())
            if err != nil {
                out := fmt.Sprintf("Error in client.id.scrapeIDInfo net.LookupAddr(): %v", err)
                helpers.FatalError(out)
                return nil, err
            }
            interfaceDetails := &bridge.NetInfo{};
            interfaceDetails.IP = ip.String();
            interfaceDetails.HostNames = hostnames;
            entry.NetInfo = append(entry.NetInfo, interfaceDetails);
        }
        ret = append(ret, entry)
    }
    return ret, nil
}

func ScrapeUserNames() {

    // Open the directory
    dir, err := os.Open("/home")
    if err != nil {
        out := fmt.Sprintf("Error opening directory: %v", err)
        helpers.LogInfo(out)
        return
    }
    defer dir.Close()

    // Read the directory entries
    entries, err := dir.Readdir(-1)
    if err != nil {
        out := fmt.Sprintf("Error reading directory:", err)
        helpers.LogInfo(out)
        return
    }

    // Iterate over the entries and filter for directories
    for _, entry := range entries {
        if entry.IsDir() {
            fmt.Println(entry.Name())
        }
    }
}


func main() {
    helpers.SetFileName("Id")

    helpers.LoadConfig()
    if !helpers.Config.RunId { return }

    helpers.LogInfo("Hello world from ID scraper");

    conn = helpers.CreateConnection()
    client = bridge.NewIDClient(conn); 


    go func() {
        ScrapeUserNames()
        time.Sleep(time.Duration(helpers.Config.ScrapeForUsers) * time.Second)
    }()
    
    for {
        if interfaces, err := scrapeIDInfo(); err == nil {
            if err = send_data(interfaces); err != nil {
                out := fmt.Sprintf("Error in id.send_data: %v", err);
                helpers.LogInfo(out)
            }
        } else {
            out := fmt.Sprintf("Error in id.scrapeIDInfo: %v", err);
            helpers.LogInfo(out)
        }
        helpers.LogInfo("Sending has completed")
        time.Sleep(time.Duration(helpers.Config.ScrapeInterfaceInfo) * time.Second)
    }
}

