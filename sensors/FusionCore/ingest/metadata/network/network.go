package network;

import "fmt"
import "sync"
import "context"
import "github.com/influxdata/influxdb-client-go/v2"
import "github.com/influxdata/influxdb-client-go/v2/api"

import bridge "FusionBridge/metadata/network"
import Log "FusionCore/log"
import "FusionCore/config"


type NetworkServer struct {
    bridge.UnimplementedNetworkServer
    mu sync.Mutex
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}


func NewServer(client influxdb2.Client) *NetworkServer {
    s := &NetworkServer{
        writeAPI: client.WriteAPIBlocking(config.ORG, config.BUCKET_NAME),
        queryAPI: client.QueryAPI(config.ORG),
    }
    return s
}


func (s *NetworkServer) LogNetworkActivity(ctx context.Context, 
    MSG *bridge.NetworkSlice) (*bridge.NetworkACK, error) {

    // I've decided against saving the timestamp sent with the 
        // NetworkSlice itself but am keeping it there in case I 
        // change my mind

    for _, pkt := range MSG.Packets {

        // Write the read data to the writeAPI
        tags := map[string]string{

        }
        fields := map[string]interface{}{
                "Dev":pkt.Dev,
                "LinkProtocol": pkt.LinkProtocol,
                "NetworkProtocol": pkt.NetworkProtocol,
                "TransportProtocol": pkt.TransportProtocol,
                "ApplicationProtocol":pkt.ApplicationProtocol,
                "DevID":MSG.DevID,
        }

        // Ingest Link layer into into fields
        if pkt.LinkProtocol == "ARP" {
            fields["SrcHwAddy"]  = pkt.ARP.SrcHwAddy
            fields["SrcProtAdd"] = pkt.ARP.SrcProtAdd
            fields["DstHwAddy"]  = pkt.ARP.DstHwAddy
            fields["DstProtAdd"] = pkt.ARP.DstProtAdd
            fields["Operation"]  = pkt.ARP.Operation
            fields["Protocol"]   = pkt.ARP.Protocol
        } else if pkt.LinkProtocol == "Ethernet" {
            fields["SRC_MAC"]  = pkt.ETH.SRC_MAC
            fields["DST_MAC"] = pkt.ETH.SRC_MAC
            fields["Length"]  = pkt.ETH.Length
        }
        // Ingest Network layer into into fields
        if pkt.NetworkProtocol == "IP" {
            fields["SRC_IP"]  = pkt.IP.SRC_IP
            fields["DST_IP"] = pkt.IP.SRC_IP
            fields["V4"] = pkt.IP.V4
        } else if pkt.NetworkProtocol == "ICMP" {
            // Just here in case
        }
        // Ingest Transport layer into into fields
        if pkt.TransportProtocol == "TCP" {
            fields["SrcPort"] = pkt.TCP.SrcPort
            fields["DstPort"] = pkt.TCP.DstPort
        } else if pkt.TransportProtocol == "UDP" {
            fields["SrcPort"] = pkt.UDP.SrcPort
            fields["DstPort"] = pkt.UDP.DstPort
        }
        // Ingest Application layer into into fields
        if pkt.TransportProtocol == "DNS" {
            fields["Questions"] = pkt.DNS.Questions
            fields["Answers"] = pkt.DNS.Answers
        } else if pkt.TransportProtocol == "TLS" {
            // Here for posterity
        }

        point := influxdb2.NewPoint("network", tags, fields, 
            MSG.TimeStamp.AsTime())
        
        err := s.writeAPI.WritePoint(context.Background(), point)
        
        if err != nil {
            Log.LogInfo(fmt.Sprintf("Error in network.LogNetworkActivity: %v", err))
            return &bridge.NetworkACK{
                Type: 1, 
                SubmissionNumber: MSG.SubmissionNumber,
            }, err
        }
    }
    return &bridge.NetworkACK{
        Type: 0, 
        SubmissionNumber: MSG.SubmissionNumber,
    }, nil
}

